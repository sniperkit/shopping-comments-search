package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var waitgroup sync.WaitGroup

func request(i int) {
	fmt.Println(i, time.Now())
	resp, err := http.Get("http://haofly.net/abc")
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	fmt.Println(string(body))
	waitgroup.Done()
	// ch := make(chan string)
	// ch <- string(body)
}

func main() {
	// 48s
	// now := time.Now()
	// for i := 0; i <= 1000; i++ {
	// 	fmt.Println(i, time.Now())
	// 	resp, err := http.Get("http://haofly.net/abc")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	defer resp.Body.Close()

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 	}
	// 	fmt.Println(string(body))
	// }
	// now1 := time.Now()
	// fmt.Println(now1.Sub(now))

	now := time.Now()
	for i := 0; i <= 1000; i++ {
		waitgroup.Add(1)
		go request(i)
	}
	waitgroup.Wait()
	now1 := time.Now()

	// for i := 0; i <= 1; i++ {
	// 	a := <- ch
	// 	fmt.Println(a)
	// }
	// // fmt.Println(ch)
	fmt.Println(now1.Sub(now))
}
