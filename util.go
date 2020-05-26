package lake

type defaultMap map[string]string

func (d defaultMap)get(key, defaultValue string) string {
    if v, ok := d[key]; ok {
        return v
    } else {
        return defaultValue
    }
}

func (d defaultMap)set(key, value string)  {
    d[key] = value
}

func panicErrNotNil(err error)  {
    if err != nil {
        panic(err)
    }
}
