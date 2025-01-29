package hashring

import (
	"fmt"
	"testing"
)

func TestNewRing(t *testing.T) {
	r := NewRing()
	fmt.Printf("ring: %v\n", r)

	r = NewRing(WithReplicas(10))
	fmt.Printf("ring: %v\n", r)

	r = NewRing(WithReplicas(10), WithHashFn(TypeHashFnMurmur3))
	fmt.Printf("ring: %v\n", r)
}

func TestAddNode(t *testing.T) {
	r := NewRing()
	r.AddNode("node1")
	r.AddNode("node2")
	r.AddNode("node3")
	cnt := make(map[string]int)
	fmt.Printf("hash ring: %v\n", r.hashRing)

	for i := 0; i < 1000; i++ {
		node, err := r.GetNodeByKey(fmt.Sprintf("key%d", i))
		if err != nil {
			t.Fatalf("GetNodeByKey failed: %v", err)
		}
		t.Logf("GetNodeByKey result: %s", node)
		cnt[node]++
	}

	fmt.Printf("cnt: %v\n", cnt)
}

func TestRemoveNode(t *testing.T) {
	r := NewRing(WithReplicas(10), WithHashFn(TypeHashFnMurmur3))
	r.AddNode("node1")
	r.AddNode("node2")
	fmt.Printf("ring: %v\n", r)
	r.RemoveNode("node2")
	fmt.Printf("ring: %v\n", r)
}

func TestHashFn(t *testing.T) {
	inputs := []string{"test1", "test2", "test3", "test4", "test5"}
	for _, input := range inputs {
		fmt.Printf("input: %s\n", input)
		fmt.Printf("murmur3 hash: %d\n", murmur3Hash(input))
		fmt.Printf("md5 hash: %d\n", md5Hash(input))
	}
}
