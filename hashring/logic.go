package hashring

import (
	"crypto/md5"
	"fmt"
	"math/big"
	"sort"

	"github.com/debidong/go-pkg/utils"
	"github.com/spaolacci/murmur3"
)

type hashFunc func(string) uint32

func murmur3Hash(str string) uint32 {
	return murmur3.Sum32([]byte(str)) // always be in range of [0, 2^32 - 1]
}

func md5Hash(str string) uint32 {
	hash := md5.Sum([]byte(str))
	// convert hash to integer
	hashInt := big.NewInt(0).SetBytes(hash[:])
	hashInt.Mod(hashInt, big.NewInt(modulo))
	return uint32(hashInt.Int64())
}

func (r *Ring) hasher() hashFunc {
	switch r.hashFunc {
	case TypeHashFnMurmur3:
		return murmur3Hash
	case TypeHashFnMd5:
		return md5Hash
	default:
		return murmur3Hash
	}
}

func (r *Ring) AddNode(node string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.nodes[node]; ok {
		return utils.NewError("node already exists", "AddNode")
	}
	nodes := make([]vNode, r.replicas)
	for i := 0; i < r.replicas; i++ {
		nodes[i] = vNode{
			hash: r.hasher()(fmt.Sprintf("%s-%d", node, i)),
			node: node,
		}
	}
	r.nodes[node] = nodes
	r.updateHashRing()
	return nil
}

func (r *Ring) RemoveNode(node string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.nodes[node]; !ok {
		return utils.NewError("node not found", "RemoveNode")
	}
	delete(r.nodes, node)
	r.updateHashRing()
	return nil
}

func (r *Ring) GetNodeByKey(key string) (string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.hashRing) == 0 {
		return "", utils.NewError("no nodes in hash ring", "GetNodeByKey")
	}

	hash := r.hasher()(key)
	idx := sort.Search(len(r.hashRing), func(i int) bool {
		return r.hashRing[i].hash >= hash
	})
	if idx == len(r.hashRing) {
		idx = 0 // ring
	}
	return r.hashRing[idx].node, nil
}

func (r *Ring) updateHashRing() {
	r.hashRing = make([]vNode, 0)
	for _, nodes := range r.nodes {
		r.hashRing = append(r.hashRing, nodes...)
	}
	sort.Slice(r.hashRing, func(i, j int) bool {
		return r.hashRing[i].hash < r.hashRing[j].hash
	})
}
