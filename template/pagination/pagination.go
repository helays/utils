package pagination

import (
	"helay.net/go/utils/v3/tools"
	"math"
)

const (
	defaultPageSize       = 10
	defaultMaxPagesToShow = 7 // 默认显示7个页码（包括折叠按钮）
)

// Pagination 分页结构体
type Pagination struct {
	TotalItems     int // 总记录数
	CurrentPage    int // 当前页码
	PageSize       int // 每页记录数
	MaxPagesToShow int // 最大显示的页码数量(包括折叠按钮)
}

// New 创建分页实例
func New(totalItems, currentPage, pageSize, maxPagesToShow int) *Pagination {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if maxPagesToShow <= 0 {
		maxPagesToShow = defaultMaxPagesToShow
	}
	if currentPage <= 0 {
		currentPage = 1
	}

	return &Pagination{
		TotalItems:     totalItems,
		CurrentPage:    currentPage,
		PageSize:       pageSize,
		MaxPagesToShow: maxPagesToShow,
	}
}

// TotalPages 计算总页数
func (p *Pagination) TotalPages() int {
	if p.PageSize == 0 {
		return 0
	}
	return int(math.Ceil(float64(p.TotalItems) / float64(p.PageSize)))
}

// PageItem 分页项类型
type PageItem struct {
	Type     string // "page", "prev", "next", "ellipsis"
	Page     int    // 页码(仅对Type="page"有效)
	Disabled bool   // 是否禁用
}

// GetPages 获取要显示的分页项
func (p *Pagination) GetPages() []PageItem {
	totalPages := p.TotalPages()
	if totalPages <= 0 {
		return []PageItem{}
	}

	var items []PageItem

	// 添加上一页按钮
	items = append(items, PageItem{
		Type:     "prev",
		Disabled: p.CurrentPage <= 1,
	})

	// 计算需要显示的页码
	if totalPages <= p.MaxPagesToShow {
		// 全部显示
		for i := 1; i <= totalPages; i++ {
			items = append(items, PageItem{
				Type: "page",
				Page: i,
			})
		}
	} else {
		// 需要折叠显示
		halfShow := (p.MaxPagesToShow - 2) / 2 // 两侧各显示多少页码

		// 总是显示第一页
		items = append(items, PageItem{
			Type: "page",
			Page: 1,
		})

		// 左侧折叠
		if p.CurrentPage > halfShow+2 {
			jumpPage := p.CurrentPage - p.MaxPagesToShow
			jumpPage = tools.Ternary(jumpPage < 1, 1, jumpPage)
			items = append(items, PageItem{
				Type: "ellipsis-prev",
				Page: jumpPage,
			})
		}

		// 中间页码
		start := maxFunc(2, p.CurrentPage-halfShow)
		end := minFunc(totalPages-1, p.CurrentPage+halfShow)

		// 调整显示范围，确保显示足够数量的页码
		if p.CurrentPage <= halfShow+1 {
			end = p.MaxPagesToShow - 2
		} else if p.CurrentPage >= totalPages-halfShow {
			start = totalPages - (p.MaxPagesToShow - 3)
		}

		for i := start; i <= end; i++ {
			items = append(items, PageItem{
				Type: "page",
				Page: i,
			})
		}

		// 右侧折叠
		if p.CurrentPage < totalPages-halfShow-1 {
			jumpPage := p.CurrentPage + p.MaxPagesToShow
			jumpPage = tools.Ternary(jumpPage > totalPages, totalPages, jumpPage)
			items = append(items, PageItem{
				Type: "ellipsis-next",
				Page: jumpPage,
			})
		}

		// 总是显示最后一页
		if totalPages > 1 {
			items = append(items, PageItem{
				Type: "page",
				Page: totalPages,
			})
		}
	}

	// 添加下一页按钮
	items = append(items, PageItem{
		Type:     "next",
		Disabled: p.CurrentPage >= totalPages,
	})

	return items
}

// HasPrev 是否有上一页
func (p *Pagination) HasPrev() bool {
	return p.CurrentPage > 1
}

// PrevPage 上一页页码
func (p *Pagination) PrevPage() int {
	if !p.HasPrev() {
		return p.CurrentPage
	}
	return p.CurrentPage - 1
}

// HasNext 是否有下一页
func (p *Pagination) HasNext() bool {
	return p.CurrentPage < p.TotalPages()
}

// NextPage 下一页页码
func (p *Pagination) NextPage() int {
	if !p.HasNext() {
		return p.CurrentPage
	}
	return p.CurrentPage + 1
}

func minFunc(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxFunc(a, b int) int {
	if a > b {
		return a
	}
	return b
}
