package parallel

import "context"

// Step
// 实现时注意，类型 fmt.Sprintf("%T", Step) 会作为 key
type Step interface {
	Process(ctx context.Context) error
}
