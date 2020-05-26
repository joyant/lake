package lake

import (
    "bytes"
    "errors"
    "fmt"
    "reflect"
)

type (
    Parameter map[string]interface{}
    nodeType  int
)

const (
    insertNode = "insert"
    selectNode = "select"
    updateNode = "update"
    deleteNode = "delete"

    ifNode        = "if"
    whereNode     = "where"
    whenNode      = "when"
    chooseNode    = "choose"
    otherwiseNode = "otherwise"
    foreachNode   = "foreach"
    setNode       = "set"
    includeNode   = "include"

    testAttr  = "test"
    refidAttr = "refid"

    nodeTypeTag     = 1
    nodeTypeContent = 2
)

var (
    errUnknownNode = "unknown node %s"
    errNeedTestAttr = "need test attr near %s"
    errNeedRefidAttr = "need refid attr near %s"
    errUnclosedNode = "unclosed node %s"
    errNeedTagNode = "need tag node near %s"
    errNeedCollectionAttr = "need collection attr near %s"
    errNeedSlice = "need slice near %s"
)

type node struct {
    kind    nodeType
    name    string
    start   bool
    attr    map[string]string
    content []byte
}

type stack struct {
    namespace string
    id        string
    nodes     []node
    action    string // select update insert delete
    params    Parameter
}

func (s *stack)push(n node)  {
    s.nodes = append(s.nodes, n)
}

func (s *stack)len() int {
    return len(s.nodes)
}

func (s *stack)Parse() (stmt []byte, err error) {
    return s.parse(0, s.len()-1)
}

func (s *stack)parse(from, to int) (sql []byte, err error) {
    buff := bytes.Buffer{}
    for i := from; i <= to; {
        v := s.nodes[i]
        if v.kind == nodeTypeContent {
            buff.Write(v.content)
            i ++
            continue
        }
        closeIndex := s.seekCloseNode(i+1, v.name)
        if closeIndex == 0 {
            err = errorF(errUnclosedNode, v.name)
            return
        }
        var bs []byte
        switch v.name {
        case ifNode:
            bs, err = s.parseIf(i, closeIndex)
        case whereNode:
            bs, err = s.parseWhere(i, closeIndex)
        case whenNode:
           bs, err = s.parseWhen(i, closeIndex)
        case chooseNode:
            bs, err = s.parseChoose(i, closeIndex)
        case foreachNode:
            bs, err = s.parseForeach(i, closeIndex)
        case setNode:
            bs, err = s.parseSet(i, closeIndex)
        case otherwiseNode:
            bs, err = s.parseOtherwise(i, closeIndex)
        case includeNode:
            bs, err = s.parseInclude(i, closeIndex)
        case insertNode, updateNode, selectNode, deleteNode:
            s.action = v.name
            bs, err = s.parse(i+1, closeIndex-1)
        default:
            err = errorF(errUnknownNode, v.name)
        }
        if err != nil {
            return
        }
        if bs != nil {
            buff.WriteByte(' ')
            buff.Write(bs)
        }
        i = closeIndex + 1
    }
    return buff.Bytes(), nil
}

func (s *stack)seekCloseNode(i int, name string) (index int) {
   count := 0
   for i < s.len() {
       v := s.nodes[i]
       if v.kind == nodeTypeTag && v.name == name {
           if v.start {
               count ++
           } else if count > 0{
               count --
           } else {
               return i
           }
       }
       i ++
   }
   return 0
}

func (s *stack)parseIf(from, to int) (sql []byte, err error) {
    v := s.nodes[from]
    if len(v.attr) == 0 || v.attr[testAttr] == "" {
        err = errorF(errNeedTestAttr, v.name)
        return
    }
    if evalTest(v.attr[testAttr], s.params) {
        return s.parse(from+1, to-1)
    }
    return nil, nil
}

func (s *stack)parseWhere(from, to int) (sql []byte, err error)  {
    b, err := s.parse(from+1, to-1)
    if err != nil {
        return nil, err
    }
    b = bytes.TrimSpace(b)
    // after parsing content inside where, see whether content start with and(AND)
    // erase and(AND) if the answer is yes
    if b != nil {
        if bytes.HasPrefix(b, []byte("AND ")) {
            b = bytes.Replace(b, []byte("AND"), nil, 1)
        } else if bytes.HasPrefix(b, []byte("and ")) {
            b = bytes.Replace(b, []byte("and"), nil, 1)
        }
        bb := make([]byte, len(b)+6)
        copy(bb, "where ")
        copy(bb[6:], b)
        b = bb
    }
    return b, nil
}

func (s *stack)parseChoose(from, to int) (sql []byte, err error) {
    var b []byte
    for i := from+1; i < to; {
        v := s.nodes[i]
        if v.kind != nodeTypeTag {
            err = errorF(errNeedTagNode, string(v.content))
            return
        }
        if v.name == whenNode {
            closeIndex := s.seekCloseNode(i+1, v.name)
            if closeIndex == 0 {
                err = errorF(errUnclosedNode, v.name)
                return
            }
            b, err = s.parseWhen(i, closeIndex)
            if err != nil {
                break
            }
            if b == nil {
                i = closeIndex+1
                continue
            }
            break
        } else if v.name == otherwiseNode {
            closeIndex := s.seekCloseNode(i+1, v.name)
            if closeIndex == 0 {
                err = errorF(errUnclosedNode, v.name)
                return
            }
            b, err = s.parseOtherwise(i, closeIndex)
            break
        } else {
            err = errorF(errUnknownNode, v.name)
            break
        }
    }
    return b, err
}

