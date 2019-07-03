package main

import (
	"sort"
)

// rankedTestTree describes a tree of ranked tests.
type rankedTestTree struct {
	// leave nil for root
	*rankedTest
	children []*rankedTestTree
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

func (r *rankedTestTree) sort() {
	sort.Slice(r.children, func(i, j int) bool {
		return r.children[i].timeRunning < r.children[j].timeRunning
	})

	for _, ch := range r.children {
		ch.sort()
	}
}
