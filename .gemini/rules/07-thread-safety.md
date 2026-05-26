# Rule: Thread Safety

> **Applies to**: any code that shares state across goroutines.

## Protect Read/Write Maps with `sync.RWMutex`

```go
type KeycloakOAuthServiceImpl struct {
    publicKeys map[string]*rsa.PublicKey
    mu         sync.RWMutex
}

// Read path — use RLock
func (k *KeycloakOAuthServiceImpl) getKey(kid string) (*rsa.PublicKey, bool) {
    k.mu.RLock()
    defer k.mu.RUnlock()
    return k.publicKeys[kid]
}

// Write path — use Lock
func (k *KeycloakOAuthServiceImpl) storeKey(kid string, key *rsa.PublicKey) {
    k.mu.Lock()
    defer k.mu.Unlock()
    k.publicKeys[kid] = key
}
```

## Use `sync.Map` for Concurrent Append-Only Key Sets

```go
var processedKeys sync.Map

if _, loaded := processedKeys.LoadOrStore(key, struct{}{}); loaded {
    return // already seen this message
}
defer processedKeys.Delete(key) // clean up after processing
```

## Goroutines Must Have a `WaitGroup` or Done Channel

```go
var wg sync.WaitGroup
for i := 0; i < workerCount; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for msg := range ch {
            handle(msg)
        }
    }()
}
// Signal workers to stop by closing ch, then:
wg.Wait()
```

## Rules

- **Never** launch an unbounded number of goroutines — use a fixed-size worker pool with a buffered channel.
- **Always** use `defer` for `mu.Unlock()` / `mu.RUnlock()` to prevent deadlocks on early returns.
- Prefer `sync.RWMutex` over `sync.Mutex` when reads vastly outnumber writes.
- A goroutine that owns shared state must be its sole writer — avoid concurrent writes from multiple goroutines without locking.
