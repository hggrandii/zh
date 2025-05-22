package deps

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
