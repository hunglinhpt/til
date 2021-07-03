package main

import (
	"fmt"
	"til/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func middle(engine *gin.Engine) {
	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s %s %s %s %d %s %s \n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
	engine.Use(gin.Recovery())
}

func firstVersion() *gin.Engine {
	//
	engine := gin.New()

	middle(engine)
	engine.GET("/ping", api.Static)
	return engine
}

func secondVersion() *gin.Engine {
	engine := gin.New()
	middle(engine)
	engine.GET("/ping", api.Static)
	return engine
}

func setup() error {
	// Setup the mgm default config
	return mgm.SetDefaultConfig(nil, "news", options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
}

// func main() {
// 	if err := setup(); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	v1 := firstVersion()
// 	go func() {
// 		v1.Run(":8000")
// 	}()
//
// 	v2 := secondVersion()
// 	v2.RunTLS(":8080", "./cert.pem", "./key.pem")
// }
