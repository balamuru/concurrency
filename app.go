package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	// sequential()
	// println()
	// goroutines_naive()
	goroutines_channels()

}

func goroutines_channels() {
	urlSizeDetailsChannel := make(chan string)
	go responseSize2("https://www.example.com", urlSizeDetailsChannel)
	go responseSize2("https://www.golang.org", urlSizeDetailsChannel)
	go responseSize2("https://www.golang.org/doc", urlSizeDetailsChannel)
	println(<- urlSizeDetailsChannel)
	println(<- urlSizeDetailsChannel)
	println(<- urlSizeDetailsChannel)

}

func responseSize2(url string, urlSizeDetailsChannel chan string) {
	urlSizeDetailsChannel <- responseSizeHelper(url)
}

func goroutines_naive() {
	go responseSize("https://www.example.com")
	go responseSize("https://www.golang.org")
	go responseSize("https://www.golang.org/doc")

	time.Sleep(5 * time.Second) //not the best way to wait for all the goroutines to complete
	println("done")

}
func sequential() {
	responseSize("https://www.example.com")
	responseSize("https://www.golang.org")
	responseSize("https://www.golang.org/doc")

}

func responseSize(url string) {
	println(responseSizeHelper(url))
}



func responseSizeHelper(url string) string {
	var sb strings.Builder
	sb.WriteString(url)
	sb.WriteString(" : ")
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	sb.WriteString(strconv.FormatInt(int64(len(body)), 10))
	return sb.String()
}
