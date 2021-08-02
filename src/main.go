package main

import (
    "fmt"
    "reflect"
    "strconv"

    "pathtree"
)


func found(n *pathtree.Node, path string, val interface{}) bool {
    leafs := n.Find(path)
    if leafs == nil {
        fmt.Printf("Didn't find: %s", path)
    }

    for _, leaf := range leafs {
        if leaf.Value == val {
            return true
        }
    }

    for _, leaf := range leafs {
        if leaf.Value != val {
            fmt.Printf("%s: value (actual) %s != %s (expected)\n", path, leaf.Value, val)
        }
    }
    return false
}

func notFound(n *pathtree.Node, path string, val interface{}) bool {
    leafs := n.Find(path)
    if leafs != nil {
        for _, leaf := range leafs {
            if leaf.Value == val {
                fmt.Printf("Should not have found: %s \n", path)
                return false
            }
        }
    }
    return true
}

func foundPath(n *pathtree.Node, path string) bool {
    node := n.FindPath(path)
    if node == nil {
        fmt.Printf("%s not found != found (expected) \n", path)
        return false
    }
    return true
}

func notFoundPath(n *pathtree.Node, path string) bool {
    node := n.FindPath(path)
    if node != nil {
        fmt.Printf("Should not have found path: %s \n", path)
        return false
    }
    return true
}

func Display(name string, x interface{}) {
    fmt.Printf("Display %s (%T):\n", name, x)
    display(name, reflect.ValueOf(x))
}

func formatAtom(v reflect.Value) string {
    switch v.Kind() {
    case reflect.Invalid:
        return "invalid"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
        reflect.Int64:
        return strconv.FormatInt(v.Int(), 10)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
        reflect.Uint64:
        return strconv.FormatUint(v.Uint(), 10)
    case reflect.Bool:
        return strconv.FormatBool(v.Bool())
    case reflect.String:
        return strconv.Quote(v.String())
    case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
        return v.Type().String() + "0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
    default:
        return v.Type().String() + "value"
    }
}

func display(path string, v reflect.Value) {
    switch v.Kind() {
    case reflect.Invalid:
        fmt.Printf("%s = invalid\n", path)
    case reflect.Slice, reflect.Array:
        for i := 0; i < v.Len(); i++ {
            display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
        }
    case reflect.Struct:
        for i := 0; i < v.NumField(); i++ {
            fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
            display(fieldPath, v.Field(i))
        }
    case reflect.Map:
        for _, key := range v.MapKeys(){
            display(fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key))
        }
    case reflect.Ptr:
        if v.IsNil() {
            fmt.Printf("%s = nil\n", path)
        } else {
            display(fmt.Sprintf("*%s", path), v.Elem())
        }
    case reflect.Interface:
        if v.IsNil() {
            fmt.Printf("%s = nil\n", path)
        } else {
            //fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
            display(path+".value", v.Elem())
        }
    default:
        fmt.Printf("%s = %s\n", path, formatAtom(v))
    }
}


func Testing() {
    fmt.Print("\n============ TestColon ============\n")
    t := pathtree.Tree()
    extraData := make(map[string]interface{})
    extraData["test"] = "just a test"
    _ = t.Add("/", "root", extraData)
    _ = t.Add("/", "root1", nil)

    _ = t.Add("/a", nil, nil)
    _ = t.Add("/a", "a", nil)
    _ = t.Add("/a", "aa", nil)
    _ = t.Add("/a", "aaa", nil)

    _ = t.Add("/a/b", "b", nil)
    _ = t.Add("/a/b", "bb", nil)
    _ = t.Add("/a/b", "bbb", nil)

    _ = t.Add("/a/b/c", "c", nil)
    _ = t.Add("/a/b/c", "cc", nil)
    _ = t.Add("/a/b/c", "ccc", nil)

    _ = t.Add("/a/d", "d", nil)
    _ = t.Add("/a/d", "dd", nil)
    _ = t.Add("/a/d", "ddd", nil)

    _ = t.Add("/a/d/e", "e", nil)
    _ = t.Add("/a/d/e", "ee", nil)
    _ = t.Add("/a/d/e", "eee", nil)

    _ = t.Add("/f/g/h", "h", nil)
    _ = t.Add("/f/g/h", "hh", nil)
    _ = t.Add("/f/g/h", "hhh", nil)

    extraData = make(map[string]interface{})
    extraData["test"] = "hi"
    t.SetPathExtraData("/a", extraData)

    found(t, "/", "root")
    found(t, "/", "root1")
    found(t, "/a", nil)
    found(t, "/a", "a")
    found(t, "/a", "aa")
    found(t, "/a", "aaa")
    found(t, "/a/b", "b")
    found(t, "/a/b", "bb")
    found(t, "/a/b", "bbb")
    found(t, "/a/b/c", "c")
    found(t, "/a/b/c", "cc")
    found(t, "/a/b/c", "ccc")
    found(t, "/a/d", "d")
    found(t, "/a/d", "dd")
    found(t, "/a/d", "ddd")
    found(t, "/a/d/e", "e")
    found(t, "/a/d/e", "ee")
    found(t, "/a/d/e", "eee")
    found(t, "/f/g/h", "h")
    found(t, "/f/g/h", "hh")
    found(t, "/f/g/h", "hhh")

    notFound(t, "/a", "!a")
    notFound(t, "/a/e", "!a")

    t.DeleteLeaf("/", "root")
    notFound(t, "/", "root")

    t.DeleteLeaf("/a", "a")
    notFound(t, "/a", "a")
    found(t, "/a", "aa")
    found(t, "/a", "aaa")

    t.DeleteLeaf("/a/b", "b")
    notFound(t, "/a/b", "b")
    found(t, "/a/b", "bb")
    found(t, "/a/b", "bbb")

    t.DeleteLeaf("/a/b", "b")
    notFound(t, "/a/g/h", "h")

    t.DeletePath("/a")
    notFoundPath(t, "/a")
    notFoundPath(t, "/a/b")
    notFoundPath(t, "/a/b/c")
    notFoundPath(t, "/a/d")
    notFoundPath(t, "/a/d/e")

    t.DeletePath("/f/g/h")

    notFoundPath(t, "/f/g/h")
    foundPath(t, "/f")
    foundPath(t, "/f/g")

    t.DeletePath("/f")
    notFoundPath(t, "/f")
    notFoundPath(t, "/f/g")
    notFoundPath(t, "/f/g/h")

}

func main() {
    Testing()
}