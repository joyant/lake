package lake

import (
    "fmt"
    "reflect"
    "strings"
)

func evalTest(s string, params Parameter) bool {
    orSlice := strings.Split(s, " or ")
    var result = false
    if len(orSlice) > 1 {
        for _, v := range orSlice {
            result = result || evalAnd(v, params)
        }
    } else {
        result = evalAnd(s, params)
    }
    return result
}

func evalAnd(s string, params Parameter) bool {
    andSlice := strings.Split(s, " and ")
    var result = true
    if len(andSlice) > 1 {
        for _, v := range andSlice {
            result = result && evalExpression(v, params)
        }
    } else {
        result = evalExpression(s, params)
    }
    return result
}

func evalExpression(s string, params Parameter) bool {
    notEqualSlice := strings.Split(s, "!=")
    if len(notEqualSlice) == 2 {
        return !equal(notEqualSlice[0], notEqualSlice[1], params)
    }
    equalSlice := strings.Split(s, "==")
    if len(equalSlice) == 2 {
        return equal(equalSlice[0], equalSlice[1], params)
    }
    panic("evalExpression error:" + s)
}

// if key exist, value will not equal to nil
// if key is 0, "0", false, '', [], map[], value is empty ,but not nil
// key is something like "name", value is
func equal(key, value string, params Parameter) bool {
    key = strings.TrimSpace(key)
    value = strings.TrimSpace(value)
    if v, ok := params[key]; ok {
        vs := fmt.Sprintf("%v", v)
        if value == "empty" {
            return vs == "" ||
                vs == "0" ||
                (reflect.TypeOf(v).Kind().String() == "bool" && v.(bool) == false) ||
                vs == "[]" ||
                vs == "map[]"
        } else if value == "nil" {
            return false
        } else {
            return vs == value
        }
    } else {
        if value == "empty" {
            return true
        }
        return value == "nil"
    }
}
