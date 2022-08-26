# from https://freshman.tech/snippets/go/cross-compile-go-programs/
# Windows
# 64-bit
GOOS=windows GOARCH=amd64 go build -o bin/accesschecker-amd64.exe main.go
# 32-bit
GOOS=windows GOARCH=386 go build -o bin/accesschecker-386.exe main.go

# Mac
# 64-bit
GOOS=darwin GOARCH=amd64 go build -o bin/accesschecker-amd64-darwin main.go
# 32-bit
# has error: "go: unsupported GOOS/GOARCH pair darwin/386"
# GOOS=darwin GOARCH=386 go build -o bin/app-386-darwin app.go

# Linux
# 64-bit
GOOS=linux GOARCH=amd64 go build -o bin/accesschecker-amd64-linux main.go
# 32-bit
GOOS=linux GOARCH=386 go build -o bin/accesschecker-386-linux main.go
