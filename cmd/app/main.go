package main

import (
	"log"
	"main/internal"
	"main/internal/model"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	application := internal.NewApplication()

	waitChan := make(chan struct{})

	sysSignal := make(chan os.Signal, 1)
	signal.Notify(sysSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer close(waitChan)

		sig := <-sysSignal
		log.Println(model.LevelInfo + "main: received sysSignal: " + sig.String())
		signal.Stop(sysSignal)

		err := application.Shutdown()
		if err != nil {
			log.Println(err.Error())
			return
		}

		log.Println(model.LevelInfo + ": application.Shutdown: successfully")
	}()

	err := application.Start()
	if err != nil {
		log.Println(err.Error())
	}
	sysSignal <- syscall.SIGTERM

	<-waitChan
}
