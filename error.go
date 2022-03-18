package lrucache

// NilErr the key not found
const NilErr = LRUCacheError("Not found")

// ExpiredErr the key was expired
const ExpiredErr = LRUCacheError("Expired")

type LRUCacheError string

func (e LRUCacheError) Error() string {
	return string(e)
}
