package collect

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkAnalyticsHandlerCollectPageView(b *testing.B) {
	handler, site := newAnalyticsHandlerTestFixture(b, AnalyticsHandlerConfig{
		MaxBodyBytes:       4096,
		MaxPropertiesBytes: 1024,
	}, nil)

	b.ReportAllocs()
	for b.Loop() {
		req := newAnalyticsCollectRequest(site.PublicKey, `{"path":"/pricing"}`)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/140.0.0.0 Safari/537.36")
		recorder := httptest.NewRecorder()

		handler.Collect(recorder, req)
		if recorder.Code != http.StatusNoContent {
			b.Fatalf("collect status = %d, want %d", recorder.Code, http.StatusNoContent)
		}
	}
}
