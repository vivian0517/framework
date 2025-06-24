package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

func FileTreeDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var dir []string
	for _, v := range entries {
		if v.IsDir() {
			dir = append(dir, v.Name())
		}
	}
	return dir, nil
}
func CopyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("open err %v", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("create err %v", err)
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	log.Print("文件复制完成")
	return err
}

func CopyLargeFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	reader := bufio.NewReader(src)
	writer := bufio.NewWriter(dst)
	buf := make([]byte, 32*1024) // 32KB 缓冲块
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("读取错误 %v", err)
		}
		if n == 0 {
			break
		}
		if _, err := writer.Write(buf[:n]); err != nil {
			log.Fatalf("写入错误 %v", err)
		}
	}
	if err := writer.Flush(); err != nil {
		log.Fatalf("flush err %v", err)
	}
	log.Print("文件复制完成")
	return err
}

func ReadLargeFile(srcPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("逐行扫描失败: %v", err)
		return err
	}
	return err
}
func RunGoFile(filename string) {
	cmd := exec.Command("go", "run", filename, "-test.v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("运行 %s 失败: %v\n", filename, err)
		os.Exit(1)
	}
}
