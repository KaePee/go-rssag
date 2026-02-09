package auth

import (
	"errors"
	"net/http"
	"strings"
)

//Extracts apikey from http requests
//Authorization: ApiKey {api_key}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authentication info")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed authentication info, expected `ApiKey api_key`")
	}
	return vals[1], nil
}