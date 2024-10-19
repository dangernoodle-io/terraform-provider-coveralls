package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	ContentType = "application/json; charset=utf-8"
)

type Client struct {
	resty    *resty.Client
	endpoint *url.URL
}

type Repository struct {
	Service               string   `json:"service,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Token                 string   `json:"token,omitempty"`
	CommentOnPullRequests bool     `json:"comment_on_pull_requests"`
	SendBuildStatus       bool     `json:"send_build_status"`
	FailThreshold         *float64 `json:"commit_status_fail_threshold"`
	FailChangeThreshold   *float64 `json:"commit_status_fail_change_threshold"`
	CreatedAt             string   `json:"created_at,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

type body struct {
	Repo *Repository `json:"repo"`
}

func NewCoveralls(endpoint, token string) (*Client, error) {
	client := resty.New()
	client.SetHeader("Accept", ContentType)
	client.SetHeader("Content-Type", ContentType)
	client.SetHeader("Authorization", fmt.Sprintf("token %s", token))

	u, _ := url.Parse(endpoint)

	return &Client{client, u}, nil
}

func (client *Client) Create(ctx context.Context, repository *Repository) (*Repository, error) {
	response, err := requestWithBody(ctx, client, repository).
		Post(fmt.Sprintf("%s/api/repos", client.endpoint.String()))

	return handlePostOrPut(201, response, err)
}

func (client *Client) Get(ctx context.Context, service, name string) (*Repository, error) {
	resp, err := client.resty.R().
		SetContext(ctx).
		SetResult(Repository{}).
		Get(fmt.Sprintf("%s/api/repos/%s/%s", client.endpoint.String(), service, name))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New(resp.String())
	}

	return resp.Result().(*Repository), nil
}

func (client *Client) Update(ctx context.Context, service, name string, repository *Repository) (*Repository, error) {
	response, err := requestWithBody(ctx, client, repository).
		Put(fmt.Sprintf("%s/api/repos/%s/%s", client.endpoint.String(), service, name))

	return handlePostOrPut(200, response, err)
}

func handlePostOrPut(status int, response *resty.Response, err error) (*Repository, error) {
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != status {
		return nil, errors.New(response.String())
	}

	// note: response doesn't currently include the token
	return response.Result().(*body).Repo, nil
}

func requestWithBody(ctx context.Context, client *Client, repository *Repository) *resty.Request {
	return client.resty.R().
		SetContext(ctx).
		SetBody(body{repository}).
		SetResult(body{})
}
