package constants

//Status
const (
	StatusCreated    = "created"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

// Context
const (
	UserIdCtx    = "userId"
	UserRoleCtx  = "userRole"
	RequestIdCtx = "requestId"
	TraceIdCtx   = "traceId"
	SpanIdCtx    = "spanId"
	ApiNameCtx   = "apiName"
)

// Errors
const (
	BindBodyError      string = "bind_body"
	BindPathError      string = "bind_path"
	UserIdTypeMismatch string = "type_mismatch"
)

