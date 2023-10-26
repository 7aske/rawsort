all: main.go exif.go util.go
	go build -o bin/rawsort .

install: bin/rawsort
	mkdir -p /usr/local/bin
	cp bin/rawsort /usr/local/bin/rawsort