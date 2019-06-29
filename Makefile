export PATH := $(GOPATH)/bin:$(PATH)
BIN=./bin
FLAGS=-mod=vendor

ifeq ($(OS),Windows_NT)
	EXT=.exe
else
	EXT=
endif

all: launcher

launcher:
	go build $(FLAGS) -o $(BIN)/launcher$(EXT) ./cmd/launcher
