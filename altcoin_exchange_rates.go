// AltoinExchangeRates package
package AltcoinExchangeRates

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

var config = new(Config)
var httpClient = &http.Client{}


// main function
func main() {
	config = loadConfig()

	httpClient = &http.Client{
		Timeout: time.Second * time.Duration(config.RequestTimeout),
	}

	// Open or create logfile
	logfile, err := os.OpenFile(config.Logfile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	sendMQTTMessage(getExchangeRates())

}

// sendMQTTMessage func is Connecting to the MQTT Server and sending coins datas message to the topic
func sendMQTTMessage(exchangeRates map[string]CoinData) {
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

// getExchangeRates func is getting and build exchange rates
func getExchangeRates() map[string]CoinData {
	exchangeRates := make(map[string]CoinData)
	for _, coin := range config.Coins {
		var tmp = CoinData{}
		tmp, err := httpGet(config.URL + coin)
		if err == nil {
			exchangeRates[coin] = tmp
		}
	}
	return exchangeRates
}

// httpGet func is getting coin data from url
func httpGet(url string) (CoinData, error) {
	response, err := httpClient.Get(url)
	if err == nil {
		buf, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err == nil {
			var coinData []CoinData
			err = json.Unmarshal(buf, &coinData)
			if err == nil {
				return coinData[0], nil
			} else {
				log.Printf("Error when parsing json: %s\n", err.Error())
			}
		} else {
			log.Println("Error when read data: " + err.Error())
		}
	} else {
		log.Println("Error when getting data: " + err.Error())
	}
	return CoinData{}, err
}

// loadConfig func loads config from config.json
func loadConfig() *Config {
	var conf = new(Config)
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	conferr := decoder.Decode(&conf)
	if conferr != nil {
		log.Fatal(conferr)
	}
	file.Close()
	return conf
}
