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

//func (t *Transcoder) Run() {
//
//	t.discoverFile()
//
//	if t.video == nil {
//		return
//	}
//
//	err := t.prepare()
//	if err != nil {
//		log.Error(err)
//		t.handleTranscodeError()
//		return
//	}
//
//	t.transcodeVideo()
//}

//func (t *Transcoder) transcodeVideo() {
//
//	log.Info("transcoding video file:" + t.video.FullPath())
//	fileManager, _ := fm.NewFm("")
//
//	timeStart := time.Now()
//
//	// loop over video settings and transcode them
//	for _, setting := range t.videoSettings {
//
//		videoTmpOut := t.cfg.TmpDir + "/" + t.video.Basename() + "_" + setting.Name() + "." + setting.OutputExtension()
//		cmd := setting.Cmd()
//
//		// make sure the tmp file does not exist
//		err := fileManager.DeleteFile(videoTmpOut)
//		if err != nil {
//			log.Error(err)
//		}
//
//		err = t.runffmpeg(cmd, videoTmpOut)
//		if err != nil {
//			log.Error("error in ffmpeg, rolling back", err)
//			t.handleTranscodeError()
//			return
//		}
//	}
//
//	// move finalized files
//
//	outPutPath := filepath.Clean(t.cfg.OutputFolder + "/" + t.relativePath)
//	err := os.MkdirAll(outPutPath, 0750)
//	if err != nil {
//		log.Error(err)
//		t.handleTranscodeError()
//		return
//	}
//
//	for _, setting := range t.videoSettings {
//
//		videoTmpOut := filepath.Clean(t.cfg.TmpDir + "/" + t.video.Basename() + "_" + setting.Name() + "." + setting.OutputExtension())
//		videoOut := filepath.Clean(outPutPath + "/" + t.video.Basename() + "_" + setting.Name() + "." + setting.OutputExtension())
//
//		// make sure the destination file does not exist
//		err := fileManager.MoveFile(videoTmpOut, videoOut, false)
//		if err != nil {
//			log.Error(err)
//			t.handleTranscodeError()
//			return
//		}
//	}
//	// move the original
//
//	err = fileManager.MoveFile(t.video.FullPath(), filepath.Clean(outPutPath+"/"+t.video.Name()), false)
//	if err != nil {
//		log.Error(err)
//		t.handleTranscodeError()
//	}
//	timeEnd := time.Now()
//	timeDiff := timeEnd.Sub(timeStart)
//	log.Info("finished processing: \"" + t.video.Name() + "\" in " + timeDiff.String())
//
//}
//
//func (t *Transcoder) handleTranscodeError() {
//	log.Warn("handling transcode Error")
//
//	fileManager, _ := fm.NewFm("")
//	for _, setting := range t.videoSettings {
//
//		videoTmpOut := t.cfg.TmpDir + "/" + t.video.Basename() + "_" + setting.Name() + "." + setting.OutputExtension()
//
//		// make sure the tmp file does not exist
//		log.Warn("deleting possible tmp file: \"" + videoTmpOut + "\"")
//		err := fileManager.DeleteFile(videoTmpOut)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//	// move the original to ignore
//
//	ignorepath := filepath.Clean(t.cfg.IgnoreDir + "/" + t.relativePath)
//	err := os.MkdirAll(ignorepath, 0750)
//	if err != nil {
//		log.Error(err)
//		return
//	}
//
//	log.Warn("moving file: " + t.video.FullPath() + " to ignore dir")
//	err = fileManager.MoveFile(t.video.FullPath(), ignorepath+"/"+t.video.Name(), false)
//	if err != nil {
//		log.Error(err)
//	}
//
//}
//
//func (t *Transcoder) runffmpeg(cmd string, videoTmpOut string) error {
//	log.Info("rendering video setting: " + t.video.Name() + " -> " + videoTmpOut)
//	binary := "/usr/bin/ffmpeg"
//
//	args := []string{}
//
//	args = append(args, "-loglevel")
//	args = append(args, "error")
//
//	args = append(args, "-i")
//	args = append(args, t.video.FullPath())
//
//	args = append(args, "-threads")
//	args = append(args, strconv.Itoa(t.cfg.ProcessorThreads))
//
//	se := strings.Split(cmd, " ")
//	for _, s := range se {
//		args = append(args, strings.Trim(s, "'"))
//	}
//
//	args = append(args, videoTmpOut)
//
//	e := exec.Command(binary, args...)
//	//e.Stdout = os.Stdout
//	e.Stderr = os.Stderr
//	err := e.Run()
//	if err != nil {
//		return err
//	}
//	return nil
//}
