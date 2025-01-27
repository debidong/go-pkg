package hashring

import (
	"fmt"
	"go-pkg/utils"

	"github.com/spaolacci/murmur3"
)

type hashFunc func(string) uint32

func murmur3Hash(str string) uint32 {
	return murmur3.Sum32([]byte(str))
}

func (r *Ring) hasher() hashFunc {
	switch r.hashFunc {
	case "murmur3":
		return murmur3Hash
	default:
		return murmur3Hash
	}
}

func (r *Ring) AddNode(node string) error {
	if _, ok := r.nodes[node]; ok {
		return utils.NewError("node already exists", "AddNode")
	}
	nodes := make([]uint32, r.replicas)
	for i := 0; i < r.replicas; i++ {
		nodes[i] = r.hasher()(fmt.Sprintf("%s-%d", node, i))
	}
	r.nodes[node] = nodes
	return nil
}

func (r *Ring) RemoveNode(node string) error {
	if _, ok := r.nodes[node]; !ok {
		return utils.NewError("node not found", "RemoveNode")
	}
	delete(r.nodes, node)
	return nil
}
