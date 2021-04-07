package transcoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Cfg struct {
	FfmpegBin  string
	FfProbeBin string // set path to ffprobe to stan the file and populate additional fields
	Template   string // command template to be populated
	InputFile  string // full path (relative or absolute) to the input file
	OutputFile string // full path (relative or absolute) to the output file, extension is stripped
}

const spaceReplace = "½ŋŋ½"

// New generates a single instance of a transcoding unit to be run once with a given configuration
func New(cfg *Cfg) (*Transcoder, error) {

	tc := Transcoder{
		ffmpeg:   cfg.FfmpegBin,
		ffprobe:  cfg.FfProbeBin,
		template: cfg.Template,
	}

	// make sure the path is in absolute notation
	inFile, err := filepath.Abs(cfg.InputFile)
	if err != nil {
		return nil, err
	}
	tc.inputFile = inFile
	tc.inputFile = strings.ReplaceAll(tc.inputFile, " ", spaceReplace)

	// make sure the path is in absolute notation
	outFile, err := filepath.Abs(cfg.OutputFile)
	if err != nil {
		return nil, err
	}
	tc.outputFile = outFile
	tc.outputFile = strings.ReplaceAll(tc.outputFile, " ", spaceReplace)

	return &tc, nil
}

type Transcoder struct {
	ffmpeg     string
	ffprobe    string
	template   string
	inputFile  string // full path (relative or absolute) to the input file
	outputFile string // full path (relative or absolute) to the output file, extension is stripped

}

// GetOutputFile returns the absolute path to the output file of the transcoded file
func (tc *Transcoder) GetOutputFile() string {
	return strings.ReplaceAll(tc.outputFile, spaceReplace, " ")
}

// args will parse the template from the profile and generate a slice of cmd arguments
func (tc *Transcoder) args() ([]string, error) {

	tpl := tc.template
	tpl = strings.ReplaceAll(tpl, "\n", " ")

	t, err := template.New("profile").Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template: %s", err)
	}

	var buf bytes.Buffer

	data, err := tc.getTmplData()
	if err != nil {
		return nil, fmt.Errorf("unable to get video details: %s", err)
	}

	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}

	result := buf.String()

	values := strings.Split(result, " ")
	values = revertSpaces(values)
	return dropEmpty(values), nil
}

// tmplData is the struct passed to the template of the profile
type tmplData struct {
	Input           string // input file
	Output          string // file output
	Width           int
	Height          int
	DurationSeconds float64
	FrameRate       int

	// todo add more fields to be used
}

// getTmplData generates a template data payload struct according to the current profile and transcoder parameters
func (tc *Transcoder) getTmplData() (*tmplData, error) {

	data := tmplData{
		Input:  tc.inputFile,
		Output: tc.outputFile,
	}

	if tc.ffprobe != "" {
		err := tc.scanF(&data)
		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}

// scanFile gets relevant information about the video vile to be used later on in the template
func (tc *Transcoder) scanF(in *tmplData) error {

	fileName := strings.ReplaceAll(tc.inputFile, spaceReplace, " ")

	var cmd []string
	cmd = append(cmd, tc.ffprobe, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", fileName)
	command := exec.Command(cmd[0], cmd[1:]...)

	// set var to get the output
	var outputBuf bytes.Buffer
	var errB bytes.Buffer

	command.Stdout = &outputBuf
	command.Stderr = &errB
	err := command.Run()
	if err != nil {
		return fmt.Errorf("error running ffprobe command: %s", err)
	}

	data := ProbeData{}
	err = json.Unmarshal(outputBuf.Bytes(), &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling ffprobe json: %s", err)
	}

	for _, stream := range data.Streams {
		if stream.CodecType == "video" {
			in.Width = stream.Width
			in.Height = stream.Height

			// calculate frame rate
			avrFrameRate := stream.AvgFrameRate
			s := strings.Split(avrFrameRate, "/")
			if len(s) != 2 {
				return fmt.Errorf("unable to calculate framerate: string does not contain two numbers")
			}
			s1, err1 := strconv.Atoi(s[0])
			s2, err2 := strconv.Atoi(s[1])

			if err1 != nil || err2 != nil {
				return fmt.Errorf("unable to calculate framerate: error converting %s to numnbers", avrFrameRate)
			}
			in.FrameRate = s1 / s2

			break // ignore possible second video streams
		}
	}
	in.DurationSeconds = data.Format.DurationSeconds

	return nil
}

// GetCmd returns the string of the ffmpeg command that would be executed
func (tc *Transcoder) GetCmd() ([]string, error) {

	var r []string
	r = append(r, tc.ffmpeg)

	args, err := tc.args()
	if err != nil {
		return nil, err
	}
	r = append(r, args...)
	return r, nil
}

// Run will execute the ffmpeg command with all the parameters
func (tc *Transcoder) Run() (string, error) {

	cmd, err := tc.GetCmd()
	if err != nil {
		return "", err
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
		return strings.Join(cmd, " "), fmt.Errorf("%v : %s", err, lines2[len(lines2)-1:][0])
	}

	return "", nil
}

// remove empty items in slice
func dropEmpty(in []string) []string {
	var out []string
	for _, v := range in {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out
}

// replaces special space placeholder (introduced earlier into file names with spaces ) by real spaces again
func revertSpaces(in []string) []string {
	var out []string
	for _, v := range in {
		s := strings.ReplaceAll(v, spaceReplace, " ")
		out = append(out, s)
	}
	return out
}

func fmtDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second
	d -= s * time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
