package service

import (
	"os"
	"sort"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kaenova/s3-scheduled-backup/pkg"
)

type BackupService struct {
	Cron           string
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

func NewBackupService(cronSchedule string, path string, maxRollbackDay int, s3 pkg.S3ObjectI, log pkg.CustomLoggerI) BackupServiceI {
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
		Cron:           cronSchedule,
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
	scheduler.Cron(b.Cron).Do(b.backup)
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

// NOTE: [WARNING] The file naming convention need to be string sortable by time!
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
	// Itterate all backed up folder and get file that has same FileName
	objStr := b.S3.ListObjectParentDir()
	sameFolderName := []string{}
	for _, val := range objStr {
		obj, err := pkg.ParseBackupFolder(val)
		if err != nil {
			b.Log.Warning("Cannot parse", val, "with an error", err.Error())
			continue
		}
		if obj.FolderName == folder {
			sameFolderName = append(sameFolderName, obj.ZipFileName)
		}
	}

	// Check if we need to delete file
	if len(sameFolderName) <= b.MaxRollbackDay {
		b.Log.Info("skipping deleting older backup for", folder)
		return
	}

	// Sort string with an assumption of the filename are time sortable
	sort.Strings(sameFolderName)

	// Get objects to be deleted
	deletedObjects := pkg.GetFirstKString(len(sameFolderName)-b.MaxRollbackDay, sameFolderName)

	// Delete the objects
	for _, v := range deletedObjects {
		err := b.S3.DeleteObject(v)
		if err != nil {
			b.Log.Error("cannot delete object with an error:", err.Error())
			return
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
		err := os.Remove(tempPath)
		if err != nil {
			b.Log.Fatal("Cannot delete temporary file of " + tempPath)
		}
		b.Log.Info("Removing temporary file " + tempPath)
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
