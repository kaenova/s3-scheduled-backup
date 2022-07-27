package pkg

import (
	"fmt"
	"regexp"
	"time"
)

type BackupFolder struct {
	// Format Zip Name
	// <foldername>--<Year-Month-Date>.zip
	ZipFileName string
	FolderName  string
	Time        BackupTime
}

type BackupTime struct {
	time.Time
}

func (b *BackupTime) String() string {
	return b.Format("2006-01-02")
}

func ParseBackupFolder(fileName string) (BackupFolder, error) {
	re := regexp.MustCompile(`(.+)--(\d+-\d+-\d+).zip$`)
	res := re.FindStringSubmatch(fileName)

	time, err := time.Parse("2006-01-02", res[2])
	if err != nil {
		return BackupFolder{}, err
	}

	return BackupFolder{
		ZipFileName: fileName,
		FolderName:  res[1],
		Time:        BackupTime{time},
	}, nil
}

func CreateBackupFolder(folderName string) BackupFolder {
	currentTime := BackupTime{time.Now()}
	fileName := fmt.Sprintf("%s--%s.zip", folderName, currentTime.String())

	return BackupFolder{
		ZipFileName: fileName,
		FolderName:  folderName,
		Time:        currentTime,
	}
}
