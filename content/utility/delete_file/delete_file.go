package delete_file

import (
	"os"
)

func RemoveFileFromDirectory(dir string) (err error) {
	if len(dir) > 0 {
		err := os.Remove(dir)

		if err != nil {
			return err
		}

	}

	return
}
