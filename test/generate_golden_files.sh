#!/bin/bash

# This generates our golden files, which we will use to compare
# expected output to actual output
# ffprobe is used to get information about the actual bitrate, sampling rate
# and vbr/cbr of the file

# We're going to have this run in parallel using xargs, we want to get the core count
# to pass to xargs -P
platform='unknown'
unamestr=$(uname)

case "$unamestr" in
    Linux)
        echo "running on linux"
        platform='linux'
        ;;
    Darwin)
        echo "running on mac"
        platform='mac'
        ;;
    *)
        echo "could not determine platform, exiting"
        exit 1
esac

function generate_golden_file {
    file=$1

    bitrate="$(ffprobe "$file" -hide_banner 2>&1 | grep -E 'Audio.+Hz' | grep 'kb/s' | awk '{print $9}')000"
    sampling_rate="$(ffprobe "$file" -hide_banner 2>&1 | grep -E 'Audio.+Hz' | awk '{print $5}')"

    vbr=''

    if echo "$file" | grep 'k\.mp3'; then
        vbr='false'
    else
        vbr='true'
    fi

    if [ $vbr = '' ]; then
        echo "could not determine vbr/cbr, exiting"
        exit 1
    fi

    frame_count=$(fq '.frames | length' "$file")

    golden_file_name=".$(echo "${file}" | cut -d. -f2).golden"

    echo "$file - $golden_file_name - $frame_count - $bitrate - $sampling_rate - $vbr"

    tabs 8

# Our mp3inspect output template
cat << EOF > "$golden_file_name"
num frames	br(bps)	vbr	sr(hz)	ID3v1	ID3v2	v2 size(b)	bad bytes	bad frames
$frame_count	$bitrate	$vbr	$sampling_rate	false	false	0	0	0	
EOF

}

export -f generate_golden_file

# parallelism in bash, oh boy
if [ $platform = 'linux' ]; then
    number_of_cores=$(grep -c ^processor /proc/cpuinfo)
    find ./files -type f -name \*.mp3 -print0 | xargs -0 -P "$number_of_cores" -n1 bash -c 'generate_golden_file "$@"' _
elif [ $platform = 'mac' ]; then
    number_of_cores=$(sysctl -n hw.ncpu)
    find ./files -type f -name \*.mp3 -print0 | gxargs -0 -P "$number_of_cores" -n1 bash -c 'generate_golden_file "$@"' _
fi
