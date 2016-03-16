package card

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Deck []Card

func (d *Deck) Shuffle() {

}

type TraitExp struct {
	Exp         string `Exp`
	Targets     string `Targets`
	Name        string `Name`
	Description string `Description`
}

type EmpireCardYML struct {
	Name        string     `Name`
	AttackPower int32      `AttackPower`
	Armor       int32      `Armor`
	Lifeforce   int32      `Lifeforce`
	Tags        []string   `Tags`
	TraitExps   []TraitExp `TraitExps`
}

func LoadEmpireDeckFromYML(path, name string) (Deck, error) {
	deck := []EmpireCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	empireDeck := Deck{}

	for _, c := range deck {
		empireCard := CreateEmpireCard(c.AttackPower, c.Armor, c.Lifeforce, c.Tags, c.TraitExps, c.Name)
		empireDeck = append(empireDeck, empireCard)
	}

	return empireDeck, nil
}

type DestinyCardYML struct {
	Name      string     `Name`
	Tags      []string   `Tags`
	TraitExps []TraitExp `TraitExps`
}

func LoadDestinyDeckFromYML(path, name string) (Deck, error) {
	deck := []DestinyCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	destinyDeck := Deck{}

	for _, c := range deck {
		destinyCard := CreateDestinyCard(c.Tags, c.TraitExps, c.Name)
		destinyDeck = append(destinyDeck, destinyCard)
	}

	return destinyDeck, nil
}

func LoadFortressCardsFromYML(path string) (Deck, error) {
	deck := []EmpireCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	fortressCards := Deck{}

	for _, c := range deck {
		fortressCard := CreateFortressCard(c.AttackPower, c.Armor, c.Lifeforce, c.Tags, c.TraitExps, c.Name)
		fortressCards = append(fortressCards, fortressCard)
	}

	return fortressCards, nil
}

func LoadGodCardsFromYML(path string) (Deck, error) {
	deck := []EmpireCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	godCards := Deck{}

	for _, c := range deck {
		godCard := CreateGodCard(c.Tags, c.TraitExps, c.Name)
		godCards = append(godCards, godCard)
	}

	return godCards, nil
}
