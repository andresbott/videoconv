package videoconv

import (
	"fmt"
	"github.com/AndresBott/videoconv/app/videoconv/config"
	"github.com/AndresBott/videoconv/internal/ffmpegtranscode"
	"github.com/AndresBott/videoconv/internal/tmpl"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Converter struct {
	Cfg        config.Conf
	DaemonMode bool
	ffmpeg     *ffmpegtranscode.Transcoder
}

// used for testing only
var processFn func(absVideo, absIn, absOut, absTmp, absFail string, profiles []config.Profile)

func New(cfg config.Conf) (*Converter, error) {

	ffmpeg, err := ffmpegtranscode.New(ffmpegtranscode.Cfg{
		FfmpegBin:  cfg.FfmpegPath,
		FfprobeBin: cfg.FfprobePath,
	})
	if err != nil {
		return nil, err
	}

	c := Converter{
		Cfg:    cfg,
		ffmpeg: ffmpeg,
	}

	// log level (not sure if I like this here)
	lv, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		lv = log.InfoLevel
	}
	log.SetLevel(lv)

	return &c, nil
}

func NewFromCfg(cfgFile string) (*Converter, error) {
	cfg, err := config.NewFromFile(cfgFile)
	if err != nil {
		return nil, err
	}
	return New(cfg)
}

func (vc *Converter) Check(create bool) error {

	for _, item := range vc.Cfg.Locations {
		e := checkLocation(vc.Cfg.ConfigLocation, item, create)
		if e != nil {
			return e
		}
	}
	return nil
}

func checkLocation(cfgLocation string, location config.Location, create bool) error {
	itemPath, err := filepath.Abs(filepath.Join(filepath.Dir(cfgLocation), location.Path))
	if err != nil {
		return fmt.Errorf("error generating absolute path for location: %v", err)
	}
	log.Debug("checking location:" + itemPath)

	if _, err := os.Stat(itemPath); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist or is not accessible", itemPath)
	}

	dirs := []string{
		location.InputDir,
		location.OutputDir,
		location.TmpDir,
		location.FailDir,
	}

	for _, d := range dirs {
		dir, err := filepath.Abs(filepath.Join(itemPath, d))
		if err != nil {
			return fmt.Errorf("error generating absolute path for location: %v", err)
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if !create {
				return fmt.Errorf("directory \"%s\" does not exits", dir)
			}

			log.Warnf("directory %s does not exist or is not accessible, trying to create", dir)
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return fmt.Errorf("unable to create dir %s, %v", dir, err)
			}
		}
	}
	return nil
}

// Run executes the main conversion loop, will exit if not run in daemon mode
func (vc *Converter) Run() {
	log.Info("starting video conversion...")
	for {
		for _, location := range vc.Cfg.Locations {
			vc.runLocation(location)
		}
		if !vc.DaemonMode {
			log.Info("finished, exiting...")
			break
		}
		log.Infof("finished, sleeping for %s", vc.Cfg.Sleep)
		time.Sleep(vc.Cfg.Sleep)
	}
}

// convert Videos on one location
func (vc *Converter) runLocation(location config.Location) {
	log.Debug("checking location:" + location.Path)
	locationPath, err := filepath.Abs(filepath.Join(filepath.Dir(vc.Cfg.ConfigLocation), location.Path))
	if err != nil {
		log.Errorf("error generating absolute path for location \"%s\", skipping location, error: %v", locationPath, err)
		return
	}
	if _, err := os.Stat(locationPath); os.IsNotExist(err) {
		log.Errorf("directory %s does not exist or is not accessible, skipping location", locationPath)
		return
	}

	err = checkLocation(vc.Cfg.ConfigLocation, location, false)
	if err != nil {
		log.Errorf("location contains error:\"%s\", sskipping...", err.Error())
	}

	videos, err := findVideos(filepath.Join(locationPath, location.InputDir), vc.Cfg.VideoExtensions)
	if err != nil {
		log.Errorf("error searching for Videos:\"%s\", sskipping...", err.Error())
	}

	for _, video := range videos {

		videoPath := filepath.Join(locationPath, location.InputDir, video)
		absInPath := filepath.Join(locationPath, location.InputDir)
		absOutPath := filepath.Join(locationPath, location.OutputDir)
		absTmpPath := filepath.Join(locationPath, location.TmpDir)
		absFailPath := filepath.Join(locationPath, location.FailDir)

		if processFn != nil {
			processFn(videoPath, absInPath, absOutPath, absTmpPath, absFailPath, location.Profiles) // used for testing purposes
		} else {
			vc.processVideo(videoPath, absInPath, absOutPath, absTmpPath, absFailPath, location.Profiles)
		}
	}
}

func renameFile(in, profileName string, overwriteExtension string) string {
	baseName := filepath.Base(in)
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	name = name + "." + profileName
	ext := filepath.Ext(baseName)
	if overwriteExtension != "" {
		ext = overwriteExtension
	}
	ext = strings.Trim(ext, ".")
	name = name + "." + ext
	return name
}

