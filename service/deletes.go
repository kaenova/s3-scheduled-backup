package service

import "github.com/kaenova/s3-scheduled-backup/pkg"

type DeleteService struct {
	S3      pkg.S3ObjectI
	folders []string
	Log     pkg.CustomLoggerI
}

func NewDeleteService() {

}
