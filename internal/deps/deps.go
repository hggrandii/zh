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

type GitProvider string

const (
	GitHub   GitProvider = "github"
	GitLab   GitProvider = "gitlab"
	Codeberg GitProvider = "codeberg"
)

type RepoInfo struct {
	Owner         string
	Repo          string
	Provider      GitProvider
	CommitHash    string
	DefaultBranch string
}

func AddDependency(url string, provider GitProvider) error {
	repo, err := parseRepoURL(url, provider)
	if err != nil {
		return err
	}

	if err := getLatestCommitInfo(repo); err != nil {
		return err
	}

	archiveURL := generateArchiveURL(repo)
	fmt.Printf("Generated archive URL: %s\n", archiveURL)

	cmd := exec.Command("zig", "fetch", "--save", archiveURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseRepoURL(url string, provider GitProvider) (*RepoInfo, error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	if !strings.Contains(url, "/") {
		return nil, fmt.Errorf("invalid repository format: %s", url)
	}

	if !strings.Contains(url, ".") {
		parts := strings.Split(url, "/")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid repository format: %s", url)
		}

		return &RepoInfo{
			Owner:    parts[0],
			Repo:     parts[1],
			Provider: provider,
		}, nil
	}

	var domain string
	switch provider {
	case GitHub:
		domain = "github.com/"
	case GitLab:
		domain = "gitlab.com/"
	case Codeberg:
		domain = "codeberg.org/"
	}

	if strings.HasPrefix(url, domain) {
		url = strings.TrimPrefix(url, domain)
	}

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repository format: %s", url)
	}

	repo := parts[1]
	repo = strings.TrimSuffix(repo, ".git")

	return &RepoInfo{
		Owner:    parts[0],
		Repo:     repo,
		Provider: provider,
	}, nil
}

func getLatestCommitInfo(repo *RepoInfo) error {
	if err := getDefaultBranch(repo); err != nil {
		return err
	}

	var url string

	switch repo.Provider {
	case GitHub:
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s",
			repo.Owner, repo.Repo, repo.DefaultBranch)
	case GitLab:
		url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s%%2F%s/repository/branches/%s",
			repo.Owner, repo.Repo, repo.DefaultBranch)
	case Codeberg:
		url = fmt.Sprintf("https://codeberg.org/api/v1/repos/%s/%s/branches/%s",
			repo.Owner, repo.Repo, repo.DefaultBranch)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch repo.Provider {
	case GitHub:
		var response struct {
			Commit struct {
				SHA string `json:"sha"`
			} `json:"commit"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
		repo.CommitHash = response.Commit.SHA

	case GitLab:
		var response struct {
			Commit struct {
				ID string `json:"id"`
			} `json:"commit"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
		repo.CommitHash = response.Commit.ID

	case Codeberg:
		var response struct {
			Commit struct {
				ID string `json:"id"`
			} `json:"commit"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
		repo.CommitHash = response.Commit.ID
	}

	return nil
}

func getDefaultBranch(repo *RepoInfo) error {
	var url string

	switch repo.Provider {
	case GitHub:
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s", repo.Owner, repo.Repo)
	case GitLab:
		url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s%%2F%s", repo.Owner, repo.Repo)
	case Codeberg:
		url = fmt.Sprintf("https://codeberg.org/api/v1/repos/%s/%s", repo.Owner, repo.Repo)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch repo.Provider {
	case GitHub:
		var response struct {
			DefaultBranch string `json:"default_branch"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
		repo.DefaultBranch = response.DefaultBranch

	case GitLab:
		var response struct {
			DefaultBranch string `json:"default_branch"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
		repo.DefaultBranch = response.DefaultBranch

	case Codeberg:
		var response struct {
			DefaultBranch string `json:"default_branch"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
		repo.DefaultBranch = response.DefaultBranch
	}

	return nil
}

func generateArchiveURL(repo *RepoInfo) string {
	switch repo.Provider {
	case GitHub:
		return fmt.Sprintf("https://github.com/%s/%s/archive/%s.tar.gz",
			repo.Owner, repo.Repo, repo.CommitHash)
	case GitLab:
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/archive/%s/%s-%s.tar.gz",
			repo.Owner, repo.Repo, repo.CommitHash, repo.Repo, repo.CommitHash)
	case Codeberg:
		return fmt.Sprintf("https://codeberg.org/%s/%s/archive/%s.tar.gz",
			repo.Owner, repo.Repo, repo.CommitHash)
	default:
		return ""
	}
}
