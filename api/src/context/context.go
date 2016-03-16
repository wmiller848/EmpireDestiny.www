package context

import (
	"ed/game"
	"ed/queue"
	"fmt"
	"net/http"
	"regexp"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

var GlobalPtr *Global

// TODO
// Make all this shit a db or cache
type Global struct {
	Session     *rethink.Session
	PlayerQueue map[string]chan *queue.Session
	Game        *game.Game
	Queue       *queue.Queue
	MatchRegex  *regexp.Regexp
}

type Context struct {
	*Global
	playerid string
}

type MatchContext struct {
	*Context
	sessionid string
}

func (c *Context) SetGlobal(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.Global = GlobalPtr
	fmt.Println("SetGlobal - Forwarding")
	rw.Header().Add("Content-Type", "application/json")
	next(rw, req)
}

func (c *Context) PlayerMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	authid, err := req.Cookie("EDAUTH")
	_authid := ""
	if err != nil {
		_authid = ""
	} else {
		_authid = authid.Value
	}
	fmt.Println("PlayerMiddleware", req.RequestURI, _authid)
	authSession, err := c.GetAuthSession(_authid)
	if (err != nil) && req.RequestURI != "/auth" {
		fmt.Println("Redirect to /auth")
		rw.Header().Set("Location", "/auth")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	} else if err == nil {
		c.playerid = authSession.PlayerID
	}
	fmt.Println("PlayerMiddleware - Forwarding")
	next(rw, req)
}

func (c *MatchContext) MatchSessionMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	sessionid, err := req.Cookie("EDSESSION")
	ssid := ""
	if err != nil {
		fmt.Println(err.Error())
		ssid = ""
	} else {
		ssid = sessionid.Value
	}

	fmt.Println("MatchSessionMiddleware", req.RequestURI, ssid)
	_, err = c.GetMatchInstance(ssid)
	if (ssid == "" || err != nil) &&
		c.MatchRegex.Match([]byte(req.RequestURI)) == false {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		c.sessionid = ssid
	}
	fmt.Println("MatchSessionMiddleware - Forwarding")
	next(rw, req)
}

func (c *Context) Root(rw web.ResponseWriter, req *web.Request) {
	// fmt.Fprint(rw, "/")
}
