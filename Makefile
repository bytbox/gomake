.PHONY: all clean

all: godep

godep: main.6
	6l -o $@ main.6

main.6: main.go
	6g -o $@ main.go

clean:
	rm godep *.6
