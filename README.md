# PCache  [![Build Status](https://github.com/storozhukBM/pcache/workflows/build/badge.svg)](https://github.com/storozhukBM/pcache/actions)  [![Go Report Card](https://goreportcard.com/badge/github.com/storozhukBM/pcache)](https://goreportcard.com/report/github.com/storozhukBM/pcache) [![codecov](https://codecov.io/gh/storozhukBM/pcache/branch/master/graph/badge.svg)](https://codecov.io/gh/storozhukBM/pcache) [![PkgGoDev](https://pkg.go.dev/badge/github.com/storozhukBM/pcache)](https://pkg.go.dev/github.com/storozhukBM/pcache)
`PCache` Cache that tries to keep data local for the goroutine and reduce synchronization overhead but keep it is safe for concurrent use.

`PCache` does its best to cache items inside and do as little synchronization as possible,
 but since it is cache, there is no guarantee that `PCache` won't evict your item after Store.
 
 Due to its implementation specifics, in some edge cases, 
  PCache can potentially restore previously-stored items after eviction, so please take into account that
  **it is possible and valid to observe "old" values of the specific key**. 
  While this behavior is unconventional, it is totally usable for:
   - immutable key-value pairs
   - keys that will always resolve into the same value
   - cases when it is easy for you to identify that
    the value is old and drop it or set to the new one

`PCache` eviction policy is based on GC cycles, so it can evict all items from time to time.

You can limit the size of the cache per goroutine, and `PCache` evicts random items
 if I goroutine local cache achieves maxSizePerGoroutine size.
 
You can also use PCache as a superfast, tiny cache in front of another globally synchronized cache.

## How to use it
You can import it as a library via Go modules or copy-paste `pcache.go` file
 it to your project and specify exact types that you need for your cache
 to achieve even better performance.

## Benches

<img width="1609" alt="Screenshot 2020-10-17 at 20 13 55" src="https://user-images.githubusercontent.com/3532750/96351739-4e142480-10b5-11eb-9509-78da3c9ef163.png">
