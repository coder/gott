package main

import (
	"sort"
	"time"
)

// rankedTestTree describes a tree of ranked tests.
type rankedTestTree struct {
	// leave nil for root
	*rankedTest
	children []*rankedTestTree

	// childrenTakeCache is the sum of time it takes for child tests to complete.
	childrenTakeCache *time.Duration
}

func (r *rankedTestTree) root() bool {
	return r.rankedTest == nil
}

func (r *rankedTestTree) insert(c rankedTest) bool {
	if !r.root() && !r.parentOf(c) {
		return false
	}

	// Direct child.
	degree := -1
	if !r.root() {
		degree = r.degree()
	}
	if c.degree() == degree+1 {
		//if !r.root() {
		//	fmt.Println("appending ", c.test, " to ", r.test)
		//}
		//if r.root() {
		//	fmt.Println("appending ", c.test, " to root")
		//}

		r.children = append(r.children, &rankedTestTree{
			rankedTest: &c,
		})
		return true
	}

	for _, tch := range r.children {
		//fmt.Printf("parent %v of %v: %v\n", tch.test, c.test, tch.parentOf(c))
		if tch.insert(c) {
			return true
		}
	}

	return false
}

// sort must be called after the test is finished being written to.
func (r *rankedTestTree) sort() {
	sort.Slice(r.children, func(i, j int) bool {
		return maxDur(
			r.children[i].timeRunning, r.children[i].childrenTake(),
		) < maxDur(
			r.children[j].timeRunning, r.children[j].childrenTake(),
		)
	})

	for _, ch := range r.children {
		ch.sort()
	}
}

func maxDur(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

// childrenTake calculates how long the test takes, along with it's children.
// It attempts to produce the most useful value, and considers parallel children to
// contribute to the parent.
// childrenTake must be ran after the tree is finished being written to.
func (r *rankedTestTree) childrenTake() time.Duration {
	if r.childrenTakeCache != nil {
		return *r.childrenTakeCache
	}

	var cTake time.Duration
	for _, ch := range r.children {
		cTake += maxDur(ch.timeRunning, ch.childrenTake())
	}

	r.childrenTakeCache = &cTake
	return cTake
}
