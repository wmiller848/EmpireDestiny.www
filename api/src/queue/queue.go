package queue

import (
	"crypto/rand"
	"ed/player"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	PlayerA   *player.Player
	PlayerB   *player.Player
	Mutex     *sync.Mutex
	Sessionid string
}

type Queue struct {
	next    *Queue
	players [2]*player.Player
}

func CreateQueue() *Queue {
	return &Queue{
		players: [2]*player.Player{},
	}
}

func (q *Queue) AddToQueue(player *player.Player) {
	if q.next != nil {
		q.next.AddToQueue(player)
	} else {

		if q.players[0] != nil {
			if q.players[0].Id() == player.Id() {
				return
			}
		} else if q.players[1] != nil {
			if q.players[1].Id() == player.Id() {
				return
			}
		}

		if q.players[0] == nil {
			q.players[0] = player
		} else if q.players[1] == nil {
			q.players[1] = player
		} else {
			q.next = CreateQueue()
			q.next.AddToQueue(player)
		}
	}
}

func (q *Queue) Consume(playerQueue *map[string]chan *Session) *Queue {
	pq := *playerQueue
	for {
		if q == nil {
			return nil
		}
		playerA := q.players[0]
		playerB := q.players[1]
		if playerA != nil && playerB != nil {
			go func() {
				for {
					time.Sleep(1000 * time.Millisecond)
					if pq == nil || playerA == nil || playerB == nil {
						return
					}
					if pq[playerA.Id()] != nil && pq[playerB.Id()] != nil {
						playerChanA := pq[playerA.Id()]
						playerChanB := pq[playerB.Id()]
						buf := make([]byte, 32)
						_, err := rand.Read(buf)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						sessionid := hex.EncodeToString(buf)
						session := Session{
							PlayerA:   playerA,
							PlayerB:   playerB,
							Sessionid: sessionid,
							Mutex:     &sync.Mutex{},
						}
						playerChanB <- &session
						playerChanA <- &session
						q.next.Consume(playerQueue)
						delete(pq, playerA.Id())
						delete(pq, playerB.Id())
					}
				}
			}()
			return q.next
		}
	}
}
