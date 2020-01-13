/* Copyright (C) 2019-2020 cmj. All right reserved. */
package tree

import (
	"fmt"
	"strings"
)

var (
	nrIndent int
	chIndent string
	chNext   string /* The symbol used for the next child */
	chNode   string /* The symbol used on the end of the child */
)

func init() {
	chIndent = "|   "
	chNext = "├── "
	chNode = "└── "
}

type treeView struct {
	Data   interface{}
	Childs []*treeView
	Parent *treeView
}

func NewTreeView(in interface{}) (out *treeView) {
	out = &treeView{
		Data:   in,
		Childs: make([]*treeView, 0),
		Parent: nil,
	}
	return
}

func (t *treeView) SetChNext(in string) (out *treeView) {
	out = t
	chNext = in

	if nrIndent = len(chNext); len(chNext) > len(chNode) {
		nrIndent = len(chNode)
	}
	return
}

func (t *treeView) SetChNode(in string) (out *treeView) {
	out = t
	chNode = in

	if nrIndent = len(chNext); len(chNext) > len(chNode) {
		nrIndent = len(chNode)
	}
	return
}

func (t *treeView) Insert(in interface{}) (out *treeView) {
	out = t.InsertTree(NewTreeView(in))
	return
}

func (t *treeView) InsertTree(in *treeView) (out *treeView) {
	t.Childs = append(t.Childs, in)
	in.Parent = t
	out = in
	return
}

func (t *treeView) Level() (out int) {
	for tmp := t; tmp.Parent != nil; tmp = tmp.Parent {
		out++
	}
	return
}

func (t *treeView) String() (out string) {
	var outlist []string

	tmp := fmt.Sprintf("%s%v", t.prefix(), t.Data)
	outlist = append(outlist, tmp)
	for _, child := range t.Childs {
		/* Save the list of the nested node */
		outlist = append(outlist, child.String())
	}

	out = strings.Join(outlist, "\n")
	return
}

/* Inner usage routine */
func (t *treeView) isLastNode() (out bool) {
	/* Return the current node is the last child of the parent */
	out = (t.Parent == nil) || (t.Parent.Childs[len(t.Parent.Childs)-1] == t)
	return
}

func (t *treeView) prefix() (out string) {
	for tmp := t; tmp.Parent != nil; tmp = tmp.Parent {
		var prefix string

		if t == tmp {
			if prefix = chNext; tmp.isLastNode() {
				/* Set the prefix of the node */
				prefix = chNode
			}
		} else {
			if prefix = chIndent; tmp.isLastNode() {
				/* Set the prefix of the node */
				prefix = strings.Repeat(" ", len(chIndent))
			}
		}

		out = prefix + out
	}

	return
}
