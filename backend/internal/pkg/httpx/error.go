package httpx

// ErrorVo 统一错误响应体（HTTP 4xx/5xx）。
type ErrorVo struct {
	ErrCode int    `json:"err_code" example:"40100"`
	ErrMsg  string `json:"err_msg" example:"未授权"`
}

type BizError struct {
	HTTPStatus int
	ErrCode    int
	ErrMsg     string
}

func (e *BizError) Error() string { return e.ErrMsg }

func NewBizError(httpStatus, errCode int, msg string) *BizError {
	return &BizError{HTTPStatus: httpStatus, ErrCode: errCode, ErrMsg: msg}
}