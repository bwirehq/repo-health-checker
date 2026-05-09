package checks

func DefaultSuite() []Check {
	return []Check{
		ArchivedCheck{},
		ReadmeCheck{},
		CommitActivityCheck{},
		IssueHealthCheck{},
		PullRequestHealthCheck{},
		CICheck{},
		LicenseCheck{},
		ReleaseCheck{},
		TestHintCheck{},
		DependencyHintCheck{},
	}
}
