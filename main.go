// A program that makes a background task for backup a folder inside a path
// To S3 Object

package main

import (
	"github.com/kaenova/s3-scheduled-backup/config"
	"github.com/kaenova/s3-scheduled-backup/pkg"
	"github.com/kaenova/s3-scheduled-backup/service"
)

func main() {
	// Init Logging
	log := pkg.NewLogger()

	config := config.MakeConfig(log)

	s3, err := pkg.NewS3Object(config.S3Config.Endpoint, config.S3Config.AccessKey,
		config.S3Config.SecretKey, config.S3Config.BucketName, config.S3Config.UseSSL)
	if err != nil {
		log.Fatal(err.Error())
	}

	backuper := service.NewBackupService(config.BackupConfig.Path, config.BackupConfig.MaxWindow, s3, log)
	backuper.StartBlocking()
}
