.PHONY: all clean

include ${GOROOT}/src/Make.${GOARCH}

all: godep gomake

gomake:

godep: src/godep.${O}
	${LD} -o $@ src/godep.${O}

MAINFILES = src/godep.go

src/godep.${O}: src/godep.go
	${GC} -o $@ $?

install: godep gomake
	cp godep ${GOBIN}

clean:
	rm -f godep src/*.${O}
