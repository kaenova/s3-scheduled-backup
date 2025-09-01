package pkg_test

import (
	"sort"
	"testing"
	"time"

	"github.com/kaenova/s3-scheduled-backup/pkg"
)

func TestNaming(t *testing.T) {
	currentTime := time.Now()
	backupTime := pkg.BackupTime{Time: currentTime}
	t.Log(backupTime.String())
}

func TestSorting(t *testing.T) {
	data := []string{"20220102", "20220202", "20220103"}
	sort.Strings(data)
	t.Log(data)
}
