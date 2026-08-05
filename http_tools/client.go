package http_tools

import (
	"PrivateComposerChecker/composer"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type Result struct {
	URI    string
	Status int
	Lib    composer.Lib
}

func createRequest(auth composer.Auth, url string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request for %s: %v\n", url, err)
		return nil
	}

	credentials := GetCredentialsForHost(auth, url)
	base64Creds := base64.StdEncoding.EncodeToString([]byte(credentials.Username + ":" + credentials.Password))
	request.Header.Add("Authorization", "Basic "+base64Creds)

	return request
}

func CheckStatus(auth composer.Auth, libs []composer.Lib) []Result {
	results := make([]Result, len(libs))
	client := &http.Client{}

	for i, lib := range libs {
		url := lib.Dist["url"]

		resp, err := client.Do(createRequest(auth, url))
		if err != nil {
			fmt.Printf("Error making GET request: %v\n", err)
			results[i] = Result{URI: url, Status: 0, Lib: lib}
			continue
		}

		if len(lib.Source["url"]) > 0 && resp.StatusCode != http.StatusOK {
			if err := resp.Body.Close(); err != nil {
				fmt.Printf("Error closing response body: %v\n", err)
			}

			url = lib.Source["url"]
			resp, err = client.Do(createRequest(auth, url))
			if err != nil {
				fmt.Printf("Error making GET request: %v\n", err)
				results[i] = Result{URI: url, Status: 0, Lib: lib}
				continue
			}
		}

		results[i] = Result{
			URI:    url,
			Status: resp.StatusCode,
			Lib:    lib,
		}

		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}

	return results
}

func GetHtml(auth composer.Auth, url string) string {
	client := &http.Client{}
	resp, err := client.Do(createRequest(auth, url))
	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)
		return ""
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\n", resp.Status)
		return ""
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return ""
	}

	return string(htmlBytes)
}
