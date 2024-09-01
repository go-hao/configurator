package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/go-hao/configurator"
)

func main() {
	filePath := "configs/default.yml"
	cfg := &Config{}
	err := configurator.Setup(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	err = configurator.Dump(cfg, filePath)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info(fmt.Sprintf("dump default config: %s", filePath))
}
