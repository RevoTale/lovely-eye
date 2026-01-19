package e2e

import (
	"context"
	"log/slog"
	"os"
	"testing"

	operations "github.com/lovely-eye/server/e2e/generated"
	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Those tests created by AI because I am tired
// TestInitialAdminFromConfig tests the INITIAL_ADMIN_USERNAME and INITIAL_ADMIN_PASSWORD
// configuration feature that allows creating an initial admin user on first startup.
func TestInitialAdminFromConfig(t *testing.T) {
	t.Run("creates admin when both username and password are set", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("initialadmin", "secure-password-123")
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// Should be able to login with configured credentials
		resp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "initialadmin",
			Password: "secure-password-123",
		})

		require.NoError(t, err)
		assert.Equal(t, "initialadmin", resp.Login.User.Username)
		assert.Equal(t, "admin", resp.Login.User.Role)
	})

	t.Run("prevents first user registration from becoming admin", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("preexistingadmin", "admin-password")
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// Try to register a new user (should be regular user, not admin)
		registerResp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "newuser",
			Password: "newpassword",
		})

		require.NoError(t, err)
		assert.Equal(t, "user", registerResp.Register.User.Role, "new users should not be admin when initial admin exists")
	})

	t.Run("does nothing when username is empty", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("", "somepassword")
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// First user should become admin (normal behavior)
		registerResp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "firstuser",
			Password: "password",
		})

		require.NoError(t, err)
		assert.Equal(t, "admin", registerResp.Register.User.Role)
	})

	t.Run("does nothing when password is empty", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("admin", "")
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// First user should become admin (normal behavior)
		registerResp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "firstuser",
			Password: "password",
		})

		require.NoError(t, err)
		assert.Equal(t, "admin", registerResp.Register.User.Role)
	})

	t.Run("only creates admin on first startup", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("onlyonce", "password123")
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// Verify initial admin exists
		loginResp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "onlyonce",
			Password: "password123",
		})
		require.NoError(t, err)
		assert.Equal(t, "admin", loginResp.Login.User.Role)

		// Register another user - should be regular user
		registerResp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "seconduser",
			Password: "password",
		})

		require.NoError(t, err)
		assert.Equal(t, "user", registerResp.Register.User.Role)
	})
}

// TestInitialAdminFromEnvironment tests that the feature works with actual environment variables,
// which is how it's used in production.
func TestInitialAdminFromEnvironment(t *testing.T) {
	t.Run("loads credentials from environment variables", func(t *testing.T) {
		// Set environment variables
		t.Setenv("INITIAL_ADMIN_USERNAME", "envadmin")
		t.Setenv("INITIAL_ADMIN_PASSWORD", "envpassword123")

		// Load config from environment
		cfg := config.Load()
		cfg.Database.Driver = "sqlite"
		cfg.Database.DSN = "file::memory:?cache=shared"
		cfg.Auth.AllowRegistration = true
		cfg.TrackerJS = []byte(`console.log("mock")`)

		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// Should be able to login with env credentials
		resp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "envadmin",
			Password: "envpassword123",
		})

		require.NoError(t, err)
		assert.Equal(t, "envadmin", resp.Login.User.Username)
		assert.Equal(t, "admin", resp.Login.User.Role)
	})

	t.Run("empty env vars result in no initial admin", func(t *testing.T) {
		// Explicitly unset or set to empty
		t.Setenv("INITIAL_ADMIN_USERNAME", "")
		t.Setenv("INITIAL_ADMIN_PASSWORD", "")

		cfg := config.Load()
		cfg.Database.Driver = "sqlite"
		cfg.Database.DSN = "file::memory:?cache=shared"
		cfg.Auth.AllowRegistration = true
		cfg.TrackerJS = []byte(`console.log("mock")`)

		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// First registration should create admin
		resp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "firstuser",
			Password: "password",
		})

		require.NoError(t, err)
		assert.Equal(t, "admin", resp.Register.User.Role)
	})
}

// TestInitialAdminWithRegistrationDisabled tests interaction between initial admin
// and ALLOW_REGISTRATION setting.
func TestInitialAdminWithRegistrationDisabled(t *testing.T) {
	t.Run("admin can login when registration is disabled", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("admin", "password")
		cfg.Auth.AllowRegistration = false
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// Admin login should work
		resp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "admin",
			Password: "password",
		})

		require.NoError(t, err)
		assert.Equal(t, "admin", resp.Login.User.Role)

		// Registration should fail
		_, err = operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
			Username: "newuser",
			Password: "password",
		})

		require.Error(t, err, "registration should be disabled")
	})
}

