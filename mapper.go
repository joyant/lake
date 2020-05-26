package lake

import (
    "bytes"
    "encoding/xml"
    "io/ioutil"
    "strings"
)

const (
    errNeedNamespaceAndId = "need namespace and id"
    errMapperNotExist = "mapper %s not exist"
    errAttrNotExist = "file %s attr %s not exist"
)

type mapper map[string][]stack

var ms = make(mapper)

// find seek a stack, s format is like blog.getById, blog is namespace, getById is ID
// find always obtain a copy stack
func (m mapper)find(s string) (sk stack, err error) {
    slice := strings.Split(s, ".")
    if len(slice) != 2 {
        err = errorF(errNeedNamespaceAndId)
        return
    }
    if sta, ok := m.get(slice[0], slice[1]); ok {
        sk = *sta
        return
    }
    err = errorF(errMapperNotExist, s)
    return
}

func (m mapper)get(namespace, id string) (*stack, bool) {
    if _, ok := m[namespace]; !ok {
        return nil, false
    }
    for k, v := range m[namespace] {
        if v.id == id {
            return &m[namespace][k], true
        }
    }
    return nil, false
}

func (m mapper)set(namespace string, sta stack)  {
    if slice, ok := m[namespace]; ok {
        m[namespace] = append(slice, sta)
    } else {
        m[namespace] = []stack{sta}
    }
}

func loadXML(filename string) (err error) {
    var b []byte
    b, err = ioutil.ReadFile(filename)
    if err != nil {
        return
    }
    var (
        t xml.Token
        decoder = xml.NewDecoder(bytes.NewReader(b))
    )
    var namespace string
    var sta stack
    for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
        switch token := t.(type) {
        case xml.StartElement:
            if token.Name.Local == "mapper" {
                if len(token.Attr) == 0 || token.Attr[0].Name.Local != "namespace" {
                    return errorF(errAttrNotExist, filename, "namespace")
                }
                namespace = token.Attr[0].Value
                continue
            }
            switch token.Name.Local {
            case insertNode, updateNode, selectNode, deleteNode:
                if len(token.Attr) == 0 || token.Attr[0].Name.Local != "id" {
                    return errorF(errAttrNotExist, filename, "id")
                }
                sta = stack{
                    id:        token.Attr[0].Value,
                    action:    token.Name.Local,
                    namespace: namespace,
                }
            }
            n := node{
                kind:  nodeTypeTag,
                name:  token.Name.Local,
                attr:  makeAttr(token.Attr),
                start: true,
            }
            sta.push(n)
        case xml.EndElement:
            if token.Name.Local == "mapper" {
                continue
            }
            n := node{
                kind: nodeTypeTag,
                name: token.Name.Local,
            }
            sta.push(n)
            switch token.Name.Local {
            case insertNode, updateNode, selectNode, deleteNode:
                ms.set(namespace, sta)
            }
        case xml.CharData:
            b := bytes.TrimSpace(token.Copy())
            if b != nil {
                n := node{
                    kind:    nodeTypeContent,
                    content: token.Copy(),
                }
                sta.push(n)
            }
        }
    }
    return nil
}

func makeAttr(attrs []xml.Attr) map[string]string {
    m := make(map[string]string)
    for _, v := range attrs {
        m[v.Name.Local] = v.Value
    }
    return m
}