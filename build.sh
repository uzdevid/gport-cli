$env:GOOS="linux"
$env:GOARCH="386"
go build -o ./build/linux/386

$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o ./build/linux/amd64

$env:GOOS="linux"
$env:GOARCH="arm64"
go build -o ./build/linux/arm64

# windows

$env:GOOS="windows"
$env:GOARCH="386"
go build -o ./build/windows/386/gport.exe

$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o ./build/windows/amd64/gport.exe

$env:GOOS="windows"
$env:GOARCH="arm64"
go build -o ./build/windows/arm64/gport.exe

# darwin

$env:GOOS="darwin"
$env:GOARCH="amd64"
go build -o ./build/darwin/amd64

$env:GOOS="darwin"
$env:GOARCH="arm64"
go build -o ./build/darwin/arm64
