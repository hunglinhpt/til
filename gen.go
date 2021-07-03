package main

import (
	"fmt"
	"til/model"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("ssssssssssssssssss")
	err := mgm.SetDefaultConfig(nil, "news", options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
	if err != nil {
		panic(err)
	}
	queue := make(chan int, 8)
	end := make(chan bool)
	for i := 0; i <= 8; i++ {
		for _ = range queue {

			go func() {
				item := &model.News{}
				fmt.Println(item)
				err := mgm.Coll(item).Create(item)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
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
