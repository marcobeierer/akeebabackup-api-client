package akeebabackup

// All the fields are optional. You can have a null data body.
type StartBackupRequestData struct {
	Profile     int      `json:"profile,omitempty"`     // Default = 1. The numeric profile ID to use for this backup job.
	Description string   `json:"description,omitempty"` // Default = Backup taken on <date> <time>. The short description to list in the Manage Backups (formerly "Administer Backup Files") page.
	Comment     string   `json:"comment,omitempty"`     // Default = Empty string. A longer comment to include in the Manage Backups (formerly "Administer Backup Files") listing for this backup record.
	Tag         string   `json:"tag,omitempty"`         // Default = json. Do not use this option! This is reserved for future use.
	Overrides   []string `json:"overrides,omitempty"`   // Configuration variable overrides.
}

func NewStartBackupRequestData() *StartBackupRequestData {
	return &StartBackupRequestData{}
}

type StartBackupResponseData struct {
	HasRun      bool     `json:""`           // When true, more steps are required to backup the site (unless Error is not empty). When false, the backup has just finished.
	Domain      string   `json:""`           // The backup domain. It can be one of "init', 'installer', 'db', 'pack', 'finalization'. It tells you which major step the backup is in, i.e. 'db' for database backup.
	Step        string   `json:""`           // Free text describing the last operation in the current domain, e.g. the name of the last folder backed up. You should only use it for verbose progress display.
	Substep     string   `json:""`           // Free text, giving more detail about the Step. You should only use it for verbose progress display.
	Error       string   `json:""`           // The last error occurred. If this is not an empty string or null, you can assume that the backup has failed, irrespectively of the HasRun value.
	Warnings    []string `json:""`           // The warnings produced during the last step. You'd better display them to the user.
	BackupIDAlt int      `json:"BackupID"`   // The numeric ID of the backup record being created.
	Archive     string   `json:""`           // The name of the backup archive. You will only get this in one of the responses of startBackup or stepBackup, when the archive is first created. In all other responses it will be empty. Only use the non-empty response.
	Progress    int      `json:""`           // The percentage of the backup process completion (0-100).
	BackupID    string   `json:"backupid"`   // The unique identifier of this backup. You need to pass this in all subsequent calls to stepBackup.
	SleepTime   int      `json:"sleepTime"`  // How many milliseconds you should wait before you run the next stepBackup.
	StepNumber  int      `json:"stepNumber"` // The sequential number of this step. Step numbers returned by stepBackup must be sequentially numbered, monotonically increasing. Anything else indicates a misbehaving server / broken backup.
	StepState   string   `json:"stepState"`  // The run state of the engine. It can be one of error, init, prepared, running, postrun and finished. This is debugging information about the engine internals and you should not act upon it.
}

type StartBackupRequest struct {
	Request
	url string
}

type StartBackupResponse struct {
	Response
}

// TODO validate url (trailing slash) and maybe key
func NewStartBackupRequest(url, frontendKey string) *StartBackupRequest {
	return &StartBackupRequest{
		Request: *newRequest(frontendKey, "startBackup", NewStartBackupRequestData()),
		url:     url,
	}
}

func NewStartBackupResponse() *StartBackupResponse {
	return &StartBackupResponse{
		Response: *newResponse(&StartBackupResponseData{}),
	}
}

func (qr *StartBackupRequest) Execute() (*StartBackupResponse, bool) {
	response := NewStartBackupResponse()
	return response, qr.Request.execute(qr.url, &response.Response, "")
}

func (qr *StartBackupResponse) Data() *StartBackupResponseData {
	return qr.Response.Body.Data().(*StartBackupResponseData)
}