// TestInitialAdminAuthentication tests that the initial admin can perform
// authenticated operations.
func TestInitialAdminAuthentication(t *testing.T) {
	cfg := testConfigWithInitialAdmin("admin", "password")
	ts := newTestServerWithConfig(t, cfg)
	ctx := context.Background()

	// Login as initial admin
	client := ts.authenticatedClient(ctx, t, "admin", "password")

	t.Run("can create sites", func(t *testing.T) {
		siteResp, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
			Domains: []string{"example.com"},
			Name:    "Example Site",
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"example.com"}, siteResp.CreateSite.Domains)
		assert.NotEmpty(t, siteResp.CreateSite.PublicKey)
	})

	t.Run("can list sites", func(t *testing.T) {
		sitesResp, err := operations.Sites(ctx, client)

		require.NoError(t, err)
		assert.Len(t, sitesResp.Sites, 1)
	})

	t.Run("me query returns correct user", func(t *testing.T) {
		meResp, err := operations.Me(ctx, client)

		require.NoError(t, err)
		require.NotNil(t, meResp.Me)
		require.NotNil(t, meResp.Me.Username)
		assert.Equal(t, "admin", *meResp.Me.Username)
	})
}

// TestInitialAdminSecurity tests security-related aspects.
func TestInitialAdminSecurity(t *testing.T) {
	t.Run("wrong password fails", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("secureadmin", "correct-password")
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		_, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "secureadmin",
			Password: "wrong-password",
		})

		require.Error(t, err)
	})

	t.Run("username is case-sensitive", func(t *testing.T) {
		cfg := testConfigWithInitialAdmin("adminuser", "password")
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		// Correct case works
		resp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "adminuser",
			Password: "password",
		})
		require.NoError(t, err)
		assert.Equal(t, "adminuser", resp.Login.User.Username)

		// Wrong case fails
		_, err = operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "AdminUser",
			Password: "password",
		})
		require.Error(t, err)
	})

	t.Run("special characters in password", func(t *testing.T) {
		specialPassword := "P@ssw0rd!#$%&*()_+-=[]{}|;:,.<>?"
		cfg := testConfigWithInitialAdmin("specialadmin", specialPassword)
		ts := newTestServerWithConfig(t, cfg)
		ctx := context.Background()

		resp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
			Username: "specialadmin",
			Password: specialPassword,
		})

		require.NoError(t, err)
		assert.Equal(t, "specialadmin", resp.Login.User.Username)
	})
}

// TestInitialAdminEdgeCases tests various edge cases for username/password handling.
func TestInitialAdminEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{"single char username", "a", "password"},
		{"unicode username", "админ", "password"},
		{"username with spaces", "admin user", "password"},
		{"emoji in username", "admin👨‍💻", "password"},
		{"long username", "administrator_with_a_very_long_username", "password"},
		{"long password", "password", "this_is_a_secure_password_with_reasonable_length"},
		{"bcrypt max (72 bytes)", "admin", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testConfigWithInitialAdmin(tt.username, tt.password)
			ts := newTestServerWithConfig(t, cfg)
			ctx := context.Background()

			resp, err := operations.Login(ctx, ts.graphqlClient(), operations.LoginInput{
				Username: tt.username,
				Password: tt.password,
			})

			require.NoError(t, err)
			assert.Equal(t, tt.username, resp.Login.User.Username)
			assert.Equal(t, "admin", resp.Login.User.Role)
		})
	}
}

// testConfigWithInitialAdmin creates a test config with initial admin credentials.
func testConfigWithInitialAdmin(username, password string) *config.Config {
	cfg := testConfig()
	cfg.Auth.InitialAdminUsername = username
	cfg.Auth.InitialAdminPassword = password
	return cfg
}

// newTestServerWithConfig creates a test server with custom config.
func newTestServerWithConfig(t *testing.T, cfg *config.Config) *testServer {
	t.Helper()

	srv, err := server.New(cfg)
	require.NoError(t, err, "failed to create server")

	httpServer := newTestHTTPServer(srv.Handler)

	t.Cleanup(func() {
		httpServer.Close()
		err := srv.Close()
		if nil != err {
			slog.Error("server close failed","error",err)
		}
	})

	return &testServer{
		Server:     srv,
		httpServer: httpServer,
	}
}

// TestInitialAdminWithUnsetEnvVars ensures the feature works when env vars are not set at all.
func TestInitialAdminWithUnsetEnvVars(t *testing.T) {
	// Save and restore original env vars
	origUsername := os.Getenv("INITIAL_ADMIN_USERNAME")
	origPassword := os.Getenv("INITIAL_ADMIN_PASSWORD")
	defer func() {
		if origUsername != "" {
			err := os.Setenv("INITIAL_ADMIN_USERNAME", origUsername)
			if nil != err {
				slog.Error("failed to set env","error",err)
			}
		}
		if origPassword != "" {
			err:=os.Setenv("INITIAL_ADMIN_PASSWORD", origPassword)
				if nil != err {
				slog.Error("failed to set env","error",err)
			}
		}
	}()

	// Unset the env vars
	os.Unsetenv("INITIAL_ADMIN_USERNAME")
	os.Unsetenv("INITIAL_ADMIN_PASSWORD")

	cfg := config.Load()
	cfg.Database.Driver = "sqlite"
	cfg.Database.DSN = "file::memory:?cache=shared"
	cfg.Auth.AllowRegistration = true
	cfg.TrackerJS = []byte(`console.log("mock")`)

	ts := newTestServerWithConfig(t, cfg)
	ctx := context.Background()

	// First registration should create admin (normal behavior)
	resp, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "firstuser",
		Password: "password",
	})

	require.NoError(t, err)
	assert.Equal(t, "admin", resp.Register.User.Role)
}
