package main

import "math"

// Vision element types
const (
	VisionOther   = 0
	Pyro          = 1
	Hydro         = 2
	Anemo         = 3
	Electro       = 4
	Dendro        = 5
	Cryo          = 6
	Geo           = 7
	VisionUnknown = 8
)

// Role attribute types
const (
	RoleTypeTravelerAether  = 0
	RoleTypeTravelerLumine  = 1
	RoleTypeCollab          = 2
	// attrib%4==3 means permanent 5-star that had UP events
	RoleTypeFourStar        = 4
	RoleTypeLimitedFiveStar = 5
	RoleTypeUnknown         = 6
	RoleTypeExcluded        = 256
)

// makePermanent5StarMeta encodes version info for permanent 5-stars that had UP events.
func makePermanent5StarMeta(major, minor, half uint32) uint32 {
	encoded := (major << 6) | ((minor & 0xF) << 2) | (half & 0x3)
	return (encoded << 3) | 3
}

// CharEntry represents a Genshin Impact character.
type CharEntry struct {
	ID     uint
	NameCN string
	Name   string
	Vision uint8
	Attrib uint32
}

func (c CharEntry) isFourStar() bool {
	return c.Attrib == RoleTypeFourStar
}

func (c CharEntry) isLimitedOrUPFiveStar() bool {
	return c.Attrib == RoleTypeLimitedFiveStar || (c.Attrib > 3 && c.Attrib%4 == 3)
}

// WishPool represents a single banner/wish pool version.
type WishPool struct {
	Up5    [24]uint
	Up4    [24]uint
	Weapon [24]uint
	Major  uint8
	Minor  uint8
	Half   uint8
	StartY uint16
	StartM uint8
	StartD uint8
	EndY   uint16
	EndM   uint8
	EndD   uint8
}

// PoolNode is a linked list node for a character's pool history.
type PoolNode struct {
	Major uint8
	Minor uint8
	Half  uint8
	Next  *PoolNode
}

// minInt32 is used as a sentinel "never appeared in UP" value.
const minInt32 = math.MinInt32
