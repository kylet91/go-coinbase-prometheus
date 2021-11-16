package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	addr        string = "192.168.1.1:8080"
	role_id     string = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	secret_id   string = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	coinspot_id string = "abc123abc123abc123abc123abc123abc123"
	vault_addr  string = "https://vault.mydomain.com:8200"
)

var (
	client         *api.Client
	scrapeInterval string
	listenPort     string
	secretKey      string = getSecret("secret")
	coinBalance    *prometheus.GaugeVec
	coinAud        *prometheus.GaugeVec
	coinRate       *prometheus.GaugeVec
)

func main() {
	// Get env vars
	scrapeInterval = os.Getenv("scrapeInterval")
	if len(scrapeInterval) == 0 {
		scrapeInterval = "10"
	}
	scrapeIntervali, _ := strconv.Atoi(scrapeInterval)
	listenPort = ":" + os.Getenv("listenPort")
	if len(listenPort) == 1 {
		listenPort = ":2113"
	}

	// HTTP
	timeout := time.Duration(2 * time.Second)
	coinClient := &http.Client{Timeout: timeout}

	log.SetFlags(0)

	// Run initial check and create all Prom gauges
	status := getStatus(coinClient)

	if status {
		fmt.Println("Working!")
	} else {
		fmt.Println("No bueno")
	}

	handler(coinClient)

	// scrapeInterval * seconds
	ticker := time.NewTicker(time.Duration(scrapeIntervali) * time.Second)

	go func() {
		for _ = range ticker.C {
			handler(coinClient)
		}
	}()

	// Start Prometheus server
	fmt.Println("Starting prom server on port ", listenPort)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(listenPort, nil)
}

func handler(coinClient *http.Client) {

	request := buildBalancesRequest(coinClient)

	// Make connection to CoinSpot
	resp, err := coinClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error:", err)
	}

	var responseBodyJson map[string]interface{}
	err = json.Unmarshal(responseBody, &responseBodyJson)
	if err != nil {
		log.Fatal("Error:", err)
	}

	out := responseBodyJson["balances"].([]interface{})

	for _, a := range out {
		for coin, b := range a.(map[string]interface{}) {
			for f, c := range b.(map[string]interface{}) {
				cString := fmt.Sprintf("%v", c)
				value, err := strconv.ParseFloat(cString, 64)
				if err != nil {
					log.Fatal("Error:", err)
				}

				switch f {
				case "audbalance":
					coinAud.With(prometheus.Labels{"coin": coin}).Set(value)
				case "rate":
					coinRate.With(prometheus.Labels{"coin": coin}).Set(value)
				case "balance":
					coinBalance.With(prometheus.Labels{"coin": coin}).Set(value)
				}
			}

		}
	}
}

func buildSha(requestBody []byte) string {
	hmacData := hmac.New(sha512.New, []byte(secretKey))
	hmacData.Write([]byte(requestBody))
	sha := hex.EncodeToString(hmacData.Sum(nil))

	return sha
}

func buildBalancesRequest(coinClient *http.Client) *http.Request {
	nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	requestBody, err := json.Marshal(map[string]string{
		"nonce": nonce,
	})
	if err != nil {
		log.Fatal("Error:", err)
	}

	sha := buildSha(requestBody)

	request, err := http.NewRequest("POST", "https://www.coinspot.com.au/api/ro/my/balances", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("key", coinspot_id)
	request.Header.Set("sign", sha)
	if err != nil {
		log.Fatal("Error:", err)
	}

	return request

}

func buildCoinRequest(coinClient *http.Client) *http.Request {
	nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	requestBody, err := json.Marshal(map[string]string{
		"nonce":    nonce,
		"cointype": "DOGE",
	})
	if err != nil {
		log.Fatal("Error:", err)
	}

	sha := buildSha(requestBody)

	request, err := http.NewRequest("POST", "https://www.coinspot.com.au/api/ro/my/balances?cointype=DOGE", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("key", coinspot_id)
	request.Header.Set("sign", sha)
	if err != nil {
		log.Fatal("Error:", err)
	}

	return request

}

func getStatus(coinClient *http.Client) bool {
	nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	requestBody, err := json.Marshal(map[string]string{
		"nonce": nonce,
	})
	if err != nil {
		log.Fatal("Error:", err)
	}

	sha := buildSha(requestBody)

	request, err := http.NewRequest("POST", "https://www.coinspot.com.au/api/ro/my/balances", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("key", coinspot_id)
	request.Header.Set("sign", sha)
	if err != nil {
		log.Fatal("Error:", err)
	}

	resp, err := coinClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error:", err)
	}

	var responseBodyJson map[string]interface{}
	err = json.Unmarshal(responseBody, &responseBodyJson)
	if err != nil {
		log.Fatal("Error:", err)
	}

	coins := responseBodyJson["balances"].([]interface{})

	for _, a := range coins {
		for coin := range a.(map[string]interface{}) {
			fmt.Println("Found:" + coin)

			coinBalance = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: "coins",
					Name:      "balance",
					Help:      "Number of coins.",
				},
				[]string{
					"coin",
				})

			coinAud = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: "coins",
					Name:      "aud",
					Help:      "Value of coins in AUD.",
				},
				[]string{
					"coin",
				})

			coinRate = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: "coins",
					Name:      "rate",
					Help:      "Value of each coin.",
				},
				[]string{
					"coin",
				})
		}
	}

	// Could probably go somewhere better than getStatus, but yolo.
	prometheus.MustRegister(coinBalance)
	prometheus.MustRegister(coinAud)
	prometheus.MustRegister(coinRate)

	out := fmt.Sprintf("%v", responseBodyJson["status"])

	if out == "ok" {
		return true
	} else {
		return false
	}

}

func getSecret(secretName string) string {
	conf := api.DefaultConfig()
	client, _ = api.NewClient(conf)
	client.SetAddress(vault_addr)
	// Get auth from Vault
	resp, err := client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   role_id,
		"secret_id": secret_id,
	})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if resp == nil {
		fmt.Println("empty response from credential provider")
		return ""
	}

	// Set client token for Vault
	client.SetToken(resp.Auth.ClientToken)
	secret, err := client.Logical().Read("coinspot/data/" + coinspot_id)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	key, _ := secret.Data["data"].(map[string]interface{})
	keyreturn := fmt.Sprintf("%v", key[secretName])
	return keyreturn
}
