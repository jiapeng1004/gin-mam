package tx

type Propagation int

const (
	Required Propagation = iota
	RequiresNew
	Mandatory
	NotSupported
	Never
)
