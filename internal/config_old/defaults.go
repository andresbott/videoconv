package config_old

const (
	Version             = "0.1.1"
	defaultThread       = 1
	defaultPollInterval = 300 // seconds
	configFileName      = "videoconv.config"
)

var defaultVideoExtension = []videoExtension{
	{
		extension: "wmv",
		videoSettingName: []string{
			"960_encoded",
			"720_encoded",
			"half_encoded",
		},
	},
	{
		extension: "flv",
		videoSettingName: []string{
			"960_encoded",
			"720_encoded",
			"half_encoded",
		},
	},
	{
		extension: "webm",
		videoSettingName: []string{
			"960_encoded",
			"720_encoded",
			"half_encoded",
		},
	},
	{
		extension: "mov",
		videoSettingName: []string{
			"960_encoded",
			"720_encoded",
			"half_encoded",
		},
	},
	{
		extension: "mp4",
		videoSettingName: []string{
			"960",
			"720",
			"half",
		},
	},
}

var defaultVideoSettings = []VideoSetting{

	{
		name:         "480_encoded",
		cmd:          `-vf 'scale=-2:min(480\,ih-mod(ih\,2))' -strict -2 -c:v libx264 -crf 23 -preset veryslow`,
		outExtension: "mp4",
	},
	{
		name:         "720_encoded",
		cmd:          `-vf 'scale=-2:min(720\,ih-mod(ih\,2))' -strict -2 -c:v libx264 -crf 23 -preset veryslow`,
		outExtension: "mp4",
	},
	{
		name:         "960_encoded",
		cmd:          `-vf 'scale=-2:min(960\,ih-mod(ih\,2))' -strict -2 -c:v libx264 -crf 23 -preset veryslow`,
		outExtension: "mp4",
	},
	{
		name:         "1280_encoded",
		cmd:          `-vf 'scale=-2:min(1280\,ih-mod(ih\,2))' -strict -2 -c:v libx264 -crf 23 -preset veryslow`,
		outExtension: "mp4",
	},
	{
		name:         "half_encoded",
		cmd:          "-vf scale=iw*.5:-1 -strict -2 -c:v libx264 -crf 23 -preset veryslow",
		outExtension: "mp4",
	},

	{
		name:         "480",
		cmd:          "-vf 'scale=-2:min(780\\,ih-mod(ih\\,2))' -strict -2",
		outExtension: "mp4",
	},
	{
		name:         "720",
		cmd:          "-vf 'scale=-2:min(720\\,ih-mod(ih\\,2))' -strict -2",
		outExtension: "mp4",
	},
	{
		name:         "960",
		cmd:          "-vf 'scale=-2:min(960\\,ih-mod(ih\\,2))' -strict -2",
		outExtension: "mp4",
	},
	{
		name:         "1280",
		cmd:          "-vf 'scale=-2:min(1280\\,ih-mod(ih\\,2))' -strict -2",
		outExtension: "mp4",
	},
	{
		name:         "half",
		cmd:          "-vf scale=iw*.5:-1 -strict -2",
		outExtension: "mp4",
	},
}
