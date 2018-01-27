GO=go
LDFLAGS=-ldflags="-w -s"
OUTPUTDIR=output

all: clean windows darwin linux

windows: directory
	GOOS=windows $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/deployer.exe github.com/nealwon/deployer
linux: directory
	GOOS=linux $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/deployer_linux github.com/nealwon/deployer
darwin: directory
	GOOS=darwin $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/deployer_darwin github.com/nealwon/deployer

directory:
	test -d $(OUTPUTDIR) || mkdir -p $(OUTPUTDIR)

clean:
	rm -f $(OUTPUTDIR)/*