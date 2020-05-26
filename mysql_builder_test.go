package lake

import (
    "reflect"
    "testing"
)

func TestBuilder(t *testing.T)  {
    b := mySQLBuilder{}
    sql, params, err := b.build([]byte("select * from user where name = #{name}"), Parameter{"name":"tom"})
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    if sql != "select * from user where name = ?" {
        t.Errorf("expect sql equal get false, sql:%s", sql)
    }
    if len(params) != 1 || !reflect.DeepEqual(params[0], "tom") {
        t.Errorf("expect [tom], get %v", params)
    }
}

func TestBuilder2(t *testing.T)  {
    b := mySQLBuilder{}
    sql, params, err := b.build([]byte("select * from user where name = #{name} or age = #{age}"), Parameter{"name":"tom", "age":11})
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    if sql != "select * from user where name = ? or age = ?" {
        t.Errorf("expect sql equal get false, sql:%s", sql)
    }
    if len(params) != 2 || !reflect.DeepEqual(params[0], "tom") || !reflect.DeepEqual(params[1], 11){
        t.Errorf("expect [tom, 11], get %v", params)
    }
}

func TestBytesIsNumber(t *testing.T)  {
    if !bytesIsInt([]byte("123")) {
        t.Errorf("expect true get false")
    }
    if bytesIsInt([]byte("abc")) {
        t.Errorf("expect false get true")
    }
}

func TestSplitKeyValue(t *testing.T)  {
    key, value := splitKeyValue([]byte("item.name"))
    if key != "item" {
        t.Errorf("expect item get %s", key)
    }
    if _, ok := value.(string); !ok {
        t.Errorf("expect value type is string")
    }
    if !reflect.DeepEqual(value, "name") {
        t.Errorf("expect true, get false")
    }
    key, value = splitKeyValue([]byte("item.0"))
    if key != "item" {
        t.Errorf("expect item get %s", key)
    }
    if _, ok := value.(int); !ok {
        t.Errorf("expect value type is int")
    }
    if !reflect.DeepEqual(value, 0) {
        t.Errorf("expect true, get false")
    }
    key, value = splitKeyValue([]byte("item.9"))
    if key != "item" {
        t.Errorf("expect item get %s", key)
    }
    if _, ok := value.(int); !ok {
        t.Errorf("expect value type is int")
    }
    if !reflect.DeepEqual(value, 9) {
        t.Errorf("expect true, get false")
    }
}

func TestLeftLike(t *testing.T) {
    v := LeftLike(1)
    if v != "%1" {
        t.Errorf("expect '%%1' get %s", v)
    }
}

func TestRightLike(t *testing.T) {
    v := RightLike(1)
    if v != "1%" {
        t.Errorf("expect '1%%' get %s", v)
    }
}

func TestLike(t *testing.T) {
    v := Like(1)
    if v != "%1%" {
        t.Errorf("expect '%%1%%' get %s", v)
    }
}
