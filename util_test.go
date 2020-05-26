package lake

import "testing"

func TestDefaultMap(t *testing.T)  {
    d := defaultMap{}
    d.set("name", "tim")
    if value := d.get("name", "tom"); value != "tim" {
        t.Errorf("expect tim, get:%s", value)
    }
    if value := d.get("age", "11"); value != "11" {
        t.Errorf("expect 11, get:%s", value)
    }
}
