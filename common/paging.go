package common

type Paging struct {
	Page  int   `json:"page" form:"page" binding:"required,min=1"`
	Limit int   `json:"limit" form:"limit" binding:"required,min=1,max=20"`
	Total int64 `json:"total" form:"-"`
}

func (c *Paging) Process() {
	if c.Page <= 0 {
		c.Page = 1
	}

	if c.Limit <= 0 {
		c.Limit = 20
	}

	if c.Limit > 100 {
		c.Limit = 20
	}

	if c.Total <= 0 {
		c.Total = 0
	}
}

func (c *Paging) Offset() int {
	return (c.Page - 1) * c.Limit
}