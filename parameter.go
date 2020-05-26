package lake

type Parameter map[string]interface{}

func (p Parameter)Get(key string) (value interface{}, exists bool) {
    value, exists = p[key]
    return
}

func (p Parameter)Set(key string, value interface{})  {
    p[key] = value
}
