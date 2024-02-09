package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/AndresBott/videoconv/app/videoconv/config"
	"github.com/AndresBott/videoconv/internal/ffprobe"
	"github.com/AndresBott/videoconv/internal/tmpl"
	"github.com/spf13/cobra"
	"strings"
)

type templateData struct {
	Args      []string `json:"args"`
	Extension string   `json:"extension"`
}

type videoData struct {
	Video   ffprobe.ProbeData
	Profile map[string]string
}

func ProbeCmd() *cobra.Command {
	configfile := "videoconv.yaml"
	template := ""

	cmd := cobra.Command{
		Use:   "probe <video>",
		Args:  cobra.ExactArgs(1),
		Short: "generate a json render of the ffmprobe output",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := config.NewFromFile(configfile)
			if err != nil {
				return err
			}

			ff, err := ffprobe.New(cfg.FfprobePath)
			if err != nil {
				return err
			}

			probeData, err := ff.Probe(args[0])
			if err != nil {
				return err
			}
			data := videoData{
				Video:   probeData,
				Profile: map[string]string{},
			}

			// print the raw data passed into the template
			if template == "" {
				b, err := json.MarshalIndent(data, "", "    ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
				return nil
			}
			// trying to render template
			tmplFile, err := tmpl.FindTemplate(cfg.TmplDirs, template)
			if err != nil {
				return err
			}

			profileTmpl, err := tmpl.NewTmplFromFile(tmplFile)
			if err != nil {
				return err
			}

			tmplData := templateData{}
			err = profileTmpl.ParseJson(data, &tmplData)
			if err != nil {
				return fmt.Errorf("error parsing template: %v", err)
			}
			tmplData.Args = dropEmpty(tmplData.Args)

			b2, err := json.MarshalIndent(tmplData, "", "    ")
			if err != nil {
				return err
			}
			fmt.Println(string(b2))
			return nil

		},
	}

	cmd.Flags().StringVarP(&configfile, "config", "c", configfile, "configuration file")
	cmd.Flags().StringVarP(&template, "template", "t", "", "optional parse a template")

	return &cmd

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
