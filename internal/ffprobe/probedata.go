package ffprobe

import (
	"math"
	"strconv"
)

// ProbeData is the root json data structure returned by an ffprobe.
type ProbeData struct {
	Format   Format     `json:"format"`
	Streams  []Stream   `json:"streams"`
	Chapters []Chapters `json:"chapters"`
	Summary  Summary    `json:"Summary"`
}

const CodecTypeVideo = "video"

func (p *ProbeData) Digest() {

	video := Video{}
	for _, s := range p.Streams {
		if s.CodecType == CodecTypeVideo {
			video.H = s.Height
			video.W = s.Width
			video.Format = s.CodecName

			break
		}
	}

	b, _ := strconv.Atoi(p.Format.BitRate)
	video.BitRate = b
	video.BitRateM = math.Round((float64(b)/1000000)*100) / 100

	p.Summary = Summary{
		Video: video,
	}

}

type Summary struct {
	Video Video
}

type Video struct {
	Format   string
	H        int
	W        int
	BitRate  int
	BitRateM float64
}

// Format is a json data structure to represent formats
type Format struct {
	Filename         string            `json:"filename"`
	NBStreams        int               `json:"nb_streams"`
	NBPrograms       int               `json:"nb_programs"`
	FormatName       string            `json:"format_name"`
	FormatLongName   string            `json:"format_long_name"`
	StartTimeSeconds float64           `json:"start_time,string"`
	DurationSeconds  float64           `json:"duration,string"`
	Size             string            `json:"size"`
	BitRate          string            `json:"bit_rate"`
	ProbeScore       int               `json:"probe_score"`
	Tags             map[string]string `json:"tags"`
}

type Chapters struct {
	ID        int64             `json:"id"`
	TimeBase  string            `json:"time_base"`
	Start     int64             `json:"start"`
	StartTime string            `json:"start_time"`
	End       int64             `json:"end"`
	EndTime   string            `json:"end_time"`
	Tags      map[string]string `json:"tags"`
}

// Stream is a json data structure to represent streams.
// A stream can be a video, audio, subtitle, etc type of stream.
type Stream struct {
	Index              *int              `json:"index"`
	ID                 string            `json:"id"`
	CodecName          string            `json:"codec_name"`
	CodecLongName      string            `json:"codec_long_name"`
	CodecType          string            `json:"codec_type"`
	CodecTimeBase      string            `json:"codec_time_base"`
	CodecTagString     string            `json:"codec_tag_string"`
	CodecTag           string            `json:"codec_tag"`
	RFrameRate         string            `json:"r_frame_rate"`
	AvgFrameRate       string            `json:"avg_frame_rate"`
	TimeBase           string            `json:"time_base"`
	StartPts           int               `json:"start_pts"`
	StartTime          string            `json:"start_time"`
	DurationTs         uint64            `json:"duration_ts"`
	Duration           string            `json:"duration"`
	BitRate            string            `json:"bit_rate"`
	BitsPerRawSample   string            `json:"bits_per_raw_sample"`
	NbFrames           string            `json:"nb_frames"`
	Disposition        StreamDisposition `json:"disposition,omitempty"`
	Tags               StreamTags        `json:"tags,omitempty"`
	Profile            string            `json:"profile,omitempty"`
	Width              int               `json:"width"`
	Height             int               `json:"height"`
	HasBFrames         int               `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string            `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string            `json:"display_aspect_ratio,omitempty"`
	PixFmt             string            `json:"pix_fmt,omitempty"`
	Level              int               `json:"level,omitempty"`
	ColorRange         string            `json:"color_range,omitempty"`
	ColorSpace         string            `json:"color_space,omitempty"`
	SampleFmt          string            `json:"sample_fmt,omitempty"`
	SampleRate         string            `json:"sample_rate,omitempty"`
	Channels           int               `json:"channels,omitempty"`
	ChannelLayout      string            `json:"channel_layout,omitempty"`
	BitsPerSample      int               `json:"bits_per_sample,omitempty"`
}

// StreamDisposition is a json data structure to represent stream dispositions
type StreamDisposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
}

// StreamTags is a json data structure to represent stream tags
type StreamTags struct {
	Rotate       int    `json:"rotate,string,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	Language     string `json:"language,omitempty"`
	Title        string `json:"title,omitempty"`
	Encoder      string `json:"encoder,omitempty"`
	Location     string `json:"location,omitempty"`
}
