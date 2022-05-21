package vidprocessing

import (
	"fmt"
	"os/exec"
)

const OutDir = "out/"

func ConvertToGifByUrl(sourceUrl string, start string, duration int, filePath string) (file []byte, err error) {
	s := fmt.Sprintf("%v", start)
	t := fmt.Sprintf("%v", duration)

	e := exec.Command("ffmpeg", "-ss", s, "-to", t, "-i", sourceUrl, "-filter_complex", "[0:v] fps=12, scale=1080:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse", filePath)
	// improve quality https://engineering.giphy.com/how-to-make-gifs-with-ffmpeg/
	stdout, err := e.CombinedOutput()
	file = stdout
	if err != nil {
		fmt.Printf("ERROR:\n%v\n", string(stdout))
		return
	}

	return
}

func ConvertToGifByUrlByStartEnd(sourceUrl string, start string, end string, filePath string) (file []byte, err error) {
	s := fmt.Sprintf("%v", start)
	t := fmt.Sprintf("%v", end)

	e := exec.Command("ffmpeg", "-ss", s, "-to", t, "-i", sourceUrl, "-filter_complex", "[0:v] fps=12, scale=1080:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse", filePath)
	// improve quality https://engineering.giphy.com/how-to-make-gifs-with-ffmpeg/
	stdout, err := e.CombinedOutput()
	file = stdout
	if err != nil {
		fmt.Printf("ERROR:\n%v\n", string(stdout))
		return
	}

	return
}

func Test() (t int) {
	t = 2 + 2
	return t
}
