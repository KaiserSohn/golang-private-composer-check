package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"project/composer"
	"project/html_tools"
	"project/http_tools"
	"strings"
)

func main() {
	lockFile, authFile := parseFlags()

	composerLock := composer.ParseLock(lockFile)

	libs := append(composerLock.Packages, composerLock.PackagesDev...)
	auth := composer.ParseAuth(authFile)

	for _, lib := range libs {
		lib.Dist["url"] = http_tools.PrepareUrlWithAuth(lib.Dist["url"], auth)

		if len(lib.Source["url"]) > 0 {
			lib.Source["url"] = http_tools.PrepareUrlWithAuth(lib.Source["url"], auth)
		}
	}

	showResult(auth, libs)
}

func showResult(auth composer.Auth, libs []composer.Lib) {
	results := http_tools.CheckStatus(auth, libs)
	success := true

	for _, result := range results {
		if result.Status != http.StatusOK {

			fmt.Println(
				result.Lib.Name+":"+result.Lib.Version,
				"Status:",
				result.Status,
				"|",
				formatAvailableVersion(html_tools.ParseVersions(auth, result.Lib)),
			)
			success = false
		}
	}

	if success {
		fmt.Println("✅ All packages is available")
	}
}

func formatAvailableVersion(versions []html_tools.AvailableVersion) string {
	var builder strings.Builder

	if len(versions) == 0 {
		return ""
	}

	for _, version := range versions {
		builder.WriteString(version.Version)
		builder.WriteString(",")
	}

	return "Available versions (" + builder.String() + ")"
}

func parseFlags() (string, string) {
	var lockFileLocation, authFileLocation string

	flag.StringVar(&lockFileLocation, "lock-file", "", "Composer.lock location")
	flag.StringVar(&authFileLocation, "auth", "", "Authentication file location")

	flag.Parse()

	if len(lockFileLocation) == 0 {
		fmt.Printf("Empty lock file location")
		os.Exit(1)
	}

	if len(authFileLocation) == 0 {
		fmt.Printf("Empty auth file location")
		os.Exit(1)
	}

	return lockFileLocation, authFileLocation
}
