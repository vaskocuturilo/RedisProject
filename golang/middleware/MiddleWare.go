package middleware

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var ratelimitSql = redis.NewScript(`
    local current = redis.call("INCR", KEYS[1])
    if current == 1 then
        redis.call("PEXPIRE", KEYS[1], ARGV[1])
    end
    return current
`)

const methodNotAllowed = "Method not allowed"

func AllowMethod(method string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				log.Printf("Method %s not allowed for %s", r.Method, r.URL.Path)
				http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

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
				slog.Warn("Rate limit exceeded", "ip", ip, "count", val)
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Too Many Requests\n"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
