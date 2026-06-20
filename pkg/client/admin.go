package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ifuryst/ard/pkg/ard"
)

type AdminListOptions struct {
	Kind      string
	Type      string
	Status    string
	PageSize  int
	PageToken string
}

type AdminReviewOptions struct {
	PageSize  int
	PageToken string
}

type AdminAuditOptions struct {
	PageSize  int
	PageToken string
}

type AdminCatalogImportResponse struct {
	Entries int `json:"entries"`
}

type AdminStatusResponse struct {
	Identifier        string `json:"identifier"`
	Status            string `json:"status"`
	Reason            string `json:"reason,omitempty"`
	Approvals         int64  `json:"approvals,omitempty"`
	RequiredApprovals int64  `json:"requiredApprovals,omitempty"`
}

type AdminAuditEvent struct {
	ID           string    `json:"id"`
	Action       string    `json:"action"`
	Identifier   string    `json:"identifier,omitempty"`
	Status       string    `json:"status,omitempty"`
	Reason       string    `json:"reason,omitempty"`
	RequestID    string    `json:"requestId,omitempty"`
	Source       string    `json:"source"`
	RemoteAddr   string    `json:"remoteAddr,omitempty"`
	PreviousHash string    `json:"previousHash,omitempty"`
	Hash         string    `json:"hash,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AdminAuditResponse struct {
	Items     []AdminAuditEvent `json:"items"`
	Total     int64             `json:"total"`
	PageToken string            `json:"pageToken,omitempty"`
}

type AdminAuditVerification struct {
	Valid               bool   `json:"valid"`
	Total               int64  `json:"total"`
	LastHash            string `json:"lastHash,omitempty"`
	FirstInvalidEventID string `json:"firstInvalidEventId,omitempty"`
	Message             string `json:"message,omitempty"`
}

func (client *Client) AdminList(ctx context.Context, options AdminListOptions) (ard.ListResponse, error) {
	var response ard.ListResponse
	err := client.doJSON(ctx, http.MethodGet, "/admin/entries", adminListQuery(options), nil, &response)
	return response, err
}

func (client *Client) AdminReviews(ctx context.Context, options AdminReviewOptions) (ard.ListResponse, error) {
	var response ard.ListResponse
	err := client.doJSON(ctx, http.MethodGet, "/admin/reviews", pageQuery(options.PageSize, options.PageToken), nil, &response)
	return response, err
}

func (client *Client) AdminExportCatalog(ctx context.Context) (ard.Catalog, error) {
	var response ard.Catalog
	err := client.doJSON(ctx, http.MethodGet, "/admin/catalog", nil, nil, &response)
	return response, err
}

func (client *Client) AdminUpsertEntry(ctx context.Context, entry ard.CatalogEntry) (ard.CatalogEntry, error) {
	var response ard.CatalogEntry
	err := client.doJSON(ctx, http.MethodPost, "/admin/entries", nil, entry, &response)
	return response, err
}

func (client *Client) AdminUpsertCatalog(ctx context.Context, catalog ard.Catalog) (AdminCatalogImportResponse, error) {
	var response AdminCatalogImportResponse
	err := client.doJSON(ctx, http.MethodPost, "/admin/catalogs", nil, catalog, &response)
	return response, err
}

func (client *Client) AdminSetStatus(ctx context.Context, identifier string, status string) (AdminStatusResponse, error) {
	var response AdminStatusResponse
	err := client.doJSON(
		ctx,
		http.MethodPatch,
		"/admin/entries/"+url.PathEscape(identifier)+"/status",
		nil,
		map[string]string{"status": status},
		&response,
	)
	return response, err
}

func (client *Client) AdminApproveReview(ctx context.Context, identifier string, reason string) (AdminStatusResponse, error) {
	return client.adminReviewDecision(ctx, identifier, "approve", reason)
}

func (client *Client) AdminRejectReview(ctx context.Context, identifier string, reason string) (AdminStatusResponse, error) {
	return client.adminReviewDecision(ctx, identifier, "reject", reason)
}

func (client *Client) AdminDeleteEntry(ctx context.Context, identifier string) error {
	return client.doJSON(ctx, http.MethodDelete, "/admin/entries/"+url.PathEscape(identifier), nil, nil, nil)
}

func (client *Client) AdminAudit(ctx context.Context, options AdminAuditOptions) (AdminAuditResponse, error) {
	var response AdminAuditResponse
	err := client.doJSON(ctx, http.MethodGet, "/admin/audit", pageQuery(options.PageSize, options.PageToken), nil, &response)
	return response, err
}

func (client *Client) AdminVerifyAudit(ctx context.Context) (AdminAuditVerification, error) {
	var response AdminAuditVerification
	err := client.doJSON(ctx, http.MethodGet, "/admin/audit/verify", nil, nil, &response)
	return response, err
}

func (client *Client) adminReviewDecision(ctx context.Context, identifier string, action string, reason string) (AdminStatusResponse, error) {
	var response AdminStatusResponse
	err := client.doJSON(
		ctx,
		http.MethodPost,
		"/admin/reviews/"+url.PathEscape(identifier)+"/"+action,
		nil,
		map[string]string{"reason": reason},
		&response,
	)
	return response, err
}

func adminListQuery(options AdminListOptions) url.Values {
	query := pageQuery(options.PageSize, options.PageToken)
	if strings.TrimSpace(options.Kind) != "" {
		query.Set("kind", options.Kind)
	}
	if strings.TrimSpace(options.Type) != "" {
		query.Set("type", options.Type)
	}
	if strings.TrimSpace(options.Status) != "" {
		query.Set("status", options.Status)
	}
	return query
}

func pageQuery(pageSize int, pageToken string) url.Values {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}
	return query
}
