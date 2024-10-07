package lrumap

import (
	"fmt"

	"github.com/breise/lrumap/lrulist"
)

type stats struct {
	puts, gets, hits, misses, evictions int
}

type LruMap struct {
	theMap  map[interface{}]*lrulist.Node
	lruList *lrulist.LruList
	stats
}

func New() *LruMap {
	return &LruMap{
		theMap:  map[interface{}]*lrulist.Node{},
		lruList: lrulist.New(),
	}
}

func (lrumap *LruMap) NItems() int {
	return lrumap.lruList.NItems()
}

func (lrumap *LruMap) MaxItems(x int) *LruMap {
	lrumap.lruList.MaxItems(x)
	return lrumap
}

/* Get retrieves the item and update its position in the lru list
 */
func (lrumap *LruMap) Get(k interface{}) (value interface{}, ok bool) {
	lrumap.gets++
	node, ok := lrumap.theMap[k]
	if !ok {
		lrumap.misses++
		return nil, ok
	}
	lrumap.hits++
	lrumap.lruList.Update(node)
	kv, ok := node.Item.(kvPair)
	if !ok {
		panic(fmt.Sprintf("the retrieved item is not a kvPair.  It is a %T", kv))
	}
	v := kv.v
	return v, ok
}

type kvPair struct{ k, v interface{} }

func (lrumap *LruMap) Put(k, v interface{}) {
	if item, ok := lrumap.theMap[k]; ok {
		// This is a reassignment of an existing key
		// TODO: 
		//    remove item from lruList
		//    continue with lruList.Add() and theMap[k] = node, below
		return
	}
	lrumap.puts++
	node, dropped := lrumap.lruList.Add(kvPair{k, v})
	lrumap.theMap[k] = node
	for i, doomed := range dropped {
		kv, ok := doomed.(kvPair)
		if !ok {
			panic(fmt.Sprintf("the %dth dropped is not a kvPair.  It is a %T", i, doomed))
		}
		delete(lrumap.theMap, kv.k)
	}
	lrumap.evictions += len(dropped)
}
