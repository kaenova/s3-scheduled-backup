package service

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kaenova/s3-scheduled-backup/pkg"
)

type BackupService struct {
	Path           string
	MaxRollbackDay int
	folders        []string
	S3             pkg.S3ObjectI
	Log            pkg.CustomLoggerI
}

type BackupServiceI interface {
	StartBlocking()
	RegisterScheduler(scheduler *gocron.Scheduler)
}

func NewBackupService(path string, maxRollbackDay int, s3 pkg.S3ObjectI, log pkg.CustomLoggerI) BackupServiceI {
	folders, err := pkg.FoldersOneLevel(path)
	if err != nil {
		log.Error(err.Error())
	}

	return &BackupService{
		Path:           path,
		MaxRollbackDay: maxRollbackDay,
		S3:             s3,
		folders:        folders,
		Log:            log,
	}
}

func (b *BackupService) StartBlocking() {
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		b.Log.Fatal("Fail to load timezone")
	}

	s := gocron.NewScheduler(tz)

	b.RegisterScheduler(s)

	b.Log.Info("Schedule started")
	s.StartBlocking()
}

func (b *BackupService) RegisterScheduler(scheduler *gocron.Scheduler) {
	scheduler.Cron(pkg.CRON_MINUTE).Do(b.backup)
	scheduler.Cron(pkg.CRON_MINUTE).Do(func() {
		b.Log.Info("Backup service is healthy")
	})
}

func (b *BackupService) backup() {
	b.Log.Info("Starting backup process")
	for _, folder := range b.folders {
		go b.backupSingleFolder(folder)
	}
}

func (b *BackupService) backupSingleFolder(folder string) {
	currentTime := time.Now()

	if b.MaxRollbackDay != 0 {
		b.Log.Info("Will backup folder", folder, "and delete older backup before",
			currentTime.Add(-time.Hour*24*time.Duration(b.MaxRollbackDay)).Format(pkg.TIME_FORMAT))
	} else {
		b.Log.Info("Will backup folder", folder)
	}

	err := b.backupFolder(folder, currentTime)
	if err != nil {
		b.Log.Warning("Won't delete old backup of folder " + folder)
		return
	}

	b.deleteOldBackup(folder, currentTime)
}

func (b *BackupService) deleteOldBackup(folder string, currentTime time.Time) {
	// Prepare maximum backed up folders
	maxTime := currentTime.Add(-time.Hour * 24 * time.Duration(b.MaxRollbackDay))

	// Itterate all backed up folder
	objStr := b.S3.ListObjectParentDir()
	for _, val := range objStr {
		obj, err := pkg.ParseBackupFolder(val)

		if err != nil {
			b.Log.Error("Cannot parse", val)
			continue
		}

		// Delete backup if have the same name folder and the time is below max time
		if (obj.FolderName == folder) && (obj.Time.Unix() < maxTime.Unix()) {
			err := b.S3.DeleteObject(val)

			if err != nil {
				b.Log.Error("Cannot delete object", val)
				continue
			}
		}
	}

	b.Log.Info("Deleting old backup of", folder, "success")
}

func (b *BackupService) backupFolder(folder string, currentTime time.Time) error {
	// Creates an abstract for folder
	backupFile := pkg.CreateBackupFolder(folder, currentTime)

	b.Log.Info("Starting zipping folder " + backupFile.FolderName)
	sourceFolderPath := b.Path + "/" + backupFile.FolderName
	tempPath := "./temp/" + backupFile.ZipFileName

	// Zip folder
	err := pkg.ZipSource(sourceFolderPath, tempPath)
	if err != nil {
		b.Log.Error("Cannot zip folder " + backupFile.FolderName + " " + err.Error())
		return err
	}

	// Zip temp cleanup
	defer func() {
		os.Remove(tempPath)
	}()

	b.Log.Info("Success zip " + backupFile.ZipFileName + " and trying to upload")

	fileName := backupFile.GenerateFileName()
	_, err = b.S3.UploadFileFromPathNamed(fileName, tempPath)
	if err != nil {
		b.Log.Error("Fail to upload " + backupFile.FolderName)
		return err
	}
	b.Log.Info("Success upload folder " + backupFile.FolderName)
	return nil
}
