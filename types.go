package main

import "math"

// 神之眼元素类型
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

// 角色属性类型
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

// 将常驻 5 星的首次 UP 池信息编码
func makePermanent5StarMeta(major, minor, half uint32) uint32 {
	encoded := (major << 6) | ((minor & 0xF) << 2) | (half & 0x3)
	return (encoded << 3) | 3
}

// CharEntry 表示一个原神角色
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

// WishPool 表示单个卡池，存储具体角色、版本、上半下半和起止日期信息
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

// PoolNode 用来表示单个卡池，作为链表节点，存储一个角色的所有 UP 历史
type PoolNode struct {
	Major uint8
	Minor uint8
	Half  uint8
	Next  *PoolNode
}

// minInt32 表示该角色从未 UP
const minInt32 = math.MinInt32
