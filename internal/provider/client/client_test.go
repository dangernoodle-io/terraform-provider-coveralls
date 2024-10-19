package client

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestCoverallsCreate(t *testing.T) {
	client := setup(t)

	want := &Repository{
		Service: "github",
		Name:    "username/reponame",
	}

	httpmock.RegisterResponder("POST", "https://coveralls.io/api/repos",
		postResponder(t, 201, map[string]*Repository{"repo": want}))

	got, err := client.Create(t.Context(), want)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestCoverallsGet(t *testing.T) {
	client := setup(t)

	want := &Repository{
		Service: "github",
		Name:    "username/reponame",
		Token:   "token",
	}

	httpmock.RegisterResponder("GET", "https://coveralls.io/api/repos/github/username/reponame",
		getResponder(t, 200, want))

	got, err := client.Get(t.Context(), "github", "username/reponame")

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestCoverallsGetNotFound(t *testing.T) {
	client := setup(t)

	httpmock.RegisterResponder("GET", "https://coveralls.io/api/repos/github/username/reponame",
		getResponder(t, 404, &Repository{}))

	_, err := client.Get(t.Context(), "github", "username/reponame")

	require.Error(t, err)
	require.Equal(t, "repository not found", err.Error())
}

func TestCoverallsUpdate(t *testing.T) {
	client := setup(t)

	ft := 5.0
	fct := 10.5

	want := &Repository{
		CommentOnPullRequests: true,
		SendBuildStatus:       true,
		FailThreshold:         &ft,
		FailChangeThreshold:   &fct,
	}

	httpmock.RegisterResponder("PUT", "https://coveralls.io/api/repos/github/username/reponame",
		putResponder(t, 200, map[string]*Repository{"repo": want}))

	got, err := client.Update(t.Context(), "github", "username/reponame", want)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestMarshalling(t *testing.T) {

}

func setup(t *testing.T) *Client {
	coveralls, _ := NewCoveralls("https://coveralls.io", "fake-token")

	httpmock.ActivateNonDefault(coveralls.resty.GetClient())
	t.Cleanup(httpmock.DeactivateAndReset)

	return coveralls
}

func getResponder(t *testing.T, status int, body any) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "GET", req.Method)
		return jsonResponder(t, status, body)(req)
	}
}

func postResponder(t *testing.T, status int, body any) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "POST", req.Method)
		return jsonResponder(t, status, body)(req)
	}
}

func putResponder(t *testing.T, status int, body any) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "PUT", req.Method)
		return jsonResponder(t, status, body)(req)
	}
}

func jsonResponder(t *testing.T, status int, body any) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		require.Equal(t, ContentType, req.Header.Get("Accept"))
		require.Equal(t, ContentType, req.Header.Get("Content-Type"))
		require.Equal(t, "token fake-token", req.Header.Get("Authorization"))

		resp, _ := httpmock.NewJsonResponse(status, body)
		return httpmock.ResponderFromResponse(resp)(req)
	}
}
