# mp3inspect

This library is still in the alpha phase, and in need of documentation (I'm working on it!). It has been tested againts CBR files but needs more testing against VBR files.  Additionally it does not currently scan for Xing or iVBR headers yet.

To use simply build the inspector utility:

    go build inspector.go

Or mp3inspect can be used as a library by calling `go get` on this repo and importing `github.com/betamike/mp3inspect/mp3`. The `Scanner` type can be used to scan the mp3 frame by frame, or the convenience method `InspectFile` will scan the whole file and provide a condensed `MP3Info` struct.

# License

This package is distributed under the MIT license.  See the LICENSE file for more details.
