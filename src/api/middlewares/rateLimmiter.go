package middlewares

import (
	"clean_api/src/api/helpers"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func RateLimiter(next http.Handler) http.Handler{
	type client struct {
		lastSeen time.Time
		limiter *rate.Limiter
	}
	var (
		mu sync.Mutex
		clients = map[string]*client{}
	)
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 10 * time.Minute{
					delete(clients,ip)
				}
			}
			mu.Unlock()

		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil{
			helpers.WriteResponse(w,helpers.GenerateResponseWithError(nil, false, 404, err))
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found{
			clients[ip] = &client{lastSeen: time.Now(), limiter: rate.NewLimiter(rate.Limit(2),3)}

		} 
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			helpers.WriteResponse(w,helpers.GenerateResponseWithError(nil, false, 404, errors.New("rate limit exceeded")))
			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}