package main

import (
	"fmt"
)

type Node struct {
	Val  interface{}
	Next *Node
}

type ListNode struct {
	Head   *Node
	Length int
}

type Method interface {
	Insert(i int, v interface{})
	Delete(i int)
	GetLength() int
	Search(v interface{}) int
	isNull() bool
}

func CreateNode(v interface{}) *Node {
	return &Node{v, nil}
}

func CreateListNode() *ListNode {
	return &ListNode{CreateNode(nil), 0}
}

func (list *ListNode) Insert(i int, v interface{}) {
	if i > list.Length+1 {
		return
	}
	s := CreateNode(v)
	cur := list.Head

	for count := 0; count <= i; count++ {
		if count == i-1 {
			s.Next = cur.Next
			cur.Next = s
			list.Length++
			break
		}
		cur = cur.Next
	}
}

func (list *ListNode) Delete(i int) {
	if i > list.Length {
		return
	}

	cur := list.Head
	for count := 1; count <= i; count++ {
		if count == i {
			cur.Next = cur.Next.Next
			list.Length--
		}
		cur = cur.Next
	}

}

func (list *ListNode) GetLength() int {
	return list.Length
}
func (list *ListNode) Search(v interface{}) int {
	cur := list.Head.Next

	for i := 1; i <= list.Length; i++ {
		if v == cur.Val {
			return i
		}
		cur = cur.Next
	}
	return 0
}
func (list *ListNode) isNull() bool {
	cur := list.Head
	if cur.Next == nil {
		return true
	}
	return false
}

func printList(list *ListNode) {
	cur := list.Head.Next

	fmt.Println("show ListNode:...")
	for i := 1; i <= list.Length; i++ {
		fmt.Println(cur.Val)
		cur = cur.Next
	}

}

func main() {
	list := CreateListNode()
	printList(list)

	list.Insert(1, 5)
	list.Insert(2, 4)
	list.Insert(1, 6)
	list.Insert(4, 9)
	printList(list)
	fmt.Println("length:", list.GetLength())

	list.Delete(4)
	printList(list)
	fmt.Println("length:", list.GetLength())

	fmt.Println("search index:", list.Search(9))
}
