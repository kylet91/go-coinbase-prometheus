# Go CoinSpot

## CoinSpot Prometheus Exporter
A Prometheus exporter that polls the [CoinSpot RO API](https://www.coinspot.com.au/api#rocoinsbalance) and makes the values available to Prometheus.

It's written to use a Hashicorp Vault approle to get a token / get the API secret key / auth key, so you'll have to rewrite this if you're not using Vault.

You will need to edit main.go to change your Vault URL / hostname.

### Environment Variables

| Variable | Description | Default |
| -------- | ----------- | ------- |
| listenPort | Port the prometheus server will listen on | 2112 |
| scrapeInterval | Interval to scrape the Remootio API for an open/close status | 30 seconds |

### To-Do:
- Add an env var to only search for a specific coin