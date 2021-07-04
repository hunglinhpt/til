package main

import (
	"bytes"
	"crypto/tls"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func Post(url string, buf *bytes.Buffer) {
	transport := &http2.Transport{DialTLS: dialT(buf), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Transport: transport}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	//res.Write(os.Stdout)

	framer := http2.NewFramer(ioutil.Discard, buf)
	for {
		f, err := framer.ReadFrame()
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		switch err.(type) {
		case nil:
			log.Println(f)
		case http2.ConnectionError:
			// Ignore. There will be many errors of type "PROTOCOL_ERROR, DATA
			// frame with stream ID 0". Presumably we are abusing the framer.
		default:
			log.Println(err, framer.ErrorDetail())
		}
	}
}

// dialT returns a connection that writes everything that is read to w.
func dialT(w io.Writer) func(network, addr string, cfg *tls.Config) (net.Conn, error) {
	return func(network, addr string, cfg *tls.Config) (net.Conn, error) {
		conn, err := tls.Dial(network, addr, cfg)
		return &TlsCon{conn, w}, err
	}
}

type TlsCon struct {
	net.Conn
	T io.Writer // receives everything that is read from Conn
}

func (w *TlsCon) Read(b []byte) (n int, err error) {
	n, err = w.Conn.Read(b)
	w.T.Write(b)
	return
}
