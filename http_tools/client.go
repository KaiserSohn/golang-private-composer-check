package http_tools

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"project/composer"
)

type Result struct {
	URI    string
	Status int
	Lib    composer.Lib
}

func createRequest(auth composer.Auth, url string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	credentials := GetCredentialsForHost(auth, url)
	base64Creds := base64.StdEncoding.EncodeToString([]byte(credentials.Username + ":" + credentials.Password))
	request.Header.Add("Authorization", "Basic "+base64Creds)

	return request
}

func CheckStatus(auth composer.Auth, libs []composer.Lib) []Result {
	results := make([]Result, len(libs))

	for i, lib := range libs {
		url := lib.Dist["url"]

		client := &http.Client{}
		resp, err := client.Do(createRequest(auth, url))

		if len(lib.Source["url"]) > 0 && resp.StatusCode != http.StatusOK {
			err := resp.Body.Close()
			if err != nil {
				return nil
			}

			url = lib.Source["url"]
			resp, err = client.Do(createRequest(auth, url))
		}

		if err != nil {
			fmt.Printf("Error making GET request: %v", err)
		}

		results[i] = Result{
			URI:    url,
			Status: resp.StatusCode,
			Lib:    lib,
		}

		err = resp.Body.Close()
		if err != nil {
			return nil
		}
	}
	return results
}

func GetHtml(auth composer.Auth, url string) string {
	client := &http.Client{}
	resp, err := client.Do(createRequest(auth, url))

	if err != nil {
		fmt.Printf("Error making GET request: %v", err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s", resp.Status)
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return ""
	}

	return string(htmlBytes)
}
