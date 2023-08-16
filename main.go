package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/irbgeo/pipline-test/accumulator"
	"github.com/irbgeo/pipline-test/controller"
	"github.com/irbgeo/pipline-test/processor"
	"github.com/irbgeo/pipline-test/publisher"
	"github.com/irbgeo/pipline-test/source"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment variables file")
	}

	sourceTimeoutENV := os.Getenv("SOURCE_TIMEOUT")
	if sourceTimeoutENV == "" {
		log.Fatal("SOURCE_TIMEOUT not found")
	}
	sourceTimeout, err := strconv.Atoi(sourceTimeoutENV)
	if err != nil {
		log.Fatalf("failed value of SOURCE_TIMEOUT (%s)", sourceTimeoutENV)
	}

	workerCountENV := os.Getenv("WORKER_COUNT")
	if workerCountENV == "" {
		log.Fatal("WORKER_COUNT not found")
	}
	workerCount, err := strconv.Atoi(workerCountENV)
	if err != nil {
		log.Fatalf("failed value of WORKER_COUNT (%s)", workerCountENV)
	}

	publisherTimeoutENV := os.Getenv("PUBLISHER_TIMEOUT")
	if publisherTimeoutENV == "" {
		log.Fatal("PUBLISHER_TIMEOUT not found")
	}
	publisherTimeout, err := strconv.Atoi(publisherTimeoutENV)
	if err != nil {
		log.Fatalf("failed value of PUBLISHER_TIMEOUT (%s)", publisherTimeoutENV)
	}

	src := source.New(time.Duration(sourceTimeout) * time.Millisecond)
	prc := processor.New(workerCount)
	accum := accumulator.New()
	pub := publisher.New()

	ctrl := controller.New(
		src,
		prc,
		accum,
		pub,
		time.Duration(publisherTimeout)*time.Second,
	)

	ctrl.Start()

	log.Println("i'm turned on")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	ctrl.Stop()

	log.Println("goodbye")
}
