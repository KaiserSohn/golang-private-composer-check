package http_tools

import (
	"PrivateComposerChecker/composer"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URI    string
	Status int
	Lib    composer.Lib
	Error  error
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	},
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

func CheckStatus(auth composer.Auth, libs []composer.Lib, workers int) []Result {
	results := make([]Result, len(libs))

	if workers < 1 {
		workers = 1
	}

	jobs := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = checkLib(auth, libs[i])
			}
		}()
	}

	for i := range libs {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	return results
}

func checkLib(auth composer.Auth, lib composer.Lib) Result {
	url := lib.Dist["url"]

	if err := ValidateUrl(url); err != nil {
		return Result{URI: url, Status: 0, Lib: lib, Error: err}
	}

	resp, err := httpClient.Do(createRequest(auth, url))
	if err != nil {
		return Result{URI: url, Status: 0, Lib: lib, Error: fmt.Errorf("GET %s: %w", url, err)}
	}

	if len(lib.Source["url"]) > 0 && resp.StatusCode != http.StatusOK {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}

		url = lib.Source["url"]

		if err := ValidateUrl(url); err != nil {
			return Result{URI: url, Status: 0, Lib: lib, Error: err}
		}

		resp, err = httpClient.Do(createRequest(auth, url))
		if err != nil {
			return Result{URI: url, Status: 0, Lib: lib, Error: fmt.Errorf("GET %s: %w", url, err)}
		}
	}

	result := Result{
		URI:    url,
		Status: resp.StatusCode,
		Lib:    lib,
	}

	if err := resp.Body.Close(); err != nil {
		fmt.Printf("Error closing response body: %v\n", err)
	}

	return result
}

func GetHtml(auth composer.Auth, url string) string {
	resp, err := httpClient.Do(createRequest(auth, url))
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
