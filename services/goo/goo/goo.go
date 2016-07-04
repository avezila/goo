package goo

import (
	"github.com/PuerkitoBio/purell"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Goo struct {
	urls  map[string]string
	hashs map[string]string
	echo  *echo.Echo
	len   int
	db    *mgo.Collection
}

func New() (*Goo, error) {
	session, err := mgo.Dial("goo-mongo")
	if err != nil {
		return nil, err
	}

	g := &Goo{
		urls:  make(map[string]string),
		hashs: make(map[string]string),
		echo:  echo.New(),
		len:   1,
		db:    session.DB("goo").C("urls"),
	}
	res := []struct {
		Url  string
		Hash string
	}{}

	if err := g.db.Find(bson.M{}).Select(bson.M{"hash": 1, "url": 1}).All(&res); err != nil {
		return nil, err
	}

	for _, row := range res {
		g.urls[row.Hash] = row.Url
		g.hashs[row.Url] = row.Hash
	}

	g.echo.Use(middleware.Logger())
	g.echo.Use(middleware.Recover())

	g.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

	g.echo.POST("/putUrl", g.Put)
	g.echo.GET("/:hash", g.Get)
	return g, nil
}

func (g *Goo) Run(addr string) {
	g.echo.Run(fasthttp.New(addr))
}

func (g *Goo) Get(ctx echo.Context) error {
	hash := ctx.Param("hash")
	url, ok := g.urls[hash]
	if ok {
		ctx.Redirect(301, url)
		return nil
	}
	ctx.String(404, "404: Not found.")
	return nil
}

func (g *Goo) Put(ctx echo.Context) error {
	req := ctx.Request()
	len := req.ContentLength()
	if len > 1024*1024 {
		ctx.String(403, "Too long url")
		return nil
	}
	buf := make([]byte, len)
	if _, err := req.Body().Read(buf); err != nil {
		ctx.String(500, err.Error())
		return nil
	}
	url := string(buf)
	if pureURL, err := purell.NormalizeURLString(url, purell.FlagsSafe); err == nil {
		url = pureURL
	}
	hash, ok := g.hashs[url]
	if ok {
		ctx.String(200, hash)
		return nil
	}
	for {
		hash = RandString(g.len)
		if _, ok = g.urls[hash]; !ok {
			break
		}
		g.len++
	}
	if err := g.db.Insert(bson.M{"hash": hash, "url": url}); err != nil {
		ctx.String(500, "error writing row "+err.Error())
		return nil
	}
	g.urls[hash] = url
	g.hashs[url] = hash
	ctx.String(200, hash)
	return nil
}
