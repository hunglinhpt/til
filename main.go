package main

import (
	"fmt"
	"til/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func middle(engine *gin.Engine, version string) {
	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s %s %s %s %s %d %s %s \n", version,
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

	middle(engine, "http1 ")
	engine.POST("/list", api.List)
	return engine
}

func secondVersion() *gin.Engine {
	engine := gin.New()
	middle(engine, "http2 ")
	engine.POST("/list", api.PushList)
	return engine
}

func setup() error {
	// Setup the mgm default config
	return mgm.SetDefaultConfig(nil, "news", options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
}

func main() {
	chErr := make(chan error)
	if err := setup(); err != nil {
		chErr <- err
	}

	go func() {
		v1 := firstVersion()
		chErr <- v1.Run(":8080")
	}()

	go func() {
		v2 := secondVersion()
		chErr <- v2.RunTLS(":8081", "./cert.pem", "./key.pem")
	}()
	select {
	case err := <-chErr:
		panic(err)
	}
}
