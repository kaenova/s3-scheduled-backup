package pkg

import (
	"os"
	"testing"

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
