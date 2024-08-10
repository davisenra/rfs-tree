package tree

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestGenerateTree(t *testing.T) {
	root := "./../../dummy"
	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}

	tree, err := GenerateTree(root)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path     string
		name     string
		isDir    bool
		hasChild bool
	}{
		{path: absRoot, name: "dummy", isDir: true, hasChild: true},
		{path: filepath.Join(absRoot, "fixtures"), name: "fixtures", isDir: true, hasChild: true},
		{path: filepath.Join(absRoot, "fixtures", "file.txt"), name: "file.txt", isDir: false, hasChild: false},
		{path: filepath.Join(absRoot, "fixtures", "deeper_fixtures"), name: "deeper_fixtures", isDir: true, hasChild: true},
		{path: filepath.Join(absRoot, "fixtures", "deeper_fixtures", "file2.txt"), name: "file2.txt", isDir: false, hasChild: false},
		{path: filepath.Join(absRoot, "fixtures", "deeper_fixtures", "file3.txt"), name: "file3.txt", isDir: false, hasChild: false},
	}

	for _, tt := range tests {
		node := findNodeByPath(tree, tt.path)

		if node == nil {
			t.Fatalf("Expected to find node for path %s, but did not", tt.path)
		}

		if node.Name != tt.name {
			t.Errorf("Expected node name %s, but got %s", tt.name, node.Name)
		}

		if node.IsDir != tt.isDir {
			t.Errorf("Expected IsDir to be %v, but got %v", tt.isDir, node.IsDir)
		}

		if (len(node.Children) > 0) != tt.hasChild {
			t.Errorf("Expected hasChild to be %v, but got %v", tt.hasChild, len(node.Children) > 0)
		}
	}
}

func TestOutputTree(t *testing.T) {
	var outputTrap bytes.Buffer

	mockTree := TreeNode{
		Name:  "fake_dir",
		Path:  "/var/www/fake_dir",
		IsDir: true,
		Children: []*TreeNode{
			{
				Name:  "nested_dir",
				Path:  "/var/www/fake_dir/nested_dir",
				IsDir: true,
				Children: []*TreeNode{
					{
						Name:     "deeply_nested_file.txt",
						Path:     "/var/www/fake_dir/nested_dir/deeply_nested_file.txt",
						IsDir:    true,
						Children: nil,
					},
					{
						Name:     "deeply_nested_file2.txt",
						Path:     "/var/www/fake_dir/nested_dir/eeply_nested_file2.txt",
						IsDir:    true,
						Children: nil,
					},
				},
			},
			{
				Name:     "nested_file.txt",
				Path:     "/var/www/fake_dir/nested_file.txt",
				IsDir:    false,
				Children: nil,
			},
		},
	}

	err := OutputTree(&mockTree, &outputTrap)
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "fake_dir\n" +
		"├── nested_dir\n" +
		"│   ├── deeply_nested_file.txt\n" +
		"│   └── deeply_nested_file2.txt\n" +
		"└── nested_file.txt\n"

	actualOutput := outputTrap.String()

	if expectedOutput != actualOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, actualOutput)
	}
}

func findNodeByPath(node *TreeNode, searchPath string) *TreeNode {
	if filepath.Clean(node.Path) == filepath.Clean(searchPath) {
		return node
	}

	for _, child := range node.Children {
		result := findNodeByPath(child, searchPath)
		if result != nil {
			return result
		}
	}

	return nil
}
