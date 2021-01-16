#!/bin/bash

OSFLAG :=
ifeq ($(OS),Windows_NT)
	OSFLAG += win32
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		OSFLAG += amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OSFLAG += ia32
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OSFLAG += linux
	endif
	ifeq ($(UNAME_S),Darwin)
		OSFLAG += osx
	endif
		UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		OSFLAG += amd64
	endif
		ifneq ($(filter %86,$(UNAME_P)),)
	OSFLAG += ia32
		endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		OSFLAG += arm
	endif
endif

help:
	@echo "config: initiate template configuration"
	@echo "build: create executable files"
	@echo "check: check os type and architechture"

config:
	@cp -n config.yaml.example config.yaml
	@echo "please set value in config.yaml"

build:
	@go mod vendor
	@./script/go-build.sh

check:
	@echo $(OSFLAG)