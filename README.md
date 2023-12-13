# Shelly Scraper

Used to scrape your shelly plugs and send the data to influxdb

## Set Up

First you will need to setup your influxdb

### start influxdb:

```bash
docker compose up -d
```

go to http://127.0.0.1:8086 and setup your influx instance (this is currently
manual but probably can be automated)

### set environment variables


You will need to set the following environment variables:

```bash
INFLUXDB_ORG=
INFLUXDB_BUCKET=
INFLUXDB_TOKEN=
# Electricty Maps token 
MAPS_TOKEN=
```

### Set Up Shelly plugs

I used the [Shelly Plug S](https://www.amazon.nl/Shelly-Huisautomatisering-Elektriciteitsmeter-Meerkleurige-LED-indicator/dp/B0BTJ1DTBX/) and set the up following the instructions in the box. You will need to find the IP address of your Shelly Plugs on the app and make sure this is reacheable


### Update the Config File

create a config file with the addresses for your shelly plugs

you can run the following:

```bash
cp config.example.yaml config.yaml
```

and set the address and name of your plugs in the created file


### Run the app

For now you can just run it with:

```bash
go run cmd/scraper/scraper.go
```