func (s *stack)parseSet(from, to int) (sql []byte, err error) {
    b, err := s.parse(from+1, to-1)
    if err != nil {
        return nil, err
    }
    b = bytes.TrimSpace(b)
    if b != nil && b[len(b)-1] == ',' {
        b = b[:len(b)-1]
    }
    // if content end with comma, erase it, and add set at begin
    if b != nil {
        bb := make([]byte, len(b)+4)
        copy(bb, "set ")
        copy(bb[4:], b)
        b = bb
    }
    return b, nil
}

// parse when is similar with parsing if, but when has short circuit effect, once pass a when, will stop parsing
// this will be charge of func parseChoose, parseWhen should not care about it.
func (s *stack)parseWhen(from, to int) (sql []byte, err error) {
    v := s.nodes[from]
    if len(v.attr) == 0 || v.attr[testAttr] == "" {
        err = errorF(errNeedTestAttr, v.name)
        return
    }
    if evalTest(v.attr[testAttr], s.params) {
        return s.parse(from+1, to-1)
    }
    return nil, nil
}

func (s *stack)parseForeach(from, to int) (sql []byte, err error) {
    a := s.nodes[from]
    attrMap := make(defaultMap)
    for name, value := range a.attr {
        attrMap.set(name, value)
    }
    co := attrMap["collection"]
    if co == "" {
        err = errorF(errNeedCollectionAttr, a.name)
        return
    }
    // return back if collection is nil, not exist is not a error
    if _, ok := s.params[co]; !ok {
        return nil, nil
    }

    likeSlice := s.params[co] // make be it is a slice, so named likeSlice
    buff := bytes.Buffer{}
    inner := s.nodes[from+1].content
    item := []byte(fmt.Sprintf("#{%s}", attrMap.get("item", "item")))
    if reflect.TypeOf(likeSlice).Kind().String() != "slice" {
        return nil, errorF(errNeedSlice, a.name)
    }
    sliceLen := reflect.ValueOf(likeSlice).Len()
    if sliceLen == 0 {
        return nil, nil
    }
    if hasMap(reflect.ValueOf(likeSlice)) {
        keys := findPlaceholder(inner)
        for i, l := 0, sliceLen; i < l; i++ {
            bs := make([]byte, len(inner))
            copy(bs, inner)
            buff.WriteString(attrMap.get("open", "("))
            for _, keyValue := range keys {
                _, value := splitKeyValue([]byte(keyValue))
                if value == nil {
                    return nil, errors.New("need format like item.key")
                }
                bs = bytes.Replace(bs, []byte("#{"+keyValue+"}"), []byte(fmt.Sprintf("#{%s.%d.%s}", co, i, value)), -1)
            }
            buff.Write(bs)
            buff.WriteString(attrMap.get("close", ")"))
            if i < l-1 {
                buff.WriteString(attrMap.get("separator", ","))
            }
        }
    } else {
        buff.WriteString(attrMap.get("open", "("))
        for i, l := 0, sliceLen; i < l; i++ {
            b := bytes.Replace(inner, item, []byte(fmt.Sprintf("#{%s.%d}", co, i)), -1)
            buff.Write(b)
            if i < l-1 {
                buff.WriteString(attrMap.get("separator", ","))
            }
        }
        buff.WriteString(attrMap.get("close", ")"))
    }
    return buff.Bytes(), nil
}

func (s *stack)parseOtherwise(from, to int) (sql []byte, err error) {
    return s.parse(from+1, to-1)
}

func (s *stack)parseInclude(from, to int) (sql []byte, err error) {
    v := s.nodes[from]
    if len(v.attr) == 0 || v.attr[refidAttr] == "" {
        err = errorF(errNeedRefidAttr, v.name)
        return
    }
    refStack, ok := ms.get(s.namespace, v.attr[refidAttr])
    if !ok {
        return nil, errorF("seek include, ref id %s not exist", v.attr[refidAttr])
    }
    ref := *refStack
    ref.params = s.params
    refBytes, err := ref.Parse()
    if err != nil {
        return nil, err
    }
    return refBytes, nil
}

func findPlaceholder(b []byte) (keys []string) {
    for i, l := 0, len(b); i < l; {
        if b[i] == '#' && i < l && b[i+1] == '{' {
            for j := i + 2; j < l; j++ {
                if b[j] == '}' {
                    keys = append(keys, string(b[i+2:j]))
                    i = j + 1
                    break
                }
            }
        } else {
            i ++
        }
    }
    return
}

func hasMap(value reflect.Value) bool {
    switch value.Index(0).Kind().String() {
    case "map":
        return true
    case "interface":
        return value.Index(0).Elem().Kind().String() == "map"
    default:
        return false
    }
}