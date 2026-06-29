package spark

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	http *resty.Client
}

type Config struct {
	BaseURL     string
	AccessToken string
	Timeout     time.Duration
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	IsCreator   bool   `json:"is_creator"`
	IsVerified  bool   `json:"is_verified"`
	CreatedAt   string `json:"created_at"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Stream struct {
	ID          string   `json:"id"`
	CreatorID   string   `json:"creator_id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Thumbnail   string   `json:"thumbnail_url,omitempty"`
	Status      string   `json:"status"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	IsLive      bool     `json:"is_live"`
	ViewerCount int      `json:"viewer_count"`
	StartedAt   string   `json:"started_at,omitempty"`
	EndedAt     string   `json:"ended_at,omitempty"`
	CreatedAt   string   `json:"created_at"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type StreamListResponse struct {
	Streams    []Stream   `json:"streams"`
	Pagination Pagination `json:"pagination"`
}

func NewClient(cfg Config) *Client {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.spark.dev/api/v1"
	}

	c := resty.New()
	c.SetBaseURL(baseURL)
	c.SetTimeout(cfg.Timeout)
	if cfg.Timeout == 0 {
		c.SetTimeout(30 * time.Second)
	}

	if cfg.AccessToken != "" {
		c.SetAuthToken(cfg.AccessToken)
	}

	return &Client{http: c}
}

func (c *Client) SetAccessToken(token string) {
	c.http.SetAuthToken(token)
}

func (c *Client) Register(email, username, password string) (*AuthResponse, error) {
	var result AuthResponse
	_, err := c.http.R().
		SetBody(map[string]string{
			"email":    email,
			"username": username,
			"password": password,
		}).
		SetResult(&result).
		Post("/auth/register")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) Login(email, password string) (*AuthResponse, error) {
	var result AuthResponse
	_, err := c.http.R().
		SetBody(map[string]string{
			"email":    email,
			"password": password,
		}).
		SetResult(&result).
		Post("/auth/login")
	if err != nil {
		return nil, err
	}
	if result.AccessToken != "" {
		c.SetAccessToken(result.AccessToken)
	}
	return &result, nil
}

func (c *Client) Me() (*User, error) {
	var result User
	_, err := c.http.R().
		SetResult(&result).
		Get("/users/me")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetUser(id string) (*User, error) {
	var result User
	_, err := c.http.R().
		SetResult(&result).
		Get("/users/" + id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ListStreams(params map[string]string) (*StreamListResponse, error) {
	var result StreamListResponse
	req := c.http.R().SetResult(&result)
	for k, v := range params {
		req.SetQueryParam(k, v)
	}
	_, err := req.Get("/streams")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetStream(id string) (*Stream, error) {
	var result Stream
	_, err := c.http.R().
		SetResult(&result).
		Get("/streams/" + id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateStream(title, description, category string, tags []string) (*Stream, error) {
	var result Stream
	_, err := c.http.R().
		SetBody(map[string]any{
			"title":       title,
			"description": description,
			"category":    category,
			"tags":        tags,
		}).
		SetResult(&result).
		Post("/streams")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetBalance() (map[string]any, error) {
	var result map[string]any
	_, err := c.http.R().
		SetResult(&result).
		Get("/wallet/balance")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Search(query string, params map[string]string) (map[string]any, error) {
	var result map[string]any
	req := c.http.R().
		SetQueryParam("q", query).
		SetResult(&result)
	for k, v := range params {
		req.SetQueryParam(k, v)
	}
	_, err := req.Get("/search")
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Ensure Client implements http.Client compatible interface
var _ http.RoundTripper = (*clientTransport)(nil)

type clientTransport struct {
	inner *resty.Client
}

func (t *clientTransport) RoundTrip(req *http.Request) (*http.Response, nil) {
	return nil, nil
}
