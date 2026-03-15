BINARY=mdev

build:
	go build -o $(BINARY)

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -f $(BINARY)

run:
	go run main.go