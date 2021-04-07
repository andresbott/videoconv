package videconv

import (
	"fmt"
	"github.com/AndresBott/videoconv/internal/transcoder"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type App struct {
	ConfigFile      string
	DaemonMode      bool
	logLevel        string
	ffmpegBin       string
	ffProbeBin      string
	videoExtensions []string
	threads         int
	sleep           time.Duration
	locations       []location
	OverlayFname    string
	profiles        map[string]profile
}

// start the video conversion app
func (vc *App) Start() error {

	err := vc.loadConfig()
	if err != nil {
		return fmt.Errorf("error while loading configuration: %v", err)
	}
	log.Info("starting video conversion...")
	vc.loop()
	return nil
}

// main execution loop, will exit if not run in daemon mode
func (vc *App) loop() {
	for true {
		for _, l := range vc.locations {
			vc.runLocation(l)
		}
		if !vc.DaemonMode {
			log.Info("finished, exiting...")
			break
		}
		log.Infof("finished, sleeping for %s", vc.sleep)
		time.Sleep(vc.sleep)
	}
}

// execute one video on one location and return true after a video conversion is finished
// in case there are no videos false is returned, errors are swallowed and only logged
// videos that contain errors are moved to fail location and true is returned anyway.
func (vc *App) runLocation(l location) {
	log.Debug("checking location:" + l.path)

	dir, _ := filepath.Abs(l.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Warnf("directory %s does not exist or is not accessible, skipping location", dir)
		return
	}

	dirs := []string{
		"./",
		l.inputDir,
		l.outputDir,
		l.tmpDir,
		l.failDir,
	}

	for _, d := range dirs {
		dir, _ := filepath.Abs(l.path + "/" + d)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Warnf("directory %s does not exist or is not accessible, trying to create", dir)
			err := os.Mkdir(dir, 0755)
			if err != nil {
				log.Fatalf("unable to create dir %s, %v", dir, err)
			}
		}
	}

	// load location specific configuration
	// todo remove overlay configuration: not needed
	newLoc, err := l.loadOverlay(vc.OverlayFname)
	if err != nil {
		log.Errorf("error while loading overlay: %s", err)
		return
	}

	for {
		// search for videos in the location
		video, err := vc.findVideoFile(newLoc.path + "/" + newLoc.inputDir)
		if err != nil {
			log.Errorf("error while searching videos: %s", err)
			return
		}
		if video == "" {
			return
		}
		vc.transcodeVideo(video, newLoc)
	}
}

// transcodeVideo takes one video and will transcode it to all the configured profiles
func (vc *App) transcodeVideo(video string, location *location) {
	start := time.Now()

	log.Infof("transcoding video: \"%s\"", video)

	absInDir, err := filepath.Abs(location.path + "/" + location.inputDir)
	if err != nil {
		r := fmt.Errorf("error getting the absolute path for input file: %s", err.Error())
		log.Error(r)
	}
	absOutDir, err := filepath.Abs(location.path + "/" + location.outputDir)
	if err != nil {
		r := fmt.Errorf("error getting the absolute path for input file: %s", err.Error())
		log.Error(r)
	}
	absTmpDir, err := filepath.Abs(location.path + "/" + location.tmpDir)
	if err != nil {
		r := fmt.Errorf("error getting the absolute path for input file: %s", err.Error())
		log.Error(r)
	}

	var toBeMoved []string

	for _, profName := range location.appliedProfiles {
		if prf, ok := vc.profiles[profName]; ok {

			log.Infof("running profile: \"%s\"", profName)
			pStart := time.Now()

			spew.Dump(prf)
			// my_original_video.<profileName>.<profileExtension>
			outputFileName := filepath.Base(video) + "." + prf.name + "." + prf.extension

			cfg := transcoder.Cfg{
				FfmpegBin:  vc.ffmpegBin,
				FfProbeBin: vc.ffProbeBin,
				Template:   prf.template,
				InputFile:  filepath.Join(absInDir, video),
				OutputFile: filepath.Join(absTmpDir, outputFileName),
			}

			spew.Dump(cfg)

			tr, err := transcoder.New(&cfg)
			if err != nil {
				r := fmt.Errorf("error with transcoder \"%s\" args: %s", profName, err.Error())
				log.Error(r)
				return
			}

			// delete a potential tmp output file before starting a new conversion
			if _, err := os.Stat(tr.GetOutputFile()); err == nil {
				log.Warn("deleting OLD tmp file: " + tr.GetOutputFile())
				e := os.Remove(tr.GetOutputFile())
				if e != nil {
					log.Fatalf("unable to delete temp file %s, error: %v ", tr.GetOutputFile(), e)
				}
			}

			logCmd, _ := tr.GetCmd()
			log.Debugf("command: %s", strings.Join(logCmd, " "))

			commads, err := tr.Run()
			if err != nil {
				log.Errorf("error while transcoding video with profile \"%s\" args: %s, command: %s", profName, err.Error(), commads)

				// todo move file to failed
				log.Warn("TODO: move file to fail location")

				log.Warn("deleting temp file: " + tr.GetOutputFile())
				e := os.Remove(tr.GetOutputFile())
				if e != nil {
					log.Fatalf("unable to delete temp file %s, error: %v ", tr.GetOutputFile(), e)
				}
			}
			toBeMoved = append(toBeMoved, tr.GetOutputFile())
			log.Infof("profile \"%s\" took %s", profName, time.Since(pStart))

		} else {
			log.Warnf("profile: \"%s\" not found", profName)
		}
	}

	log.Infof("video processing took %s", time.Since(start))
	log.Infof("movig files to final destination: %s", absOutDir)

	destPath := filepath.Join(absOutDir, filepath.Dir(video))
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err = os.MkdirAll(destPath, 0755)
		if err != nil {
			log.Fatalf("unable to create folder: %s, error: %v ", destPath, err)
		}
	}

	// move the converted files
	for _, f := range toBeMoved {
		err := os.Rename(f, filepath.Join(destPath, filepath.Base(f)))
		if err != nil {
			log.Fatalf("unable to move file %s, error: %v ", f, err)
		}
	}

	// move the original file
	err = os.Rename(filepath.Join(absInDir, video), filepath.Join(absOutDir, video))
	if err != nil {
		log.Fatalf("unable to move file %s, error: %v ", filepath.Join(absInDir, video), err)
	}
}

// findVideoFiles recurses the provided root dir and searches video files
// returns the first found video file
func (vc *App) findVideoFile(rootPath string) (string, error) {

	var matches string

	err := filepath.Walk(rootPath, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fInfo.IsDir() {
			return nil
		}
		ext := filepath.Ext(fPath)
		ext = ext[1:]

		if !vc.isVideo(ext) {
			return nil
		}

		rel, err := filepath.Rel(rootPath, fPath)
		if err != nil {
			return err
		}
		matches = rel
		return io.EOF // exit early after the first match
	})
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return "", err
	}
	return matches, nil
}

// used to check if file is a video based on extansion
func (vc *App) isVideo(val string) bool {
	c := strings.TrimSpace(val)
	c = strings.ToLower(val)

	for _, item := range vc.videoExtensions {
		if strings.ToLower(item) == c {
			return true
		}
	}
	return false
}
