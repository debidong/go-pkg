package hashring

import "sync"

type typeHashFn string

const (
	TypeHashFnMurmur3 typeHashFn = "murmur3"

	DefaultReplicas = 1600
	DefaultHashFunc = TypeHashFnMurmur3
)

type Ring struct {
	replicas int
	hashFunc typeHashFn
	nodes    map[string][]vNode
	hashRing []vNode

	sync.RWMutex
}

// vNode is the struct for virtual node.
type vNode struct {
	hash uint32
	node string
}

type RingOption func(*Ring)

func NewRing(opts ...RingOption) *Ring {
	r := new(Ring)
	r.Lock()
	defer r.Unlock()
	r.nodes = make(map[string][]vNode)
	for _, opt := range opts {
		opt(r)
	}

	// default
	r.replicas = DefaultReplicas
	r.hashFunc = DefaultHashFunc
	return r
}

func (r *Ring) WithReplicas(replicas int) *Ring {
	if replicas <= 0 {
		panic("WithReplicas: replicas must be greater than 0")
	}
	r.replicas = replicas
	return r
}

func (r *Ring) WithHashFn(hashFn typeHashFn) *Ring {
	r.hashFunc = hashFn
	return r
}
