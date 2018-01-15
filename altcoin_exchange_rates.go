package main

import (
	"strconv"
	"os"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

// MQTT server config data structure
type MqttConfig struct {
	Server	string
	Port	int
	Topic	string
	ClientID	string
}

// Program configuration data structure
type Config struct {
	Mqtt	MqttConfig `json:"mqtt"`
	Message	string `json:"message"`
	Url	string `json:"url"`
	Request_timeout	int `json:"request_timeout"`
	Coins	[]string `json:"coins"`
	Logfile	string `json:"logfile"`
}

// Coin data structure
type CoinData struct {
	Id	string `json:"id"`
	Name	string `json:"name"`
	Symbol	string `json:"symbol"`
	Rank	string `json:"rank"`
	Price_usd	string `json:"price_usd"`
	Price_btc	string `json:"price_btc"`
	Aday_volume_usd	string `json:"24h_volume_usd"`
	Market_cap_usd	string `json:"market_cap_usd"`
	Total_supply	string `json:"total_supply"`
	Max_supply	string `json:"max_supply"`
	Percent_change_1h	string `json:"percent_change_1h"`
	Percent_change_24h	string `json:"percent_change_24h"`
	Percent_change_7d	string `json:"percent_change_7d"`
	Last_updated	string `json:"last_updated"`
}

// Main function
func main() {
	// Load config from config.json
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := new(Config)
	conferr := decoder.Decode(&config)
	if conferr != nil {
		log.Fatal(conferr)
	}

	// Open or create logfile
	logfile, err := os.OpenFile(config.Logfile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	// Get exchange rates
	var httpClient = &http.Client{
		Timeout: time.Second * time.Duration(config.Request_timeout),
	}
	exchange_rates := make(map[string]CoinData)
	for _, coin := range config.Coins {
		var coinData []CoinData
		response, err := httpClient.Get(config.Url + coin)
		if err == nil {
			buf, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err == nil {
				err = json.Unmarshal(buf, &coinData)
				if err == nil {
					exchange_rates[coin] = coinData[0]
				} else {
					log.Printf("Error when parsing json: %s\n", err.Error())
				}
			} else {
				log.Println("Error when read data: " + err.Error())
			}
		} else {
			log.Println("Error when getting data: " + err.Error())
		}
	}

	// Connect to the MQTT Server.
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + config.Mqtt.Server + ":" + strconv.Itoa(config.Mqtt.Port)).SetClientID(config.Mqtt.ClientID)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}

	// Publish data
	mqttMessage, err := json.Marshal(exchange_rates)
	if err == nil {
		if token := client.Publish(config.Mqtt.Topic, 1, true, mqttMessage); token.Wait() && token.Error() != nil {
			log.Printf("Error when publish MQTT message: %s\n", token.Error())
		}
	} else {
		log.Println("Error when creating JSON: " + err.Error())
	}

	// Disconnect the Network Connection.
	client.Disconnect(250)

}
