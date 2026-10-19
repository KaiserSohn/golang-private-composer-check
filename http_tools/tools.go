package http_tools

import (
	"PrivateComposerChecker/composer"
	"errors"
	"fmt"
	"net/url"
)

func Parse(uriString string) *url.URL {
	uri, err := url.Parse(uriString)
	if err != nil {
		fmt.Println("Bad uri", uriString)
		return nil
	}
	return uri
}

func GetCredentialsForHost(auth composer.Auth, hostname string) composer.UsernamePasswordCredentials {
	uri := Parse(hostname)
	if uri == nil {
		return composer.UsernamePasswordCredentials{}
	}

	return auth.HttpBasic[uri.Host]
}

func ValidateUrl(uriString string) error {
	if len(uriString) == 0 {
		return errors.New("empty url")
	}

	uri, err := url.Parse(uriString)
	if err != nil {
		return fmt.Errorf("invalid url %q: %w", uriString, err)
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		return fmt.Errorf("invalid url %q: unsupported scheme %q", uriString, uri.Scheme)
	}

	if len(uri.Host) == 0 {
		return fmt.Errorf("invalid url %q: missing host", uriString)
	}

	return nil
}

func PrepareUrlWithAuth(uriString string, auth composer.Auth) string {
	uri := Parse(uriString)
	if uri == nil {
		return uriString
	}

	credentials := GetCredentialsForHost(auth, uri.Host)
	if len(credentials.Username) == 0 && len(credentials.Password) == 0 {
		return uriString
	}

	return uri.Scheme + "://" + credentials.Username + ":" + credentials.Password + "@" + uri.Host + uri.RequestURI()
}
