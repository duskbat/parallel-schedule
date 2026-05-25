package parallel

import "context"

// Step
// 实现时注意，类型 fmt.Sprintf("%T", Step) 会作为 key
// Note: the type name via fmt.Sprintf("%T", Step) is used as the scheduling key
type Step interface {
	Process(ctx context.Context) error
}
