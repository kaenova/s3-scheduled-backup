// A program that makes a background task for backup a folder inside a path
// To S3 Object

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kaenova/s3-scheduled-backup/pkg"
	"github.com/kaenova/s3-scheduled-backup/service"
)

func main() {
	// Input var
	var pathToBackup, s3Endpoint, bucketName, accessKey, secretKey string
	var useSSL bool

	pathToBackup = pkg.InputString("Input a parent folder for the children folder to be backed up: ")
	folders, err := pkg.FoldersOneLevel(pathToBackup)
	if err != nil {
		log.Fatal(err.Error())
	}
	isContinue := pkg.InputBool(func() string {
		finalString := "This folder(s) will be backed up every 23:59\n"
		for i, v := range folders {
			finalString += fmt.Sprintf("%d. %s\n", i, v)
		}
		finalString += "Are you sure? (y/[n]) "
		return finalString
	}(), func(s string) bool {
		if s == "y" {
			return true
		}
		return false
	})

	if !isContinue {
		log.Fatal("User terminated the program")
		os.Exit(1)
	}

	s3Endpoint = pkg.InputString("Input S3 Enpoint (ex. is3.cloudhost.id): ")
	bucketName = pkg.InputString("Input Bucket Name: ")
	accessKey = pkg.InputString("Input AccessKey: ")
	secretKey = pkg.InputString("Input SecretKey: ")
	useSSL = pkg.InputBool("Do you want to use ssl? ([y]/n) : ", func(s string) bool {
		if s == "n" {
			return false
		}
		return true
	})

	log := pkg.NewLogger()
	s3, err := pkg.NewS3Object(s3Endpoint, accessKey, secretKey, bucketName, useSSL)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	backuper := service.NewBackupService(pathToBackup, s3, log)
	backuper.StartBlocking()
}
