package vidprocessing

import (
	"fmt"
	"os/exec"
)

const OutDir = "out/"

// ConvertToGifCutByDuration takes in a url to a video (eg. http(s)://.*.mp4), start time, duration, filepath in the local machine, and calls Command to execute FFMPEG (installed on the host machine).
// The resulting Gif file created is assigned to file of type []byte else err of type error
func ConvertToGifCutByDuration(sourceUrl string, start string, duration int, filePath string) (file []byte, err error) {
	durationStr := fmt.Sprintf("%v", duration)

	// improve quality https://engineering.giphy.com/how-to-make-gifs-with-ffmpeg/
	e := exec.Command("ffmpeg", "-ss", start, "-t", durationStr, "-i", sourceUrl, "-filter_complex", "[0:v] fps=12, scale=1080:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse", filePath)
	stdout, err := e.CombinedOutput()
	file = stdout

	if err != nil {
		fmt.Printf("ERROR:\n%v\n", string(stdout))
		return
	}

	return
}

// ConvertToGifCutByDuration takes in a url to a video (eg. http(s)://.*.mp4), start time, end time, filepath in the local machine, and calls Command to execute FFMPEG (installed on the host machine).
// The resulting Gif file created is assigned to file of type []byte else err of type error
func ConvertToGifCutByStartEnd(sourceUrl string, start string, end string, filePath string) (file []byte, err error) {
	// improve quality https://engineering.giphy.com/how-to-make-gifs-with-ffmpeg/
	e := exec.Command("ffmpeg", "-ss", start, "-to", end, "-i", sourceUrl, "-filter_complex", "[0:v] fps=12, scale=1080:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse", filePath)
	stdout, err := e.CombinedOutput()
	file = stdout

	if err != nil {
		fmt.Printf("ERROR:\n%v\n", string(stdout))
		return
	}

	return
}
