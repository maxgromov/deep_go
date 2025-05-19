package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	nameMaxLength = 42

	// attributeMask битовые позиции
	manaBitsStart   = 0
	manaBitsLen     = 10
	healthBitsStart = 10
	healthBitsLen   = 10

	flagHouseBit  = 20
	flagGunBit    = 21
	flagFamilyBit = 22
	typeBitsStart = 23
	typeBitsLen   = 2

	// statsMask битовые позиции
	respectBitsStart    = 0
	strengthBitsStart   = 4
	experienceBitsStart = 8
	levelBitsStart      = 12
	bitsPerStat         = 4
)

// Маски
const (
	manaMask       = (1 << manaBitsLen) - 1
	healthMask     = (1 << healthBitsLen) - 1
	statMask       = (1 << bitsPerStat) - 1
	personTypeMask = (1 << typeBitsLen) - 1
)

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z       int32
	gold          uint32
	attributeMask uint32 // 0-9 бит мана, 10-19 здоровье, 20 дом, 21 оружие, 22 семья, 23-24 тип персонажа
	statsMask     uint16 // 0-3 бит уважение, 4-7 сила, 8-11 опыт, 12-15 уровень
	name          [nameMaxLength]byte
}

// В задании нужно упаковать данные игрока в структуру таким образом, чтобы ее размер был не более, чем 64 байта.
type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		n := copy(person.name[:], name)
		for i := n; i < len(person.name); i++ {
			person.name[i] = 0
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributeMask = (person.attributeMask &^ (manaMask << manaBitsStart)) | (uint32(mana&manaMask) << manaBitsStart)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributeMask = (person.attributeMask &^ (healthMask << healthBitsStart)) | (uint32(health&healthMask) << healthBitsStart)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.statsMask = (person.statsMask &^ (statMask << respectBitsStart)) | (uint16(respect&statMask) << respectBitsStart)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.statsMask = (person.statsMask &^ (statMask << strengthBitsStart)) | (uint16(strength&statMask) << strengthBitsStart)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.statsMask = (person.statsMask &^ (statMask << experienceBitsStart)) | (uint16(experience&statMask) << experienceBitsStart)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.statsMask = (person.statsMask &^ (statMask << levelBitsStart)) | (uint16(level&statMask) << levelBitsStart)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributeMask |= 1 << flagHouseBit
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributeMask |= 1 << flagGunBit
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributeMask |= 1 << flagFamilyBit
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributeMask = (person.attributeMask &^ (personTypeMask << typeBitsStart)) | ((uint32(personType) & personTypeMask) << typeBitsStart)
	}
}

func NewGamePerson(options ...Option) GamePerson {
	p := new(GamePerson)
	for _, opt := range options {
		opt(p)
	}
	return *p
}

func (p *GamePerson) Name() string {
	n := 0
	for n < len(p.name) && p.name[n] != 0 {
		n++
	}

	return string(p.name[:n])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int((p.attributeMask >> manaBitsStart) & manaMask)
}

func (p *GamePerson) Health() int {
	return int((p.attributeMask >> healthBitsStart) & healthMask)
}

func (p *GamePerson) Respect() int {
	return int((p.statsMask >> 12) & 0xF)
}

func (p *GamePerson) Strength() int {
	return int((p.statsMask >> 8) & 0xF)
}

func (p *GamePerson) Experience() int {
	return int((p.statsMask >> 4) & 0xF)
}

func (p *GamePerson) Level() int {
	return int(p.statsMask & 0xF)
}

func (p *GamePerson) HasHouse() bool {
	return p.attributeMask&(1<<flagHouseBit) != 0
}

func (p *GamePerson) HasGun() bool {
	return p.attributeMask&(1<<flagGunBit) != 0
}

func (p *GamePerson) HasFamily() bool {
	return p.attributeMask&(1<<flagFamilyBit) != 0
}

func (p *GamePerson) Type() int {
	return int((p.attributeMask >> typeBitsStart) & personTypeMask)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
