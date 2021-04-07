# Internals

### internal/transcode

This package contains two items needed to transcode video correctly:
* The profile struct uses a template, loaded on creation, to transform that in a slice of args to be passed to ffmpeg
* The Transcoder represents a single execution of a transcoding action for one profile, the creates a new transcoder
every time it wants to apply a new transcoding profile to a video.

