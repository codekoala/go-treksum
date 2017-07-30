package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

var Global conf

type conf struct {
	DbHost     string `default:"localhost"`
	DbPort     uint16 `default:"5432"`
	DbName     string `default:"treksum"`
	DbUser     string
	DbPassword string
}

func init() {
	if err := envconfig.Process("treksum", &Global); err != nil {
		fmt.Printf("unable to parse config: %s", err)
		os.Exit(1)
	}
}
