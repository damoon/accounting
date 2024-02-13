#!env bash

set -euxo pipefail

go run ./bin/pdf-to-txt/main.go | while read line 
do
    pdftotext -layout "$line"
done
