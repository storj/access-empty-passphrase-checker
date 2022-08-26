# based partially on https://freshman.tech/snippets/go/cross-compile-go-programs/
# Windows
# 64-bit
GOOS=windows GOARCH=amd64 go build -o bin/accesschecker-amd64.exe main.go
# 32-bit
GOOS=windows GOARCH=386 go build -o bin/accesschecker-386.exe main.go

# Mac
# 64-bit
GOOS=darwin GOARCH=amd64 go build -o bin/accesschecker-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o bin/accesschecker-darwin-arm64 main.go
# 32-bit
# has error: "go: unsupported GOOS/GOARCH pair darwin/386"
# GOOS=darwin GOARCH=386 go build -o bin/accesschecker-386-darwin main.go

# Linux
# 64-bit
GOOS=linux GOARCH=amd64 go build -o bin/accesschecker-linux-amd64 main.go
GOOS=linux GOARCH=arm64 go build -o bin/accesschecker-linux-arm64 main.go
# 32-bit
GOOS=linux GOARCH=386 go build -o bin/accesschecker-linux-386 main.go
GOOS=linux GOARCH=arm go build -o bin/accesschecker-linux-arm main.go

# Free BSD
GOOS=freebsd GOARCH=amd64 go build -o bin/accesschecker-freebsd-amd64 main.go
