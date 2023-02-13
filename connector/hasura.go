package connector

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

type HasuraSQLError struct {
	Path         string `json:"path,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (h *HasuraSQLError) Error() string {
	return fmt.Sprintf("error: %s, path: %s", h.ErrorMessage, h.Path)
}

type HasuraSQLRequest struct {
	Type string               `json:"type"`
	Args HasuraSQLRequestArgs `json:"args"`
}

type HasuraSQLRequestArgs struct {
	Source string `json:"source"`
	SQL    string `json:"sql"`
}

type HasuraSQLResponse struct {
	ResultType string `json:"result_type,omitempty"`
	Result     any    `json:"result,omitempty"`
}

type HasuraPostgresConnector struct {
	hasuralURL  string
	adminSecret string
	cl          *http.Client
}

func NewHasuraPostgresConnector(hasuraURL string, adminSecret string) *HasuraPostgresConnector {
	return &HasuraPostgresConnector{
		hasuralURL:  hasuraURL,
		adminSecret: adminSecret,
		cl:          &http.Client{},
	}
}

func (h *HasuraPostgresConnector) makeHTTPRequest(ctx context.Context, sql string) (HasuraSQLResponse, error) {
	// TODO: allow setting read_only as that improves performance
	reqBody := HasuraSQLRequest{
		Type: "run_sql",
		Args: HasuraSQLRequestArgs{
			Source: "default",
			SQL:    sql,
		},
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return HasuraSQLResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.hasuralURL, bytes.NewReader(b))
	if err != nil {
		return HasuraSQLResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hasura-Admin-Secret", h.adminSecret)

	resp, err := h.cl.Do(req)
	if err != nil {
		return HasuraSQLResponse{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp := &HasuraSQLError{}
		if err := json.NewDecoder(resp.Body).Decode(errResp); err != nil {
			panic(err)
		}
		return HasuraSQLResponse{}, errResp
	}

	var respBody HasuraSQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return HasuraSQLResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return respBody, nil
}

func (h *HasuraPostgresConnector) TestConnection(ctx context.Context) error {
	_, err := h.makeHTTPRequest(ctx, "SELECT 1")
	return err
}

type UserByRefreshTokenResponse struct {
	ID          string
	DefaultRole string
	IsAnonymous bool
	Roles       []string
}

func (h *HasuraPostgresConnector) GetUserByRefreshToken(ctx context.Context, refreshToken string) (UserByRefreshTokenResponse, error) {
	sha256 := sha256.Sum256([]byte(refreshToken))
	hashedToken := fmt.Sprintf(`\x%x`, sha256)

	// TODO: insecure. only for testing, sql injection possible
	sql := fmt.Sprintf(`SELECT u.id, u.default_role, u.is_anonymous,r.role
            FROM auth.refresh_tokens AS rt
            JOIN auth.users AS u ON rt.user_id = u.id
            JOIN auth.user_roles AS r on rt.user_id = r.user_id
            WHERE rt.refresh_token_hash = '%s'
              AND u.disabled = false
              AND rt.expires_at > NOW();`, hashedToken)
	res, err := h.makeHTTPRequest(ctx, sql)
	if err != nil {
		return UserByRefreshTokenResponse{}, err
	}

	resp, ok := res.Result.([]any)
	if !ok {
		return UserByRefreshTokenResponse{}, fmt.Errorf("failed to cast result to []any")
	}

	if len(resp) == 1 {
		return UserByRefreshTokenResponse{}, fmt.Errorf("user not found")
	}

	user := UserByRefreshTokenResponse{
		Roles: make([]string, 0, 1),
	}
	headers := resp[0].([]any)
	for _, r := range resp[1:] {
		row := r.([]any)
		for i, header := range headers {
			switch header.(string) {
			case "id":
				user.ID = row[i].(string)
			case "default_role":
				user.DefaultRole = row[i].(string)
			case "is_anonymous":
				user.IsAnonymous = row[i].(string) == "t"
			case "role":
				user.Roles = append(user.Roles, row[i].(string))
			}
		}
	}

	return user, nil
}

func (h *HasuraPostgresConnector) DeleteExpiredRefreshToken(ctx context.Context) error {
	sql := `DELETE FROM "auth"."refresh_tokens" WHERE expires_at < NOW();`
	_, err := h.makeHTTPRequest(ctx, sql)
	return err
}

func (h *HasuraPostgresConnector) UpdateUserLastSeen(ctx context.Context, userID string) error {
	// TODO: insecure. only for testing, sql injection possible
	sql := fmt.Sprintf(`UPDATE "auth"."users" SET last_seen = NOW() WHERE id = '%s';`, userID)
	_, err := h.makeHTTPRequest(ctx, sql)
	return err
}
