env GOOS=windows GOARCH=amd64 go build -o builds/cupt_windows64.exe
env GOOS=darwin GOARCH=amd64 go build -o builds/cupt_darwin64
env GOOS=linux GOARCH=amd64 go build -o builds/cupt_linux64