package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

// app struct
// 	tree map orphans
//	tree map pages
//		pdf-document file name
//		effects
//		list index

type Key struct{}

type PageBase struct{}

func main() {
	m := treemap.NewWithStringComparator() // empty (keys are of type int)
	m.Put("1", "x")                        // 1->x
	m.Put("2", "b")                        // 1->x, 2->b (in order)

	mm := treemap.NewWithStringComparator() // empty (keys are of type int)
	h := hey{42, mm}

	m.Put("3", h) // 1->a, 2->b (in order)

	//	fmt.Println(m)

	s, _ := m.Get("3") // b, true
	fmt.Println(s.(hey).Tree)

}
