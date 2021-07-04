package main

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func BenchmarkHttpListApi(t *testing.B) {
	url := "http://localhost:8080/list"
	start := time.Now()
	// some computation

	for i := 0; i < 200; i++ {
		payload := strings.NewReader(fmt.Sprintf(`{"page":%d}`, i))
		req, err := http.NewRequest("POST", url, payload)
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		tr := &http.Transport{
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
	fmt.Println("Http", elapsed)
}
