package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"./broker"
)

func main() {
	message := `Welcome to 
  _____       __  __  ____ _______ _______ ____            _             
 / ____|     |  \/  |/ __ \__   __|__   __|  _ \          | |            
| |  __  ___ | \  / | |  | | | |     | |  | |_) |_ __ ___ | | _____ _ __ 
| | |_ |/ _ \| |\/| | |  | | | |     | |  |  _ <| '__/ _ \| |/ / _ \ '__|
| |__| | (_) | |  | | |__| | | |     | |  | |_) | | | (_) |   <  __/ |   
 \_____|\___/|_|  |_|\___\_\ |_|     |_|  |____/|_|  \___/|_|\_\___|_| 
 !!!! WIP !!!!
`

	fmt.Printf("%s", message)
	runtime.GOMAXPROCS(runtime.NumCPU())
	config, err := broker.ConfigureConfig(os.Args[1:])
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatal("New Broker error: ", err)
	}
	b.Start()

	s := waitForSignal()
	log.Println("signal received, broker closed.", s)
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
