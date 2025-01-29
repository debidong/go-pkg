package hashring

import "sync"

type typeHashFn string

const (
	TypeHashFnMurmur3 typeHashFn = "murmur3"
	TypeHashFnMd5     typeHashFn = "md5"

	DefaultReplicas = 1600
	DefaultHashFunc = TypeHashFnMurmur3

	modulo int64 = 1 << 32
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
	// default
	r.nodes = make(map[string][]vNode)
	r.replicas = DefaultReplicas
	r.hashFunc = DefaultHashFunc

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithReplicas(replicas int) RingOption {
	return func(r *Ring) {
		if replicas <= 0 {
			panic("WithReplicas: replicas must be greater than 0")
		}
		r.replicas = replicas
	}
}

func WithHashFn(hashFn typeHashFn) RingOption {
	return func(r *Ring) {
		r.hashFunc = hashFn
	}
}
