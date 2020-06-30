package week1

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type node struct {
	name         string
	path         string
	nestingLevel int
	isDir        bool
	size         string
	content      []*node
	parent       *node
}

var (
	prefixNotLast = "├───"
	prefixSkipped = "│"
	prefixLast    = "└───"
)

func (l *node) String() string {
	var inner string
	for _, el := range l.content {
		inner += el.name
		inner += " "
	}

	return fmt.Sprintf("{name:%s; path:%s; nestingLevel:%d; size:%s; inner:[%s]; isDir:%t;}", l.name, l.path,
		l.nestingLevel, l.size, inner, l.isDir)
}

// NewLeaf ...
func NewLeaf(path string, parentNode *node, defaultNestingLevel int, withFiles bool) *node {
	path, _ = filepath.Abs(path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileDescription, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var l node
	l.isDir = fileDescription.IsDir()
	l.path = path
	l.nestingLevel = len(strings.Split(l.path, "/")) - defaultNestingLevel
	l.name = fileDescription.Name()
	l.path = path
	l.parent = parentNode

	if !l.isDir {
		if fileDescription.Size() == 0 {
			l.size = "empty"
		} else {
			l.size = strconv.FormatInt(fileDescription.Size(), 10)
		}
	} else {
		l.size = "[]"
		allEntries, _ := ioutil.ReadDir(path)

		for _, el := range allEntries {
			newPath := fmt.Sprintf("%s/%s", path, el.Name())
			tmp := NewLeaf(newPath, &l, defaultNestingLevel, withFiles)
			if !withFiles {
				if tmp.isDir {
					l.content = append(l.content, tmp)
				}
			} else {
				l.content = append(l.content, tmp)
			}
		}

		sort.Slice(l.content, func(i, j int) bool {
			return l.content[i].name < l.content[j].name
		})
	}

	return &l
}

func DirTree(out io.Writer, path string, withFiles bool) error {
	path, _ = filepath.Abs(path)
	defaultNestingLevel := len(strings.Split(path, "/"))
	l := NewLeaf(path, nil, defaultNestingLevel, withFiles)
	treeTraversal(out, l)
	printLeaf(out, l, "", withFiles)

	return nil
}

func treeTraversal(out io.Writer, l *node) {
	for _, el := range l.content {
		treeTraversal(out, el)
	}
}

func printLeaf(out io.Writer, l *node, indent string, printFiles bool) {
	var nodePrefix, newIndent string

	for _, file := range l.content {
		if l.content[len(l.content)-1].name == file.name {
			nodePrefix = prefixLast
			newIndent = indent + "\t"
		} else {
			nodePrefix = prefixNotLast
			newIndent = indent + prefixSkipped + "\t"
		}

		if file.isDir {
			fmt.Fprintf(out, "%s%s%s\n", indent, nodePrefix, file.name)
			printLeaf(out, file, newIndent, printFiles)
		} else {
			var sizeStr string
			if unicode.IsDigit([]rune(file.size)[0]) {
				sizeStr = fmt.Sprintf("%sb", file.size)
			} else {
				sizeStr = "empty"
			}
			fmt.Fprintf(out, "%s%s%s (%s)\n", indent, nodePrefix, file.name, sizeStr)
		}
	}
}