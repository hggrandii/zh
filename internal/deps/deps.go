package deps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type BranchResponse struct {
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

type RepoResponse struct {
	DefaultBranch string `json:"default_branch"`
}

func AddDependency(url string) error {
	owner, repo, err := parseGitHubURL(url)
	if err != nil {
		return err
	}

	commitHash, err := getLatestCommitHash(owner, repo)
	if err != nil {
		return err
	}

	archiveURL := fmt.Sprintf("https://github.com/%s/%s/archive/%s.tar.gz", owner, repo, commitHash)
	fmt.Printf("Generated archived url %s\n", archiveURL)

	cmd := exec.Command("zig", "fetch", "--save", archiveURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()

}

func parseGitHubURL(url string) (string, string, error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	if strings.HasPrefix(url, "github.com/") {
		url = strings.TrimPrefix(url, "github.com/")
	}

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid Github format %s", url)
	}

	repo := parts[1]
	repo = strings.TrimSuffix(repo, ".git")
	return parts[0], repo, nil
}

func getLatestCommitHash(owner, repo string) (string, error) {
	defaultBranch, err := getDefaultBranch(owner, repo)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", owner, repo, defaultBranch)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Github API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var branch BranchResponse
	if err := json.NewDecoder(resp.Body).Decode(&branch); err != nil {
		return "", err
	}

	return branch.Commit.SHA, nil
}

func getDefaultBranch(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Github API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var repoData RepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return "", err
	}

	return repoData.DefaultBranch, nil
}
