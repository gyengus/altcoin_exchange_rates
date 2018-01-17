// altoin_exchange_rates package
package altcoin_exchange_rates

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

// MqttConfig type is server config data structure
type MqttConfig struct {
	// Server url or ip address
	Server	string
	// Port is the MQTT servers port
	Port	int
	// Topic
	Topic	string
	// ClientID
	ClientID	string
}

// Config type is a struct to store program configuration data
type Config struct {
	Mqtt	MqttConfig `json:"mqtt"`
	Message	string `json:"message"`
	URL	string `json:"url"`
	RequestTimeout	int `json:"request_timeout"`
	Coins	[]string `json:"coins"`
	Logfile	string `json:"logfile"`
}

// CoinData structure is store the coins datas
type CoinData struct {
	ID	string `json:"id"`
	Name	string `json:"name"`
	Symbol	string `json:"symbol"`
	Rank	string `json:"rank"`
	PriceUSD	string `json:"price_usd"`
	PriceBTC	string `json:"price_btc"`
	AdayVolumeUSD	string `json:"24h_volume_usd"`
	MarketCapUSD	string `json:"market_cap_usd"`
	TotalSupply	string `json:"total_supply"`
	MaxSupply	string `json:"max_supply"`
	PercentChange1h	string `json:"percent_change_1h"`
	PercentChange24h	string `json:"percent_change_24h"`
	PercentChange7d	string `json:"percent_change_7d"`
	LastUpdated	string `json:"last_updated"`
}

// main function
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
		Timeout: time.Second * time.Duration(config.RequestTimeout),
	}
	exchangeRates := make(map[string]CoinData)
	for _, coin := range config.Coins {
		var coinData []CoinData
		response, err := httpClient.Get(config.URL + coin)
		if err == nil {
			buf, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err == nil {
				err = json.Unmarshal(buf, &coinData)
				if err == nil {
					exchangeRates[coin] = coinData[0]
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
	mqttMessage, err := json.Marshal(exchangeRates)
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
