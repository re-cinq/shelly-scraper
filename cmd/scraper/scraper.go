package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/klog"

	"github.com/re-cinq/shelly-scraper/pkg/shelly"
)

func main() {
	// termination Handeling
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	// Create a new context, with its cancellation function
	// from the original context
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		Tick(ctx)
	}()

	<-termChan

	cancel()
	<-ctx.Done()
}

func Tick(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	c := shelly.New("192.168.2.51")
	klog.Info("starting ticker")
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			klog.Info("stopping ticker")
			return
		case <-ticker.C:
			klog.Info("fetching energy consumption")
			r, err := c.GetSwitchStatus("0")
			if err != nil {
				klog.Error(err)
				continue
			}

			klog.Info(r.AEnergy.Total)
		}
	}
}
