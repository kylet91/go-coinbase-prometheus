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

### Example Prometheus Output
```
# HELP coins_aud Value of coins in AUD.
# TYPE coins_aud gauge
coins_aud{coin="AGIX"} 9001.00
coins_aud{coin="CND"} 9001.00
coins_aud{coin="DOGE"} 9001.00
coins_aud{coin="DOT"} 9001.00
coins_aud{coin="FET"} 9001.00
coins_aud{coin="HIVE"} 9001.00
coins_aud{coin="LOOM"} 9001.00
coins_aud{coin="SCRT"} 9001.00
coins_aud{coin="SHIB"} 9001.00
coins_aud{coin="SNT"} 9001.00
coins_aud{coin="SOL"} 9001.00
# HELP coins_balance Number of coins.
# TYPE coins_balance gauge
coins_balance{coin="AGIX"} 69
coins_balance{coin="CND"} 69
coins_balance{coin="DOGE"} 69
coins_balance{coin="DOT"} 69
coins_balance{coin="FET"} 69
coins_balance{coin="HIVE"} 69
coins_balance{coin="LOOM"} 69
coins_balance{coin="SCRT"} 69
coins_balance{coin="SHIB"} 69
coins_balance{coin="SNT"} 69
coins_balance{coin="SOL"} 69
# HELP coins_rate Value of each coin.
# TYPE coins_rate gauge
coins_rate{coin="AGIX"} 0.420
coins_rate{coin="CND"} 0.420
coins_rate{coin="DOGE"} 0.420
coins_rate{coin="DOT"} 0.420
coins_rate{coin="FET"} 0.420
coins_rate{coin="HIVE"} 0.420
coins_rate{coin="LOOM"} 0.420
coins_rate{coin="SCRT"} 0.420
coins_rate{coin="SHIB"} 0.420
coins_rate{coin="SNT"} 0.420
coins_rate{coin="SOL"} 0.420

```

### To-Do:
- Add an env var to only search for a specific coin
