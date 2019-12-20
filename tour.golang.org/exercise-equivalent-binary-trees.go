package main

import "golang.org/x/tour/tree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if (t == nil) {
		return
	}
	Walk(t.Left, ch)
	fmt.Println(t.Value)
	ch <- t.Value
	Walk(t.Right, ch)
	
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int, 10)
	c2 := make(chan int, 10)
	go Walk(t1, c1)
	go Walk(t2, c2)
	for i:=0; i < 10; i++ {
		v1 := <- c1
		v2 := <- c2
		
		if (v1 != v2) {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(Same(tree.New(5), tree.New(5)))
}
