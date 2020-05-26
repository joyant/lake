package main

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "testing"
)

var testSession Session

func init()  {
    var err error
    testSession, err = NewSession(MySQLDrive, "panda_dev:Panda_dev@tcp(rm-2zer9e7ki71n903rrvo.mysql.rds.aliyuncs.com:3306)/erp_dev?charset=utf8&parseTime=True&loc=Local", nil)
    if err != nil {
       panic(err)
    }
}

func TestSession(t *testing.T) {
    results, err := testSession.Select("user.getByLimit", Parameter{"limit":3})
    if err != nil {
        t.Errorf("select err:%v", err)
    }
    fmt.Println(results)
}

func TestSession2(t *testing.T) {
    results, err := testSession.Select("user.selectInID", Parameter{"ids":[]int{1, 2, 3}})
    if err != nil {
       t.Errorf("select err:%v", err)
    }
    fmt.Println(results)
}

func TestSessionOne(t *testing.T) {
    result, err := testSession.SelectOne("user.selectInID", Parameter{"ids":[]int{1, 2, 3}})
    if err != nil {
        t.Errorf("select err:%v", err)
    }
    fmt.Println(result)
}

func TestSession3(t *testing.T)  {
    p := Parameter{
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
    results, err := testSession.Select("opn.search", p)
    if err != nil {
        t.Errorf("select err:%v", err)
    }
    fmt.Println(results)
}

func TestSession4(t *testing.T)  {
    p := Parameter{
        "users": []map[string]interface{}{
            {"name":"a", "age":1},
            {"name":"b", "age":2},
            {"name":"c", "age":3},
        },
    }
    act, err := testSession.Insert("user.ifBatchSave", p)
    fmt.Println(act, err)
}

func TestUpdateSession(t *testing.T)  {
    p := Parameter{
        "id":3,
        "name":"aa",
    }
    rows, err := testSession.Update("user.updateNameById", p)
    fmt.Println(rows, err)
}

func TestUpdate2Session(t *testing.T)  {
    p := Parameter{
        "id":3,
        "name":"aaaa",
        "eq":true,
    }
    rows, err := testSession.Update("user.updateNameById2", p)
    fmt.Println(rows, err)
}

func TestDeleteSession(t *testing.T)  {
    p := Parameter{
        "id":3,
        "name":"aaaa", //  param name is redundant, but not effect execution
    }
    rows, err := testSession.Delete("user.deleteById", p)
    fmt.Println(rows, err)
}

func TestDelete2Session(t *testing.T)  {
    p := Parameter{
        "like_name":Like("aaaa"),
    }
    rows, err := testSession.Delete("user.like_name", p)
    fmt.Println(rows, err)
}

func TestSelectInclude(t *testing.T)  {
    p := Parameter{
        "name": Like("m"),
    }
    result, err := testSession.SelectOne("user.selectCount", p)
    if err != nil {
        t.Errorf("not expect err:%v", err)
    }
    if result["count"].(int64) != 2 {
        t.Errorf("expect 2 get %v", result)
    }
}