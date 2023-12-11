package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
	"k8s.io/klog"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/re-cinq/shelly-scraper/pkg/shelly"
)

func main() {

	options := loadConfig()
	// termination Handeling
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()

	// Create a new context, with its cancellation function
	// from the original context
	ctx, cancel := context.WithCancel(ctx)

	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(
		"http://localhost:8086",
		os.Getenv("INFLUXDB_TOKEN"),
	)

	api := client.WriteAPIBlocking(
		os.Getenv("INFLUXDB_ORG"),
		os.Getenv("INFLUXDB_BUCKET"),
	)

	for i := range options {
		go func(opt *Option) {
			Tick(ctx, api, *opt)
		}(&options[i])
	}

	<-termChan
	cancel()
}

type Option struct {
	Addr string
	Name string
}

func Tick(ctx context.Context, db api.WriteAPIBlocking, o Option) {
	ticker := time.NewTicker(3 * time.Minute)
	c := shelly.New(o.Addr)
	klog.Infof("starting to track on plug: %s", o.Addr)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			klog.Info("stopping ticker")
			return
		case <-ticker.C:
			klog.Infof("fetching energy consumption for: %s", o.Addr)
			r, err := c.GetSwitchStatus("0")
			if err != nil {
				klog.Error(err)
				continue
			}
			now := time.Now()
			for i := range r.AEnergy.ByMinute {
				p := influxdb2.NewPoint("energy_consumption",
					map[string]string{
						"unit":   "milliwatt-hours",
						"device": o.Addr,
						"label":  o.Name,
					},
					map[string]interface{}{
						"consumption": r.AEnergy.ByMinute[i],
					},
					now.Add(time.Duration(-i)*time.Minute))
				err = db.WritePoint(context.Background(), p)
				if err != nil {
					klog.Error(err)
					continue
				}
			}
		}
	}
}

type config struct {
	Plugs []Option
}

func loadConfig() []Option {

	var c config

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		klog.Fatalf("failed loading config: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return c.Plugs
}
