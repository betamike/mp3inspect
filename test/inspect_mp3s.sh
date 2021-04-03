#!/bin/bash

find ./files -name \*.mp3 -type f -exec sh -c 'echo {} && mp3inspect {}' \;
