#!/bin/bash
__main() {
    _name="deepseek-transfer"

    GOOS=linux GOARCH=amd64 go build -o "${_name}_linux" main.go

    GOOS=darwin GOARCH=amd64 go build -o "${_name}_mac" main.go

    GOOS=windows GOARCH=amd64 go build -o "${_name}_windows.exe" main.go

}
__main