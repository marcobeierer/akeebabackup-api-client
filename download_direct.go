package akeebabackup

// Do note that the data is always returned unencrypted, no mater what the requested encapsulation was in your request.
type DownloadDirectRequestData struct {
	BackupID int `json:"backup_id"` // The numeric ID of the backup record whose files you want to download
	PartID   int `json:"part_id"`   // The backup part you wish to download. For example, if we have a multi-part archive with the base name test.jpa and six parts (.j01 through .j05 and .jpa) then part_id=2 means that we want to download test.j02 and part_id=6 means that we want to download the last part, i.e. test.jpa. If the backup record is not a multi-part archive just use 1.
}

func NewDownloadDirectRequestData(backupID, partID int) *DownloadDirectRequestData {
	return &DownloadDirectRequestData{
		BackupID: backupID,
		PartID:   partID,
	}
}

// Unencrypted binary stream containing the raw file data. You will also receive standard HTTP headers setting the content disposition to Attachment, specifying an application/octet-stream MIME type and notifying you of the size of the download. In fact, you can generate a URL and pass it to any third party download tool (e.g. cURL, Wget) to download the archive part.
type DownloadDirectResponseData struct {
}

type DownloadDirectRequest struct {
	Request
	url string
}

type DownloadDirectResponse struct {
	Response
}

// TODO validate url (trailing slash) and maybe key
func NewDownloadDirectRequest(url, frontendKey string, backupID int) *DownloadDirectRequest {
	return &DownloadDirectRequest{
		Request: *newRequest(frontendKey, "downloadDirect", NewDownloadDirectRequestData(backupID, 1)),
		url:     url,
	}
}

func NewDownloadDirectResponse() *DownloadDirectResponse {
	return &DownloadDirectResponse{
		Response: *newResponse(&DownloadDirectResponseData{}),
	}
}

func (qr *DownloadDirectRequest) Execute(filepath string) (*DownloadDirectResponse, bool) {
	response := NewDownloadDirectResponse()
	return response, qr.Request.execute(qr.url, &response.Response, filepath)
}
