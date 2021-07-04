package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
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

func PushList(ctx *gin.Context) {
	if !ctx.Request.ProtoAtLeast(2, 1) {
		ctx.JSON(400, Response{"message": "not support http 2"})
	}
	flush, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.JSON(400, Response{"message": "not support http 2"})
	}
	flush.Flush()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, ctx.Request.Body)
	if err != nil {
		ctx.JSON(400, Response{"message": "not support http 2"})
	}
	buf, err = getData(buf.Bytes())
	if err != nil {
		ctx.JSON(500, Response{"message": err})
	}
	_, _ = io.Copy(CustomWriter{w: ctx.Writer}, buf)
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
