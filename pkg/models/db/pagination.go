package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
)

// Pagination ...
type Pagination struct {
	DB      *sql.DB         `json:"-"`
	Default int             `json:"-"`
	Max     int             `json:"-"`
	Min     int             `json:"-"`
	Data    *PaginationData `json:"data"`
}

// kind of context for pagination
type PaginationData struct {
	CurrentPage int   `json:"_current,omitempty"`
	LastPage    int   `json:"_last,omitempty"`
	Next        int   `json:"_next,omitempty"`
	Prev        int   `json:"_prev,omitempty"`
	Total       int   `json:"_total,omitempty"`
	Limit       int   `json:"_limit,omitempty"`
	ShowPage    int   `json:"-"`
	Pages       []int `json:"-"`
}

// Query extracts `page` and `limit` from quey
func (p *Pagination) Query(q *url.URL) *Pagination {
	// create new data for every call to Query
	p.Data = &PaginationData{
		Limit: p.Default,
	}

	page, err := strconv.Atoi(q.Query().Get("page"))
	if err != nil || page == 0 {
		page = 1
	}
	limit, err := strconv.Atoi(q.Query().Get("limit"))
	if err != nil || limit == 0 || limit < p.Min || limit > p.Max {
		limit = p.Default
	}

	showpage, err := strconv.Atoi(q.Query().Get("showpage"))
	if err != nil || showpage == 0 {
		showpage = 5
	}

	p.Data.Limit = limit
	p.Data.CurrentPage = page
	p.Data.ShowPage = showpage

	return p
}

// Generate return the limit and offset sql keyword
func (p *Pagination) Generate(q string) string {
	// make sure current page could not be 0
	if p.Data.CurrentPage == 0 {
		p.Data.CurrentPage = 1
	}

	offset := 0
	if p.Data.Total > 0 {
		offset = p.Data.Limit * (p.Data.CurrentPage - 1)
	}
	// fmt.Println("pagination in Generate", p)
	return fmt.Sprintf("%s limit %d offset %d", q, p.Data.Limit, offset)
}

// SetCount when we know the total number without query
func (p *Pagination) SetCount(n int) {
	p.Data.Total = n
	var lastPage int

	// limit can be very high or negative, to fully load result in 1 page
	if p.Data.Limit >= p.Data.Total || p.Data.Limit == -1 {
		p.Data.CurrentPage, p.Data.LastPage = 1, 1
		return
	}
	lastPage = p.Data.Total / p.Data.Limit

	if lastPage*p.Data.Limit < p.Data.Total {
		p.Data.LastPage = lastPage + 1
	} else {
		p.Data.LastPage = lastPage
	}

	if p.Data.CurrentPage < p.Data.LastPage {
		p.Data.Next = p.Data.CurrentPage + 1
	} else {
		p.Data.CurrentPage = p.Data.LastPage
		p.Data.Next = 0
	}

	if p.Data.CurrentPage > 1 {
		p.Data.Prev = p.Data.CurrentPage - 1
	} else {
		p.Data.Prev = 0
	}

	// ShowPage default 1,3,5,7,9...
	p.Data.Pages = []int{}
	var limitPages = p.Data.CurrentPage + p.Data.ShowPage/2 + 1
	if p.Data.ShowPage >= p.Data.LastPage {
		for i := 1; i < p.Data.LastPage+1; i++ {
			p.Data.Pages = append(p.Data.Pages, i)
		}
	} else {
		if limitPages > p.Data.LastPage {
			for i := p.Data.LastPage - p.Data.ShowPage; i <= p.Data.LastPage; i++ {
				p.Data.Pages = append(p.Data.Pages, i)
			}
		} else {
			if p.Data.CurrentPage <= p.Data.ShowPage/2 {
				for i := 1; i < p.Data.ShowPage+1; i++ {
					p.Data.Pages = append(p.Data.Pages, i)
				}
			} else {
				for i := p.Data.CurrentPage - p.Data.ShowPage/2; i < limitPages; i++ {
					p.Data.Pages = append(p.Data.Pages, i)
				}
			}
		}
	}

}

// Count execute the count query to find the data.total
func (p *Pagination) Count(query string, dest ...interface{}) error {
	row := p.DB.QueryRow(query, dest...)
	if err := row.Scan(&p.Data.Total); err != nil {
		return err
	}
	p.SetCount(p.Data.Total)
	return nil
}

// MaxLimit set Limit to Max
func (p *Pagination) MaxLimit() {
	p.Data.Limit = p.Max
}
