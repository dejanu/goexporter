package main

import (
	"log"
	"net/http"
)


func curlENDPOINT(url string) {
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	log.Println(res.StatusCode)
	log.Println(string(res.Body))
}



func main() {
	url := "ifconfig.me"
	curlENDPOINT(url)
}

