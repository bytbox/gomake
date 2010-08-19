.PHONY: all clean

include ${GOROOT}/src/Make.${GOARCH}

all: godep gomake gorules

gomake: src/gomake.${O}
	${LD} -o $@ src/gomake.${O}

godep: src/godep.${O}
	${LD} -o $@ src/godep.${O}

gorules: src/gorules.${O}
	${LD} -o $@ src/gorules.${O}

src/godep.${O}: src/godep.go src/common.go
	${GC} -o $@ src/godep.go src/common.go

src/gomake.${O}: src/gomake.go src/common.go
	${GC} -o $@ src/gomake.go src/common.go

src/gorules.${O}: src/gorules.go src/common.go
	${GC} -o $@ src/gorules.go src/common.go

install: all
	cp godep gomake gorules ${GOBIN}
	cp doc/*.1 /usr/local/share/man/man1

format:
	gofmt -w src/*.go

clean:
	rm -f godep gomake getgo src/*.${O}
