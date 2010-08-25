.PHONY: all clean

include ${GOROOT}/src/Make.inc

all: godep gomake gorules goinfo

gomake: src/gomake.${O}
	${LD} -o $@ src/gomake.${O}

godep: src/godep.${O}
	${LD} -o $@ src/godep.${O}

gorules: src/gorules.${O}
	${LD} -o $@ src/gorules.${O}

goinfo: src/goinfo.${O}
	${LD} -o $@ src/goinfo.${O}

src/godep.${O}: src/godep.go src/common.go
	${GC} -o $@ src/godep.go src/common.go

src/gomake.${O}: src/gomake.go src/common.go
	${GC} -o $@ src/gomake.go src/common.go

src/gorules.${O}: src/gorules.go src/common.go
	${GC} -o $@ src/gorules.go src/common.go

src/goinfo.${O}: src/goinfo.go src/common.go
	${GC} -o $@ src/goinfo.go src/common.go

install: all
	cp godep gomake gorules goinfo /usr/local/bin
	cp doc/*.1 /usr/local/share/man/man1

format:
	gofmt -w src/*.go

clean:
	rm -f godep gomake getgo goinfo src/*.${O}
