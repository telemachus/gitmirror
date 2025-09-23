package cli

type resultKind int

const (
	resultCloned resultKind = iota
	resultUpdated
	resultUpToDate
	resultError
)

type result struct {
	repo string
	kind resultKind
}

func (cmd *cmdEnv) collectResult(res result) {
	switch res.kind {
	case resultCloned:
		cmd.results.cloned = append(cmd.results.cloned, res.repo)
	case resultUpdated:
		cmd.results.updated = append(cmd.results.updated, res.repo)
	case resultUpToDate:
		cmd.results.upToDate = append(cmd.results.upToDate, res.repo)
	case resultError:
		cmd.results.errors = append(cmd.results.errors, res.repo)
	}
}
