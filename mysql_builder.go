package lake

import (
    "bytes"
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

type mySQLBuilder struct {}

func (b *mySQLBuilder)build(stmt []byte, argv Parameter) (sql string, params []interface{}, err error) {
    buff := bytes.Buffer{}
    for i, l := 0, len(stmt); i < l; {
        s := stmt[i]
        if (s == '#' || s == '$') && i+1 < l && stmt[i+1] == '{' {
            for j := i+2; j < l; j++ {
                if stmt[j] == '}' {
                    key, value := splitKeyValue(stmt[i+2:j])
                    if _, ok := argv[key]; !ok {
                        err = errors.New("param " + key + " not exist")
                        return
                    }
                    switch value.(type) {
                    case nil:
                        if s == '#' {
                            params = append(params, argv[key])
                        } else {
                            buff.WriteString(fmt.Sprintf("%v", argv[key]))
                        }
                    case int:
                        v := reflect.ValueOf(argv[key]).Index(value.(int)).Interface()
                        if s == '#' {
                            params = append(params, v)
                        } else {
                            buff.WriteString(fmt.Sprintf("%v", v))
                        }
                    case string: // item.0.key
                        v := value.(string)
                        if dot := strings.IndexByte(v, '.'); dot > -1 {
                            index, err := strconv.Atoi(v[:dot])
                            panicErrNotNil(err)
                            vv := reflect.ValueOf(argv[key]).Index(index).Interface().(map[string]interface{})[v[dot+1:]]
                            if s == '#' {
                                params = append(params, vv)
                            } else {
                                buff.WriteString(fmt.Sprintf("%v", vv))
                            }
                        }
                    }
                    if s == '#' {
                        buff.WriteByte('?')
                    }
                    i = j + 1
                    break
                }
            }
        } else {
            buff.WriteByte(stmt[i])
            i ++
        }
    }
    return buff.String(), params, nil
}

func (*mySQLBuilder)lastSQL(sql string, params []interface{}) string {
    for _, v := range params {
        sql = strings.Replace(sql, "?", fmt.Sprintf("'%v'", v), 1)
    }
    return sql
}

// splitKeyValue separate key and value in string like "item.name" or "item.0"
// name is string, 0 is int, so returning value is interface type
func splitKeyValue(b []byte) (key string, value interface{}) {
    if index := bytes.IndexByte(b, '.'); index > -1 {
        if bytesIsInt(b[index+1:]) {
            i, _ := strconv.Atoi(string(b[index+1:]))
            return string(b[:index]), i
        }
        return string(b[:index]), string(b[index+1:])
    }
    return string(b), nil
}

func bytesIsInt(bs []byte) bool {
    for _, v := range bs {
        if v < '0' || v > '9' {
            return false
        }
    }
    return true
}

func LeftLike(v interface{}) string {
    return fmt.Sprintf("%%%v", v)
}

func RightLike(v interface{}) string {
    return fmt.Sprintf("%v%%", v)
}

func Like(v interface{}) string {
    return fmt.Sprintf("%%%v%%", v)
}

