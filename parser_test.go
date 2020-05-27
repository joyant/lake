package lake

import (
    "fmt"
    "reflect"
    "testing"
)

func init()  {
    err := Load("test/mappers.xml")
    if err != nil {
        panic(err)
    }
}

func TestSelectIf(t *testing.T)  {
    sk, err := ms.find("user.getByLimit")
    if err != nil {
        t.Error("find err", err.Error())
        return
    }
    sk.params = Parameter{
        "limit":10,
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Error("parse err:", err.Error())
        return
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Error("build err:", err)
        return
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestSelectIn(t *testing.T)  {
    sk, err := ms.find("user.selectInID")
    if err != nil {
        t.Errorf("find user.selectInId err:%v", err)
    }
    sk.params = Parameter{
        "ids":[]int{1, 2, 3},
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestSelectChoose(t *testing.T)  {
    sk, err := ms.find("user.selectByNameOrAge")
    if err != nil {
        t.Error("find err", err.Error())
        return
    }
    sk.params = Parameter{
        "minID":0,
        "age":12,
        "name":"jim",
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Error("parse err:", err)
        return
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Error("build err:", err)
        return
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestSelectWhen(t *testing.T)  {
    sk, err := ms.find("user.getById")
    if err != nil {
        t.Error("find err", err.Error())
        return
    }
    sk.params = Parameter{
        "age":11,
        "name":"tom",
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Error("parse err:", err.Error())
        return
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Error("build err:", err)
        return
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestSelectIn2(t *testing.T)  {
    sk, err := ms.find("user.selectInID")
    if err != nil {
        t.Errorf("find user.selectInId err:%v", err)
    }
    sk.params = Parameter{
        "ids":[]interface{}{"1", "2", "3"},
        //"ids":[]interface{}{map[string]interface{}{"a":1}},
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    //fmt.Println("stmt", string(stmt))
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestSelectIn3(t *testing.T)  {
    sk, err := ms.find("user.selectInID")
    if err != nil {
        t.Errorf("find user.selectInId err:%v", err)
    }
    sk.params = Parameter{
        "ids":[]interface{}{"1", "2", "3"},
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestSelectComplex(t *testing.T)  {
    sk, err := ms.find("user.selectByMultipleCondition")
    if err != nil {
        t.Errorf("find user.selectByMultipleCondition err:%v", err)
    }
    sk.params = Parameter{
        "table":"user",
        "ids":[]string{"1", "2", "3"},
        "age":11,
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestInsert1(t *testing.T)  {
    sk, err := ms.find("user.saveUser")
    if err != nil {
        t.Errorf("find user.saveUser err:%v", err)
    }
    sk.params = Parameter{
        "name":"tom",
        "age":11,
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestInsert2(t *testing.T)  {
    sk, err := ms.find("user.saveUser")
    if err != nil {
        t.Errorf("find user.saveUser err:%v", err)
    }
    sk.params = Parameter{
        "name":"tom",
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestFindPlaceholder(t *testing.T)  {
    keys := findPlaceholder([]byte("#{name}, #{age}, #{gender}"))
    if !reflect.DeepEqual(keys, []string{"name", "age", "gender"}) {
        t.Errorf("expect true get false, keys are:%v", keys)
    }
}

func TestBatchSave(t *testing.T)  {
    sk, err := ms.find("user.batchSave")
    if err != nil {
        t.Errorf("find user.batchSave err:%v", err)
    }
    sk.params = Parameter{
        "users": []map[string]interface{}{
            {"name":"a", "age": 1},
            {"name":"b", "age": 2},
            {"name":"c", "age": 3},
        },
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
       t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestBatchSave2(t *testing.T)  {
    sk, err := ms.find("user.batchSave")
    if err != nil {
        t.Errorf("find user.batchSave err:%v", err)
    }
    sk.params = Parameter{
        "users": []interface{}{
            map[string]interface{}{"name":"a", "age": 1},
            map[string]interface{}{"name":"b", "age": 2},
            map[string]interface{}{"name":"c", "age": 3},
        },
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestUpdate(t *testing.T)  {
    sk, err := ms.find("user.updateNameById")
    if err != nil {
        t.Errorf("find user.updateNameById err:%v", err)
    }
    sk.params = Parameter{
        "id":1,
        "name":"tom",
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestOpn(t *testing.T)  {
    sk, err := ms.find("opn.search")
    if err != nil {
        t.Errorf("find opn.search err:%v", err)
    }
    sk.params = Parameter{
        "pro_id":1,
        "sub_type":2,
        "auditing":3,
        "publish":4,
        "pro_ver":5,
        "disable":6,
        "collection_id":7,
        "withCollection":true,
        "keyword":Like("小汤"),
        "searchAuthor":Like("汤普森"),
        "lessonIds":[]int{100, 101, 102, 103, 104, 105, 106, 107},
        "newTypeFlagArr":[]interface{}{"newType1", "newType2", "newType3", "newType4"},
        "collectionKeyword":LeftLike("a"),
        "searchCollection":Like("b"),
        "startRow":0,
        "pageSize":20,
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestIfBatchSave(t *testing.T)  {
    sk, err := ms.find("user.ifBatchSave")
    if err != nil {
        t.Errorf("find user.ifBatchSave err:%v", err)
    }
    sk.params = Parameter{
        "users": []map[string]interface{}{
            {"name":"a", "age":1},
            {"name":"b", "age":2},
            {"name":"c", "age":3},
        },
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestSelectCount(t *testing.T)  {
    sk, err := ms.find("user.selectCount")
    if err != nil {
        t.Errorf("find user.selectCount err:%v", err)
    }
    sk.params = Parameter{
        "name": "tom",
    }
    stmt, err := sk.Parse()
    if err != nil {
        t.Errorf("parse err:%v", err)
    }
    builder := mySQLBuilder{}
    sql, params, err := builder.build(stmt, sk.params)
    if err != nil {
        t.Errorf("build err:%v", err)
    }
    fmt.Println("sql:", sql)
    fmt.Println("params:", params)
}

func TestRestrainTrim(t *testing.T)  {
    b := restrainTrim([]byte("        test \n     abd    \r       test      \t"))
    if string(b) != "test abd test" {
        t.Errorf("expect test abd test get %s", b)
    }
    b = restrainTrim([]byte("        \t\t\ntest\n\n\t\t     abd    \r       test      \t"))
    if string(b) != "test abd test" {
        t.Errorf("expect test abd test get %s", b)
    }
}
