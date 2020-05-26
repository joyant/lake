package lake

import (
    "testing"
)

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

func TestParseExpression(t *testing.T)  {
    o, left, right, err := parseExpression("a >= b")
    if o != operatorEgt || left != "a" || right != "b" || err != nil {
        t.Error("expect true get false", o, left, right, err)
    }
    o, left, right, err = parseExpression("a > 1")
    if o != operatorGt || left != "a" || right != "1" || err != nil {
        t.Error("expect true get false", o, left, right, err)
    }
    o, left, right, err = parseExpression("a == 1")
    if o != operatorEqual || left != "a" || right != "1" || err != nil {
        t.Error("expect true get false", o, left, right, err)
    }
    o, left, right, err = parseExpression("a != 1")
    if o != operatorNotEqual || left != "a" || right != "1" || err != nil {
        t.Error("expect true get false", o, left, right, err)
    }
    o, left, right, err = parseExpression("a <= 1")
    if o != operatorElt || left != "a" || right != "1" || err != nil {
        t.Error("expect true get false", o, left, right, err)
    }
    o, left, right, err = parseExpression("a < 1")
    if o != operatorLt || left != "a" || right != "1" || err != nil {
        t.Error("expect true get false", o, left, right, err)
    }
}

func TestGetEvalFunc(t *testing.T)  {
    fn, key := getEvalFunc("len(sports)")
    if fn == nil || key != "sports" {
        t.Errorf("expect get evalLen and sports, get %v", key)
    }

    fn, key = getEvalFunc("xyz")
    if fn != nil || key != "" {
        t.Errorf("expect get nil, get %v", key)
    }
}

func TestGetNumber(t *testing.T)  {
    a, b := getNumber("len(sports)", "6", Parameter{"sports":[]string{"football", "baseball", "otherball"}})
    if a != 3 || b != 6 {
        t.Errorf("expect 3 and 6, get %v %v", a, b)
    }
    a, b = getNumber("sports", "6", Parameter{})
    if a != 0 || b != 6 {
        t.Errorf("expect 0 and 6, get %v %v", a, b)
    }
    a, b = getNumber("sports", "abc", Parameter{})
    if a != 0 || b != 0 {
        t.Errorf("expect 0 and 0, get %v %v", a, b)
    }
    a, b = getNumber("sports", "abc", Parameter{"sports": 5})
    if a != 5 || b != 0 {
        t.Errorf("expect 5 and 0, get %v %v", a, b)
    }
}

func TestEqual(t *testing.T)  {
    if !equal("a", "1", Parameter{"a":1}) {
        t.Errorf("expect true get false")
    }
    if !equal("len(a)", "3", Parameter{"a":[]interface{}{1, 2, 3}}) {
        t.Errorf("expect true get false")
    }
    if !equal("a", "0", Parameter{}) {
        t.Errorf("expect true get false")
    }
}

func TestEvalTest(t *testing.T)  {
    if !evalTest("len(a) > 1", Parameter{"a":[]int{1, 2}}) {
        t.Errorf("expect true get false")
    }
    if evalTest("a >= 1", Parameter{"a": -99}) {
        t.Errorf("expect false get true")
    }
    if !evalTest("a >= 1", Parameter{"a": 1}) {
        t.Errorf("expect true get false")
    }
    if !evalTest("len(a) <= 1", Parameter{"a":[]int{1}}) {
        t.Errorf("expect true get false")
    }
    if !evalTest("len(a) <= 1", Parameter{}) {
        t.Errorf("expect true get false")
    }
    if !evalTest("len(a) < 2", Parameter{"a":[]int{1}}) {
        t.Errorf("expect true get false")
    }
}
