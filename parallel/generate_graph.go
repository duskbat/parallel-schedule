package parallel

import (
	"fmt"
	"os"
	"strings"
)

// GenerateGraphLR 生成 markdown 流程图，从左到右
func (d *Scheduler) GenerateGraphLR(fileName string) error {
	return d.GenerateGraph(fileName, "LR")
}

// GenerateGraphTB 生成 markdown 流程图，从上到下
func (d *Scheduler) GenerateGraphTB(fileName string) error {
	return d.GenerateGraph(fileName, "TB")
}

// GenerateGraph 调用该方法生成 markdown 文件，可以直接渲染成依赖图。本地生成完后请删除该方法调用
func (d *Scheduler) GenerateGraph(fileName string, direction string) error {
	defer func() {
		fmt.Println("请删除 parallel.GenerateGraph 方法调用")
		os.Exit(1)
	}()
	if len(fileName) == 0 {
		fileName = "graph.md" // 默认当前路径
	}
	f, err1 := os.Create(fileName)
	if err1 != nil {
		return err1
	}
	defer func(f *os.File) {
		if er := f.Close(); er != nil {
			fmt.Println(er)
		}
	}(f)
	f.WriteString("```mermaid\n")
	f.WriteString(fmt.Sprintf("flowchart %s\n", direction))
	for _, dependency := range d.dependencies {
		f.WriteString(fmt.Sprintf("    %s --> %s\n", strings.TrimLeft(fmt.Sprintf("%T", dependency[0]), "*"), strings.TrimLeft(fmt.Sprintf("%T", dependency[1]), "*")))
	}
	f.WriteString("```\n")
	err2 := f.Sync()
	if err2 != nil {
		return err2
	}
	return nil
}
