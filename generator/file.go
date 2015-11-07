package generator

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var fileNameRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)

func fileify(t string) string {
	return strings.ToLower(fileNameRegex.ReplaceAllString(t, ""))
}

func deleteIfExists(file string) error {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return nil
	}

	return os.Remove(file)
}

func write(file string, code []byte) error {
	if err := deleteIfExists(file); err != nil {
		return err
	}

	return ioutil.WriteFile(file, code, 0644)
}
