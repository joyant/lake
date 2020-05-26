package lake

import "fmt"

var logger logInterface

func init()  {
    logger = &log{}
}

type logInterface interface{
    Infof(format string, args ...interface{})
}

type log struct {}

// use stdout as output by default
func (l *log)Infof(format string, args ...interface{})  {
    fmt.Printf(format, args...)
}

func SetLogger(lg logInterface)  {
    logger = lg
}