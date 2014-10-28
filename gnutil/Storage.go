package gnutil

type Storage struct {
	space	map[string]interface{}
}

func NewStorage() *Storage {
	return &Storage {
		space:	make(map[string]interface{}),
	}
}

func (this *Storage) Get(key string) (interface{}, bool) {
	val, ok := this.space[key]
	return val, ok
}

func (this *Storage) Set(key string, val interface{}) {
	this.space[key] = val
}
