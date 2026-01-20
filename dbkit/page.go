package dbkit

type Page struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type Pageable interface {
	IsValid() bool
	GetOffset() int
	GetLimit() int
}

func (p *Page) IsValid() bool {
	return p != nil && p.PageNum > 0 && p.PageSize > 0
}

func (p *Page) GetOffset() int {
	if !p.IsValid() {
		return 0
	}
	return (p.PageNum - 1) * p.PageSize
}

func (p *Page) GetLimit() int {
	if !p.IsValid() {
		return 0
	}
	return p.PageSize
}
