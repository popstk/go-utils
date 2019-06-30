export PATH := $(GOPATH)/bin:$(PATH)
BIN=./bin
FLAGS=-mod=vendor

ifeq ($(OS),Windows_NT)
	EXT=.exe
else
	EXT=
endif

all: launcher dataSync

launcher:
	go build $(FLAGS) -o $(BIN)/launcher$(EXT) ./cmd/launcher

dataSync:
	go build $(FLAGS) -o $(BIN)/dataSync$(EXT) ./cmd/dataSync

