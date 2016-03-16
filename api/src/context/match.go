package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"ed/card"
	"ed/match"
	"ed/player"
	"ed/queue"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

type EnemyRsp struct {
	Field    []card.Card
	Destrict []card.Card
	God      card.Card
	Fortress card.Card
}

type PlayerRsp struct {
	Field    []card.Card
	Destrict []card.Card
	Hand     []card.Card
	God      card.Card
	Fortress card.Card
	Gold     int32
}

func (c *Context) SetMatchInstance(mtch *match.Match) error {
	resp, err := rethink.DB("ed").Table("matches").Insert(mtch).Run(c.Session)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func (c *Context) GetMatchInstance(id string) (*match.Match, error) {
	cursor, err := rethink.DB("ed").Table("matches").Get(id).Run(c.Session)
	if err != nil {
		return nil, err
	}
	matchInstance := match.Match{}
	err = cursor.One(&matchInstance)
	if err != nil {
		return nil, err
	}
	return &matchInstance, nil
}

func (c *MatchContext) GetMatch(rw web.ResponseWriter, req *web.Request) {

	_, err := c.GetMatchInstance(c.sessionid)
	if err == nil {
		rw.Header().Set("Location", "/match/info")
		rw.WriteHeader(http.StatusFound)
		return
	}
	fmt.Println(err.Error())
	// deckid := req.PathParams["deckid"]
	//
	//
	plr, err := player.LoadPlayerDB(c.playerid, c.Session)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	deckid := req.PathParams["deckid"]
	fmt.Println(deckid)
	dck := plr.Decks[deckid]
	if dck == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	plrInstance, err := plr.LoadPlayer(c.Session)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.PlayerQueue[c.playerid] = nil
	c.Queue.AddToQueue(plrInstance)
	pchan := make(chan *queue.Session)
	c.PlayerQueue[c.playerid] = pchan
	// BLOCK on matching in the queue
	session, open := <-pchan
	close(pchan)
	// match has started
	if session != nil && open == true {
		session.Mutex.Lock()
		_, err = c.GetMatchInstance(c.sessionid)
		if err != nil {
			c.SetMatchInstance(match.CreateMatch(session.PlayerA.Id(), session.PlayerB.Id(), session.Sessionid))
		}
		session.Mutex.Unlock()
		runtime.Gosched()
		cookie := &http.Cookie{
			Name:  "EDSESSION",
			Value: session.Sessionid,
			Path:  "/",
		}
		http.SetCookie(rw, cookie)
	}
	rw.WriteHeader(http.StatusCreated)
}

func (c *MatchContext) GetInfo(rw web.ResponseWriter, req *web.Request) {
	mtch, err := c.GetMatchInstance(c.sessionid)
	fmt.Println(mtch, err)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err = mtch.LoadPlayers(c.Session)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	var plr *player.Player
	var enemy *player.Player
	if mtch.PlayerA.Id() == c.playerid {
		plr = mtch.PlayerA
		enemy = mtch.PlayerB
	} else {
		plr = mtch.PlayerB
		enemy = mtch.PlayerA
	}

	fmt.Println(plr, enemy)

	rsp := make(map[string]interface{})
	rsp["Session"] = c.sessionid
	rsp["Enemy"] = EnemyRsp{
		Field:    enemy.Field(),
		Destrict: enemy.Destrict(),
		God:      enemy.God(),
		Fortress: enemy.Fortress(),
	}
	rsp["Player"] = PlayerRsp{
		Field:    plr.Hand(),
		Destrict: plr.Destrict(),
		God:      enemy.God(),
		Fortress: enemy.Fortress(),
		// Hand:     plr.Hand(),
		// Gold:     plr.Gold(),
	}
	fmt.Println(rsp)
	d, err := json.Marshal(rsp)
	if err != nil {
		d = []byte{}
	}
	fmt.Fprint(rw, string(d))
}

func (c *MatchContext) GetCommit(rw web.ResponseWriter, req *web.Request) {
	// m := c.Matches[c.sessionid]
	// fmt.Println("GetCommit", req.PathParams)
	// fmt.Println(m, req)
}

func (c *MatchContext) GetReset(rw web.ResponseWriter, req *web.Request) {
	// m := c.Matches[c.sessionid]
	// fmt.Println("GetReset", req.PathParams)
	// fmt.Println(m, req)
}

func (c *MatchContext) GetMoveAttempt(rw web.ResponseWriter, req *web.Request) {
	// m := c.Matches[c.sessionid]
	// fmt.Println(m, req)
}

func (c *MatchContext) GetMatchFinish(rw web.ResponseWriter, req *web.Request) {
	// buf := make([]byte, 1024)
	// _, err := req.Body.Read(buf)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(buf)
	// m := c.Matches[c.sessionid]
	// fmt.Println(m)
}
