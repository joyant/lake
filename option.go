package main

import "time"

const (
    ModeDebug = "debug"
    ModeRelease = "release"
)

type Mode string

type Option struct {
    Mode Mode
    ConnMaxLifetime time.Duration
    MaxIdleConns int
    MaxOpenConns int
}
