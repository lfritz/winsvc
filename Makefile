.PHONY: all clean

all: service.exe installer.exe

clean:
	rm -f service.exe installer.exe

service.exe: cmd/service/main.go
	GOOS=windows GOARCH=amd64 go build -o service.exe cmd/service/main.go

installer.exe: cmd/installer/main.go
	GOOS=windows GOARCH=amd64 go build -o installer.exe cmd/installer/main.go
