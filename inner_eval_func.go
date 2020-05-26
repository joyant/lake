package lake

import (
    "fmt"
    "reflect"
    "unicode/utf8"
)

func init()  {
    RegisterEvalFunc("len", evalLen)
}

func evalLen(key, value string, params Parameter) interface{} {
    if v, ok := params[key]; ok {
        switch reflect.TypeOf(v).Kind() {
        case reflect.Slice, reflect.Map:
            return reflect.ValueOf(v).Len()
        default:
            return utf8.RuneCountInString(fmt.Sprintf("%v", v))
        }
    } else {
        return 0
    }
}
