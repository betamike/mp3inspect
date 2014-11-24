# mp3inspect

This library is still in the alpha phase, and in need of documentation (I'm working on it!). It has been tested against
a decent number of CBR files but needs more testing against VBR files. Additionally it does not currently scan for Xing
or iVBR headers.

To use simply install the `mp3inspect` utility:

    go get github.com/betamike/mp3inspect
    go install github.com/betamike/mp3inspect

Or mp3inspect can be used as a library by importing `github.com/betamike/mp3inspect/mp3`. The `Scanner` type can be used
to scan an mp3 frame by frame, or the convenience method `InspectFile` will scan an entire file and provide a condensed
`MP3Info` struct.

# Todo

- Add tests for all currently implemented features
- Handle Xing and iVBR headers
- Detect APE tags
- Parse metadata from tags ?
- Allow editing of tags ?

# License

This package is distributed under the MIT license.  See the LICENSE file for more details.
