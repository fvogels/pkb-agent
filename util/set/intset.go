package set

import "pkb-agent/util"

type IntSet struct {
	members []bool
}

func NewIntSet() *IntSet {
	return &IntSet{
		members: nil,
	}
}

func NewIntSetWithInitialCapacity(capacity int) *IntSet {
	return &IntSet{
		members: make([]bool, capacity),
	}
}

func (set *IntSet) Add(item int) {
	set.growIfNecessary(item)
	set.members[item] = true
}

func (set *IntSet) Contains(item int) bool {
	if item < len(set.members) {
		return set.members[item]
	} else {
		return false
	}
}

func (set *IntSet) growIfNecessary(item int) {
	if item >= len(set.members) {
		newMembers := make([]bool, (item+1)*2)
		copy(newMembers, set.members)
		set.members = newMembers
	}
}

func (set *IntSet) IntersectWith(other *IntSet) {
	imax := util.MinInt(len(set.members), len(other.members))

	i := 0
	for i != imax {
		set.members[i] = set.members[i] && other.members[i]
		i++
	}

	for i < len(set.members) {
		set.members[i] = false
		i++
	}
}

func (set *IntSet) UnionWith(other *IntSet) {
	imax := util.MinInt(len(set.members), len(other.members))
	for i := range imax {
		set.members[i] = set.members[i] || other.members[i]
	}

	if len(other.members) > len(set.members) {
		set.members = append(set.members, other.members[len(set.members):]...)
	}
}
