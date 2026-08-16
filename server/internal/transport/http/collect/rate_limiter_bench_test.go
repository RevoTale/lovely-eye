package collect

import (
	"strconv"
	"testing"
)

func BenchmarkRateLimiter(b *testing.B) {
	b.Run("existing_key", func(b *testing.B) {
		limiter := NewRateLimiter(true, 120, 240)
		b.ReportAllocs()
		for b.Loop() {
			limiter.Allow("ip|203.0.113.10")
		}
	})

	b.Run("new_keys", func(b *testing.B) {
		const keyCount = 1024
		keys := make([]string, keyCount)
		for index := range keys {
			keys[index] = "ip|" + strconv.Itoa(index)
		}
		limiter := NewRateLimiter(true, 120, 240)

		b.ReportAllocs()
		b.ResetTimer()
		for index := range b.N {
			if index%keyCount == 0 {
				limiter = NewRateLimiter(true, 120, 240)
			}
			limiter.Allow(keys[index%keyCount])
		}
	})
}
