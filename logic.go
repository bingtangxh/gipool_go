package main

import (
	"math"
	"sort"
	"time"
)

var daysPassedSinceLastUP []int
var arrangedInOrderOfDays []int
var poolLinkLists []*PoolNode

func initDynamicThings() {
	n := len(charMap)
	daysPassedSinceLastUP = make([]int, n)
	poolLinkLists = make([]*PoolNode, n)
	arrangedInOrderOfDays = make([]int, n)
	for i := range arrangedInOrderOfDays {
		arrangedInOrderOfDays[i] = i
	}
	getDaysPassedSinceLastUp()
	arrangeByDaysPassedSinceLastUp()
}

// poolEndHour returns the hour (local time) at which a pool ends based on its half value.
func poolEndHour(half uint8) int {
	switch half % 10 {
	case 1:
		return 18
	case 2:
		return 15
	case 3:
		return 0
	default:
		return math.MinInt32
	}
}

func makeTimeFromYMD(y uint16, m, d uint8, hour, min, sec int) time.Time {
	return time.Date(int(y), time.Month(m), int(d), hour, min, sec, 0, time.Local)
}

func daysSinceSinglePoolEnds(pool WishPool) int {
	hour := poolEndHour(pool.Half)
	if hour == math.MinInt32 {
		return minInt32
	}
	endTime := makeTimeFromYMD(pool.EndY, pool.EndM, pool.EndD, hour, 0, 0)
	diff := time.Since(endTime).Hours() / 24
	if diff < 0 {
		return int(diff) - 1
	} else if diff > 0 {
		return int(diff) + 1
	}
	return 0
}

func getDaysPassedSinceLastUp() {
	for c := range charMap {
		daysPassedSinceLastUP[c] = minInt32
		ch := charMap[c]
		for p := len(wishPool) - 1; p >= 0; p-- {
			pool := wishPool[p]
			found := false
			if ch.isFourStar() {
				for _, id := range pool.Up4 {
					if id == 0 {
						break
					}
					if id == ch.ID {
						found = true
						break
					}
				}
			} else {
				for _, id := range pool.Up5 {
					if id == 0 {
						break
					}
					if id == ch.ID {
						found = true
						break
					}
				}
			}
			if found {
				daysPassedSinceLastUP[c] = daysSinceSinglePoolEnds(pool)
				break
			}
		}
	}
}

func arrangeByDaysPassedSinceLastUp() {
	sort.Slice(arrangedInOrderOfDays, func(i, j int) bool {
		return daysPassedSinceLastUP[arrangedInOrderOfDays[i]] > daysPassedSinceLastUP[arrangedInOrderOfDays[j]]
	})
}

func checkIntegrity() int {
	errors := 0
	for i, c := range charMap {
		if c.ID != uint(i) {
			errors++
		}
	}
	return errors
}

func buildPoolLinkList(index int) {
	ch := charMap[index]
	var head, tail *PoolNode
	for _, pool := range wishPool {
		found := false
		if ch.isFourStar() {
			for _, id := range pool.Up4 {
				if id == 0 {
					break
				}
				if id == ch.ID {
					found = true
					break
				}
			}
		} else {
			for _, id := range pool.Up5 {
				if id == 0 {
					break
				}
				if id == ch.ID {
					found = true
					break
				}
			}
		}
		if found {
			node := &PoolNode{Major: pool.Major, Minor: pool.Minor, Half: pool.Half}
			if head == nil {
				head = node
				tail = node
			} else {
				tail.Next = node
				tail = node
			}
		}
	}
	poolLinkLists[index] = head
}
