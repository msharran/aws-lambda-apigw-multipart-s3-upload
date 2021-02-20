package multipart

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

//DecodeFromBase64String - decodes form data from base64 encoded string and returns as map of key and value
func DecodeFromBase64String(request events.APIGatewayProxyRequest, keys []string) (map[string]string, error) {
	decodedTxt := make([]byte, base64.StdEncoding.DecodedLen(len(request.Body)))
	_, err := base64.StdEncoding.Decode(decodedTxt, []byte(request.Body))
	if err != nil {
		return nil, err
	}

	boundryStr := strings.Split(request.Headers["Content-Type"], "; ")[1]
	boundry := strings.Split(boundryStr, "=")[1]

	body := make(map[string]string, 0)
	var rdr = multipart.NewReader(bytes.NewBuffer(decodedTxt), boundry)
	for len(body) < len(keys) {
		part, err := rdr.NextPart()
		if err == io.EOF {
			break
		}
		for i := range keys {
			if part.FormName() == keys[i] {
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				body[part.FormName()] = buf.String()
			}
		}
	}

	return body, nil
}
