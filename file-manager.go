package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"path/filepath"
)

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func OpenDir(target *tview.TreeNode, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name())).
			SetSelectable(true)

		if file.IsDir() {
			node.SetColor(tcell.ColorDarkGreen)
		}
		target.AddChild(node)
	}
}

func SelectNode(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return
	}
	children := node.GetChildren()
	if len(children) == 0 {
		path := reference.(string)
		if IsDir(path) {
			OpenDir(node, path)
		} else {
			PlayFile(OpenFile(path))
		}

	} else {
		node.SetExpanded(!node.IsExpanded())
	}
}

func ListFiles() {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	OpenDir(root, rootDir)
	tree.SetSelectedFunc(SelectNode)

	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}

}

func OpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
