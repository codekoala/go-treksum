cmds := api scraper

build: api

bin:
	mkdir -p ./bin/

$(cmds): bin
	go build -o ./bin/treksum-$(@) ./cmd/treksum-$(@)
	cd ./bin/
	sha256sum treksum-$(@) > SHA256SUMS
