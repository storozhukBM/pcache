# PCache
`PCache` is goroutine local cache that is safe for concurrent use.

`PCache` does its best to cache items inside and do as little synchronization as possible,
 but since it is cache, there is no guarantee that `PCache` won't evict your item after Store.

`PCache` eviction policy is based on GC cycles, so it can evict all items from time to time.

You can limit the size of the cache per goroutine, and `PCache` evicts random items
 if I goroutine local cache achieves maxSizePerGoroutine size.
 
You can also use PCache as a superfast, tiny cache in front of another globally synchronized cache.

## How to use it
You can import it as a library via Go modules or copy-paste `pcache.go` file
 it to your project and specify exact types that you need for your cache
 to achieve even better performance.

## Benches

<img width="1914" alt="Screenshot 2020-09-14 at 23 25 26" src="https://user-images.githubusercontent.com/3532750/93144547-26aaf000-f6e2-11ea-99ec-cd8f43c20cec.png">