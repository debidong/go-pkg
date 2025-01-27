package hashring

type typeHashFn string

const (
	TypeHashFnMurmur3 typeHashFn = "murmur3"

	DefaultReplicas = 1600
	DefaultHashFunc = TypeHashFnMurmur3
)

type Ring struct {
	replicas int
	hashFunc typeHashFn
	nodes    map[string][]uint32 // physical node -> virtual nodes
	hashRing []uint32
}

type RingOption func(*Ring)

func NewRing(opts ...RingOption) *Ring {
	r := new(Ring)
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
