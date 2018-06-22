package changelog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
)

var (
	// ErrSHANotLongEnough is the error returned from various functions in this package when a
	// given SHA is not long enough for use
	ErrSHANotLongEnough = errors.New("SHA not long enough")
)

// SingleRepoVals generates a changelog entry from vals.OldRelease to sha. It returns the commits that were unparseable (and had to be skipped) or any error encountered during the process. On a nil error, vals is filled in with all of the sorted changelog entries. Note that any nil commits will not be in the returned string slice
func SingleRepoVals(client *github.Client, vals *Values, sha, name string, includeRepoName bool) ([]string, error) {
	var skippedCommits []string
	commitCompare, resp, err := client.Repositories.CompareCommits(context.Background(), "teamhephy", name, vals.OldRelease, sha)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errTagNotFoundForRepo{repoName: name, tagName: vals.OldRelease}
		}
		return nil, errCouldNotCompareCommits{old: vals.OldRelease, new: sha, err: err}
	}
	for _, commit := range commitCompare.Commits {
		if commit.Commit.Message == nil {
			continue
		}
		if commit.SHA == nil {
			continue
		}
		commitMessage := strings.Split(*commit.Commit.Message, "\n")[0]
		shortSHA, err := shortSHATransform(*commit.SHA)
		if err != nil {
			return nil, err
		}
		focus := commitFocus(*commit.Commit.Message)
		title := commitTitle(*commit.Commit.Message)
		shortSHALink := fmt.Sprintf("[`%s`](https://github.com/teamhephy/%s/commit/%s)", shortSHA, name, *commit.SHA)
		changelogMessage := fmt.Sprintf("%s %s: %s", shortSHALink, focus, title)
		if includeRepoName {
			changelogMessage = fmt.Sprintf("%s (%s) - %s: %s", shortSHALink, name, focus, title)
		}

		skippedCommits = appendToValues(vals, commitMessage, changelogMessage, skippedCommits)
	}
	return skippedCommits, nil
}

// shortSHATransform returns the shortened version of the given SHA given in s. If the given
// string is not long enough, returns the empty string and ErrSHANotLongEnough
func shortSHATransform(s string) (string, error) {
	if len(s) < 7 {
		return "", ErrSHANotLongEnough
	}
	return s[:7], nil
}

// appendToValues appends a changelogMessage to vals depending on commitMessage,
// returning any skipped commits (which will be appended to the provided skippedCommits)
func appendToValues(vals *Values, commitMessage, changelogMessage string, skippedCommits []string) []string {
	if strings.HasPrefix(commitMessage, "feat(") {
		vals.Features = append(vals.Features, changelogMessage)
	} else if strings.HasPrefix(commitMessage, "ref(") {
		vals.Refactors = append(vals.Refactors, changelogMessage)
	} else if strings.HasPrefix(commitMessage, "fix(") {
		vals.Fixes = append(vals.Fixes, changelogMessage)
	} else if strings.HasPrefix(commitMessage, "docs(") || strings.HasPrefix(commitMessage, "doc(") {
		vals.Documentation = append(vals.Documentation, changelogMessage)
	} else if strings.HasPrefix(commitMessage, "test(") || strings.HasPrefix(commitMessage, "tests(") {
		vals.Tests = append(vals.Tests, changelogMessage)
	} else if strings.HasPrefix(commitMessage, "chore(") {
		vals.Maintenance = append(vals.Maintenance, changelogMessage)
	} else {
		if !strings.HasPrefix(commitMessage, "Merge pull request") {
			skippedCommits = append(skippedCommits, changelogMessage)
		}
	}
	return skippedCommits
}
