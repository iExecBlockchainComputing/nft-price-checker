package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var sum float64
var temp float64

type Collection struct {
	Stats struct {
		Floor_price float64 `json:"floor_price"`
	}
}

type Price struct {
	Ethereum struct {
		Usd float64 `json:"usd"`
	}
}

// Reading input.txt file
func readInput(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}
	return entries, scanner.Err()
}

// Asking for Opensea's colections Floor Prices
func floorPrice(collection_name string) string {
	var col1 Collection

	url := fmt.Sprintf("https://api.opensea.io/api/v1/collection/%s/stats", collection_name)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &col1)

	temp = col1.Stats.Floor_price //Temporary storing the value of the floor price

	if col1.Stats.Floor_price == 0 {
		return "  x " + collection_name + " cannot be found on Opensea"
	} else {
		return "--> " + collection_name + " Floor price = " + fmt.Sprintf("%f", col1.Stats.Floor_price) + " eth"
	}
}

// Asking for Eth price
func ethPrice() float64 {
	var pr1 Price

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &pr1)

	return pr1.Ethereum.Usd
}

func main() {

	iexec_out := os.Getenv("IEXEC_OUT")
	iexec_in := os.Getenv("IEXEC_IN")
	iexec_input_file := os.Getenv("IEXEC_INPUT_FILE_NAME_1")

	temp = 0

	entries, readErr := readInput(iexec_in + "/" + iexec_input_file)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// Append some results in /iexec_out/
	fr, err := os.OpenFile(iexec_out+"/result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	for _, col := range entries {
		sum += temp //Incrementing the total value
		if _, err := fr.Write([]byte(floorPrice(col) + "\n")); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := fr.Write([]byte("------------- \n The estimate total value of your portfolio is : " + fmt.Sprintf("%f", sum) + " eth\n Or " + fmt.Sprintf("%f", sum*ethPrice()) + " Usd")); err != nil {
		log.Fatal(err)
	}
	if err := fr.Close(); err != nil {
		log.Fatal(err)
	}

	// Declare everything is computed
	fc, err := os.OpenFile(iexec_out+"/computed.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := fc.Write([]byte("{ \"deterministic-output-path\" : \"" + iexec_out + "/result.txt\" }")); err != nil {
		log.Fatal(err)
	}
	if err := fc.Close(); err != nil {
		log.Fatal(err)
	}
}