package lake

import "testing"

func TestEvalLen(t *testing.T)  {
    v := evalLen("a", "1", Parameter{"a": []int{1, 2, 3}})
    if v.(int) != 3 {
        t.Errorf("expect 3 get %v", v)
    }
}

func TestEvalLen2(t *testing.T)  {
    v := evalLen("a", "1", Parameter{"b": []int{1, 2, 3}})
    if v.(int) != 0 {
        t.Errorf("expect 0 get %v", v)
    }
}

func TestEvalLen3(t *testing.T)  {
    v := evalLen("a", "1", Parameter{"a": 12345})
    if v.(int) != 5 {
        t.Errorf("expect 5 get %v", v)
    }
}

func TestEvalLen4(t *testing.T)  {
    v := evalLen("a", "1", Parameter{"a": "abcdef"})
    if v.(int) != 6 {
        t.Errorf("expect 6 get %v", v)
    }
}
