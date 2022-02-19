package main

import (
	"fmt"
	"sort"

	xxhash "github.com/cespare/xxhash/v2"
)

type ConsistentHash struct {
	ring           map[uint64]Server
	sortedRingKeys []uint64
}

func (c *ConsistentHash) addServer(serverName string) {
	for i := 0; i < 3; i++ {
		serverName_v := fmt.Sprintf("%s%d", serverName, i)
		serverKey := xxhash.Sum64([]byte(serverName_v))
		c.ring[serverKey] = Server{make(map[uint64]string)}
		c.sortedRingKeys = append(c.sortedRingKeys, serverKey)
	}
	sort.Slice(c.sortedRingKeys, func(i int, j int) bool {
		return c.sortedRingKeys[i] < c.sortedRingKeys[j]
	})
}

func (c *ConsistentHash) insertKey(key string, value string) (bool, error) {
	keyHash := xxhash.Sum64([]byte(key))
	fmt.Println("KeyHash: ", keyHash)
	i := sort.Search(len(c.sortedRingKeys), func(i int) bool {
		return c.sortedRingKeys[i] > keyHash
	})
	fmt.Println("index: ", i)
	c.ring[c.sortedRingKeys[i]].data[keyHash] = value
	return true, nil
}

func (c *ConsistentHash) getKey(key string) (string, error) {
	keyHash := xxhash.Sum64([]byte(key))
	fmt.Println("KeyHash: ", keyHash)
	i := sort.Search(len(c.sortedRingKeys), func(i int) bool {
		return c.sortedRingKeys[i] > keyHash
	})
	value := c.ring[c.sortedRingKeys[i]].data[keyHash]
	return value, nil

}

func New() *ConsistentHash {
	return &ConsistentHash{make(map[uint64]Server), make([]uint64, 0)}
}

type Server struct {
	data map[uint64]string
}

func main() {
	chash := New()

	chash.addServer("server1")

	chash.addServer("server2")

	chash.addServer("server3")

	chash.insertKey("Emin", "FF")

	chash.insertKey("Ali", "FF2")

	chash.insertKey("Dersdfsdfsdf", "FF3")

	val, _ := chash.getKey("Ali")

	fmt.Println(val)

}
