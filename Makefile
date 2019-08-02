export PATH := $(GOPATH)/bin:$(PATH)
BIN=./bin
FLAGS=-mod=vendor

ifeq ($(OS),Windows_NT)
	EXT=.exe
else
	EXT=
endif

all: launcher dataSync autoLogin c4ldr

launcher:
	go build $(FLAGS) -o $(BIN)/launcher$(EXT) ./cmd/launcher

dataSync:
	go build $(FLAGS) -o $(BIN)/dataSync$(EXT) ./cmd/dataSync

autoLogin:
	go build $(FLAGS) -o $(BIN)/autoLogin$(EXT) ./cmd/autoLogin

staffAttendance:
	go build $(FLAGS) -o $(BIN)/staffAttendance$(EXT) ./cmd/staffAttendance

c4ldr:
	go build $(FLAGS) -o $(BIN)/c4ldr$(EXT) ./cmd/c4ldr
