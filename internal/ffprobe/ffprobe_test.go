package ffprobe

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestProbeVideo(t *testing.T) {

	tcs := []struct {
		name   string
		in     string
		expect ProbeData
	}{
		{
			name:   "mp4",
			in:     "testdata/video.mp4",
			expect: kodakProbe,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ff, err := New()
			if err != nil {
				t.Fatal(err)
			}

			got, err := ff.Probe(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}

var index0 = 0
var index1 = 1
var kodakProbe = ProbeData{
	Streams: []Stream{
		{
			Index:            &index0,
			ID:               "0x1",
			CodecName:        "h264",
			CodecLongName:    "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
			CodecType:        "video",
			CodecTimeBase:    "",
			CodecTagString:   "avc1",
			CodecTag:         "0x31637661",
			RFrameRate:       "25/1",
			AvgFrameRate:     "25/1",
			TimeBase:         "1/12800",
			StartPts:         0,
			StartTime:        "0.000000",
			DurationTs:       0x7000,
			Duration:         "2.240000",
			BitRate:          "190710",
			BitsPerRawSample: "8",
			NbFrames:         "56",
			Disposition: StreamDisposition{
				Default: 1,
			},
			Tags: StreamTags{
				Language: "und",
			},
			Profile:            "Constrained Baseline",
			Width:              320,
			Height:             180,
			HasBFrames:         0,
			SampleAspectRatio:  "",
			DisplayAspectRatio: "",
			PixFmt:             "yuv420p",
			Level:              12,
			ColorRange:         "tv",
			ColorSpace:         "bt709",
			SampleFmt:          "",
			SampleRate:         "",
			Channels:           0,
			ChannelLayout:      "",
			BitsPerSample:      0,
		},
		{
			Index:            &index1,
			ID:               "0x2",
			CodecName:        "aac",
			CodecLongName:    "AAC (Advanced Audio Coding)",
			CodecType:        "audio",
			CodecTimeBase:    "",
			CodecTagString:   "mp4a",
			CodecTag:         "0x6134706d",
			RFrameRate:       "0/0",
			AvgFrameRate:     "0/0",
			TimeBase:         "1/48000",
			StartPts:         0,
			StartTime:        "0.000000",
			DurationTs:       0x01a400,
			Duration:         "2.240000",
			BitRate:          "131331",
			BitsPerRawSample: "",
			NbFrames:         "106",
			Disposition: StreamDisposition{
				Default: 1,
			},
			Tags: StreamTags{
				Language: "und",
			},
			Profile:            "LC",
			Width:              0,
			Height:             0,
			HasBFrames:         0,
			SampleAspectRatio:  "",
			DisplayAspectRatio: "",
			PixFmt:             "",
			Level:              0,
			ColorRange:         "",
			ColorSpace:         "",
			SampleFmt:          "fltp",
			SampleRate:         "48000",
			Channels:           2,
			ChannelLayout:      "stereo",
			BitsPerSample:      0,
		},
	},
	Format: Format{
		Filename:         "testdata/video.mp4",
		NBStreams:        2,
		NBPrograms:       0,
		FormatName:       "mov,mp4,m4a,3gp,3g2,mj2",
		FormatLongName:   "QuickTime / MOV",
		StartTimeSeconds: 0,
		DurationSeconds:  2.262,
		Size:             "93142",
		BitRate:          "329414",
		ProbeScore:       100,
		Tags: FormatTags{
			MajorBrand:       "isom",
			MinorVersion:     "512",
			CompatibleBrands: "isomiso2avc1mp41",
			CreationTime:     "",
		},
	},
}
