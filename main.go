package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileInfo 定义了一个包含文件路径和大小的结构体
type FileInfo struct {
	Path string
	Size int64
}

// BySize 实现了 sort.Interface，用于根据文件大小对 FileInfo 切片进行排序
type BySize []FileInfo

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size > a[j].Size }

func main() {
	startTime := time.Now()
	var slice1 []string
	var filePaths []string

	// 执行wmic命令获取磁盘分区信息
	cmd := exec.Command("wmic", "logicaldisk", "get", "caption")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing wmic:", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过标题行
		if strings.HasPrefix(line, "Caption") {
			continue
		}
		// 处理盘符信息
		if len(line) > 0 {
			driveLetter := strings.TrimSpace(line)
			slice1 = append(slice1, driveLetter)

		}

	}

	fmt.Println(slice1)

	root := "D:/"
	// 使用filepath.Walk遍历文件夹
	err1 := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil // 跳过这个错误，继续遍历
			}
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		if !f.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	if err1 != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err1)
	}

	sorts(filePaths)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// 将耗时转换为秒
	seconds := duration.Seconds()

	// 输出耗时
	fmt.Printf("执行耗时: %.2f 秒\n", seconds)

}

func sorts(filePaths []string) {
	// 创建一个 FileInfo 切片来存储文件和它们的大小
	var files []FileInfo

	// 遍历文件路径切片，获取每个文件的大小，并添加到 files 切片中
	for _, filePath := range filePaths {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", filePath, err)
			continue // 跳过错误
		}
		files = append(files, FileInfo{
			Path: filePath,
			Size: fileInfo.Size(),
		})
	}

	// 使用自定义的 BySize 类型对 files 切片进行排序（按大小降序）
	sort.Sort(BySize(files))

	// 输出最大的前20个文件
	fmt.Println("Top 20 largest files:")
	for i, file := range files {
		if i >= 20 {
			break // 只输出前20个
		}
		fmt.Printf("%s - %d mb\n", file.Path, (file.Size / 1024 / 1024))
	}

}
