package pathtree

import (
    "errors"
    "strings"
)

type Node struct {
    Edges     map[string]*Node
    ExtraData interface{}
    Leafs   []*Leaf
    LeafCnt int
}

type Leaf struct {
    Value     interface{}
    ExtraData interface{}
    Order     int
}


func Tree() *Node {
    return &Node{Edges: make(map[string]*Node), ExtraData: nil, Leafs: []*Leaf{}, LeafCnt: 0}
}


func (t *Node) Add(path string, val interface{}, extraData interface{}) error {
    if path == "" || path[0] != '/' {
        return errors.New("path must begin with /")
    }
    t.LeafCnt++
    return t.add(t.LeafCnt, splitPath(path), val, extraData)
}


func (t *Node) addLeaf(leaf *Leaf) error {
    if t.Leafs == nil {
        t.Leafs = []*Leaf{}
    }
    t.Leafs = append(t.Leafs, leaf)
    return nil
}

func (t *Node) add(order int, elements []string, value, extraData interface{}) error {
    if len(elements) == 0 {
        leaf := &Leaf{
            Value:     value,
            ExtraData: extraData,
            Order:     order,
        }
        return t.addLeaf(leaf)
    }

    var el string
    el, elements = elements[0], elements[1:]
    if el == "" {
        return errors.New("empty path elements are not allowed")
    }

    e, ok := t.Edges[el]
    if !ok {
        e = Tree()
        t.Edges[el] = e
    }

    return e.add(order, elements, value, extraData)
}


func (t *Node) Find(path string) []*Leaf {
    if len(path) == 0 || path[0] != '/' {
        return nil
    }
    return t.find(splitPath(path))
}

func (t *Node) find(elements []string) []*Leaf {
    if len(elements) == 0 {
        return t.Leafs
    }

    var el string
    el, elements = elements[0], elements[1:]

    var leafs []*Leaf
    if nextNode, ok := t.Edges[el]; ok {
        if nextNode != nil {
            leafs = nextNode.find(elements)
        }
    }
    return leafs
}

func (t *Node) FindPath(path string) *Node {
    if len(path) == 0 || path[0] != '/' {
        return nil
    }
    return t.findPath(splitPath(path))
}

func (t *Node) findPath(elements []string) *Node {
    if len(elements) == 0 {
        return nil
    }
    if len(elements) == 1 && t.Edges != nil {
        if targetNode, ok := t.Edges[elements[0]]; ok {
            return targetNode
        }
    }

    var el string
    el, elements = elements[0], elements[1:]

    var targetNode *Node
    if nextNode, ok := t.Edges[el]; ok {
        if nextNode != nil {
            targetNode = nextNode.findPath(elements)
        }
    }
    return targetNode
}

func (t *Node) DeleteLeaf(path string, value interface{}) ([]*Leaf, *Node) {
    if len(path) == 0 || path[0] != '/' {
        return nil, nil
    }
    leafs, targetNode := t.deleteLeaf(splitPath(path), value)
    return leafs, targetNode
}

func (t *Node) deleteLeaf(elements []string, value interface{}) ([]*Leaf, *Node) {
    if len(elements) == 0 {
        leaf := matchLeaf(t.Leafs, value)
        if leaf != nil {
            popLeaf(t.Leafs, leaf)
        }
        return t.Leafs, t
    }

    var el string
    el, elements = elements[0], elements[1:]

    var targetNode *Node
    var leafs []*Leaf
    if nextNode, ok := t.Edges[el]; ok {
        if nextNode != nil {
            leafs, targetNode = nextNode.deleteLeaf(elements, value)
            matched := matchLeaf(leafs, value)
            if matched != nil {
                targetNode = nextNode
                popLeaf(nextNode.Leafs, matched)
                popLeaf(leafs, matched)
                return leafs, targetNode
            }
        }
    }
    return leafs, targetNode
}


func (t *Node) DeletePath(path string) (map[string]*Node, *Node) {
    if len(path) == 0 || path[0] != '/' {
        return nil, nil
    }
    targetNode, fatherNode := t.deletePath(splitPath(path))
    return targetNode, fatherNode
}

func (t *Node) deletePath(elements []string)(map[string]*Node, *Node) {
    fatherNode := t
    if len(elements) == 1 && t.Edges != nil {
        if _, ok := t.Edges[elements[0]]; ok {
            delete(t.Edges, elements[0])
        }
        return fatherNode.Edges, fatherNode
    }

    var el string
    el, elements = elements[0], elements[1:]

    var targetNode map[string]*Node
    if nextNode, ok := t.Edges[el]; ok {
        if nextNode != nil {
            targetNode, fatherNode = nextNode.deletePath(elements)
        }
    }
    return targetNode, fatherNode
}

func (t *Node) SetPathExtraData(path string, extraData interface{}) bool {
    node := t.FindPath(path)
    if node != nil {
        node.ExtraData = extraData
        return true
    }
    return false
}

func matchLeaf(leafs []*Leaf, value interface{}) *Leaf {
    for i := 0; i < len(leafs); i++ {
        leaf := leafs[i]
        if leaf !=nil && leaf.Value == value {
            return leaf
        }
    }
    return nil
}

func popLeaf(leafs []*Leaf, targetLeaf *Leaf) {
    if leafs == nil || targetLeaf == nil{
        return
    }
    for i := 0; i < len(leafs); i++ {
        if leafs[i] == targetLeaf {
            leafs = append(leafs[:i], leafs[i+1:]...)
        }
    }
}

func splitPath(path string) []string {
    elements := strings.Split(path, "/")
    if elements[0] == "" {
        elements = elements[1:]
    }
    if elements[len(elements)-1] == "" {
        elements = elements[:len(elements)-1]
    }
    return elements
}

