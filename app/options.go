package app

// CmdLineOptioner abstracts configuration options for reading parameters
// from command line.
type CmdLineOptioner interface {
	Flags() (fs NamedFlagSets)
	// 验证命令行参数是否合法、参数补全、设置默认值等
	Validate() []error
}

// CompleteableOptions abstracts options which can be completed.
type CompleteableOptions interface {
	Complete() error
}
