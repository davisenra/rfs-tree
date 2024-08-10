package tree

import (
	"io"
	"os"
	"path/filepath"
)

type TreeNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*TreeNode
}

func GenerateTree(root string) (*TreeNode, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	node := TreeNode{
		Name:     filepath.Base(root),
		Path:     root,
		IsDir:    info.IsDir(),
		Children: []*TreeNode{},
	}

	paths, err := filepath.Glob(root + "/*")
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		info, err = os.Lstat(path)
		if err != nil {
			return nil, err
		}

		symLink, err := isSymlink(path)
		if err != nil {
			return nil, err
		}

		if symLink {
			return &node, nil
		}

		if info.IsDir() {
			pathChildren, err := GenerateTree(path)
			if err != nil {
				return nil, err
			}

			node.Children = append(node.Children, pathChildren)
		} else {
			node.Children = append(node.Children, &TreeNode{
				Name:     filepath.Base(path),
				Path:     path,
				IsDir:    false,
				Children: nil,
			})
		}
	}

	return &node, nil
}

func isSymlink(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}

const VERTICAL_LINE = "│   "
const BRANCH_WITH_ITEMS = "├── "
const BRANCH_LAST_ITEM = "└── "

func OutputTree(node *TreeNode, w io.Writer) error {
	_, err := w.Write([]byte(node.Name + "\n"))
	if err != nil {
		return err
	}

	return recursivelyDescendTree(node.Children, w, "")
}

func recursivelyDescendTree(nodes []*TreeNode, w io.Writer, prefix string) error {
	for i, node := range nodes {
		branch := BRANCH_WITH_ITEMS
		if i == len(nodes)-1 {
			branch = BRANCH_LAST_ITEM
		}

		_, err := w.Write([]byte(prefix + branch + node.Name + "\n"))
		if err != nil {
			return err
		}

		newPrefix := prefix
		if i == len(nodes)-1 {
			newPrefix += "    "
		} else {
			newPrefix += VERTICAL_LINE
		}

		if len(node.Children) > 0 {
			err = recursivelyDescendTree(node.Children, w, newPrefix)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
