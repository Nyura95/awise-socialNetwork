# Messenger Appointrip

This package is the http/socket server for the chat of https://appointrip.com

---

- [Install](#install)
- [Explanation](#Explanation)
- [Run](#Run)
- [Dependances](#Dependances)
- [Build](#Build)

---

## Install

First, install golang on your computer [install GOlang](https://golang.org/doc/install)

After install, go to '/Users/username/go/src' and copy/paste into this folder the project

run a terminal and :

```sh
go install
go run main.go
```

## Run

For run the project, open a cmd with Go installed and run

```sh
go run main.go
```

## Explanation

See the package for read the doc

## Build

Build for ubuntu run

```
GOOS=linux GOARCH=amd64 go build -o socialnetwork main.go
```

Build for macos run

```
GOOS=darwin GOARCH=amd64 go build -o socialnetwork main.go
```

Build for window run

```
GOOS=windows GOARCH=amd64 go build -o socialnetwork.exe main.go
```

This command create a binary file named 'socialnetwork' or 'socialnetwork.exe'
