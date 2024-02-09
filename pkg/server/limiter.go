package server

import (
	"encoding/json"
	"github.com/jonboulle/clockwork"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter interface {
	IpRateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler
}

type rateLimiter struct {
	config  Config
	IpItems map[string]*ipItem
	Mu      sync.Mutex
	clock   clockwork.Clock
}

type ipItem struct {
	bucket      chan bool
	lastSeen    time.Time
	lastFilling time.Time
	ip          string
}

type Message struct {
	Status string
	Body   string
}

func NewRateLimit(c Config, clock clockwork.Clock) RateLimiter {
	if clock == nil {
		clock = clockwork.NewRealClock()
	}

	r := &rateLimiter{config: c, IpItems: make(map[string]*ipItem), clock: clock}

	go removeExpiredIpItems(r)

	return r
}

func (r *rateLimiter) IpRateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Mu.Lock()

		if !r.allow(ip) {
			slog.Info("To many requests for IP", "ip", ip)

			r.Mu.Unlock()

			message := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later.",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				slog.Error("Failed to write response:" + err.Error())
			}

			return
		}

		r.Mu.Unlock()
		next(w, req)
	})
}

func (r *rateLimiter) allow(ip string) bool {
	item, exists := r.IpItems[ip]

	if !exists {
		slog.Debug("the bucket is not exists. Creating...",
			"ip", ip,
			"bucketSize", r.config.RequestsLimit,
		)

		item = &ipItem{bucket: r.preparedBucket(), ip: ip, lastFilling: r.clock.Now()}
		r.IpItems[ip] = item

		slog.Debug("Created ipItem bucket", "ip", ip)
	}

	r.fillBucket(item)
	r.updateIpItemTimes(item)

	select {
	case <-item.bucket:
		return true
	default:
		return false
	}
}

func (r *rateLimiter) preparedBucket() chan bool {
	bucket := make(chan bool, r.config.RequestsLimit)
	for i := 0; i < r.config.RequestsLimit; i++ {
		bucket <- true
	}

	return bucket
}

// Refill "tokens", if the needed time has passed
func (r *rateLimiter) fillBucket(item *ipItem) {
	now := r.clock.Now()

	if now.After(item.lastFilling.Add(time.Duration(r.config.Within) * time.Second)) {
		tokensToAdd := r.config.RequestsLimit - len(item.bucket)
		for i := 0; i < tokensToAdd; i++ {
			select {
			case item.bucket <- true:
			default:
			}
		}
		item.lastFilling = now
	}
}

func (r *rateLimiter) updateIpItemTimes(item *ipItem) {
	item.lastSeen = r.clock.Now()
}

func removeExpiredIpItems(r *rateLimiter) {
	expired := time.Duration(r.config.IpBucketLifeTimeSeconds) * time.Second
	ticker := r.clock.NewTicker(30 * time.Second)

	for range ticker.Chan() {
		r.Mu.Lock()
		now := r.clock.Now()
		slog.Info("Checking the IpItems state")

		for ip, item := range r.IpItems {
			if now.After(item.lastSeen.Add(expired)) {
				slog.Info("Removing unused element", "ip", ip, "expiredSeconds", expired)
				delete(r.IpItems, ip)
			}
		}

		r.Mu.Unlock()
	}
}
