package parallel

import (
	"context"
	"fmt"
	"runtime/debug"
)

type Scheduler struct {
	dependencies  [][2]Step                      // 依赖关系
	table         map[string]Step                // 映射
	adjacencyList map[string]map[string]struct{} // 邻接表
	inDegrees     map[string]int                 // 入度
}

func InitScheduler() *Scheduler {
	return &Scheduler{
		dependencies:  [][2]Step{},
		adjacencyList: map[string]map[string]struct{}{},
		table:         map[string]Step{},
		inDegrees:     map[string]int{},
	}
}

func (d *Scheduler) AddDependency(pre Step, next Step) *Scheduler {
	d.dependencies = append(d.dependencies, [2]Step{pre, next})
	return d
}

func (d *Scheduler) init() error {
	for _, edge := range d.dependencies {
		k1 := fmt.Sprintf("%T", edge[0])
		k2 := fmt.Sprintf("%T", edge[1])
		d.table[k1] = edge[0]
		d.table[k2] = edge[1]
		if _, ok := d.adjacencyList[k1]; !ok {
			d.adjacencyList[k1] = make(map[string]struct{})
		}
		d.adjacencyList[k1][k2] = struct{}{}

		if _, ok := d.inDegrees[k1]; !ok {
			d.inDegrees[k1] = 0
		}
		d.inDegrees[k2]++
	}
	visit := make(map[string]int)
	for k := range d.table {
		if err := d.checkEndlessLoop(k, visit); err != nil {
			return err
		}
	}
	return nil
}

// 死循环自检 防呆傻
func (d *Scheduler) checkEndlessLoop(k string, visit map[string]int) error {
	visit[k] = -1
	if children, ok := d.adjacencyList[k]; ok {
		for c := range children {
			switch visit[c] {
			case 0:
				if err := d.checkEndlessLoop(c, visit); err != nil {
					return err
				}
			case -1:
				return fmt.Errorf("endless loop: %s", k)
			}
		}
	}
	visit[k] = 1
	return nil
}

func (d *Scheduler) Launch(ctx context.Context) (res error) {
	err0 := d.init()
	if err0 != nil {
		return err0
	}
	l := len(d.table)
	finishChan := make(chan string, l)
	errorChan := make(chan error, l)
	// 分析依赖关系，拓扑排序
	for k, in := range d.inDegrees {
		if len(errorChan) > 0 {
			res = <-errorChan
			break
		}
		if in == 0 {
			go func() {
				defer deferRecoverFunc(ctx, errorChan)
				err1 := d.table[k].Process(ctx)
				if err1 != nil {
					errorChan <- err1
				}
				finishChan <- k
			}()
		}
	}
	for i := 0; res == nil && i < l; i++ {
		select {
		case res = <-errorChan:
		case k := <-finishChan:
			for child := range d.adjacencyList[k] {
				if len(errorChan) > 0 {
					res = <-errorChan
					break
				}
				d.inDegrees[child]--
				if d.inDegrees[child] == 0 {
					go func() {
						defer deferRecoverFunc(ctx, errorChan)
						err2 := d.table[child].Process(ctx)
						if err2 != nil {
							errorChan <- err2
						}
						finishChan <- child
					}()
				}
			}
		}
	}
	return res
}

func deferRecoverFunc(ctx context.Context, c chan error) {
	if r := recover(); r != nil {
		c <- PanicError{
			ErrMsg: fmt.Sprintf("[PANIC] %v %s", r, debug.Stack()),
		}
	}
}
