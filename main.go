package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	currency := flag.String("currency", "BTC-USD", "Currency type to retrieve. Ex: BTC-USD, ETH-USD, LTC-USD")
	start := flag.String("start", time.Now().Add(-24*time.Hour).Format(time.RFC3339), "Start of date range")
	interval := flag.Duration("interval", time.Hour, "Duration specified in go format. Ex: 1h")
	out := flag.String("out", "out.log", "File to write output to")
	flag.Parse()

	fmt.Println("currency:", *currency)
	fmt.Println("start:", *start)
	fmt.Println("interval:", interval.Seconds())

	resp, err := http.Get(fmt.Sprintf("https://api.gdax.com/products/%s/candles?start=%s&end=%s&granularity=%v", *currency, *start, time.Now().Format(time.RFC3339), interval.Seconds()))
	if err != nil {
		log.Fatal("http error", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("http error", resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	ticks := [][]float64{}
	json.Unmarshal(body, &ticks)
	fmt.Printf("count: %v", len(ticks))
	file, err := os.OpenFile(*out, os.O_CREATE, os.ModePerm)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range ticks {
		if tick, err := json.Marshal(newTick(*currency, v)); err == nil {
			io.WriteString(file, fmt.Sprintln(string(tick)))
		}
	}
}

type tick struct {
	Currency string
	Time     string
	Low      float64
	High     float64
	Open     float64
	Close    float64
	Volume   float64
}

func newTick(currency string, props []float64) tick {
	return tick{
		Currency: currency,
		Time:     time.Unix(int64(props[0]), 0).Format(time.RFC3339),
		Low:      props[1],
		High:     props[2],
		Open:     props[3],
		Close:    props[4],
		Volume:   props[5],
	}
}
