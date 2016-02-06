package actor

import "fmt"

type Path interface {
	fmt.Stringer
	Name() string

	join(name string) Path
}

type path struct {
	path string
	name string
}

func (p *path) String() string {
	return p.path
}

func (p *path) Name() string {
	return p.name
}

func (p *path) join(name string) Path {
	return &path{
		path: p.path + "/" + name,
		name: name,
	}
}

var rootPath = &path{
	path: "",
	name: "",
}
