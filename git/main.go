package git

type Commit struct {
	Hash        string
	Author      string
	Descritpion string
	Message     string
	Prevous     string
	Tree        []Tree
}

type Tree struct {
	Mode int
	Type string
	Hash string
	Path string
}

type Repo struct {
	Path    string
	Commits []Commit
}

func (r *Repo) PullCommits() {
	// git --no-pager log --pretty=format:'{"commit": "%h", "author": "%ae", "date": "%ah", "Description": "%s", "Message": "%b"},'
}

func (r *Repo) GetTree(c Commit) Tree {
	// git ls-tree -r --full-name f2ba1e3

	return Tree{}
}

func (r *Repo) GetFile(hash, path string) string {}

func (r *Repo) GetDiff(hash1, hash2, path string) string {}

git --no-pager log -L:TestSelector:selector.go --pretty=format:'%h %ad | %s%d [%an]'
