package main

import (
	"ed/context"
	"ed/game"
	"ed/queue"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

var global context.Global

func main() {

	rethinkSession, err := rethink.Connect(rethink.ConnectOpts{
		Address: "localhost:28015",
		// Addresses: []string{"localhost:28015", "localhost:28016"},
		// Database: "test",
		// AuthKey:       "14daak1cad13dj",
		DiscoverHosts: true,
		MaxIdle:       10,
		MaxOpen:       10,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	gameInstance, err := game.LoadGameFromYML("game/data/game.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	chn := gameInstance.Sync()
	for {
		crdRevison, closed := <-chn
		if closed == false {
			break
		}
		fmt.Println("Card Revision", crdRevison)
		resp, err := rethink.DB("ed").Table("cards").Insert(crdRevison, rethink.InsertOpts{
			Conflict: "replace",
		}).RunWrite(rethinkSession)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(resp)
		}
	}

	exp, err := regexp.Compile(`^/match/queue/[0-9a-zA-z]*`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	global := context.Global{
		Session:     rethinkSession,
		PlayerQueue: make(map[string]chan *queue.Session),
		Game:        gameInstance,
		Queue:       queue.CreateQueue(),
		MatchRegex:  exp,
	}
	context.GlobalPtr = &global

	fmt.Println(global)

	go func() {
		for {
			gq := global.Queue.Consume(&global.PlayerQueue)
			if gq == nil {
				gq = queue.CreateQueue()
			}
			global.Queue = gq
			time.Sleep(500 * time.Millisecond)
		}
	}()

	rootRouter := web.New(context.Context{}). // Create your router
							Middleware(web.LoggerMiddleware). // included logging middleware
							Middleware(web.ShowErrorsMiddleware).
							Middleware((*context.Context).SetGlobal).        // Set the global indexes
							Middleware((*context.Context).PlayerMiddleware). // User Session Validation
							Get("/", (*context.Context).Root).               //
							Get("/auth", (*context.Context).GetAuth).        //
							Head("/auth", (*context.Context).GetAuth).       //
							Get("/player", (*context.Context).GetPlayer).    //
							Post("/decks", (*context.Context).PostDecks).    //
							Put("/decks/:id", (*context.Context).PutDecks).  //
							Get("/game", (*context.Context).GetGame)         //

	matchRouter := rootRouter.Subrouter(context.MatchContext{}, "/match").
		Middleware((*context.MatchContext).MatchSessionMiddleware). // Match Session Validation
		Get("/queue/:deckid", (*context.MatchContext).GetMatch).    //
		Get("/finish", (*context.MatchContext).GetMatchFinish).     //
		Get("/info", (*context.MatchContext).GetInfo).              //
		// Get("/move/engage/:id/taget/:targetid", (*context.MatchContext).GetEngage). //
		// Get("/move/disengage/:id", (*context.MatchContext).GetDisengage).           //
		// Get("/move/bow/:id/target/:targetid", (*context.MatchContext).GetBow).      //
		// Get("/move/unbow/:id", (*context.MatchContext).GetUnbow).                   //
		Get("/move/commit", (*context.MatchContext).GetCommit).
		Get("/move/reset", (*context.MatchContext).GetReset)

	fmt.Println(rootRouter, matchRouter)
	http.ListenAndServe("localhost:3000", rootRouter)
}
