cmds := api scraper

all: api scraper checksums

bin:
	mkdir -p ./bin/

$(cmds): bin
	go build -o ./bin/treksum-$(@) ./cmd/treksum-$(@)

checksums:
	cd ./bin/; sha256sum treksum-* > SHA256SUMS

clean:
	rm -rf ./bin

docker:
	docker build -t codekoala/treksum .
