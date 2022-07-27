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
	err := b.backupFolder(folder)
	if err != nil {
		b.Log.Warning("Won't delete old backup of folder " + folder)
		return
	}
	b.deleteOldBackup(folder)
}

func (b *BackupService) deleteOldBackup(folder string) {

}

func (b *BackupService) backupFolder(folder string) error {
	// Creates an abstract for folder
	backupFile := pkg.CreateBackupFolder(folder)

	b.Log.Info("Backuping folder " + backupFile.FolderName)
	sourceFolderPath := b.Path + "/" + backupFile.FolderName
	tempPath := "./temp/" + backupFile.ZipFileName

	// Zip folder
	err := pkg.ZipSource(sourceFolderPath, tempPath)
	if err != nil {
		b.Log.Error("Cannot Backup Folder " + folder + " " + err.Error())
		return err
	}

	// Zip temp cleanup
	defer func() {
		os.Remove(tempPath)
	}()

	_, err = b.S3.UploadFileFromPathNamed(backupFile.FolderName+"--"+backupFile.Time.String(), tempPath)
	if err != nil {
		b.Log.Error("Fail to upload " + backupFile.FolderName)
		return err
	}
	b.Log.Info("Success upload folder " + backupFile.FolderName)
	return nil
}
