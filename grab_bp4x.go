package scrape

import (
	"fmt"
	"log"

	"github.com/javscrape/go-scrape/query"
)

const bp4xJavURL = "https://www.bp4x.com/?q=%s"
const bp4xAmateurURL = "https://www.bp4x.com/?c=amateur&q=%s"
const bp4xIVURL = "https://www.bp4x.com/?c=iv&q=%s"
const bp4xHentaiURL = "https://www.bp4x.com/?c=hentai&q=%s"

// GrabBP4XType ...
type GrabBP4XType int

// JAV ...
const (
	JAV GrabBP4XType = iota
	AMATEUR
	IV
	HENTAI
)

var bp4xGrabList = []string{
	JAV:     bp4xJavURL,
	AMATEUR: bp4xAmateurURL,
	IV:      bp4xIVURL,
	HENTAI:  bp4xHentaiURL,
}

type grabBP4X struct {
	grabType GrabBP4XType
}

// Find ...
func (g *grabBP4X) Find(name string) error {
	url := bp4xGrabList[g.grabType]
	url = fmt.Sprintf(url, name)
	document, e := query.New(url)
	if e != nil {
		return e
	}
	ret, e := document.Html()
	log.Println(ret)
	return nil
}

// Decode ...
func (g *grabBP4X) Decode(*Message) error {
	panic("implement me")
}

// NewGrabBP4X ...
func NewGrabBP4X(grabType GrabBP4XType) IGrab {
	return &grabBP4X{
		grabType: grabType,
	}
}
