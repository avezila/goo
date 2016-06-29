package goo

import (
	"fmt"

	"github.com/PuerkitoBio/purell"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

type Goo struct {
	urls  map[string]string
	hashs map[string]string
	iris  *iris.Framework
	len   int
}

func New() (*Goo, error) {
	g := &Goo{
		urls:  make(map[string]string),
		hashs: make(map[string]string),
		iris:  iris.New(),
		len:   1,
	}

	g.iris.Use(cors.Default())
	g.iris.Post("/putUrl", g.Put)
	g.iris.Get("/:hash", g.Get)
	return g, nil
}

func (g *Goo) Start(addr string) {
	g.iris.Listen(addr)
}

func (g *Goo) Stop() {
	g.iris.Close()
}

func (g *Goo) Get(ctx *iris.Context) {
	hash := ctx.Param("hash")
	fmt.Println(hash)
	url, ok := g.urls[hash]
	if ok {
		ctx.Redirect(url, 301)
		return
	}
	ctx.EmitError(404)
}

func (g *Goo) Put(ctx *iris.Context) {
	url := string(ctx.Request.Body())
	if pureURL, err := purell.NormalizeURLString(url, purell.FlagsSafe); err == nil {
		url = pureURL
	}
	fmt.Println(url)
	hash, ok := g.hashs[url]
	if ok {
		ctx.WriteString(hash)
		return
	}
	for {
		hash = RandString(g.len)
		if _, ok = g.urls[hash]; !ok {
			break
		}
		g.len++
	}
	g.urls[hash] = url
	g.hashs[url] = hash
	ctx.WriteString(hash)
}
