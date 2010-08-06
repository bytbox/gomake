.PHONY: all clean

include ${GOROOT}/src/Make.${GOARCH}

all: godep gomake getgo

getgo: src/getgo.${O}
	${LD} -o $@ src/getgo.${O}

gomake: src/gomake.${O}
	${LD} -o $@ src/gomake.${O}

godep: src/godep.${O}
	${LD} -o $@ src/godep.${O}

src/godep.${O}: src/godep.go src/common.go
	${GC} -o $@ src/godep.go src/common.go

src/getgo.${O}: src/getgo.go src/common.go
	${GC} -o $@ src/getgo.go src/common.go

src/gomake.${O}: src/gomake.go src/common.go
	${GC} -o $@ src/gomake.go src/common.go

install: all
	cp godep gomake ${GOBIN}

clean:
	rm -f godep gomake getgo src/*.${O}
