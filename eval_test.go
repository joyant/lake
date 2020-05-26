package main

import "testing"

func TestEval(t *testing.T)  {
    if !evalTest("name != nil", Parameter{"name":0}) {
       t.Errorf("expect true get false, name is not nil")
    }
    if !evalTest("name == nil", Parameter{}) {
       t.Errorf("expect true get false, name is nil")
    }
    if !evalTest("name == empty", Parameter{"name":0}) {
       t.Errorf("expect true get false, name(0) is empty")
    }
    if !evalTest("name == empty", Parameter{"name":""}) {
       t.Errorf("expect true get false, name(\"\") is empty")
    }
    if !evalTest("name == empty", Parameter{"name":false}) {
       t.Errorf("expect true get false, name(false) is empty")
    }
    if !evalTest("name == tom", Parameter{"name":"tom"}) {
       t.Errorf("expect true get false, name equal to tom")
    }
    if !evalTest("name == 2", Parameter{"name":"2"}) {
       t.Errorf("expect true get false, name equal to 2")
    }
    if !evalTest("name != tom", Parameter{"name":"jim"}) {
       t.Errorf("expect true get false, tom not equal to jim")
    }
    if !evalTest("name != 3", Parameter{"name":"2"}) {
       t.Errorf("expect true get false, 3 not equal to 2")
    }
    if !evalTest("name == tom and age == 11", Parameter{"name":"tom", "age":11}) {
       t.Errorf("expect true get false, name equal to tom and age equal to 11")
    }
    if !evalTest("name == tom or age == 11", Parameter{"name":"jim", "age":11}) {
       t.Errorf("expect true get false, tom != jim, but age == 11")
    }
    if evalTest("name == tom or age == 12", Parameter{"name":"jim", "age":11}) {
       t.Errorf("expect true get false, tom != jim, and 12 != 11")
    }
}

func TestEval2(t *testing.T)  {
    if !evalTest("withCollection != empty and searchCollection != empty and collectionKeyword != empty", Parameter{
        "withCollection":1,
        "searchCollection":1,
        "collectionKeyword":1,
    }) {
        t.Errorf("expect true get false")
    }
}

func TestEval3(t *testing.T)  {
    if !evalTest("x == nil or withCollection != empty and searchCollection != empty and collectionKeyword != empty", Parameter{

    }) {
        t.Errorf("expect true get false")
    }
}
