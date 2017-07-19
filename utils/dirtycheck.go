/*
Package common 是公用的接口包.
*/
package utils

import (
	"bufio"
	"io"
	"os"

	"github.com/huayuego/wordfilter/trie"
)

var DirtyWord *trie.Trie

func InitDirtyWord(fileNames ...string) (dirtyWord *trie.Trie, err error) {
	dirtyWord = trie.NewTrie()
	for _, fileName := range fileNames {
		if lines, err := ReadDirtyFile(fileName); err == nil {
			for _, line := range lines {
				dirtyWord.Add(line)
			}
		}
	}
	return dirtyWord, nil
}

func ReadDirtyFile(fileName string) (lines []string, err error) {
	fd, err := os.Open(fileName)
	if err != nil {
		Info("", "打开文件失败, file_name: %s, error: %s\n", fileName, err)
		return nil, err
	}
	defer fd.Close()
	fileReader := bufio.NewReader(fd)
	lines = make([]string, 0)

	for {
		line, err := fileReader.ReadSlice(byte('\n'))
		if err != nil {
			if err != io.EOF {
				Info("", "读取文件失败，filename: %s, error: %s\n", fileName, err)
				return lines, err
			}
			break
		}
		word := string(line[:len(line)-1])
		if word != "" {
			lines = append(lines, word)
		}
	}
	return lines, nil
}
