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

type CmdArgs struct {
	Ffmpeg   string
	InitArgs []string
	Input    string
	Args     []string
	Output   string
}

func (cmd CmdArgs) String() string {
	return fmt.Sprintf("%s %s -i \"%s\" %s \"%s\"",
		cmd.Ffmpeg, strings.Join(cmd.InitArgs, " "), cmd.Input,
		strings.Join(cmd.Args, " "), cmd.Output,
	)

}

func (cmd CmdArgs) Slice() []string {
	var r []string
	r = append(r, cmd.Ffmpeg)
	r = append(r, cmd.InitArgs...)
	r = append(r, "-i", cmd.Input)
	r = append(r, cmd.Args...)
	r = append(r, cmd.Output)
	return r
}

// GetCmd returns the string of the ffmpeg command that would be executed
func (tc *Transcoder) GetCmd(input, output string, init, args []string) (CmdArgs, error) {
	if input == output {
		return CmdArgs{}, fmt.Errorf("input cannot be the same as the output")
	}

	// make sure the path is in absolute notation
	inFile, err := filepath.Abs(input)
	if err != nil {
		return CmdArgs{}, err
	}

	// make sure the path is in absolute notation
	outFile, err := filepath.Abs(output)
	if err != nil {
		return CmdArgs{}, err
	}

	r := CmdArgs{
		Ffmpeg:   tc.ffmpeg,
		InitArgs: init,
		Input:    inFile,
		Args:     args,
		Output:   outFile,
	}

	return r, nil
}

// Run will execute the ffmpeg command with all the parameters
func (tc *Transcoder) Run(input, output string, init, args []string) (CmdArgs, error) {

	cmd, err := tc.GetCmd(input, output, init, args)
	if err != nil {
		return CmdArgs{}, err
	}
	cmdSlice := cmd.Slice()
	command := exec.Command(cmdSlice[0], cmdSlice[1:]...)

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