// processVideo is responsible for taking one video and generate all the renditions as per profile configuration
// it returns the list of done Videos, if something goes wrong an error is returned as well
func (vc *Converter) processVideo(absVideo, absIn, absOut, absTmp, absFail string, profiles []config.Profile) {
	log.Infof("procesing video: \"%s\"", filepath.Base(absVideo))

	cmd := []string{}
	err := func() error {

		doneVideos := []string{}
		for _, profile := range profiles {
			tmplFile, err := tmpl.FindTemplate(vc.Cfg.TmplDirs, profile.Template)
			if err != nil {
				return err
			}
			profileTmpl, err := tmpl.NewTmplFromFile(tmplFile)
			if err != nil {
				return err
			}
			if profile.Name == "" {
				return fmt.Errorf("profile name cannot be empty")
			}

			// todo put in ffmpeg probe + profile data
			tmplData, err := profileTmpl.Parse("")
			if err != nil {
				return fmt.Errorf("error parsing template: %v", err)
			}
			outFileName := renameFile(filepath.Base(absVideo), profile.Name, tmplData.FileExt)

			tmpFilePath := filepath.Join(absTmp, outFileName)
			// delete a potential tmp output file before starting a new conversion
			if _, err := os.Stat(tmpFilePath); err == nil {
				log.Warn("deleting OLD tmp file: " + outFileName)
				e := os.Remove(tmpFilePath)
				if e != nil {
					return fmt.Errorf("unable to delete temp file %s, error: %v ", tmpFilePath, e)

				}
			}

			// todo test error in ffmpeg
			cmd, err = vc.ffmpeg.Run(absVideo, tmpFilePath, tmplData.Args)
			if err != nil {
				return fmt.Errorf("error trancoding video: %v", err)
			}
			doneVideos = append(doneVideos, tmpFilePath)
		}

		relativePath, err := filepath.Rel(absIn, absVideo)
		if err != nil {
			return fmt.Errorf("error getting the relative path: %v", err)
		}
		// create output directories
		destPath := filepath.Join(absOut, filepath.Dir(relativePath))
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			err = os.MkdirAll(destPath, 0755)
			if err != nil {
				return fmt.Errorf("unable to create folder: %s, error: %v ", destPath, err)
			}
		}

		//move the converted files
		for _, f := range doneVideos {
			outFile := filepath.Join(absOut, filepath.Dir(relativePath), filepath.Base(f))
			//spew.Dump(fmt.Sprintf("move video from %s to %s", f, outFile))

			err := os.Rename(f, outFile)
			if err != nil {
				return fmt.Errorf("unable to move file %s to %s, error: %v ", f, outFile, err)
			}
		}

		outFile := filepath.Join(absOut, filepath.Dir(relativePath), filepath.Base(filepath.Base(absVideo)))
		err = os.Rename(absVideo, outFile)
		if err != nil {
			return fmt.Errorf("unable to move file \"%s\", error: %v ", filepath.Base(absVideo), err)
		}

		return nil

	}()
	if err != nil {
		relativePath, err2 := filepath.Rel(absIn, absVideo)
		if err2 != nil {
			panic(fmt.Errorf("failure in getting relative path during error handling: %s", err))
		}

		log.Errorf("Error transcoding video: \"%s\" %s", relativePath, err)
		if len(cmd) > 0 {
			log.Errorf("command run: \"%s\"", strings.Join(cmd, " "))
		}

		// create output directories
		failPath := filepath.Join(absFail, filepath.Dir(relativePath))
		if _, err3 := os.Stat(failPath); os.IsNotExist(err3) {
			err3 = os.MkdirAll(failPath, 0755)
			if err3 != nil {
				panic(fmt.Errorf("unable to create folder \"%s\" during error handling, error: %v ", failPath, err3))
			}
		}
		// move failed video
		failOut := filepath.Join(failPath, filepath.Base(absVideo))
		err = os.Rename(absVideo, failOut)
		if err != nil {
			panic(fmt.Errorf("unable to move file \"%s\" to failed location, error: %v ", filepath.Base(absVideo), err))
		}

	}

	//// todo move to failed in case of error
	//if err != nil {
	//	log.Errorf("error processing video: %s, %s", video, err.Error())
	//}
	//
	//return nil

}

// findVideos recursively searches Videos in the rootPath and returns an array of relative paths of Videos
func findVideos(rootPath string, videoExtensions []string) ([]string, error) {
	var videos []string

	err := filepath.Walk(rootPath, func(fPath string, fInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if fInfo.IsDir() {
			return nil
		}
		ext := filepath.Ext(fPath)
		ext = ext[1:]
		if !isVideo(ext, videoExtensions) {
			return nil
		}

		rel, err := filepath.Rel(rootPath, fPath)
		if err != nil {
			return err
		}
		videos = append(videos, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return videos, nil
}

// used to check if file is a video based on extension
func isVideo(val string, videoExtensions []string) bool {
	c := strings.TrimSpace(val)
	c = strings.ToLower(c)

	for _, item := range videoExtensions {
		if strings.ToLower(item) == c {
			return true
		}
	}
	return false
}
