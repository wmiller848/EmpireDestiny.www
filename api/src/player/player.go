package player

import (
	"ed/card"
	"errors"
	"fmt"

	rethink "github.com/dancannon/gorethink"
)

type Move struct {
}

func (m *Move) CardsToPlay() []string {
	return []string{}
}

type DeckDB struct {
	EmpireDeck  []int32
	DestinyDeck []int32
	FortessCard int32
	GodCard     int32
}

type PlayerDB struct {
	Id    string             `Id`
	Decks map[string]*DeckDB `Decks`
	Cards []string           `Cards`
}

func CreatePlayerDB(id string) *PlayerDB {
	return &PlayerDB{
		Id: id,
	}
}

func LoadPlayerDB(id string, session *rethink.Session) (*PlayerDB, error) {
	cursor, err := rethink.DB("ed").Table("players").Get(id).Run(session)
	if err != nil {
		return nil, err
	}
	pdb := PlayerDB{}

	if cursor.IsNil() == false {
		err = cursor.One(&pdb)
		fmt.Println(cursor, pdb, err)
		if err != nil {
			return nil, err
		} else {
			return &pdb, nil
		}
	} else {
		return nil, errors.New("No user found")
	}
}

func (p *PlayerDB) Save(session *rethink.Session) error {
	resp, err := rethink.DB("ed").Table("players").Insert(p, rethink.InsertOpts{
		Conflict: "replace",
	}).RunWrite(session)
	if err != nil {
		return err
	}
	fmt.Println("Save", resp)
	return nil
}

func (p *PlayerDB) LoadPlayer(session *rethink.Session) (*Player, error) {
	cursor, err := rethink.DB("ed").Table("players").Get(p.Id).Run(session)
	if err != nil {
		return nil, err
	}
	pdb := PlayerDB{}
	suc := cursor.Next(&pdb)
	fmt.Println(suc, pdb)
	// dck := card.GetDeckFromIds(plr.Deck(deckid))
	// dck := card.GetDeckFromIds(plr.Deck(deckid))
	// p.deckName = dck.Name()
	// p.deck = &Deck{}
	return CreatePlayer(p.Id), nil
}

func (p *PlayerDB) AddDeck(session *rethink.Session, deckid string, deck *DeckDB) error {
	if p.Decks == nil {
		p.Decks = make(map[string]*DeckDB)
	}
	p.Decks[deckid] = deck
	return p.Save(session)
}

type Deck struct {
	EmpireDeck  card.Deck
	DestinyDeck card.Deck
	FortessCard card.FortessCard
	GodCard     card.GodCard
}

type Player struct {
	id        string
	deckName  string
	deck      *Deck
	hand      []card.Card
	field     []card.Card
	districts []card.Card
	gold      int32
}

func CreatePlayer(id string) *Player {
	return &Player{
		id:        id,
		hand:      []card.Card{},
		field:     []card.Card{},
		districts: []card.Card{},
		deck:      &Deck{},
	}
}

func (p *Player) Cards() (card.Deck, card.Deck) {
	return p.deck.EmpireDeck, p.deck.DestinyDeck
}

func (p *Player) Hand() []card.Card {
	return p.hand
}

func (p *Player) Field() []card.Card {
	return p.field
}

func (p *Player) Destrict() []card.Card {
	return p.districts
}

func (p *Player) God() card.Card {
	return &p.deck.GodCard
}

func (p *Player) Fortress() card.Card {
	return &p.deck.FortessCard
}

func (p *Player) SetGold(value int32) {
	p.gold = value
}

func (p *Player) Gold() int32 {
	return p.gold
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) SetId(id string) {
	p.id = id
}

func (p *Player) Shuffle() {
	p.deck.EmpireDeck.Shuffle()
	p.deck.DestinyDeck.Shuffle()
}
