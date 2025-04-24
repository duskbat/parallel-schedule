package parallel

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	myDataBus := &MyDataBus{}
	s1 := &MyStep1{Input: myDataBus}
	s2 := &MyStep2{Input: myDataBus}
	s3 := &MyStep3{Input: myDataBus}
	s4 := &MyStep4{Input: myDataBus}
	s5 := &MyStep5{Input: myDataBus}
	s6 := &MyStep6{Input: myDataBus}
	scheduler := InitScheduler().
		AddDependency(s1, s2).
		AddDependency(s1, s5).
		AddDependency(s2, s3).
		AddDependency(s2, s4).
		AddDependency(s5, s6).
		AddDependency(s3, s6).
		AddDependency(s4, s6)
	// scheduler.GenerateGraphTB("")
	err := scheduler.Launch(context.Background()) // 启动！
	if err != nil {
		fmt.Println("error:", err)
	}
}

/*
step 由用户定义，不需要对 Runner 可见
*/

type MyDataBus struct {
	cellPhone string
}

type MyStep1 struct {
	Step
	Input *MyDataBus
}

func (s *MyStep1) Process(ctx context.Context) error {
	fmt.Println("start s1")
	fmt.Println("end s1")
	time.Sleep(1 * time.Second)
	panic("panic")
	return nil
}

type MyStep2 struct {
	Step
	Input *MyDataBus
}

func (s *MyStep2) Process(ctx context.Context) error {
	fmt.Println("start s2")
	fmt.Println("end s2")
	return nil
}

type MyStep3 struct {
	Step
	Input *MyDataBus
}

func (s *MyStep3) Process(ctx context.Context) error {
	fmt.Println("start s3")
	fmt.Println("end s3")
	return fmt.Errorf("s3")
}

type MyStep4 struct {
	Step
	Input *MyDataBus
}

func (s *MyStep4) Process(ctx context.Context) error {
	fmt.Println("start s4")
	fmt.Println("end s4")
	return nil
}

type MyStep5 struct {
	Step
	Input *MyDataBus
}

func (s *MyStep5) Process(ctx context.Context) error {
	fmt.Println("start s5")
	fmt.Println("end s5")
	return fmt.Errorf("s5")
}

type MyStep6 struct {
	Step
	Input *MyDataBus
}

func (s *MyStep6) Process(ctx context.Context) error {
	fmt.Println("start s6")
	fmt.Println("end s6")
	return nil
}
