package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

//func TestGen(t *testing.T) {
//	err := mgm.SetDefaultConfig(nil, "news", options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
//	if err != nil {
//		panic(err)
//	}
//	queue := make(chan int, 8)
//	end := make(chan bool)
//	for i := 0; i <= 16; i++ {
//		go func() {
//			for _ = range queue {
//				item := &model.News{}
//				_ = faker.FakeData(item)
//				fmt.Println(item)
//				err := mgm.Coll(item).Create(item)
//				if err != nil {
//					fmt.Println(err)
//				}
//			}
//		}()
//
//	}
//	for i := 0; i < 2<<20; i++ {
//		queue <- i
//	}
//	close(queue)
//	end <- true
//
//	select {
//	case <-end:
//		return
//	}
//}

func TestHTTP2(t *testing.T) {
	url := "https://localhost:8081/list"
	start := time.Now()
	// some computation

	//for i := 0; i < 200; i++ {
	payload := `{"page":200}`
	Post(url, bytes.NewBuffer([]byte(payload)))
	elapsed := time.Since(start)
	fmt.Println("Https", elapsed)
}
