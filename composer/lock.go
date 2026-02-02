package composer

import (
	"PrivateComposerChecker/json_tools"
)

type ComposerLock struct {
	Packages    []Lib `json:"packages"`
	PackagesDev []Lib `json:"packages-dev"`
}

type Lib struct {
	Name    string
	Version string
	Source  map[string]string
	Dist    map[string]string
}

func ParseLock(fileLocation string) ComposerLock {
	var composer ComposerLock

	json_tools.Parse(fileLocation, &composer)

	return composer
}
