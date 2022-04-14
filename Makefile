SRCS	:= $(shell find . -type d -name archive -prune -o -type f -name '*.go')
LDFLAGS	:= -ldflags="-s -w -extldflags \"-static\""
NAME    := bin/yaml-merge

build: $(SRCS)
	CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o $(NAME) pkg/main.go
