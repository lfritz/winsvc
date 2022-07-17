.PHONY: all clean

all: service.exe

clean:
	rm -f service.exe

service.exe: cmd/service/main.go
	GOOS=windows GOARCH=amd64 go build -o service.exe cmd/service/main.go
