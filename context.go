package hasaki

import "context"

const ctxRequestLatency = "REQUEST_LATENCY"

// 从 ctx 获得请求延迟，一般用在 AfterFunc 中
func GetRequestLatency(ctx context.Context) int64 {
	if v := ctx.Value(ctxRequestLatency); v != nil {
		if latency, ok := v.(int64); ok {
			return latency
		}
	}
	return 0
}
