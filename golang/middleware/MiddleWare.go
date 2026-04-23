package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
)

var ratelimitSql = redis.NewScript(`
    local current = redis.call("INCR", KEYS[1])
    if current == 1 then
        redis.call("PEXPIRE", KEYS[1], ARGV[1])
    end
    return current
`)

var (
	ratelimitEvents = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_ratelimit_exceeded_total",
		Help: "Total number of requests blocked by rate limiter",
	}, []string{"method", "path"})
)

func RateLimiter(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ip := r.RemoteAddr
			key := fmt.Sprintf("ratelimit:%s", ip)

			val, err := ratelimitSql.Run(ctx, rdb, []string{key}, window.Milliseconds()).Int64()

			if err != nil {
				slog.Error("Redis Lua script error", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, int64(limit)-val)))

			if val > int64(limit) {
				ratelimitEvents.WithLabelValues(r.Method, r.URL.Path).Inc()
				slog.Warn("Rate limit exceeded", "ip", ip, "count", val)
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Too Many Requests\n"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
