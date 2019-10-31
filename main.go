package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/marema31/kamino/cmd"
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

	end := make(chan error, 1)

	//This function will stop itself at context cancellation or at end of main goroutine
	go func() {
		end <- cmd.Execute(ctx)
	}()

	//Waiting for signal on the channel and call cancel on the context
	select {
	case <-sigChan: //Received CTRL+C
		cancel() // Cancellation of context, that will propagate to all function that listen ctx.Done
		log.Println("Aborting ....")
		log.Println("Waiting for all sub task abortion...")
		<-end
	case <-ctx.Done(): // the context has been cancelled
		log.Println("Waiting for all sub task abortion...")
		<-end
	case executeError := <-end: //The goroutine executing the action is finished we can stop here
		if executeError != nil {
			os.Exit(1)
		}
	}
	os.Exit(0)
}

/* func main() {

/*
	fmt.Printf("Will run the sync %s\n", strings.Join(os.Args[i:], ", "))

	for _, syncName := range os.Args[i:] {
		err := kaminoSync.Do(ctx, config, syncName, environment, instances)
		if err != nil {
			log.Println(err)
			cancel()                    // Cancellation of context, that will propagate to all function that listen ct.Done
			time.Sleep(5 * time.Second) // TODO:Wait on waitgroup ??
			os.Exit(1)
		}
	}
	time.Sleep(1 * time.Second) // TODO:Wait on waitgroup ??
}
*/
