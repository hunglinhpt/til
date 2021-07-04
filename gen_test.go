package main

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"til/model"
)

func TestGen(t *testing.T) {
	fmt.Println("ssssssssssssssssss")
	err := mgm.SetDefaultConfig(nil, "news", options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
	if err != nil {
		panic(err)
	}
	queue := make(chan int, 8)
	end := make(chan bool)
	for i := 0; i <= 16; i++ {
		go func() {
			for _ = range queue {
				item := &model.News{}
				_ = faker.FakeData(item)
				fmt.Println(item)
				err := mgm.Coll(item).Create(item)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()

	}
	var i int
	for i < 2<<10 {
		queue <- i
	}
	close(queue)
	end <- true

	select {
	case <-end:
		return
	}
}
