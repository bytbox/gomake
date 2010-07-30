.PHONY: all clean

include ${GOROOT}/src/Make.${GOARCH}

all: godep

godep: main.${O}
	${LD} -o $@ main.${O}

MAINFILES = main.go

main.${O}: ${MAINFILES}
	${GC} -o $@ ${MAINFILES}

install: godep
	cp godep ${GOBIN}

clean:
	rm -f godep *.${O}
