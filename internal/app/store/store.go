package store

var mapStore map[string]string

func init() {
	mapStore = make(map[string]string)
}

func Get(k string) (string, bool) {
	v, ok := mapStore[k]
	return v, ok
}

func Set(k string, v string) {
	mapStore[k] = v
}
