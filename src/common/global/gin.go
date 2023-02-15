package global

const (
	GinRespHeaderTracerIdKey = "x-tracer-id"

	GinParamQueryKey    = "query"
	GinParamOrderByKey  = "order_by"
	GinParamPageKey     = "page"
	GinParamPageSizeKey = "page_size"

	GinCtxTracerIdKey = "__gin_tracer_id"
	GinCtxQueryKey    = "__gin_query"
	GinCtxOrderByKey  = "__gin_order_by"
	GinCtxPageKey     = "__gin_page"
	GinCtxPageSizeKey = "__gin_page_size"
)
