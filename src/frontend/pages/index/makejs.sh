#!/bin/bash.exe
export GOOS=js
export GOARCH=ecmascript
#export GOPHERJS_GOROOT="/c/Progs/Go1.12.16"

gopherjs build -v -m -o ../../../../Dist/index.js