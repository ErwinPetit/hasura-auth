package connector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nhost/hasura-auth/connector"
)

func TestTestConnection(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		url    string
		secret string
	}{
		{
			name:   "works",
			url:    "http://localhost:8080/v1/query",
			secret: "hello123",
		},
		{
			name:   "wrong pass",
			url:    "http://localhost:8080/v1/query",
			secret: "hello124",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc := tc

			hasura := connector.NewHasuraPostgresConnector(tc.url, tc.secret)
			err := hasura.TestConnection(context.Background())
			if err != nil {
				t.Fatalf("failed to test connection: %v", err)
			}
		})
	}
}

func TestGetUserByRefreshToken(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		url    string
		secret string
	}{
		{
			name:   "works",
			url:    "http://localhost:8080/v1/query",
			secret: "hello123",
		},
		// {
		// 	name:   "wrong pass",
		// 	url:    "http://localhost:8080/v1/query",
		// 	secret: "hello124",
		// },
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc := tc

			hasura := connector.NewHasuraPostgresConnector(tc.url, tc.secret)
			id, err := hasura.GetUserByRefreshToken(
				context.Background(),
				"0ec5ec69-13d4-4048-a526-ed3df102a85e",
			)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(id)
		})
	}
}

func TestDeleteExpiredRefreshToken(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		url    string
		secret string
	}{
		{
			name:   "works",
			url:    "http://localhost:8080/v1/query",
			secret: "hello123",
		},
		// {
		// 	name:   "wrong pass",
		// 	url:    "http://localhost:8080/v1/query",
		// 	secret: "hello124",
		// },
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc := tc

			hasura := connector.NewHasuraPostgresConnector(tc.url, tc.secret)
			err := hasura.DeleteExpiredRefreshToken(
				context.Background(),
			)
			if err != nil {
				t.Fatalf("%va", err)
			}
		})
	}
}

func TestUpdateUserLastSeen(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		url    string
		secret string
	}{
		{
			name:   "works",
			url:    "http://localhost:8080/v1/query",
			secret: "hello123",
		},
		// {
		// 	name:   "wrong pass",
		// 	url:    "http://localhost:8080/v1/query",
		// 	secret: "hello124",
		// },
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc := tc

			hasura := connector.NewHasuraPostgresConnector(tc.url, tc.secret)
			err := hasura.UpdateUserLastSeen(
				context.Background(),
				"bc6ca4cb-d826-4369-b0f3-3e960b9bb8e8",
			)
			if err != nil {
				t.Fatalf("%va", err)
			}
		})
	}
}
