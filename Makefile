cmds := api scraper

all: api scraper checksums

bin:
	mkdir -p ./bin/

$(cmds): bin generate
	go build -o ./bin/treksum-$(@) ./cmd/treksum-$(@)

generate:
	go generate

compress:
	upx ./bin/treksum-*
	$(MAKE) checksums

checksums:
	cd ./bin/; sha256sum treksum-* > SHA256SUMS

clean:
	rm -rf ./bin

docker:
	docker build -t codekoala/treksum .
