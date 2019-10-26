package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/marema31/kamino/cmd"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/recipe"
	//	"github.com/marema31/kamino/config"
	//	kaminoSync "github.com/marema31/kamino/sync"
)

func main() {
	// Create the application context for correct sub-task abortion if CTRL+C
	ctx := context.Background()
	// trap Ctrl+C
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// At end of function we remove signal trapping
	defer signal.Stop(sigChan)

	//Waiting for signal on the channel and call cancel on the context
	//This function will stop itself at context cancellation or at end of main goroutine
	go func() {
		select {
		case <-sigChan:
			cancel() // Cancellation of context, that will propagate to all function that listen cts.Done
			log.Println("Aborting ....")
		case <-ctx.Done(): // the context has been cancelled
		}
		log.Println("Waiting for all sub task abortion...")
		time.Sleep(5 * time.Second) //TODO: Wait on waitgroup ??
		os.Exit(1)
	}()
	cmd.Execute(ctx)
	//TODO: not the right action, just for testing
	_ = recipe.Load(context.Background(), "testdata", "pokemon", datasource.New(), &provider.KaminoProvider{})
}

/* func main() {

/*
	//TODO: During CLI review add options for environment and instances
	environment := ""
	var instances []string

	//	instances = append(instances, "sch1")
	//	instances = append(instances, "1")

	if len(os.Args) > 2 && os.Args[1] == "-d" {
		i = 3
		configPath = os.Args[2]
		//	configFile = filepath.Base(os.Args[2])
	}

	config, err := config.New(configPath, configFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Will run the sync %s\n", strings.Join(os.Args[i:], ", "))

	for _, syncName := range os.Args[i:] {
		err := kaminoSync.Do(ctx, config, syncName, environment, instances)
		if err != nil {
			log.Println(err)
			cancel()                    // Cancellation of context, that will propagate to all function that listen cts.Done
			time.Sleep(5 * time.Second) // TODO:Wait on waitgroup ??
			os.Exit(1)
		}
	}
	time.Sleep(1 * time.Second) // TODO:Wait on waitgroup ??
}
*/
