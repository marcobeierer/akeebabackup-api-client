package akeebabackup

import (
	"encoding/json"
	"log"
)

// WARNING: The JSON string is prefixed and suffixed with a triple hash mark (###). The client program must throw away anything before and including the first occurrence of the triple has mark, as well as anything after and including the last occurrence of the triple hash mark. This is done for compatibility reasons, as some servers inject PHP warnings, or mistakenly serve their output as a full HTML page instead of the requested raw JSON format.
type Response struct {
	Encapsulation int           `json:"encapsulation"` // The encapsulation the server decided to use for its response
	Body          *ResponseBody `json:"body"`
}

type ResponseBody struct {
	Status     int         `json:"status"` // The status of the request
	DataString string      `json:"data"`   // The JSON-encoded Response Data, in the specified encapsulation
	data       interface{} `json:"-"`
}

func newResponse(data interface{}) *Response {
	return &Response{
		Body: &ResponseBody{
			data: data,
		},
	}
}

func (qr *ResponseBody) Data() interface{} {
	err := json.Unmarshal([]byte(qr.DataString), qr.data)
	if err != nil {
		log.Println(err)
	}

	return qr.data
}
