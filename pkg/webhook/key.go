package webhook

type Key struct {
	// Org is the organization name into which a commit was pushed.
	Org string

	// Rep is the repository name into which a commit was pushed.
	Rep string

	// Bra is the branch name into which a commit was pushed.
	Bra string
}
