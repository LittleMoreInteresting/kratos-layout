package pager

type Pager struct {
	PageNumber int
	PageSize   int
}

func (p *Pager) Offsite() int {
	if p.PageNumber == 0 {
		return 0
	}
	return (p.PageNumber - 1) * p.PageSize
}

func (p *Pager) Size() int {
	if p.PageSize <= 0 || p.PageSize > 200 {
		return 20
	}
	return p.PageSize
}
