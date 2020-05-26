package main

type builder interface {
    build(stmt []byte, argv Parameter)(sql string, params []interface{}, err error)
    lastSQL(sql string, params []interface{}) string
}