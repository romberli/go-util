package common

import (
	"os"

	"github.com/pkg/sftp"
)

func StringInSlice(str string, slice []string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}

	return false
}

func StringInMap(str string, m map[string]string) bool {
	if _, ok := m[str]; ok {
		return true
	}

	return false
}

func PathExistsLocal(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func PathExistsRemote(path string, client *sftp.Client) (bool, error) {
	if _, err := client.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
