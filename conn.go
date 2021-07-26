package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type SubRes struct {
	Message    string `json:"message"`
	UUID       string `json:"uuid"`
	Result     string `json:"result"`
	API        string `json:"api"`
	Visibility string `json:"visibility"`
	Options    struct {
		Useragent string `json:"useragent"`
	} `json:"options"`
	URL string `json:"url"`
}

type Payload struct {
	URL        string `json:"url"`
	Visibility string `json:"visibility"`
}

func conn(ss string) {
	data := Payload{
		URL:        os.Args[1],
		Visibility: "private",
	}
	payloadBytes, err := json.Marshal(data)
	checkError(err)

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://urlscan.io/api/v1/scan/", body)
	checkError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", "$APIKEY") // Provide own key

	resp, err := http.DefaultClient.Do(req)
	checkError(err)

	defer resp.Body.Close()

	time.Sleep(15 * time.Second) // Necessary, to let urlscan finish scanning before getting results

	respBody, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	urlSubResUnmarshal(respBody)
}

func urlSubResUnmarshal(sb []byte) {
	subRes := SubRes{}
	json.Unmarshal(sb, &subRes)
	subRes.printSubRes()
}

func (s SubRes) printSubRes() {
	// Will use API res in future to gather specific results to pivot off of
	fmt.Println(s.Message)
	fmt.Println(s.Result)
	fmt.Println(s.API)
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
