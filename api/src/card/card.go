package card

import (
	"ed/util"
	"encoding/json"
)

type Event struct {
}

// Card Basic Agreement Between all cards
type Card interface {
	Id() string
	Name() string
	Event() Event
	Tags() []string
	TraitExps() []TraitExp //

	SetProp(int32, int32, int32)
	Prop() (int32, int32, int32) //

	Engage()
	Disengage()
	Engaged() bool //

	Bow()
	Unbow()
	Bowed() bool

	Kill()
	Revive()
	Dead() bool

	Lock()
	Unlock()
	Locked() bool

	MarshalJSON() ([]byte, error)
	Type() string
}

//
type BaseCard struct {
	tags      []string
	traitExps []TraitExp
	name      string
	id        string
	engaged   bool
	bowed     bool
	dead      bool
	locked    bool
}

func (b *BaseCard) Tags() []string {
	return b.tags
}

func (b *BaseCard) TraitExps() []TraitExp {
	return b.traitExps
}

func (b *BaseCard) SetProp(attack, armor, life int32) {
}

func (b *BaseCard) Prop() (int32, int32, int32) {
	return 0, 0, 0
}

func (b *BaseCard) Name() string {
	return b.name
}

func (b *BaseCard) Id() string {
	return b.id
}

func (b *BaseCard) Event() Event {
	return Event{}
}

func (b *BaseCard) Engage() {
	if b.locked == false {
		b.engaged = true
		b.Bow()
	}
}
func (b *BaseCard) Disengage() {
	if b.locked == false {
		b.engaged = false
		b.Unbow()
	}
}
func (b *BaseCard) Engaged() bool {
	return b.engaged
}

func (b *BaseCard) Bow() {
	if b.locked == false {
		b.bowed = true
	}
}
func (b *BaseCard) Unbow() {
	if b.locked == false {
		b.bowed = false
	}
}
func (b *BaseCard) Bowed() bool {
	return b.bowed
}

func (b *BaseCard) Kill() {
	if b.locked == false {
		b.dead = true
	}
}
func (b *BaseCard) Revive() {
	if b.locked == false {
		b.dead = false
	}
}
func (b *BaseCard) Dead() bool {
	return b.dead
}

func (b *BaseCard) Lock() {
	b.locked = true
}
func (b *BaseCard) Unlock() {
	b.locked = false
}
func (b *BaseCard) Locked() bool {
	return b.locked
}

func (b *BaseCard) MarshalMap() map[string]interface{} {
	return map[string]interface{}{
		"Name":      b.Name(),
		"Id":        b.Id(),
		"Dead":      b.Dead(),
		"Engaged":   b.Engaged(),
		"Bowed":     b.Bowed(),
		"Locked":    b.Locked(),
		"Tags":      b.Tags(),
		"TraitExps": b.TraitExps(),
	}
}

//
type EmpireCard struct {
	BaseCard
	attackPower int32
	armor       int32
	lifeforce   int32
}

func (e *EmpireCard) SetProp(attack, armor, life int32) {
	e.attackPower = attack
	e.armor = armor
	e.lifeforce = life
}

func (e *EmpireCard) Prop() (int32, int32, int32) {
	return e.attackPower, e.armor, e.lifeforce
}

func (e *EmpireCard) MarshalJSON() ([]byte, error) {
	att, arm, lif := e.Prop()
	cmap := e.MarshalMap()
	cmap["Attackpower"] = att
	cmap["Armor"] = arm
	cmap["Lifeforce"] = lif
	cmap["Type"] = e.Type()
	return json.Marshal(cmap)
}

func (e *EmpireCard) MarshalRQL() (interface{}, error) {
	att, arm, lif := e.Prop()
	cmap := e.MarshalMap()
	cmap["Attackpower"] = att
	cmap["Armor"] = arm
	cmap["Lifeforce"] = lif
	cmap["Type"] = e.Type()
	return cmap, nil
}

func (e *EmpireCard) Type() string {
	return "empire"
}

func CreateEmpireCard(attack, armor, life int32, tags []string, traitExps []TraitExp, name string) *EmpireCard {
	card := &EmpireCard{
		BaseCard: BaseCard{
			tags:      tags,
			traitExps: traitExps,
			name:      name,
			id:        util.Hash(name),
		},
		attackPower: attack,
		armor:       armor,
		lifeforce:   life,
	}
	return card
}

//
type DestinyCard struct {
	BaseCard
}

func (d *DestinyCard) MarshalJSON() ([]byte, error) {
	cmap := d.MarshalMap()
	return json.Marshal(cmap)
}

func (d *DestinyCard) MarshalRQL() (interface{}, error) {
	cmap := d.MarshalMap()
	cmap["Type"] = d.Type()
	return cmap, nil
}

func (d *DestinyCard) Type() string {
	return "destiny"
}

func CreateDestinyCard(tags []string, traitExps []TraitExp, name string) *DestinyCard {
	card := &DestinyCard{
		BaseCard{
			tags:      tags,
			traitExps: traitExps,
			name:      name,
			id:        util.Hash(name),
		},
	}
	return card
}

//
type FortessCard struct {
	EmpireCard
}

func (f *FortessCard) Type() string {
	return "fortess"
}

func CreateFortressCard(attack, armor, life int32, tags []string, traitExps []TraitExp, name string) *FortessCard {
	card := &FortessCard{
		EmpireCard{
			BaseCard: BaseCard{
				tags:      tags,
				traitExps: traitExps,
				name:      name,
				id:        util.Hash(name),
			},
			attackPower: attack,
			armor:       armor,
			lifeforce:   life,
		},
	}
	return card
}

//
type GodCard struct {
	DestinyCard
}

func (g *GodCard) Type() string {
	return "god"
}

func CreateGodCard(tags []string, traitExps []TraitExp, name string) *GodCard {
	card := &GodCard{
		DestinyCard{
			BaseCard{
				tags:      tags,
				traitExps: traitExps,
				name:      name,
				id:        util.Hash(name),
			},
		},
	}
	return card
}
