package main

import (
	"log"

	"github.com/vrischmann/envconfig"

	"github.com/LarsFox/pow-tcp/src/hashcash"
	"github.com/LarsFox/pow-tcp/src/quotes"
	"github.com/LarsFox/pow-tcp/src/server"
)

type config struct {
	Addr       string `required:"true"`
	TargetBits int64  `envconfig:"default=24"`
}

func main() {
	cfg := &config{}
	if err := envconfig.InitWithPrefix(cfg, "pow_tcp"); err != nil {
		log.Fatal(err)
	}

	hc := hashcash.New(cfg.TargetBits)
	server := server.NewPOWServer(quotes.BookDostoevsky, hc)

	log.Println("Listening on", cfg.Addr)
	server.Listen(cfg.Addr)
}
