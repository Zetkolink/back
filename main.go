package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"
)

var (
	b   *back
	cfg *config
)

func init() {
	confPath := os.Getenv("STOCK_CONFPATH")

	if confPath == "" {
		confPath = "./etc/config.yml"
	}

	yamlFile, err := ioutil.ReadFile(confPath)

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)

	if err != nil {
		log.Fatal(err)
	}

	err = initBack()

	if err != nil {
		log.Fatal(err)
	}
}

func main()  {
	b.Run()

	log.Println("Success start")

	listenSignals()
}

func listenSignals() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	for sig := range signals {
		log.Println("Got signal: " + sig.String())

		_ = destroyBack()

		return
	}
}

func initBack() error {
	if b != nil {
		return nil
	}

	var err error
	b, err = newBack()

	if err != nil {
		return err
	}

	return nil
}

func destroyBack() error {
	if b != nil {
		b.Stop()
		b = nil
	}

	return nil
}
