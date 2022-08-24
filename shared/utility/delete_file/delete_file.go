package delete_file

import (
	"os"
)

// RemoveFileFromDirectory deletes the local host file created when converting from video to gif.
// An improvement would be to make FFMPEG not create this local file
func RemoveFileFromDirectory(dir string) (err error) {
	if len(dir) > 0 {
		err := os.Remove(dir)

		if err != nil {
			return err
		}

	}

	return
}
