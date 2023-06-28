package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type response struct {
	body       string
	statusCode int
}

func curlEndpoint(endpoint string, ch chan<- response) {
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("cache-control", "no-cache")

	client := &http.Client{}
	// res, err := http.DefaultClient.Do(req)
	res, err := client.Do(req)
	statusCode := res.StatusCode
	if err != nil {
		fmt.Println("Error: ", err)
		ch <- response{"", statusCode}
	}
	fmt.Println("Status: ", res.Status)
	defer res.Body.Close()
	// convert []byte to string
	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)
	// send the response struct over the channel
	ch <- response{bodyString, statusCode}
}

// readConfig reads a file and returns an array of strings
func readConfig(filename string) ([]string, error) {
	file, err := os.Open(filename)
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return lines, err
}

func main() {

	// create sync group and channel. A WaitGroup waits for a collection of goroutines to finish.
	var wg sync.WaitGroup
	ch := make(chan response)

	lines, _ := readConfig("test")
	//create empty slice of endpoints that we want to call concurrently
	endpoints := []string{}
	println(endpoints)
	// create map of endpoints with line string as key and endpoint array of string as value
	endpoints_map := map[string][]string{}

	// iterate over string array lines
	for _, value := range lines {
		// append each endpoint to the endpoints slice
		// endpoints = append(endpoints, fmt.Sprintf("https://vault.admin.%s/v1/sys/health", value))
		// endpoints = append(endpoints, fmt.Sprintf("https://registry.admin.%s/api/v2.0/health", value))
		// endpoints = append(endpoints, fmt.Sprintf("https://idp.admin.%s/auth/admin/synergy/console/", value))
		// add value to endpoints_map and endpoints slice if not already present
		endpoints_map[value] = []string{fmt.Sprintf("https://vault.admin.%s/v1/sys/health", value), fmt.Sprintf("https://registry.admin.%s/api/v2.0/health", value)}
	}

	// add all endpoints to endpoints slice
	for _, value := range endpoints_map {
		for _, endpoint := range value {
			endpoints = append(endpoints, endpoint)
		}
	}

	// iterate over endpoints slice
	for _, endpoint := range endpoints {
		wg.Add(1)
		// launch a goroutine for each endpoint using an anonymous function
		go func(ep string) {
			defer wg.Done()
			// call curlEndpoint function
			curlEndpoint(endpoint, ch)
		}(endpoint)
	}

	//separe goroutine to wait for all the goroutines to finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	for resp := range ch {
		fmt.Println(resp.statusCode)
		fmt.Println(resp.body)
	}

	fmt.Println(endpoints_map)
	fmt.Println(endpoints)
	// for range endpoints {
	// 	// receive the response struct from the channel
	// 	res := <-ch
	// 	// fmt.Println(res.statusCode)
	// 	fmt.Println(res.body)
	// }

	// http.Handle("/metrics", promhttp.Handler())
	// http.ListenAndServe(":2112", nil)
}
