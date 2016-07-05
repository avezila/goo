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
	echo   *echo.Echo
	db     *mgo.Collection
	hasher *Hasher
}

func New() (*Goo, error) {
	session, err := mgo.Dial("goo-mongo")
	if err != nil {
		return nil, err
	}

	g := &Goo{
		echo:   echo.New(),
		db:     session.DB("goo").C("urls"),
		hasher: NewHasher(),
	}

	res := []struct {
		Url  string
		Hash string
	}{}

	if err := g.db.Find(bson.M{}).Select(bson.M{"hash": 1, "url": 1}).All(&res); err != nil {
		return nil, err
	}

	for _, row := range res {
		g.hasher.Set(row.Hash, row.Url)
	}

	g.echo.Use(middleware.Logger())
	g.echo.Use(middleware.Recover())

	g.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

	g.echo.POST("/putUrl", g.Put)
	g.echo.GET("/:hash", g.Get)

	return g, nil
}

func (g *Goo) Run(addr string) {
	g.hasher.Start()
	defer g.hasher.Stop()
	g.echo.Run(fasthttp.New(addr))
}

func (g *Goo) Get(ctx echo.Context) error {
	hash := ctx.Param("hash")
	url := g.hasher.Get(hash)
	if url != "" {
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
	hash, exists := g.hasher.Insert(url)
	if !exists {
		if err := g.db.Insert(bson.M{"hash": hash, "url": url}); err != nil {
			g.hasher.Delete(hash)
			ctx.String(500, "error writing db "+err.Error())
			return err
		}
	}
	ctx.String(200, hash)
	return nil
}
