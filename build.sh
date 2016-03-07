#!/bin/bash

version=v1.0.0

rm -rf dist

darwin() {
  # Mac Os X
  GOARCH=amd64 GOOS=darwin go build -ldflags "-X main.version=$version" -o dist/darwin/ft
  zip -q -j dist/darwin.zip dist/darwin/ft
}

amd64() {
  # Linux AMD64
  GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version=$version" -o dist/amd64/ft
  zip -q -j dist/amd64.zip dist/amd64/ft
}

linux386() {
  # Edison
  GOARCH=386 GOOS=linux go build -ldflags "-X main.version=$version" -o dist/386/ft
  zip -q -j dist/386.zip dist/386/ft
}

arm7() {
  # Beaglebone
  GOARCH=arm GOOS=linux GOARM=7 go build -ldflags "-X main.version=$version" -o dist/arm7/ft
  zip -q -j dist/arm7.zip dist/arm7/ft
}

arm6() {
  # RaspberryPi
  GOARCH=arm GOOS=linux GOARM=6 go build -ldflags "-X main.version=$version" -o dist/arm6/ft
  zip -q -j dist/arm6.zip dist/arm6/ft
}

darwin &
amd64 &
linux386 &
arm7 &
arm6 &
