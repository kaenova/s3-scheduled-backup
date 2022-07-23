package service

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kaenova/s3-scheduled-backup/pkg"
)

// Do every midnight
const CRON_MIDNIGHT = "59 23 * * *"

// Do every minute
const CRON_MINUTE = "*/1 * * * *"

// Do every 5 minute
const CRON_5_MINUTE = "*/5 * * * *"

type BackupService struct {
	Path    string
	S3      pkg.S3ObjectI
	folders []string
	Log     pkg.CustomLoggerI
}

type BackupServiceI interface {
	StartBlocking()
}

func NewBackupService(path string, s3 pkg.S3ObjectI, log pkg.CustomLoggerI) BackupServiceI {
	folders, err := pkg.FoldersOneLevel(path)
	if err != nil {
		log.Error(err.Error())
	}

	// Clear terminal screen
	fmt.Println("\033[2J")

	return &BackupService{
		Path:    path,
		S3:      s3,
		folders: folders,
		Log:     log,
	}
}

func (b *BackupService) StartBlocking() {
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		b.Log.Error("Fail to load timezone")
		os.Exit(1)
	}

	s := gocron.NewScheduler(tz)

	s.Cron(CRON_MIDNIGHT).Do(b.backup)
	s.Cron(CRON_5_MINUTE).Do(func() {
		b.Log.Info("Health check: Normal")
	})

	b.Log.Info("Schedule started")
	s.StartBlocking()
}

func (b *BackupService) backup() {
	b.Log.Info("Backuping")
	for _, folder := range b.folders {
		go b.backupSingleFolder(folder)
	}
}

func (b *BackupService) backupSingleFolder(folder string) {
	// <foldername>--<Year-Month-Date>
	currentTime := time.Now()
	fileName := fmt.Sprintf("%s--%s", folder, currentTime.Format("2006-01-02"))

	b.Log.Info("Backuping folder " + folder)
	sourceFolderPath := b.Path + "/" + folder
	zipFolderPath := "./temp/" + fileName + ".zip"
	err := pkg.ZipSource(sourceFolderPath, zipFolderPath)
	if err != nil {
		b.Log.Error("Cannot Backup Folder " + folder + " " + err.Error())
		return
	}

	// Zip temp cleanup
	defer func() {
		os.Remove(zipFolderPath)
	}()

	_, err = b.S3.UploadFileFromPathNamed(fileName, zipFolderPath)
	if err != nil {
		b.Log.Error("Fail to upload " + folder)
		return
	}
	b.Log.Info("Success upload folder " + folder)
}
