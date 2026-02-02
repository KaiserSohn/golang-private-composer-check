package composer

import (
	"project/json_tools"
)

type Auth struct {
	HttpBasic map[string]UsernamePasswordCredentials `json:"http-basic"`
}

type UsernamePasswordCredentials struct {
	Username string
	Password string
}

func ParseAuth(fileLocation string) Auth {
	var auth Auth

	json_tools.Parse(fileLocation, &auth)

	return auth
}
