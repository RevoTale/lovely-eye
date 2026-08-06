package dashboard

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testIndexHTML = `<!doctype html>
<html>
<head>
  <base href="{{BASE_PATH}}/" />
  <script src="{{BASE_PATH}}/config.js"></script>
</head>
<body>dashboard</body>
</html>`

func TestHandlerAppliesRuntimeBasePathAndSPAFallback(t *testing.T) {
	tests := []struct {
		name             string
		basePath         string
		expectedBaseHref string
		expectedEnvPath  string
	}{
		{name: "root", basePath: "", expectedBaseHref: `<base href="/" />`, expectedEnvPath: ""},
		{
			name:             "single segment",
			basePath:         "/lovely-eye",
			expectedBaseHref: `<base href="/lovely-eye/" />`,
			expectedEnvPath:  "/lovely-eye",
		},
		{
			name:             "nested",
			basePath:         "/tools/lovely-eye",
			expectedBaseHref: `<base href="/tools/lovely-eye/" />`,
			expectedEnvPath:  "/tools/lovely-eye",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newTestDashboardHandler(t, tt.basePath)

			for _, route := range []string{"/", "/sites/new"} {
				response := serveDashboardRequest(handler, route)
				require.Equal(t, http.StatusOK, response.Code)
				require.Contains(t, response.Body.String(), tt.expectedBaseHref)
			}

			configResponse := serveDashboardRequest(handler, "/config.js")
			require.Equal(t, http.StatusOK, configResponse.Code)
			require.Contains(t, configResponse.Body.String(), "BASE_PATH: '"+tt.expectedEnvPath+"'")

			assetResponse := serveDashboardRequest(handler, "/assets/app.js")
			require.Equal(t, http.StatusOK, assetResponse.Code)
			require.Equal(t, "console.log('ok');", assetResponse.Body.String())

			missingAssetResponse := serveDashboardRequest(handler, "/assets/missing.js")
			require.Equal(t, http.StatusNotFound, missingAssetResponse.Code)
		})
	}
}

func newTestDashboardHandler(t *testing.T, basePath string) http.Handler {
	t.Helper()
	dashboardPath := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dashboardPath, "assets"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dashboardPath, "index.html"), []byte(testIndexHTML), 0o600))
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(dashboardPath, "assets", "app.js"),
			[]byte("console.log('ok');"),
			0o600,
		),
	)
	return Handler(Config{
		BasePath:      basePath,
		APIUrl:        basePath + "/api",
		GraphQLUrl:    basePath + "/graphql",
		DashboardPath: dashboardPath,
	})
}

func serveDashboardRequest(handler http.Handler, requestPath string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, requestPath, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
