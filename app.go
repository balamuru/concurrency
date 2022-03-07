package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// https://notes.shichao.io/gopl/ch8/#channels
func main() {

	//note: see app_test.go for demo drivers
	// goroutines_naive()
	// goroutines_channels()
	// goroutines_spinner()
	// goroutines_timeserver()
	// goroutines_waitGroups()
	// goroutines_syncDemoUsingAtomic()
	// goroutines_syncDemoUsingMutex()
	// goroutines_channelsDemo()
	// goroutines_unbufferedChannels()
	goroutines_bufferedChannels()

}

func goroutines_bufferedChannels() {
	ch := make(chan int, 3)
	defer close(ch)
	var waitGroup sync.WaitGroup

	produceSingleItemToBufferedChannel := func(i int, ch chan int) {
		ch <- i
	}

	consumeAllValuesSlowlyFromBufferedChannel := func( ch chan int, maxRead int) {		
		waitGroup.Add(1)
		read := 0	
		for {
			//take a nap
			fmt.Println("Zzzzzzzzzz")
			time.Sleep(time.Millisecond*500)

			//read all items currently in channel
			for i := 0; i < len(ch); i++ {
				if read >= maxRead {
					fmt.Printf("Done reading %v values\n", read)
					waitGroup.Done()
					return
				}
				read++
				fmt.Printf("Read: %v\n", <- ch)	
			}			
		}		
	}

	go consumeAllValuesSlowlyFromBufferedChannel(ch, 7)

	for i := 0; i < 10; i++ {
		go produceSingleItemToBufferedChannel(i, ch)
	}

	//shut down program 
	waitGroup.Wait()



}

func goroutines_unbufferedChannels() {
	ch := make(chan int)
	defer close(ch)
	var waitGroup sync.WaitGroup

	produceSingleItemToUnbufferedChannel := func(i int, ch chan int) {
		ch <- i
	}

	consumeAllValuesSlowlyFromUnbufferedChannel := func( ch chan int, maxRead int) {		
		waitGroup.Add(1)
		read := 0	
		for {
			if read >= maxRead {
				fmt.Printf("Done reading %v values\n", read)
				waitGroup.Done()
				return
			}
			read++
			//take a nap
			fmt.Println("Zzzzzzzzzz")
			time.Sleep(time.Millisecond*500)
			i := <- ch	
			fmt.Printf("Read: %v\n", i)	
		}

	}

	go consumeAllValuesSlowlyFromUnbufferedChannel(ch, 7)

	for i := 0; i < 10; i++ {
		go produceSingleItemToUnbufferedChannel(i, ch)
	}

	//shut down program 
	waitGroup.Wait()

}

func goroutines_channelsDemo() {
	ch := make(chan int)


	mult := func(x, y int) {
		res := x*y
		ch <- res
	}


	//pass value of channel as parameter
	div := func(x, y int, channel chan int) {
		res := x/y
		channel <- res
	}

	go mult(10,2)
	val, ok := <- ch
	fmt.Printf("type of channel = %T\n", ch)	
	fmt.Printf("type of channel = %v\n", ch)	
	fmt.Printf("result of value = %v\n", val)	
	fmt.Printf("result of ok = %v\n", ok)


	ch2 := make(chan int)
	go div(10,2, ch2)
	val2, ok2 := <- ch2
	fmt.Printf("type of channel2 = %T\n", ch2)	
	fmt.Printf("type of channel2 = %v\n", ch2)	
	fmt.Printf("result of value2 = %v\n", val2)	
	fmt.Printf("result of ok2 = %v\n", ok2)
}

func goroutines_syncDemoUsingMutex() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	f := func(waitGroup *sync.WaitGroup, mutex *sync.Mutex, num *int64) {
		defer waitGroup.Done() //decrement

		for i := 0; i < 100; i++ {

			func() {
				defer mutex.Unlock() //ensure unlock

				mutex.Lock() //lock

				*num++
			}()
		}
	}

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	var num int64 = 0
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)               //increment
		go f(&waitGroup, &mutex, &num) //pass in ref
	}

	waitGroup.Wait() //wait till 0

	println(num)

	println("done")
}

func goroutines_syncDemoUsingAtomic() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	f := func(waitGroup *sync.WaitGroup, num *int64) {
		defer waitGroup.Done() //decrement

		for i := 0; i < 100; i++ {
			atomic.AddInt64(num, 1)
			// *num++
		}
	}

	var waitGroup sync.WaitGroup

	var num int64 = 0
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)       //increment
		go f(&waitGroup, &num) //pass in ref
	}

	waitGroup.Wait() //wait till 0

	println(num)

	println("done")
}

func goroutines_waitGroups() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	f := func(waitGroup *sync.WaitGroup, i int) {
		defer waitGroup.Done() //decrement
		println(i)
	}

	var waitGroup sync.WaitGroup

	for i := 0; i < 5; i++ {
		waitGroup.Add(1) //increment
		// go f(&waitGroup, i) - bad - don't pass in copy - will deadlock
		go f(&waitGroup, i) //pass in ref
	}

	waitGroup.Wait() //wait till 0

	println("done")
}

//client can connect via nc eg nc localhost 8000
func goroutines_timeserver() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}
		//handle multiple clients
		// go handleConn(conn)
		go handleEchoConn(conn)
	}
}

func handleEchoConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		//a nested goroutine to make each echo async
		go echo(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t\t\t", strings.ToUpper(shout))
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		t := time.Now().Format("15:04:05\n")
		println("serving time: " + t)
		_, err := io.WriteString(c, t)
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
func goroutines_spinner() {
	go spinner(100 * time.Millisecond)
	time.Sleep(5 * time.Second)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func goroutines_channels() {
	urlSizeDetailsChannel := make(chan string)
	defer close(urlSizeDetailsChannel)
	go responseSize2("https://www.example.com", urlSizeDetailsChannel)
	go responseSize2("https://www.golang.org", urlSizeDetailsChannel)
	go responseSize2("https://www.golang.org/doc", urlSizeDetailsChannel)
	println(<-urlSizeDetailsChannel) //channel sends to println
	println(<-urlSizeDetailsChannel) //channel sends to println
	println(<-urlSizeDetailsChannel) //channel sends to println

}

func responseSize2(url string, urlSizeDetailsChannel chan string) {
	urlSizeDetailsChannel <- responseSizeHelper(url) //channel receives
}

func goroutines_naive() {
	go responseSize("https://www.example.com")
	go responseSize("https://www.golang.org")
	go responseSize("https://www.golang.org/doc")

	time.Sleep(5 * time.Second) //not the best way to wait for all the goroutines to complete
	println("done")

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
