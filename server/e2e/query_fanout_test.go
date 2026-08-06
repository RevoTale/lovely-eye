package e2e

import (
	"context"
	"sync/atomic"
	"testing"

	operations "github.com/lovely-eye/server/e2e/generated"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

type queryCounter struct {
	count atomic.Int64
}

func (counter *queryCounter) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (counter *queryCounter) AfterQuery(context.Context, *bun.QueryEvent) {
	counter.count.Add(1)
}

func TestDashboardQueryFanOut(t *testing.T) {
	ts := newTestServer(t)
	ctx := context.Background()

	_, err := operations.Register(ctx, ts.graphqlClient(), operations.RegisterInput{
		Username: "fanout-admin",
		Password: "password123",
	})
	require.NoError(t, err)
	client := ts.authenticatedClient(ctx, t, "fanout-admin", "password123")
	siteResponse, err := operations.CreateSite(ctx, client, operations.CreateSiteInput{
		Domains: []string{"fanout.example"},
		Name:    "Fan-out",
	})
	require.NoError(t, err)
	postPageView(
		t,
		ts.httpServer.Client(),
		collectURL(ts.httpServer.URL, siteResponse.CreateSite.PublicKey),
		"https://fanout.example",
		"/measured",
	)

	counter := new(queryCounter)
	ts.DB.AddQueryHook(counter)
	_, err = operations.Dashboard(
		ctx,
		client,
		siteResponse.CreateSite.Id,
		nil,
		nil,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		defaultPaging,
		nil,
		defaultPaging,
	)
	require.NoError(t, err)

	queryCount := counter.count.Load()
	require.LessOrEqual(t, queryCount, int64(11))
	t.Logf("dashboard GraphQL operation executed %d SQL queries", queryCount)
}
