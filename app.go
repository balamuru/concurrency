package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)
func main() {

	// sequential()
	println()
	concurrent()

}

func sequential() {
	responseSize("https://www.example.com")
	responseSize("https://www.golang.org")
	responseSize("https://www.golang.org/doc")



}


func concurrent() {
	go responseSize("https://www.example.com")
	go responseSize("https://www.golang.org")
	go responseSize("https://www.golang.org/doc")

	time.Sleep(3 * time.Second)
	println("done")

}

func responseSize(url string) {
	response, err := http.Get(url)
	if (err != nil) {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if (err != nil) {
		log.Fatal(err)
	}
	fmt.Printf("%s: page size = %d\n", url, len(body))
}