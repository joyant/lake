package main

import (
    "testing"
)

func TestMappers(t *testing.T)  {
    err := Load("xxx.xml") // xxx.xml is not exit
    if err == nil {
        t.Errorf("expect err get nil")
    }
}

func TestMappers2(t *testing.T)  {
    err := Load("test/mappers.xml")
    if err != nil {
        t.Errorf("expect nil get:%v", err)
    }
    namespace := "user"
    ids := []string{"getByLimit", "getById", "selectByNameOrAge"}
    for _, id := range ids {
        s, ok := ms.get(namespace, id)
        if !ok{
            t.Errorf("expect %s exit get false", id)
        }
        if s == nil {
            t.Errorf("expect not nil get nil")
        }
    }
}
