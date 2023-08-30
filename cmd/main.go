package main

import (
	fm "golang-transpiler/pkg/filemanager"
	tr "golang-transpiler/pkg/transpiler"
)

func main() {
	var manager = fm.FileManager{}
	var source = manager.ReadFile("../source.txt")
	output := tr.Transpile(source)
	manager.WriteFile(output, "../output.md")
}
