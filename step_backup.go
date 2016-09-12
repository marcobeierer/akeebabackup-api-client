package akeebabackup

type StepBackupRequestData struct {
	Tag      string `json:"tag,omitempty"`      // Optional. Default = json. Do not use this option! This is reserved for future use.
	BackupID string `json:"backupid,omitempty"` // MANDATORY. The unique identifier of this backup, as returned by the previous call to startBackup or stepBackup.
}

func NewStepBackupRequestData(backupID string) *StepBackupRequestData {
	return &StepBackupRequestData{
		BackupID: backupID,
	}
}

type StepBackupResponseData struct {
	HasRun      bool     `json:""`           // When true, more steps are required to backup the site (unless Error is not empty). When false, the backup has just finished.
	Domain      string   `json:""`           // The backup domain. It can be one of "init', 'installer', 'db', 'pack', 'finalization'. It tells you which major step the backup is in, i.e. 'db' for database backup.
	Step        string   `json:""`           // Free text describing the last operation in the current domain, e.g. the name of the last folder backed up. You should only use it for verbose progress display.
	Substep     string   `json:""`           // Free text, giving more detail about the Step. You should only use it for verbose progress display.
	Error       string   `json:""`           // The last error occurred. If this is not an empty string or null, you can assume that the backup has failed, irrespectively of the HasRun value.
	Warnings    []string `json:""`           // The warnings produced during the last step. You'd better display them to the user.
	Archive     string   `json:""`           // The name of the backup archive. You will only get this in one of the responses of startBackup or stepBackup, when the archive is first created. In all other responses it will be empty. Only use the non-empty response.
	Progress    int      `json:""`           // The percentage of the backup process completion (0-100).
	BackupID    string   `json:"backupid"`   // The unique identifier of this backup. You need to pass this in all subsequent calls to stepBackup.
	SleepTime   int      `json:"sleepTime"`  // How many milliseconds you should wait before you run the next stepBackup.
	StepNumber  int      `json:"stepNumber"` // The sequential number of this step. Step numbers returned by stepBackup must be sequentially numbered, monotonically increasing. Anything else indicates a misbehaving server / broken backup.
	StepState   string   `json:"stepState"`  // The run state of the engine. It can be one of error, init, prepared, running, postrun and finished. This is debugging information about the engine internals and you should not act upon it.
	BackupIDAlt int      `json:"BackupID"`   // The numeric ID of the backup record being created.
}

type StepBackupRequest struct {
	Request
	url string
}

type StepBackupResponse struct {
	Response
}

// TODO validate url (trailing slash) and maybe key
func NewStepBackupRequest(url, frontendKey, backupID string) *StepBackupRequest {
	return &StepBackupRequest{
		Request: *newRequest(frontendKey, "stepBackup", NewStepBackupRequestData(backupID)),
		url:     url,
	}
}

func NewStepBackupResponse() *StepBackupResponse {
	return &StepBackupResponse{
		Response: *newResponse(&StepBackupResponseData{}),
	}
}

func (qr *StepBackupRequest) Execute() (*StepBackupResponse, bool) {
	response := NewStepBackupResponse()
	return response, qr.Request.execute(qr.url, &response.Response, "")
}

func (qr *StepBackupResponse) Data() *StepBackupResponseData {
	return qr.Response.Body.Data().(*StepBackupResponseData)
}
