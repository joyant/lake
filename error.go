package lake

import (
    "fmt"
)

func errorF(s string, argv ...interface{}) error {
    return fmt.Errorf(s, argv...)
}

func panicIfNotNil(err interface{})  {
    if err != nil {
        panic(err)
    }
}
