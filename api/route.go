package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"til/model"
)

type CustomWriter struct {
	w io.Writer
}

func (t CustomWriter) Write(p []byte) (n int, err error) {
	// Đưa data vào writer của request
	n, err = t.w.Write(p)
	if f, ok := t.w.(http.Flusher); ok {
		// Chuyền data cho client
		f.Flush()
	}
	return
}

type Body struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
}

type Response map[string]interface{}

// http1 protocol
func List(ctx *gin.Context) {
	b := &Body{
		PerPage: 200,
	}
	err := ctx.Bind(b)
	if err != nil {
		ctx.JSON(400, Response{"message": "Invalid body request"})
		return
	}
	if b.PerPage < 200 {
		b.PerPage = 200
	}
	if b.Page < 0 {
		b.Page = 0
	}
	result := make([]*model.News, 0)
	skip := b.Page * b.PerPage
	err = mgm.Coll(&model.News{}).SimpleFind(&result, bson.M{}, &options.FindOptions{
		Limit: &b.PerPage,
		Skip:  &skip,
	})
	if err != nil {
		ctx.JSON(500, Response{"message": ""})
	}
	ctx.JSON(200, Response{"data": result, "message": "success"})
}

// Support http2 protocol
func PushList(ctx *gin.Context) {
	if !(ctx.Request.ProtoMajor == 2) {
		ctx.JSON(400, Response{"message": fmt.Sprintf("not support http %d", ctx.Request.ProtoMajor)})
		return
	}

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(400, Response{"message": "not support http 2"})
		return
	}
	buf, err = getData(buf.Bytes())
	if err != nil {
		ctx.JSON(500, Response{"message": err})
		return
	}
	_, _ = io.Copy(CustomWriter{w: ctx.Writer}, buf)
	ctx.Writer.Flush()
}

func getData(bits []byte) (*bytes.Buffer, error) {
	b := &Body{
		PerPage: 200,
	}
	err := json.Unmarshal(bits, b)
	if err != nil {
		return nil, err
	}
	if b.PerPage < 200 {
		b.PerPage = 200
	}
	if b.Page < 0 {
		b.Page = 0
	}
	result := make([]*model.News, 0)
	skip := b.Page * b.PerPage
	err = mgm.Coll(&model.News{}).SimpleFind(&result, bson.M{}, &options.FindOptions{
		Limit: &b.PerPage,
		Skip:  &skip,
	})
	if err != nil {
		return nil, err
	}
	tmp, err := json.Marshal(&result)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(tmp), nil
}
