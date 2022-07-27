package pkg_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/kaenova/s3-scheduled-backup/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseBackupFile(t *testing.T) {
	var expected, parsedFileName []string
	var parsedFile []pkg.BackupFolder
	totalDummy := 10000

	randomString := func(n int) (string, error) {
		const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 -_@#$%^&*(!"
		b := make([]byte, n)
		for i := range b {
			idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
			if err != nil {
				return "", err
			}
			b[i] = letterBytes[idx.Int64()]
		}
		return string(b), nil
	}

	// Create Expected File Name
	for i := 0; i < totalDummy; i++ {
		fileName, err := randomString(20)
		if err != nil {
			t.Fatal(err.Error())
		}
		currentTime := time.Now()

		expected = append(expected, fmt.Sprintf("%s--%s.zip", fileName, currentTime.Format(pkg.TIME_FORMAT)))
	}

	// Parse File Name
	for _, val := range expected {
		obj, err := pkg.ParseBackupFolder(val)
		if err != nil {
			t.Fatal(err.Error())
		}
		parsedFile = append(parsedFile, obj)
	}

	// Create parsed
	for _, val := range parsedFile {
		parsedFileName = append(parsedFileName, fmt.Sprintf("%s--%s.zip", val.FolderName, val.Time.String()))
	}
	assert.Equal(t, expected, parsedFileName)
}
