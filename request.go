package akeebabackup

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ugomo/recruitertool/utils/crypto"
)

type Request struct {
	Encapsulation int          `json:"encapsulation"` // An Encapsulation Identifier constant (integer), defining the encapsulation method of the body.
	BodyString    string       `json:"body"`          // A JSON string representing the Request Body, encoded in the specified encapsulation
	Body          *RequestBody `json:"-"`
}

// TODO is there a more type safe way to do this? so that a change of Request struct does not affect this method
func (qr *Request) MarshalJSON() ([]byte, error) {
	bodyData, err := json.Marshal(qr.Body)
	if err != nil {
		log.Println(err)
		return bodyData, err
	}

	qr.BodyString = string(bodyData)

	bodyData, err = json.Marshal(qr.BodyString)
	if err != nil {
		log.Println(err)
		return bodyData, err
	}

	jsonString := fmt.Sprintf("{\"encapsulation\": %d, \"body\": %s}", qr.Encapsulation, string(bodyData))
	return []byte(jsonString), nil
}

type RequestBody struct {
	// This field is required if and only if the encapsulation is 1 (ENCAPSULATION_RAW). It consists of a salt string and an MD5 hash, separated by a colon, like this: salt:md5.
	// The salt can be an arbitrary length alphanumeric string.
	// The md5 part of the challenge is the result of the MD5 hash of the concatenated string of the salt and the Akeeba Backup front-end secret key, as configured in the component's Parameters.
	// For example, if the salt is foo and the secret key is bar, the md5 is md5(foobar) = 3858f62230ac3c915f300c664312c63f, therefore the challenge is foo:3858f62230ac3c915f300c664312c63f.
	//
	// If the encapsulation is higher than 1, the challenge field is completely ignored.
	// The authentication is implicitly performed as soon as the body is successfully deciphered, since the knowledge of the server's key is a prerequisite for successful encryption.
	Challenge string      `json:"challenge"`
	Key       string      `json:"key,omitempty"` // This field must be present and is required if and only if the encapsulation is higher than 1. The server's response will be encrypted using this key, not the Secret Key defined in the component's Parameters.
	Method    string      `json:"method"`        // The name of the method you want the server to execute
	Data      interface{} `json:"data"`          // A JSON object containing the parameters to be passed to the method.
}

/*func (qr *RequestBody) MarshalJSON() ([]byte, error) {

}*/

func newRequest(frontendKey, method string, data interface{}) *Request {
	return &Request{
		Encapsulation: 1, // ENCAPSULATION_RAW // TODO implement enum
		Body: &RequestBody{
			Challenge: challenge(frontendKey),
			Method:    method, // TODO use an enum
			Data:      data,
		},
	}
}

func (qr *Request) execute(url string, response *Response, filepath string) bool {
	jsonData, err := json.Marshal(qr)
	if err != nil {
		log.Println(err)
		return false
	}

	responsex, err := http.Get(url + string(jsonData))
	if err != nil {
		log.Println(err)
		return false
	}
	defer responsex.Body.Close()

	if filepath != "" { // download
		// TODO check if file already exists

		file, err := os.Create(filepath)
		if err != nil {
			log.Println(err)
			return false
		}
		defer file.Close()

		_, err = io.Copy(file, responsex.Body)
		if err != nil {
			log.Println(err)
			return false
		}
	} else {
		responseData, err := ioutil.ReadAll(responsex.Body)
		if err != nil {
			log.Println(err)
			return false
		}

		responseString := string(responseData)

		// necessary because of an odd compatibility issue fix; see documentation of server response object:
		indexOfFirstTripleHashtag := strings.Index(responseString, "###")
		indexOfLastTripleHashtag := strings.LastIndex(responseString, "###")
		responseString = responseString[indexOfFirstTripleHashtag+3 : indexOfLastTripleHashtag]

		if err := json.Unmarshal([]byte(responseString), response); err != nil {
			log.Println(err)
			log.Println(responseString)
			log.Println(string(jsonData))
			log.Println(url)

			return false
		}
	}

	return true
}

func challenge(frontendKey string) string {
	salt, err := cryptoutils.RandomStringAsHex(8)
	if err != nil || salt == "" {
		log.Println("FATAL:", err)
		salt = "dce49ca2ab195ee2" // TODO better cancel the request instead of using a fallback
	}

	data := []byte(salt + frontendKey)
	return fmt.Sprintf("%s:%x", salt, md5.Sum(data))
}
