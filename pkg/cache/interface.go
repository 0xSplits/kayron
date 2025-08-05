package cache

// Interface defines a generic concurrency safe key-value cache used to
// communicate between packages of varying business logic.
type Interface[K comparable, V any] interface {
	Create(K, V) bool
	Delete(...K)
	Exists(K) bool
	Length() int
	Search(K) (V, bool)
	Update(K, V)
}
