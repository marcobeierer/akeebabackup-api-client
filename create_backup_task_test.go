package akeebabackup

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestCreateBackupTask(qt *testing.T) {
	websiteURL := "https://"
	frontendKey := ""

	task := NewCreateBackupTask(websiteURL, frontendKey)
	if !task.Execute() {
		qt.Fatal("execution failed")
	}
}
