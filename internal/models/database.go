package models

import "time"

type ListOptions struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
	Search string `form:"search"`

	MinClicks int       `form:"min_clicks"`
	MaxClicks int       `form:"max_clicks"`
	MinDate   time.Time `form:"min_date"`
	MaxDate   time.Time `form:"max_date"`
}

func (o *ListOptions) Normalize() {
	if o.Page < 1 {
		o.Page = 1
	}
	if o.Limit < 1 || o.Limit > 100 {
		o.Limit = 20
	}
}
