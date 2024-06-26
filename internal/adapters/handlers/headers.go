package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

func GetHeaders(headers http.Header) (header map[string]interface{}, err error) {
	UserInfo, ok := headers["X-Endpoint-Api-Userinfo"]
	if !ok {
		return nil, err
	}
	b64String := UserInfo[0]
	b64String += strings.Repeat("=", (4-len(b64String)%4)%4)
	rawDecodedText, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawDecodedText, &header)
	return header, err
}
