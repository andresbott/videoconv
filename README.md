# videoconv

Batch video conversion based on directory observation

videoconv searches for video files in a given directory and applies transcoding templates to those video files to
move the finalized files to an output directory

Multiple locations with different templates are supported.

videoconv can run as cli or daemon regularly looking for new files that have been dropped in the observation directory


# Getting started

    # create a simple directory stucture
    $ videoconv init

    ./videoconv.yaml
    ./main/in
    ./main/fail
    ./main/out
    ./main/tmp
    

# Improvement ideas
* add profile flag -hwaccel for hw decoding

* add profile tempaltes that are popilated with key/value pairs from profiles

* us a task runner to allow multiple executions to run in parallel, e.g when using GPU encoding, and you have more than one gpu

* add parameter npp scale,
  
* use scale paramteres to add or remove the arg to the cmd based on the video size

# build

    goreleaser release --rm-dist --skip-publish --skip-validate
