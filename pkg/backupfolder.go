package pkg

import (
	"fmt"
	"regexp"
	"time"
)

const (
	// Time Format
	// <Year-Month-Date>
	TIME_FORMAT = "2006-01-02"

	// Format Zip Name
	// <foldername>--<Year-Month-Date>.zip
	REGEX_FILE_FORMAT = `(.+)--(\d+-\d+-\d+).zip$`
)

type BackupFolder struct {
	ZipFileName string
	FolderName  string
	Time        BackupTime
}

func (b *BackupFolder) GenerateFileName() string {
	return b.FolderName + "--" + b.Time.String()
}

type BackupTime struct {
	time.Time
}

func (b *BackupTime) String() string {
	return b.Time.Format(TIME_FORMAT)
}

func (b *BackupTime) Unix() int64 {
	return b.Time.Unix()
}

func ParseBackupFolder(fileName string) (BackupFolder, error) {
	re := regexp.MustCompile(REGEX_FILE_FORMAT)
	res := re.FindStringSubmatch(fileName)

	time, err := time.Parse(TIME_FORMAT, res[2])
	if err != nil {
		return BackupFolder{}, err
	}

	return BackupFolder{
		ZipFileName: fileName,
		FolderName:  res[1],
		Time:        BackupTime{time},
	}, nil
}

func CreateBackupFolder(folderName string, time time.Time) BackupFolder {
	currentTime := BackupTime{time}
	fileName := fmt.Sprintf("%s--%s.zip", folderName, currentTime.String())

	return BackupFolder{
		ZipFileName: fileName,
		FolderName:  folderName,
		Time:        currentTime,
	}
}
