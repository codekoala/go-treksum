cmds := api scraper

build: api checksums

all: api scraper checksums

bin:
	mkdir -p ./bin/

$(cmds): bin
	go build -o ./bin/treksum-$(@) ./cmd/treksum-$(@)

checksums:
	cd ./bin/; sha256sum treksum-* > SHA256SUMS

clean:
	rm -rf ./bin
