package webgo

type Resource interface {
	Head()
	Get()
	Post()
	Put()
	Delete()
	Patch()
}
