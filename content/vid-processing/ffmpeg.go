package vidprocessing

import (
	"fmt"
	"os/exec"
)

const OutDir = "out/"

func ConvertToGifByUrl(sourceUrl string, start int, duration int, filePath string) (file []byte, err error) {
	s := fmt.Sprintf("%v", start)
	t := fmt.Sprintf("%v", duration)

	e := exec.Command("ffmpeg", "-t", t, "-ss", s, "-i", sourceUrl, "-filter_complex", "[0:v] fps=12, scale=1080:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse", filePath)
	// improve quality https://engineering.giphy.com/how-to-make-gifs-with-ffmpeg/
	stdout, err := e.CombinedOutput()
	file = stdout
	if err != nil {
		fmt.Printf("ERROR:\n%v\n", string(stdout))
		return
	}

	return
}
