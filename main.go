package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marema31/kamino/config"
	kaminoSync "github.com/marema31/kamino/sync"
)

// Theses variables will be provided by goreleaser via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {

	ctx := context.Background()

	i := 1
	configPath := "."
	configFile := ".kamino"

	if len(os.Args) > 2 && os.Args[1] == "-c" {
		i = 3
		configPath = filepath.Dir(os.Args[2])
		configFile = filepath.Base(os.Args[2])
	}

	if (len(os.Args)-i) < 1 || os.Args[i] == "-h" {
		fmt.Printf("kamino %v, commit %v, built at %v\n\n", version, commit, date)
		fmt.Println("usage: kamino [-c configFile] syncName [syncName ... syncName]")
		fmt.Println()
		fmt.Println("The config file name must be provided without the file extension, kamino will try json, toml and yaml")
		os.Exit(0)
	}

	config, err := config.New(configPath, configFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Will run the sync %s\n", strings.Join(os.Args[i:], ", "))

	for _, syncName := range os.Args[i:] {
		err := kaminoSync.Do(ctx, config, syncName)
		if err != nil {
			log.Fatal(err)
		}
	}
}
