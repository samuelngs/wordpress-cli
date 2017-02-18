package app

type container struct {
	name    string
	created bool
	running bool
}

func (v *container) Name() string {
	return v.name
}

func (v *container) Created() bool {
	return v.created
}

func (v *container) Running() bool {
	return v.running
}
