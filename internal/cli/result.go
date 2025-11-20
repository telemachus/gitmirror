package cli

type resultKind int

const (
	resultCloned resultKind = iota
	resultUpdated
	resultUpToDate
	resultFailed
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
	case resultFailed:
		cmd.results.failed = append(cmd.results.failed, res.repo)
	}
}
