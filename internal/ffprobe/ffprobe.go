package ffprobe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type FfProbe struct {
	binary string
}

const defaultFfporbe = "/usr/bin/ffprobe"

func New(bin ...string) (FfProbe, error) {

	binary := defaultFfporbe
	if len(bin) > 0 && bin[0] != "" {
		binary = bin[0]
	}

	if _, err := os.Stat(binary); err != nil {
		return FfProbe{}, fmt.Errorf("ffprobe not found at %s", binary)
	}

	f := FfProbe{
		binary: binary,
	}
	return f, nil
}

func (ff FfProbe) Probe(file string) (ProbeData, error) {

	var cmd []string
	cmd = append(cmd, ff.binary, "-v", "quiet", "-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_chapters",
		file)
	command := exec.Command(cmd[0], cmd[1:]...)

	// set var to get the output
	var outputBuf bytes.Buffer
	var errB bytes.Buffer

	command.Stdout = &outputBuf
	command.Stderr = &errB
	err := command.Run()
	if err != nil {
		return ProbeData{}, fmt.Errorf("error running ffprobe command: %s", err)
	}

	data := ProbeData{}
	err = json.Unmarshal(outputBuf.Bytes(), &data)
	if err != nil {
		return ProbeData{}, fmt.Errorf("error unmarshaling ffprobe json: %s", err)
	}
	data.Digest()

	return data, nil
}
