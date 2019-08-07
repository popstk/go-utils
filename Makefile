export PATH := $(GOPATH)/bin:$(PATH)
BIN=./bin
FLAGS=-mod=vendor

ifeq ($(OS),Windows_NT)
	EXT=.exe
else
	EXT=
endif

all: launcher autoLogin zookeeper

launcher:
	go build $(FLAGS) -o $(BIN)/launcher$(EXT) ./cmd/launcher

autoLogin:
	go build $(FLAGS) -o $(BIN)/autoLogin$(EXT) ./cmd/autoLogin

staffAttendance:
	go build $(FLAGS) -o $(BIN)/staffAttendance$(EXT) ./cmd/staffAttendance

zookeeper:
	go build $(FLAGS) -o $(BIN)/zookeeper$(EXT) ./cmd/zookeeper
