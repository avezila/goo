package goo

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var sources = make(chan rand.Source, runtime.NumCPU())

func init() {
	for i := 0; i < cap(sources); i++ {
		sources <- rand.NewSource(time.Now().UnixNano() + int64(i))
	}
}

func RandString(n int) string {
	src := <-sources
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	sources <- src
	return string(b)
}

type Hasher struct {
	generator chan string
	closer    chan chan struct{}
	vals      map[string]string
	hashs     map[string]string
	mutex     sync.RWMutex
	nowLength int
}

func (h *Hasher) Get(hash string) string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.vals[hash]
}

func (h *Hasher) Set(hash, val string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.vals[hash] = val
	h.hashs[val] = hash
}

func (h *Hasher) Insert(val string) (string, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	hash, ok := h.hashs[val]
	if ok {
		return hash, true
	}
	select {
	case hash = <-h.generator:
	default:
		hash = h.generate(true)
	}
	h.hashs[val] = hash
	h.vals[hash] = val
	return hash, false
}

func (h *Hasher) Delete(hash string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	val, ok := h.vals[hash]
	if !ok {
		return
	}
	delete(h.vals, hash)
	delete(h.hashs, val)
}
func (h *Hasher) Start() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.closer != nil {
		return
	}
	h.closer = make(chan chan struct{})
	for i := 0; i < cap(h.generator); i++ {
		go h.loop()
	}
}
func (h *Hasher) Stop() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.closer == nil {
		return
	}
	for i := 0; i < cap(h.generator); i++ {
		ch := make(chan struct{})
		h.closer <- ch
		<-ch
	}
}

func NewHasher() *Hasher {
	return &Hasher{
		vals:      make(map[string]string),
		hashs:     make(map[string]string),
		generator: make(chan string, runtime.NumCPU()),
		nowLength: 1,
	}
}

func (h *Hasher) loop() {
	hash := h.generate(false)
	for {
		select {
		case ch := <-h.closer:
			delete(h.vals, hash)
			ch <- struct{}{}
			return
		case h.generator <- hash:
			hash = h.generate(false)
		}
	}
}

func (h *Hasher) generate(locked bool) string {
	var hash string
	for {
		hash = RandString(h.nowLength)
		if !locked {
			h.mutex.Lock()
		}
		_, ok := h.vals[hash]
		if !ok {
			h.vals[hash] = ""
		}
		if !locked {
			h.mutex.Unlock()
		}
		if !ok {
			break
		}
		h.nowLength++
	}
	return hash
}
