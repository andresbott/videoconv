package ffmpegtranscode

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Transcoder struct {
	ffmpeg string
}

type Cfg struct {
	FfmpegBin string
}

// New creates a transcoder capable of running ffmpeg
func New(cfg Cfg) (*Transcoder, error) {
	if _, err := os.Stat(cfg.FfmpegBin); err != nil {
		return nil, fmt.Errorf("ffmpeg not found at %s", cfg.FfmpegBin)
	}
	tc := Transcoder{
		ffmpeg: cfg.FfmpegBin,
	}
	return &tc, nil
}

// GetCmd returns the string of the ffmpeg command that would be executed
func (tc *Transcoder) GetCmd(input, output string, args []string) ([]string, error) {
	if input == output {
		return nil, fmt.Errorf("input cannot be the same as the output")
	}

	// make sure the path is in absolute notation
	inFile, err := filepath.Abs(input)
	if err != nil {
		return nil, err
	}

	// make sure the path is in absolute notation
	outFile, err := filepath.Abs(output)
	if err != nil {
		return nil, err
	}

	var r []string
	r = append(r, tc.ffmpeg)
	r = append(r, "-i", inFile)
	r = append(r, args...)
	r = append(r, outFile)

	return r, nil
}

// Run will execute the ffmpeg command with all the parameters
func (tc *Transcoder) Run(input, output string, args []string) ([]string, error) {

	cmd, err := tc.GetCmd(input, output, args)
	if err != nil {
		return nil, err
	}
	command := exec.Command(cmd[0], cmd[1:]...)

	// set var to get the output
	var out bytes.Buffer
	var errB bytes.Buffer

	// set the output to our variable
	command.Stdout = &out
	command.Stderr = &errB
	err = command.Run()
	if err != nil {

		lines := errB.String()
		lines = strings.TrimSpace(lines)

		lines2 := strings.Split(lines, "\n")
		return cmd, fmt.Errorf("%v : %s", err, lines2[len(lines2)-1:][0])
	}

	return cmd, nil
}
