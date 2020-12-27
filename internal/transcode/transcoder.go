package transcoder

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Cfg struct {
	FfmpegBin  string
	FfmpegOpts FfmpegOpts

	InputFile  string
	OutputFile string
}

type Transcoder struct {
	ffmpeg     string
	opts       FfmpegOpts
	outputFile string
	inputFile  string
}

// New() generates a single instance of a transcoding unit to be run once with a given configuration
func New(cfg *Cfg) (*Transcoder, error) {

	tc := Transcoder{
		ffmpeg: cfg.FfmpegBin,
		opts:   cfg.FfmpegOpts,
	}

	inFile, err := filepath.Abs(cfg.InputFile)
	if err != nil {
		return nil, err
	}

	tc.inputFile = inFile
	outfName := strings.TrimSuffix(filepath.Base(cfg.OutputFile), filepath.Ext(filepath.Base(cfg.OutputFile)))

	out, err := filepath.Abs(filepath.Dir(cfg.OutputFile) + "/" + outfName + "." + cfg.FfmpegOpts.Name + "." + cfg.FfmpegOpts.VideoExt())
	if err != nil {
		return nil, err
	}
	tc.outputFile = out

	return &tc, nil
}
func (tc *Transcoder) GetOutputFile() string {
	return tc.outputFile
}

func (tc *Transcoder) GetCmd() ([]string, error) {

	var r []string
	r = append(r, tc.ffmpeg)

	r = append(r, "-i", "\""+tc.inputFile+"\"")

	args, err := tc.opts.Args()
	if err != nil {
		return nil, err
	}
	r = append(r, args...)

	r = append(r, "\""+tc.outputFile+"\"")

	return r, nil

}

func (tc *Transcoder) Run() (string, error) {

	var cmd []string
	cmd = append(cmd, tc.ffmpeg)
	cmd = append(cmd, "-i", tc.inputFile)
	args, err := tc.opts.Args()
	if err != nil {
		return "", err
	}
	cmd = append(cmd, args...)
	cmd = append(cmd, tc.outputFile)

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
		return strings.Join(cmd, " "), fmt.Errorf("%v : %s", err, lines2[len(lines2)-1:][0])
	}

	return "", nil
}
