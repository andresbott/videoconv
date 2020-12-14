package transcoder

import (
	"github.com/AndresBott/videoconv/internal/videconv"
	ffmpeg "github.com/floostack/transcoder/ffmpeg"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

type Cfg struct {
	FfmpedBin string
	VideoFile string
	OutputDir string
	TmpDir    string
	FailDir   string
}

type Transcoder struct {
	finalOutputFile string
	tmpOutputFile   string
	inputFile       string
	profile         videconv.Profile
	opts            ffmpeg.Options
	ffmpegCfg       ffmpeg.Config
}

func New(cfg *Cfg) (*Transcoder, error) {

	tc := Transcoder{}

	tc.ffmpegCfg = ffmpeg.Config{
		FfmpegBinPath:   cfg.FfmpedBin,
		FfprobeBinPath:  "/usr/local/bin/ffprobe",
		ProgressEnabled: false,
	}

	tc.inputFile = cfg.VideoFile

	// ############# Deal with destination dirs
	tmpOut, err := filepath.Abs(cfg.TmpDir + "/" + filepath.Base(cfg.VideoFile) + ".Ext")
	if err != nil {
		return nil, err
	}
	out, err := filepath.Abs(cfg.OutputDir + "/" + filepath.Base(cfg.VideoFile) + ".Ext")
	if err != nil {
		return nil, err
	}

	if cfg.TmpDir != "" {
		tc.finalOutputFile = out
		tc.tmpOutputFile = tmpOut
	} else {
		tc.finalOutputFile = out
		tc.tmpOutputFile = out
	}

	tc.opts = ffmpeg.Options{}

	format := "mp4"
	tc.opts.OutputFormat = &format

	overwrite := true
	tc.opts.Overwrite = &overwrite

	//coder := ffmpeg.
	//	New(ffmpegConf).
	//	Input("/tmp/avi").
	//	Output("/tmp/mp4").
	//	WithOptions(opts)

	return &tc, nil
}

func (tc *Transcoder) getCmd() ([]string, error) {

	var r []string
	r = append(r, tc.ffmpegCfg.FfmpegBinPath)

	r = append(r, "-i", tc.inputFile)

	r = append(r, tc.opts.GetStrArguments()...)

	r = append(r, tc.tmpOutputFile)

	return r, nil

}

func (tc *Transcoder) Run() error {

	progress, err := ffmpeg.
		New(&tc.ffmpegCfg).
		Input("/tmp/avi").
		Output("/tmp/mp4").
		Start(tc.opts)

	if err != nil {
		return err
	}

	for msg := range progress {
		log.Printf("%+v", msg)
	}
	return nil
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
