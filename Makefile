SHELL := /bin/bash
PRODUCT_NAME := net-test
BIN := dist/${PRODUCT_NAME}
DNAME := ${PRODUCT_NAME}_
ENTRY_POINT := ./cmd/${PRODUCT_NAME}
HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
VERSION := ${HASH}
LD_FLAGS := -s -w -X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.buildName=${PRODUCT_NAME}'
COMP_BIN := go

ifeq ($(OS),Windows_NT)
	BIN := $(BIN).exe
	DNAME := $(DNAME).exe
endif

.PHONY: build dev test clean deploy

build:
	mkdir -p dist/
	${COMP_BIN} build -o "${BIN}" -ldflags="${LD_FLAGS}" ${ENTRY_POINT}

dev:
	mkdir -p dist/
	${COMP_BIN} build -tags debug -o "${BIN}" -ldflags="${LD_FLAGS}" ${ENTRY_POINT}
	./${BIN}

test:
	${COMP_BIN} test

clean:
	${COMP_BIN} clean
	rm -r -f dist/

deploy: build
	sudo systemctl stop "${PRODUCT_NAME}"
	sudo cp "${BIN}" /usr/local/bin
	sudo systemctl start "${PRODUCT_NAME}"

