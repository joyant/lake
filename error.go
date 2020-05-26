package main

import (
    "fmt"
)

func errorF(s string, argv ...interface{}) error {
    return fmt.Errorf(s, argv...)
}
