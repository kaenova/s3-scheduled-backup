package service

import "github.com/kaenova/s3-scheduled-backup/pkg"

type BackupService struct {
	s3 pkg.S3ObjectI
}

type BackupServiceI interface {
}

func NewBackupService(s3 pkg.S3ObjectI) BackupServiceI {
	return &BackupService{}
}
