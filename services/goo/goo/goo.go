package goo

import (
	"github.com/PuerkitoBio/purell"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

type Goo struct {
	urls  map[string]string
	hashs map[string]string
	echo  *echo.Echo
	len   int
}

func New() (*Goo, error) {
	g := &Goo{
		urls:  make(map[string]string),
		hashs: make(map[string]string),
		echo:  echo.New(),
		len:   1,
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
	ctx.NoContent(404)
	return nil
}

func (g *Goo) Put(ctx echo.Context) error {
	req := ctx.Request()
	res := ctx.Response()
	len := req.ContentLength()
	if len > 1024*1024 {
		res.WriteHeader(403)
		res.Write([]byte("Too long url"))
		return nil
	}
	buf := make([]byte, len)
	if _, err := req.Body().Read(buf); err != nil {
		res.WriteHeader(500)
		return nil
	}
	url := string(buf)
	if pureURL, err := purell.NormalizeURLString(url, purell.FlagsSafe); err == nil {
		url = pureURL
	}
	hash, ok := g.hashs[url]
	if ok {
		res.Write([]byte(hash))
		return nil
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
	res.Write([]byte(hash))
	return nil
}
