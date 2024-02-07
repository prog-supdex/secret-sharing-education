package server

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter interface {
	IpRateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler
	allow(ip string) bool
	fillBucket(ip string)
}

type rateLimiter struct {
	config ConfigRequest
	ips    map[string]chan bool
	mu     sync.Mutex
}

type Message struct {
	Status string
	Body   string
}

func NewRateLimit(c ConfigRequest) RateLimiter {
	return &rateLimiter{config: c, ips: make(map[string]chan bool)}
}

func (r *rateLimiter) IpRateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.mu.Lock()

		if !r.allow(ip) {
			slog.Info("To many requests for IP", "ip", ip)

			r.mu.Unlock()

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

		r.mu.Unlock()
		next(w, req)
	})
}

func (r *rateLimiter) allow(ip string) bool {
	bucket, exists := r.ips[ip]

	if !exists {
		slog.Debug("the bucket is not exists. Creating...",
			"ip", ip,
			"bucketSize", r.config.RequestsLimit,
		)
		bucket = r.preparedBucket()

		r.ips[ip] = bucket

		go r.fillBucket(ip)
	}

	select {
	case <-bucket:
		return true
	default:
		return false
	}
}

// add "token" to the bucket every "within / requestLimit". It allows us to add "tokens" gradually.
// For example, if "within" is 60 and "requestLimit" is 2, it means that "tokens" will be added every 30 seconds, and
// the user will be able to spend "token" every 30 seconds (or two tokens at once if tokens are collected)
func (r *rateLimiter) fillBucket(ip string) {
	ticker := time.NewTicker(time.Duration(r.config.Within) * time.Second / time.Duration(r.config.RequestsLimit))

	for range ticker.C {
		r.mu.Lock()

		if _, exists := r.ips[ip]; !exists {
			r.mu.Unlock()
			ticker.Stop()
			return
		}

		slog.Debug("Add the token to the bucket", "ip", ip)

		select {
		case r.ips[ip] <- true:
		default:
		}

		r.mu.Unlock()
	}
}

func (r *rateLimiter) preparedBucket() chan bool {
	bucket := make(chan bool, r.config.RequestsLimit)
	for i := 0; i < r.config.RequestsLimit; i++ {
		bucket <- true
	}

	return bucket
}
