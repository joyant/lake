package lake

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "testing"
)

var testSession Session

func init() {
    var err error
    testSession, err = NewSession(MySQLDrive, "panda_pre:Panda_pre@tcp(x.com:3306)/erp_pre?charset=utf8&parseTime=True&loc=Local", nil)
    if err != nil {
        panic(err)
    }
}

func TestSession(t *testing.T) {
    results, err := testSession.Select("user.getByLimit", Parameter{"limit": 3})
    if err != nil {
        t.Errorf("select err:%v", err)
    }
    fmt.Println(results)
}

func TestSession2(t *testing.T) {
    results, err := testSession.Select("user.selectInID", Parameter{"ids": []int{1, 2, 3}})
    if err != nil {
        t.Errorf("select err:%v", err)
    }
    fmt.Println(results)
}

func TestSessionOne(t *testing.T) {
    result, err := testSession.SelectOne("user.selectInID", Parameter{"ids": []int{1, 2, 3}})
    if err != nil {
        t.Errorf("select err:%v", err)
    }
    fmt.Println(result)
}

func TestSession3(t *testing.T) {
    p := Parameter{
        "pro_id":            1,
        "sub_type":          2,
        "auditing":          3,
        "publish":           4,
        "pro_ver":           5,
        "disable":           6,
        "collection_id":     7,
        "withCollection":    true,
        "keyword":           Like("小汤"),
        "searchAuthor":      Like("汤普森"),
        "lessonIds":         []int{100, 101, 102, 103, 104, 105, 106, 107},
        "newTypeFlagArr":    []interface{}{"newType1", "newType2", "newType3", "newType4"},
        "collectionKeyword": LeftLike("a"),
        "searchCollection":  Like("b"),
        "startRow":          0,
        "pageSize":          20,
    }
    results, err := testSession.Select("opn.search", p)
    if err != nil {
        t.Errorf("select err:%v", err)
    }
    fmt.Println(results)
}

func TestSession4(t *testing.T) {
    p := Parameter{
        "users": []map[string]interface{}{
            {"name": "a", "age": 1},
            {"name": "b", "age": 2},
            {"name": "c", "age": 3},
        },
    }
    act, err := testSession.Insert("user.ifBatchSave", p)
    fmt.Println(act, err)
}

func TestUpdateSession(t *testing.T) {
    p := Parameter{
        "id":   3,
        "name": "aa",
    }
    rows, err := testSession.Update("user.updateNameById", p)
    fmt.Println(rows, err)
}

func TestUpdate2Session(t *testing.T) {
    p := Parameter{
        "id":   3,
        "name": "aaaa",
        "eq":   true,
    }
    rows, err := testSession.Update("user.updateNameById2", p)
    fmt.Println(rows, err)
}

func TestDeleteSession(t *testing.T) {
    p := Parameter{
        "id":   3,
        "name": "aaaa", //  param name is redundant, but not effect execution
    }
    rows, err := testSession.Delete("user.deleteById", p)
    fmt.Println(rows, err)
}

func TestDelete2Session(t *testing.T) {
    p := Parameter{
        "like_name": Like("aaaa"),
    }
    rows, err := testSession.Delete("user.like_name", p)
    fmt.Println(rows, err)
}

func TestSelectInclude(t *testing.T) {
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

func TestSelectApp(t *testing.T) {
    p := Parameter{
        "pro_id":         1,
        "publish":        1,
        "disable":        0,
        "parent_id":      0,
        "auditing":       0,
        "sub_type":       1,
        "pro_ver":        "5.1.0",
        "newTypeFlagArr": []string{"mmusic", "mmusicconfig"},
    }
    results, err := testSession.Select("user.getAppList", p)
    if err != nil {
        t.Error(err.Error())
    }
    for k, v := range results {
        fmt.Println(k, "=>", v)
    }
}

func TestSelectDict(t *testing.T) {
    p := Parameter{
        "type":     "opern_src_types2",
        "key_code": "d_video",
    }
    result, err := testSession.SelectOne("opn.dict", p)
    if err != nil {
        t.Error(err.Error())
    }
    fmt.Println(result)
}

func TestSelectDict2(t *testing.T) {
    p := Parameter{
        "type":     "opern_src_types2",
        "key_code": "d_video",
    }
    result, err := testSession.Select("opn.dict", p)
    if err != nil {
        t.Error(err.Error())
    }
    fmt.Println(result)
}

func TestExportDB_InsertMap(t *testing.T) {
    p := Parameter{
        "type":     "my_test",
        "key_code": "123",
    }
    a, _, err := testSession.DB().InsertTable("erp_dict", p)
    if err != nil {
        t.Error(err.Error())
    }
    if a != 1 {
        t.Error("expected 1 got", a)
    }
}

func TestExportDB_UpdateTable(t *testing.T) {
    a, err := testSession.DB().UpdateTable("a", Parameter{"id": 4}, Parameter{"name": "b"})
    if err != nil {
        t.Error(err.Error())
    }
    if a != 1 {
        t.Error("expected 1 got", a)
    }
}

func TestExportDB_NamedExec(t *testing.T) {
    result, err := testSession.DB().NamedExec("update a set name = :name where id = :id", map[string]interface{}{"name": "c", "id": 4})
    if err != nil {
        t.Error(err)
    }
    a, err := result.RowsAffected()
    if err != nil {
        t.Error(err)
    }
    if a != 1 {
        t.Error("expected 1 got ", a)
    }
}

func TestExportDB_BatchInsertTable(t *testing.T) {
    ps := []Parameter{{
        "name": "x",
        "age":  1,
    }, {
        "name": "y",
        "age":  2,
    }, {
        "name": "z",
        "age":  3,
    }}
    a, b, err := testSession.DB().BatchInsertTable("a", ps)
    if err != nil {
        t.Error(err)
    }
    if a != 3 {
        t.Error("expected 3 got", a)
    }
    _ = b // b should be max id
}

func TestExportDB_QueryMapSlice(t *testing.T) {
    records, err := testSession.DB().QueryMapSlice(`select id, crc, url
from opn_resource r
where (r.crc is null or r.crc = '')
  and r.url != ''
  and r.id not in (0, 13095)
order by id
limit ?`, 10)
    if err != nil {
        t.Error(err)
    } else {
        for _, record := range records {
            t.Logf("%+v", record)
        }
    }
}

func TestExportDB_QueryMap(t *testing.T) {
    records, err := testSession.DB().QueryMap("select * from opn_resource where id < 0 limit ?", 1)
    if err != nil {
        t.Error(err)
    } else {
        fmt.Println(records)
    }
}
