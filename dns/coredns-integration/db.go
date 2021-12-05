package coredns_integration

type BlockedHostsList interface {
	Get(string) (bool, error)
}

var globalBlockedHostsList BlockedHostsList

type dotRemover struct {
	internal BlockedHostsList
}

func (dr *dotRemover) Get(host string) (bool, error) {
	return dr.internal.Get(host[:len(host)-1])
}

func NewAdapter(adapter BlockedHostsList) BlockedHostsList {
	return &dotRemover{
		internal: adapter,
	}
}
