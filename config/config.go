package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kaenova/s3-scheduled-backup/pkg"
)

type Config struct {
	ApplicationConfig
	S3Config
	BackupConfig
}

// Mode Constant
const (
	DOCKER_MODE = "docker"
	LOCAL_MODE  = "local"
)

// Default Docker Mode Path
const DOCKER_VOL_PATH = "/dockervol"

type ApplicationConfig struct {
	Mode string
}

type BackupConfig struct {
	// Maximum of days of an old backup can be keep. 0 will never delete old backup
	MaxWindow int

	// Path of the folders that the children will be backed up
	Path string

	// Cron Job
	Cron string

	// Exclude list
	ExcludeFolders []string
}

type S3Config struct {
	Endpoint   string
	BucketName string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
}

func MakeConfig(log pkg.CustomLoggerI) Config {
	// Check dotenv
	err := godotenv.Load()
	if err != nil {
		log.Warning("Cannot load .env file")
	}

	app := MakeApplicationConfig(log)
	s3 := MakeS3Config(log)

	backup := MakeBackupConfig(log, app)

	return Config{
		ApplicationConfig: app,
		BackupConfig:      backup,
		S3Config:          s3,
	}
}

func MakeApplicationConfig(log pkg.CustomLoggerI) ApplicationConfig {
	mode := LOCAL_MODE
	if os.Getenv("MODE") == DOCKER_MODE {
		mode = DOCKER_MODE
	}
	return ApplicationConfig{
		Mode: mode,
	}
}

func MakeBackupConfig(log pkg.CustomLoggerI, app ApplicationConfig) BackupConfig {
	window, err := strconv.Atoi(os.Getenv("MAXIMUM_BACKUP_WINDOW"))
	if err != nil || window == 0 {
		log.Warning("Maximum time window not detected, defaulted to 0. Which never deleted old backup")
		window = 0
	}

	var path string
	if app.Mode != DOCKER_MODE {
		path = os.Getenv("PATH_BACKUP")
		if path == "" {
			path = pkg.InputString("Input a parent folder for the children folder to be backed up: ")
		}
	} else {
		path = DOCKER_VOL_PATH
	}

	folders, err := pkg.FoldersOneLevel(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	cronTime := os.Getenv("CRON_SCHEDULE")
	if cronTime == "" {
		log.Warning("CRON_SCHEDULE environment not specified, using default value of 5 minutes")
		cronTime = pkg.CRON_5_MINUTE
	}

	excludeFolderStr := os.Getenv("EXCLUDE_FOLDERS")
	excludeFolderDirty := strings.Split(excludeFolderStr, ",")
	excludeFolderClean := []string{}
	for _, v := range excludeFolderDirty {
		v = strings.TrimSpace(v)
		if v != "" {
			excludeFolderClean = append(excludeFolderClean, v)
		}
	}

	folders = pkg.FilterFolders(folders, excludeFolderClean)

	finalString := "These folder(s) will be backed based on this CRON " + cronTime + ": "
	for _, v := range folders {
		finalString += fmt.Sprintf(" %s, ", v)
	}

	return BackupConfig{
		MaxWindow:      window,
		Path:           path,
		Cron:           cronTime,
		ExcludeFolders: excludeFolderClean,
	}
}

func MakeS3Config(log pkg.CustomLoggerI) S3Config {
	var err error
	config := S3Config{}

	config.Endpoint = os.Getenv("S3_ENDPOINT")
	if config.Endpoint == "" {
		config.Endpoint = pkg.InputString("Input S3 Enpoint (ex. is3.cloudhost.id): ")
	}

	config.BucketName = os.Getenv("S3_BUCKET_NAME")
	if config.Endpoint == "" {
		config.Endpoint = pkg.InputString("Input Bucket Name: ")
	}

	config.AccessKey = os.Getenv("S3_ACCESS_KEY")
	if config.Endpoint == "" {
		config.Endpoint = pkg.InputString("Input Access Key: ")
	}

	config.SecretKey = os.Getenv("S3_SECRET_KEY")
	if config.Endpoint == "" {
		config.Endpoint = pkg.InputString("Input Secret Key: ")
	}

	config.UseSSL, err = strconv.ParseBool(os.Getenv("S3_USE_SSL"))
	if err != nil {
		config.UseSSL = pkg.InputBool("Do you want to use ssl? ([y]/n) : ", func(s string) bool {
			return s != "n"
		})
	}

	return config
}
