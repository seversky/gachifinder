all: windows macosx linux

windows:
	cd cmd/gachifinder && GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o windows/gachifinder.exe
mac:
	cd cmd/gachifinder && GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -o osx/gachifinder
linux:
	cd cmd/gachifinder && GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o linux/gachifinder
