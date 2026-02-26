package cmd

import (
	"flag"
	"fmt"
	"os"

	"go.deepl.dev/mealie-webhook-handler/pkg/api"
	"go.deepl.dev/mealie-webhook-handler/pkg/appcontext"
	"go.deepl.dev/mealie-webhook-handler/pkg/configuration"
)

func Execute(noticeContent []byte) {
	args := flag.NewFlagSet("mealie-webhook-handler", flag.ContinueOnError)
	var configPath string
	var isLicense bool
	args.StringVar(&configPath, "config-file", "webhooks.toml", "Path to the config file")
	args.BoolVar(&isLicense, "license", false, "Show license information")
	err := args.Parse(os.Args[1:])
	if err != nil {
		return
	}

	if isLicense {
		fmt.Printf("%s", noticeContent)
		os.Exit(0)
	}

	configVal, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Failed to open config file: %s\n", err.Error())
		os.Exit(1)
	}

	config, err := configuration.ParseConfiguration(configVal)
	if err != nil {
		fmt.Printf("Failed to parse config file: %s\n", err.Error())
		os.Exit(1)
	}
	if err := config.Init(); err != nil {
		fmt.Printf("Failed to validate config: %s\n", err.Error())
		os.Exit(1)
	}

	ctx := appcontext.AppContext{
		Config: config,
	}

	srv := api.NewServer(ctx)
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
