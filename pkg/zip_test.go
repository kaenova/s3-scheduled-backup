package pkg

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateZipFile(t *testing.T) {
	os.Mkdir("./zip-test", 0777)
	os.Mkdir("./zip-test/test", 0777)
	err := ZipSource("./zip-test", "./zip-test.zip")
	if err != nil {
		os.Remove("./zip-test")
		t.Fatal("Error on Zipping: " + err.Error())
	}
	_, err = os.Stat("./zip-test.zip")
	if err != nil {
		os.Remove("./zip-test")
		t.Fatal("Error on checking zip-test.zip : " + err.Error())
	}
	os.Remove("./zip-test.zip")
	os.RemoveAll("./zip-test")
}

func TestWalkDir(t *testing.T) {
	paths := []string{"a", "b", "c"}
	expected := []string{}
	for _, v := range paths {
		expected = append(expected, v)
		err := os.MkdirAll("./folder-test/"+v, 0777)
		if err != nil {
			t.Fatal("Cannot create test dir")
		}
		for _, vi := range paths {
			err := os.MkdirAll("./folder-test/"+v+"/"+vi, 0777)
			if err != nil {
				t.Fatal("Cannot create test dir")
			}
		}
	}

	defer func() {
		os.RemoveAll("./folder-test")
	}()

	folders, err := FoldersOneLevel("./folder-test")
	if err != nil {
		t.Fatal("Error on Walking Path")
	}

	assert.Equal(t, expected, folders)
}

func TestParseBackupFile(t *testing.T) {
	var expected, parsedFileName []string
	var parsedFile []BackupFolder
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

		expected = append(expected, fmt.Sprintf("%s--%s.zip", fileName, currentTime.Format("2006-01-02")))
	}

	// Parse File Name
	for _, val := range expected {
		obj, err := ParseBackupFolder(val)
		if err != nil {
			t.Fatal(err.Error())
		}
		parsedFile = append(parsedFile, obj)
	}

	// Create parsed
	for _, val := range parsedFile {
		parsedFileName = append(parsedFileName, fmt.Sprintf("%s--%s.zip", val.FolderName, val.Time.String()))
	}
	t.Log("Actual:", expected)
	t.Log("Parsed:", parsedFileName)
	assert.Equal(t, expected, parsedFileName)
}
