package dbkit

type PageResponse[T any] struct {
	Code  int       `json:"code"`
	Msg   string    `json:"msg"`
	Data  []T       `json:"data"`
	Total int64     `json:"total"`
	Page  *PageInfo `json:"page,omitempty"`
}

type PageInfo struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func Success[T any](data T) Response[T] {
	return Response[T]{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func SuccessWithPage[T any](data []T, page *Page, total int64) PageResponse[T] {
	resp := PageResponse[T]{
		Code:  200,
		Msg:   "success",
		Data:  data,
		Total: total,
	}

	if page != nil && page.IsValid() {
		resp.Page = &PageInfo{
			PageNum:  page.PageNum,
			PageSize: page.PageSize,
		}
	}

	return resp
}

func Error(msg string) Response[interface{}] {
	return Response[interface{}]{
		Code: 500,
		Msg:  msg,
		Data: nil,
	}
}
