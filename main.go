package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tothszabi/argonone/internal/fan"
	"github.com/tothszabi/argonone/internal/log"
	"github.com/tothszabi/argonone/internal/temperature"
)

const sleepDuration = 2 * time.Second

func main() {
	os.Exit(run())
}

func run() int {
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, os.Kill, syscall.SIGABRT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	var readyChan = make(chan struct{}, 1)
	var errChan = make(chan error, 1)

	go controlFan(ctx, readyChan, errChan)

	statusCode := 0

	select {
	case sig := <-killSig:
		log.Info("Received signal:", sig)
		cancel()

	case err := <-errChan:
		log.Error("Failure:", err)
		cancel()
		statusCode = 1
	}

	select {
	case <-readyChan:
		return statusCode
	}
}

func controlFan(ctx context.Context, readyChan chan<- struct{}, errChan chan<- error) {
	log.Info("Starting fan control")

	var previousSpeed = -1

	for {
		temp, err := temperature.CurrentTemperature()
		if err != nil {
			errChan <- err
		}

		speed := calculateFanSpeed(temp)

		if previousSpeed != speed {
			err = setFanSpeed(speed)
			if err != nil {
				errChan <- err
			}

			previousSpeed = speed
		}

		select {
		case <-time.After(sleepDuration):
			// repeats loop
		case <-ctx.Done():
			stopFan()

			log.Info("Fan control stopping")

			readyChan <- struct{}{}
			return
		}
	}
}

func calculateFanSpeed(temp int) int {
	switch {
	case temp < 55:
		return 0
	case temp < 60:
		return 10
	case temp < 65:
		return 50
	default:
		return 100
	}
}

func setFanSpeed(speed int) error {
	if speed == 0 {
		log.Info("Stopping fan")
	} else {
		log.Info("Setting fan speed to", fmt.Sprintf("%d%%", speed))
	}

	return fan.SetSpeed(speed)
}

func stopFan() {
	setFanSpeed(0)
}
