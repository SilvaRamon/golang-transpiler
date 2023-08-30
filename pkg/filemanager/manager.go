package filemanager

import (
	"bufio"
	"os"
	"golang-transpiler/pkg/types"
)

type FileManager struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func (f *FileManager) newReader(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	f.Scanner = bufio.NewScanner(file)
}

func (f *FileManager) ReadFile(file string) []types.SourceCode {
	f.newReader(file)
	var source []types.SourceCode
	var lineNumber int = 1
	
	for f.Scanner.Scan() {
		source = append(source, types.SourceCode{Line: f.Scanner.Text(), LineNumber: lineNumber})
		lineNumber++
	}

	return source
}

func (f *FileManager) newWriter(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	f.Writer = bufio.NewWriter(file)
}

func (f *FileManager) writeLine(line string) {
	_, err := f.Writer.WriteString(line + "\n")
	if err != nil {
		panic(err)
	}
}

func (f *FileManager) WriteFile(source []string, outputPath string) {
	f.newWriter(outputPath)

	f.writeLine("```mermaid")
	f.writeLine("flowchart LR")
	for _, line := range source {
		f.writeLine(line)
	}
	f.writeLine("```")
	f.Writer.Flush()
}