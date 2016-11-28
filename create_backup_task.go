package akeebabackup

import "log"

type CreateBackupTask struct {
	websiteURL  string
	frontendKey string
}

func NewCreateBackupTask(websiteURL, frontendKey string) *CreateBackupTask {
	return &CreateBackupTask{
		websiteURL:  websiteURL,
		frontendKey: frontendKey,
	}
}

// websiteURL has to have a trailing slash
func (qt *CreateBackupTask) Execute() bool {
	url := qt.websiteURL + "index.php?option=com_akeeba&view=json&format=component&json=" // just Joomla

	startBackupRequest := NewStartBackupRequest(url, qt.frontendKey)

	startBackupResponse, ok := startBackupRequest.Execute()
	if !ok {
		return false
	}

	startBackupData := startBackupResponse.Data()

	if startBackupResponse.Body.Status != 200 || startBackupData.BackupID == "" {
		log.Println("something went wrong")
		return false
	}

	finished := false
	filename := ""
	backupID := 0

	for !finished {
		stepBackupRequest := NewStepBackupRequest(url, qt.frontendKey, startBackupData.BackupID)

		stepBackupResponse, ok := stepBackupRequest.Execute()
		if !ok {
			return false
		}

		stepBackupData := stepBackupResponse.Data()

		if stepBackupResponse.Body.Status != 200 || stepBackupData.Error != "" {
			log.Println(stepBackupData.Error)
			return false
		}

		if !stepBackupData.HasRun {
			finished = true
		} else {
			backupID = stepBackupData.BackupIDAlt
		}

		filename = stepBackupData.Archive
	}

	if filename != "" && backupID != 0 {
		downloadDirectRequest := NewDownloadDirectRequest(url, qt.frontendKey, backupID)

		_, ok := downloadDirectRequest.Execute("/tmp/" + filename) // TODO use os.TmpDir
		if !ok {
			return false
		}
	} else {
		log.Println("filename or backupID not set")
		log.Println(filename)
		log.Println(backupID)
	}

	log.Println("finished")
	return true
}
