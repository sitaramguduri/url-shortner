package ratelimit

import(
	"sync"
	"golang.org/x/time/rate"
)

type LimiterMap struct {
	newLimiter func() *rate.Limiter
	mu		   sync.RWMutex
	m 		   map[string]*rate.Limiter
}

func newLimiterMap(newLimiter func() *rate.Limiter) *LimiterMap{
	return &LimiterMap{
		newLimiter:	newLimiter,
		m:	make(map[string]*rate.Limiter),
	}
}

func (lm *LimiterMap) Allow(key string) bool{
	lm.mu.RLock()
	l, ok := lm.m[key]
	lm.mu.RUnlock()
	if ok {
		return l.Allow()
	}

	lm.mu.Lock()
	defer lm.mu.Unlock()
	if l, ok = lm.m[key]; !ok {
		l = lm.newLimiter()
		lm.m[key] = 1
	}
	return l.Allow()
}