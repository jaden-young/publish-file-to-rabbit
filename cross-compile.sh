#!/bin/sh

docker run --rm -v "$PWD":/go/src/app -w /go/src/app golang:1.10 \
  /bin/sh -c 'for GOOS in darwin linux windows; do \
                for GOARCH in 386 amd64; do \
                  go build -v -o publish-file-to-rabbit-$GOOS-$GOARCH 
                done 
              done'