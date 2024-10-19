package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	ctx = tflog.SetField(ctx, "repository", repository)
	tflog.Debug(ctx, "Creating coveralls repository")

	response, err := requestWithBody(ctx, client, repository).
		Post(fmt.Sprintf("%s/api/repos", client.endpoint.String()))

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, handleErrorResponse(ctx, response)
	}

	// note: response doesn't currently include the token
	result, ok := response.Result().(*body)
	if !ok {
		return nil, errors.New("unexpected response format: couldn't convert to body type")
	}
	return result.Repo, nil
}

func (client *Client) Get(ctx context.Context, service, name string) (*Repository, error) {
	ctx = tflog.SetField(ctx, "service", service)
	ctx = tflog.SetField(ctx, "name", name)
	tflog.Debug(ctx, "Retrieving coveralls repository")

	response, err := client.resty.R().
		SetContext(ctx).
		SetResult(Repository{}).
		Get(fmt.Sprintf("%s/api/repos/%s/%s", client.endpoint.String(), service, name))

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, handleErrorResponse(ctx, response)
	}

	result, ok := response.Result().(*Repository)
	if !ok {
		return nil, errors.New("unexpected response format: couldn't convert to body type")
	}
	return result, nil
}

func (client *Client) Update(ctx context.Context, service, name string, repository *Repository) (*Repository, error) {
	response, err := requestWithBody(ctx, client, repository).
		Put(fmt.Sprintf("%s/api/repos/%s/%s", client.endpoint.String(), service, name))

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, handleErrorResponse(ctx, response)
	}

	// note: response doesn't currently include the token
	result, ok := response.Result().(*body)
	if !ok {
		return nil, errors.New("unexpected response format: couldn't convert to body type")
	}
	return result.Repo, nil
}

func handleErrorResponse(ctx context.Context, response *resty.Response) error {
	statusCode := response.StatusCode()

	ctx = tflog.SetField(ctx, "status_code", statusCode)
	ctx = tflog.SetField(ctx, "error_message", response.String())
	tflog.Debug(ctx, "Error response received")

	if statusCode == 404 {
		return errors.New("repository not found")
	}

	return errors.New(response.String())
}

func requestWithBody(ctx context.Context, client *Client, repository *Repository) *resty.Request {
	return client.resty.R().
		SetContext(ctx).
		SetBody(body{repository}).
		SetResult(body{})
}
