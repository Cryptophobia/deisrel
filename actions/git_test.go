package actions

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/deis/deisrel/testutil"
	"github.com/google/go-github/github"
)

// TestCreateGitTag tests that creating a tag should be ok
func TestCreateGitTag(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	ts.Mux.HandleFunc("/repos/deis/controller/git/refs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "GET" {
			t.Errorf("Request method: %v, want GET or POST", r.Method)
		}
		tag := `
		{
		  "ref": "refs/heads/b",
		  "url": "https://api.github.com/repos/deis/controller/git/refs/heads/b",
		  "object": {
		    "type": "commit",
		    "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		    "url": "https://api.github.com/repos/deis/controller/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		  }
		}`
		if r.Method == "GET" {
			// return a list of tags
			fmt.Fprintf(w, "[%s]", tag)
		} else {
			fmt.Fprint(w, tag)
		}
	})

	repos := []repoAndSha{
		{
			repoName: "controller",
			sha:      "aa218f56b14c9653891f9e74264a383fa43fefbd",
		},
	}

	if err := createGitTag(ts.Client, repos, "b"); err != nil {
		t.Errorf("createGitTag returned error: %v", err)
	}

	refs, _, err := ts.Client.Git.ListRefs("deis", "controller", nil)
	if err != nil {
		t.Errorf("Git.ListRefs returned error: %v", err)
	}

	want := []github.Reference{
		{
			Ref: github.String("refs/heads/b"),
			URL: github.String("https://api.github.com/repos/deis/controller/git/refs/heads/b"),
			Object: &github.GitObject{
				Type: github.String("commit"),
				SHA:  github.String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
				URL:  github.String("https://api.github.com/repos/deis/controller/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
			},
		},
	}

	if !reflect.DeepEqual(refs, want) {
		t.Errorf("createGitTag returned %+v, want %+v", refs, want)
	}
}

// TestCreateGitTagAlreadyExists tests that creating a tag that already exists should be fine
func TestCreateGitTagAlreadyExists(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	ts.Mux.HandleFunc("/repos/deis/controller/git/refs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Request method: %v, want POST", r.Method)
		}
		w.WriteHeader(github.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{
		  "message": "Reference already exists",
		  "documentation_url": "https://developer.github.com/v3/git/refs/#create-a-reference"
		}`)
	})

	repos := []repoAndSha{
		{
			repoName: "controller",
			sha:      "aa218f56b14c9653891f9e74264a383fa43fefbd",
		},
	}

	if err := createGitTag(ts.Client, repos, "b"); err != nil {
		t.Errorf("createGitTag returned error: %v", err)
	}
}

// TestCreateGitTagBadSHA tests that creating a tag with a bad SHA should fail
func TestCreateGitTagBadSHA(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	ts.Mux.HandleFunc("/repos/deis/controller/git/refs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Request method: %v, want POST", r.Method)
		}
		w.WriteHeader(github.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{
		  "message": "Object does not exist",
		  "documentation_url": "https://developer.github.com/v3/git/refs/#create-a-reference"
		}`)
	})

	repos := []repoAndSha{
		{
			repoName: "controller",
			sha:      "aa218f56b14c9653891f9e74264a383fa43fefbd",
		},
	}

	if err := createGitTag(ts.Client, repos, "b"); err == nil {
		t.Error("createGitTag didn't return an error when the API says it's an invalid SHA")
	}
}
