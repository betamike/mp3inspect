#!/bin/bash

# See https://en.wikipedia.org/wiki/MP3#Bit_rate
read -r -d '' valid_bit_rates_mpeg1_layer3 << EOF
32
40
48
56
64
80
96
112
128
160
192
224
256
320
EOF

read -r -d '' valid_bit_rates_mpeg2_layer3 << EOF
8
16
24
32
40
48
56
64
80
96
112
128
144
160
EOF

read -r -d '' valid_bit_rates_mpeg2_5_layer3 << EOF
8
16
24
32
40
48
56
64
EOF

read -r -d '' valid_sampling_rates_mpeg1_layer3 << EOF
32000
44100
48000
EOF

read -r -d '' valid_sampling_rates_mpeg2_layer3 << EOF
16000
22050
24000
EOF

read -r -d '' valid_sampling_rates_mpeg2_5_layer3 << EOF
8000
11025
12000
EOF

if [ ! -d "files" ]; then
    mkdir files
    cd files || exit 1
else
    cd files || exit 1
    rm -rf ./*
fi

while read -r sample_rate; do
    ffmpeg \
        -f lavfi \
        -i "sine=frequency=1000:sample_rate=$sample_rate:duration=3" \
        -c:a pcm_s16le \
        -ar "$sample_rate" \
        -ac 2 \
        "test_stereo_${sample_rate}Hz.wav" < /dev/null

    for quality in $(seq 0 9);
    do
        lame \
            -V "$quality" \
            "test_stereo_${sample_rate}Hz.wav" \
            "test_stereo_${sample_rate}Hz_v${quality}.mp3" < /dev/null
    done

    while read -r bit_rate; do
        lame \
            -b "$bit_rate" \
            "test_stereo_${sample_rate}Hz.wav" \
            "test_stereo_${sample_rate}Hz_${bit_rate}k.mp3" < /dev/null
    done <<< "$valid_bit_rates_mpeg1_layer3"
done <<< "$valid_sampling_rates_mpeg1_layer3"

while read -r sample_rate; do
    ffmpeg \
        -f lavfi \
        -i "sine=frequency=1000:sample_rate=$sample_rate:duration=3" \
        -c:a pcm_s16le \
        -ar "$sample_rate" \
        -ac 2 \
        "test_stereo_${sample_rate}Hz.wav" < /dev/null

    for quality in $(seq 0 9);
    do
        lame \
            -V "$quality" \
            "test_stereo_${sample_rate}Hz.wav" \
            "test_stereo_${sample_rate}Hz_v${quality}.mp3" < /dev/null
    done

    while read -r bit_rate; do
        lame \
            -b "$bit_rate" \
            "test_stereo_${sample_rate}Hz.wav" \
            "test_stereo_${sample_rate}Hz_${bit_rate}k.mp3" < /dev/null
    done <<< "$valid_bit_rates_mpeg2_layer3"
done <<< "$valid_sampling_rates_mpeg2_layer3"

while read -r sample_rate; do
    ffmpeg \
        -f lavfi \
        -i "sine=frequency=1000:sample_rate=$sample_rate:duration=3" \
        -c:a pcm_s16le \
        -ar "$sample_rate" \
        -ac 2 \
        "test_stereo_${sample_rate}Hz.wav" < /dev/null

    for quality in $(seq 0 9);
    do
        lame \
            -V "$quality" \
            "test_stereo_${sample_rate}Hz.wav" \
            "test_stereo_${sample_rate}Hz_v${quality}.mp3" < /dev/null
    done

    while read -r bit_rate; do
        lame \
            -b "$bit_rate" \
            "test_stereo_${sample_rate}Hz.wav" \
            "test_stereo_${sample_rate}Hz_${bit_rate}k.mp3" < /dev/null
    done <<< "$valid_bit_rates_mpeg2_5_layer3"
done <<< "$valid_sampling_rates_mpeg2_5_layer3"

rm -f *.wav
