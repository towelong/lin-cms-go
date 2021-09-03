package vo

type Page struct {
	Page  int         `json:"page"`
	Count int         `json:"count"`
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type Option func(*Page)

func NewPage(page, count int, options ...Option) *Page{
	p := &Page{
		Page: page,
		Count: count,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *Page) SetItems(items interface{}) {
	p.Items = items
}

func (p *Page) SetTotal(total int) {
	p.Total = total
}

func WithTotal(total int) Option{
	return func(page *Page) {
		page.Total = total
	}
}

func WithItems(items interface{}) Option{
	return func(page *Page) {
		page.Items = items
	}
}