package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	ConfigPath string `short:"c" long:"config" description:"config path" required:"true"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	parser.Usage = "[OPTIONS] address"

	args, err := parser.Parse()

	if err != nil {
		return
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		return
	}

	database, err := NewDatabase(opts.ConfigPath)
	if err != nil {
		log.Fatalln(err)
	}
	dnsHandler := NewDNSHandler(args[0], database)

	if err := dnsHandler.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
