package httpx

type ErrorVo struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
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