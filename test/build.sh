#!/bin/bash

echo "(Flogo OSS) Build Executable for local platform !!"

export GO111MODULE=on

rm -rf app

flogo create -f $1.json app

cd ./app

flogo build -e --verbose

