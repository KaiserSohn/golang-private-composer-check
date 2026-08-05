package http_tools

import (
	"PrivateComposerChecker/composer"
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
