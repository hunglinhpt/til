package main

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"net/http"
	"strings"
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

func BenchmarkHttpsListApi(t *testing.B) {
	url := "https://localhost:8081/list"
	start := time.Now()
	// some computation

	for i := 0; i < 200; i++ {
		payload := strings.NewReader(fmt.Sprintf(`{"page":%d}`, i))
		req, err := http.NewRequest("POST", url, payload)
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		tr := &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		res, err := client.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()
		//assert.NoError(t, err)
		//bits, _ := ioutil.ReadAll(res.Body)
		//fmt.Println(string(bits))
	}
	elapsed := time.Since(start)
	fmt.Println("Https", elapsed)
}
