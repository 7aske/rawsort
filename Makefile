PREFIX = /usr/local
CMD = rawsort

default: cmd/rawsort/main.go
	go build -o bin/$(CMD) cmd/rawsort/main.go

install: default
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp bin/$(CMD) $(PREFIX)/bin/$(CMD)

clean:
	rm -f bin/rawsort