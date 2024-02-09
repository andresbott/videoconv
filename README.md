# videoconv

Batch video conversion based on directory observation

videoconv searches recursively for video files in a given (input) directory and applies transcoding 
profiles to those videos, finally it moves the done files to an output directory.

Multiple locations with different templates are supported.

videoconv can run as cli or daemon regularly looking for new files that have been dropped in the observation directory.


## Getting started

    # create a simple directory stucture
    $ videoconv init

    ./videoconv.yaml
    ./sample/in
    ./sample/fail
    ./sample/out
    ./sample/tmp
    
## For developers

### Improvement ideas

* use a task runner to allow multiple executions to run in parallel, 
e.g when using GPU encoding, and you have more than one gpu

### TODOS
* ffprobe already provides some precalculated data, e.g. 
  * simplified informationa bout audio stremas
* use json5 to allow comments in json
* probe and run use copied code to render the template => unify
* add dry run to print actions but not execute
### Build

    goreleaser release --rm-dist --skip-publish --skip-validate
