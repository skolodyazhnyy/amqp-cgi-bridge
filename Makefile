VERSION?=unknown
COMMIT?=unknown
BUILDARGS?=
BUILDLDFLAGS?=
BUILDOUT?=amqp-cgi-bridge

build:
	go build $(BUILDARGS) -ldflags '$(BUILDLDFLAGS) -X main.version=$(VERSION) -X main.commit=$(COMMIT)' -o ${BUILDOUT}

crossbuild:
	GOOS=linux BUILDOUT=amqp-cgi-bridge-linux BUILDLDFLAGS='-extldflags "-static"' make build
	GOOS=darwin BUILDOUT=amqp-cgi-bridge-darwin make build
	GOOS=windows BUILDOUT=amqp-cgi-bridge-windows.exe make build

clean:
	rm amqp-cgi-bridge*
