package lake

import (
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

type EvalFunc func(left, right string, params Parameter) interface{}

var evalFuncMap = make(map[string]EvalFunc)

func RegisterEvalFunc(fnName string, fn EvalFunc)  {
    evalFuncMap[fnName] = fn
}

type operator string

const (
    operatorEqual = operator("==")
    operatorNotEqual = operator("!=")
    operatorEgt = operator(">=")
    operatorGt = operator(">")
    operatorElt = operator("<=")
    operatorLt = operator("<")
)

func evalTest(s string, params Parameter) bool {
    orSlice := strings.Split(s, " or ")
    var result = false
    if len(orSlice) > 1 {
        for _, v := range orSlice {
            result = result || evalAnd(v, params)
            if result == true {
                return true // short circuit
            }
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
            if result == false {
                return false // short circuit
            }
        }
    } else {
        result = evalExpression(s, params)
    }
    return result
}

func parseExpression(s string) (operator operator, left, right string, err error) {
    for i, l := 0, len(s); i < l; i++ {
        a := s[i]
        b := uint8(' ')
        if i < l {
            b = s[i+1]
        }
        switch {
        case a == '!' && b == '=':
            return operatorNotEqual, strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+2:]), nil
        case a == '=' && b == '=':
            return operatorEqual, strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+2:]), nil
        case a == '>' && b == '=':
            return operatorEgt, strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+2:]), nil
        case a == '<' && b == '=':
            return operatorElt, strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+2:]), nil
        case a == '>':
            return operatorGt, strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:]), nil
        case a == '<':
            return operatorLt, strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:]), nil
        }
    }
    err = errors.New("seek none operator")
    return
}

func evalExpression(s string, params Parameter) bool {
    op, left, right, err := parseExpression(s)
    panicIfNotNil(err)
    switch op {
    case operatorEqual:
        return equal(left, right, params)
    case operatorNotEqual:
        return !equal(left, right, params)
    case operatorEgt:
        return egt(left, right, params)
    case operatorGt:
        return gt(left, right, params)
    case operatorElt:
        return elt(left, right, params)
    case operatorLt:
        return lt(left, right, params)
    default:
        panicIfNotNil("unknown operator near:" + s)
        return false
    }
}

// if key exist, value will not equal to nil
// if key is 0, "0", false, '', [], map[], value is empty ,but not nil
// key is something like "name", value is
func equal(key, value string, params Parameter) bool {
    switch value {
    case "nil", "empty":
        return equalNilAndEmpty(key, value, params)
    default:
        return equalValue(key, value, params)
    }
}

func equalValue(key, value string, params Parameter) bool {
    a, b := getNumber(key, value, params)
    return a == b
}

func equalNilAndEmpty(key, value string, params Parameter) bool {
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

func egt(key, value string, params Parameter) bool {
    a, b := getNumber(key, value, params)
    return a >= b
}

func gt(key, value string, params Parameter) bool {
    a, b := getNumber(key, value, params)
    return a > b
}

func elt(key, value string, params Parameter) bool {
    a, b := getNumber(key, value, params)
    return a <= b
}

func lt(key, value string, params Parameter) bool {
    a, b := getNumber(key, value, params)
    return a < b
}

// getNumber return two float64 value according to parse string key and value
// assume key is "age", value is "11", params is {"age":12}, getNumber will return float64(12) and float(11)
// if key not exists, will return 0, if key is not numeric, will return 0
// if value is not numeric, will return 0
func getNumber(key, value string, params Parameter) (a float64, b float64) {
    var err error
    fn, k := getEvalFunc(key)
    if fn != nil {
        result := fn(k, value, params)
        a, err = strconv.ParseFloat(fmt.Sprintf("%v", result), 64)
    } else if v, ok := params[key]; ok{
        a, err = strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
    }
    b, err = strconv.ParseFloat(value, 64)
    _ = err
    return
}

func getEvalFunc(s string) (fn EvalFunc, key string) {
    for i, l := 0, len(s); i < l; i++ {
        if s[i] == '(' {
            for j := i+1; j < l; j++ {
                if s[j] == ')' {
                    var ok bool
                    fnName := s[:i]
                    fn, ok = evalFuncMap[fnName]
                    if ok {
                        key = s[i+1:j]
                        return
                    }
                }
            }
        }
    }
    return
}
