package main

import (
	"PrivateComposerChecker/composer"
	"PrivateComposerChecker/html_tools"
	"PrivateComposerChecker/http_tools"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	lockFile, authFile, workers := parseFlags()

	composerLock := composer.ParseLock(lockFile)

	libs := append(composerLock.Packages, composerLock.PackagesDev...)
	auth := composer.ParseAuth(authFile)

	for _, lib := range libs {
		if lib.Dist == nil {
			continue
		}

		if len(lib.Dist["url"]) > 0 {
			lib.Dist["url"] = http_tools.PrepareUrlWithAuth(lib.Dist["url"], auth)
		}

		if lib.Source != nil && len(lib.Source["url"]) > 0 {
			lib.Source["url"] = http_tools.PrepareUrlWithAuth(lib.Source["url"], auth)
		}
	}

	showResult(auth, libs, workers)
}

func showResult(auth composer.Auth, libs []composer.Lib, workers int) {
	results := http_tools.CheckStatus(auth, libs, workers)
	success := true

	versions := make([][]html_tools.AvailableVersion, len(results))

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
				versions[i] = html_tools.ParseVersions(auth, results[i].Lib)
			}
		}()
	}

	for i, result := range results {
		if result.Status != http.StatusOK && result.Error == nil {
			jobs <- i
		}
	}
	close(jobs)
	wg.Wait()

	for i, result := range results {
		if result.Status != http.StatusOK {
			if result.Error != nil {
				fmt.Println(result.Lib.Name+":"+result.Lib.Version, "|", result.Error)
			} else {
				fmt.Println(
					result.Lib.Name+":"+result.Lib.Version,
					"Status:",
					result.Status,
					"|",
					formatAvailableVersion(versions[i]),
				)
			}
			success = false
		}
	}

	if success {
		fmt.Println("✅ All packages is available")
	}
}

func formatAvailableVersion(versions []html_tools.AvailableVersion) string {
	if len(versions) == 0 {
		return ""
	}

	parts := make([]string, 0, len(versions))
	for _, version := range versions {
		parts = append(parts, version.Version)
	}

	return "Available versions (" + strings.Join(parts, ", ") + ")"
}

func parseFlags() (string, string, int) {
	var lockFileLocation, authFileLocation string
	var workers int

	flag.StringVar(&lockFileLocation, "lock-file", "", "Composer.lock location")
	flag.StringVar(&authFileLocation, "auth", "", "Authentication file location")
	flag.IntVar(&workers, "workers", 10, "Number of concurrent workers")
	flag.Parse()

	if len(lockFileLocation) == 0 {
		fmt.Printf("Empty lock file location")
		os.Exit(1)
	}

	if len(authFileLocation) == 0 {
		fmt.Printf("Empty auth file location")
		os.Exit(1)
	}

	return lockFileLocation, authFileLocation, workers
}
