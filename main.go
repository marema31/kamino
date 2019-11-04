package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/marema31/kamino/cmd"
)

func main() {
	log := cmd.GetLogger()

	// Synchro beetween the main goroutin and the Execution goroutine
	end := make(chan error, 1)

	// Create the application context for correct sub-task abortion if CTRL+C
	ctx := context.Background()
	// trap Ctrl+C
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// At end of function we remove signal trapping
	defer signal.Stop(sigChan)

	//This function will stop itself at context cancellation or at end of main goroutine
	go func() {
		end <- cmd.Execute(ctx)
	}()

	//Waiting for signal on the channel and call cancel on the context
	select {
	case <-sigChan: //Received CTRL+C
		cancel() // Cancellation of context, that will propagate to all function that listen ctx.Done
		log.Warn("Aborting ....")
		log.Info("Waiting for all sub task abortion...")
		<-end // We wait for the main goroutine to end
		log.Info("exiting on CTRL+C")
		os.Exit(1)
	case <-ctx.Done(): // the context has been cancelled (It should not happen since only the previous case fire a cancellation)
		log.Info("Waiting for all sub task abortion...")
		<-end // We wait for the main goroutine to end
		log.Info("exiting on cancel")
		os.Exit(1)
	case executeError := <-end: //The goroutine executing the action is finished we can stop here
		if executeError != nil {
			log.Info("exiting on error")
			os.Exit(1)
		}
	}
	log.Info("bye")
	log.Debug("Debug")
	os.Exit(0)
}
