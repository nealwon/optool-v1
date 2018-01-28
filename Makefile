GO=go
LDFLAGS=-ldflags="-w -s"
OUTPUTDIR=output
SHELL=/bin/bash --posix

all: clean windows darwin linux done

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

pack: clean
	GOOS=windows $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/optool.exe github.com/nealwon/optool
	/usr/bin/zip output/optool-windows.zip $(OUTPUTDIR)/optool.exe
	GOOS=darwin $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/optool github.com/nealwon/optool
	chmod +x $(OUTPUTDIR)/optool
	/usr/bin/zip output/optool-darwin.zip $(OUTPUTDIR)/optool
	GOOS=linux $(GO) build $(LDFLAGS) -o $(OUTPUTDIR)/optool github.com/nealwon/optool
	chmod +x $(OUTPUTDIR)/optool
	/usr/bin/zip output/optool-linux.zip $(OUTPUTDIR)/optool
	rm -f $(OUTPUTDIR)/optool.exe $(OUTPUTDIR)/optool
