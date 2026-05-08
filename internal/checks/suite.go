package checks

func DefaultSuite() []Check {
	return []Check{
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
