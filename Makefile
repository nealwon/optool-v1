GO=go
LDFLAGS=-ldflags="-w -s"
OUTPUTDIR=output

all: clean windows darwin linux

windows: directory
	GOOS=windows $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/optool.exe github.com/nealwon/optool
linux: directory
	GOOS=linux $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/optool_linux github.com/nealwon/optool
darwin: directory
	GOOS=darwin $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/optool_darwin github.com/nealwon/optool

directory:
	test -d $(OUTPUTDIR) || mkdir -p $(OUTPUTDIR)

clean:
	rm -f $(OUTPUTDIR)/*