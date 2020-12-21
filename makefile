#!/bin/bash

help:
	@echo "build: create executable files"
	@echo "build: create executable files"

build:
	@./go-build.sh

config:
	@cp -n config.yaml.example config.yaml

