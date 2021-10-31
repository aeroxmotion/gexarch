package util

import (
	"os"
	"sync"

	"golang.org/x/mod/modfile"
)

var (
	once          sync.Once
	parsedModFile *modfile.File
)

func ParseModfile() *modfile.File {
	once.Do(func() {
		modFileBytes, err := os.ReadFile("go.mod")
		PanicIfError(err)

		parsedModFile, err = modfile.Parse("go.mod", modFileBytes, func(_, version string) (string, error) {
			return version, nil
		})
		PanicIfError(err)
	})

	return parsedModFile
}
