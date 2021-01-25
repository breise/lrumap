package lrulist

import "fmt"

const (
	oneMeg = 1 << 20
)

type Node struct {
	Item interface{}
	prev *Node
	next *Node
}

type LruList struct {
	head     *Node
	tail     *Node
	nItems   int
	maxItems int
}

func New() *LruList {
	rv := &LruList{
		head:     &Node{}, // to make removals consistent
		tail:     &Node{}, // to make removals consistent
		maxItems: oneMeg,
	}
	rv.head.next = rv.tail
	rv.tail.prev = rv.head
	return rv
}

func (lrulist *LruList) MaxItems(x int) *LruList {
	if x < lrulist.nItems {
		panic(fmt.Sprintf("cannot set MaxItems (%d) less than the current number of items (%d)", x, lrulist.nItems))
	}
	lrulist.maxItems = x
	return lrulist
}

func (lrulist *LruList) NItems() int { return lrulist.nItems }

func (lrulist *LruList) ToSlice() []interface{} {
	rv := []interface{}{}
	for cur := lrulist.head.next; /*this is the tail node*/ cur.next != nil; cur = cur.next {
		rv = append(rv, cur.Item)
	}
	return rv
}

func (lrulist *LruList) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
	lrulist.nItems--
}

// drop the first node; return the item
func (lrulist *LruList) drop() interface{} {
	node := lrulist.head.next
	lrulist.remove(node)
	return node.Item
}

func (lrulist *LruList) Add(item interface{}) (node *Node, dropped []interface{}) {
	dropped = []interface{}{}
	for lrulist.nItems >= lrulist.maxItems {
		dropped = append(dropped, lrulist.drop())
	}
	tailNode := lrulist.tail
	node = &Node{
		Item: item,
		prev: tailNode.prev,
		next: tailNode,
	}
	tailNode.prev.next = node
	tailNode.prev = node
	lrulist.nItems++
	return node, dropped
}

func (lrulist *LruList) Update(node *Node) {
	lrulist.remove(node)
	_, dropped := lrulist.Add(node.Item)
	if len(dropped) > 0 {
		panic("unexpectedly dropped a node while updating")
	}
}