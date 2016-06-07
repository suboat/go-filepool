package utils

import (
	"bufio"
	"io"
	"os"
)

type Reader struct {
	Path string // 文件路径
	file *os.File
	rd   *bufio.Reader
	EOF  error
}

// 打开文件
func (reader *Reader) Open() (err error) {
	reader.file, err = os.Open(reader.Path)
	if err == nil {
		reader.rd = bufio.NewReader(reader.file)
	}
	return
}

// 关闭文件
func (reader *Reader) Close() (err error) {
	err = reader.file.Close()
	return
}

// 初始化按操作
func (reader *Reader) init() (err error) {
	err = reader.Open()
	return
}

// 读行，每次返回完整的一行
func (reader *Reader) ReadOneline() (line string, err error) {
	line, err = reader.rd.ReadString('\n')

	if len(line) > 0 {
		if line[len(line)-1] == '\n' {
			drop := 1
			if len(line) > 1 && line[len(line)-2] == '\r' {
				drop = 2
			}
			line = line[:len(line)-drop]
		}
	}

	return line, err
}

func NewReader(path string) (*Reader, error) {
	result := &Reader{
		Path: path,
		EOF:  io.EOF,
	}
	err := result.init()
	return result, err
}
