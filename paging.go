package gpaging

import (
	"errors"
	glist "github.com/og/x/list"
	"log"
	"math"
)

type Gen struct {
	Page int `json:"page"`
	PerPage int `json:"perPage"`
	Total int `json:"total"`
	ClosestPagesLength int `json:"closestPagesLength"`
	JumpBatchPageInterval  int `json:"jumpBatchPageInterval"`
}
var WarningLog = func(msg string) {
	log.Print(msg)
}
func genCheckAndFix(gen *Gen) {
	if gen.Page < 0 {
		panic(errors.New(messagePrefix + "gen.Page can not less 0"))
	}
	if gen.Page == 0 {
		gen.Page = 1
		WarningLog(messagePrefix + "gen.Page can not be 0, gpaging will set page 1, but you need check your code")
	}
	if gen.Total < 0 {
		panic(errors.New(messagePrefix + "gen.Total can not less 0"))
	}
	if gen.PerPage < 0 {
		panic(errors.New(messagePrefix + "gen.PerPage can not less 0"))
	}
	if gen.PerPage == 0 {
		panic(errors.New(messagePrefix + "gen.PerPage can not be 0"))
	}
}
type Render struct {
	HasPaging bool `json:"hasPaging"`
	LastPage int `json:"lastPage"`
	IsFirstPage bool `json:"isFirstPage"`
	IsLastPage bool `json:"isLastPage"`
	ClosestPages RenderClosestPages `json:"closestPages"`
	JumpBatchPage RenderJumpBatchPage `json:"jumpBatchPage"`
}
type RenderClosestPages struct {
	Prev []int `json:"prev"`
	Next []int `json:"next"`
}
type RenderJumpBatchPage struct{
	HasPrev bool `json:"hasPrev"`
	PrevPage int `json:"prevPage"`
	HasNext bool `json:"hasNext"`
	NextPage int `json:"nextPage"`
}
const messagePrefix = "og/go-paging: gpaging.CreateData(gen) "
func CreateData(gen *Gen) ( Render) {
	genCheckAndFix(gen)
	render := Render{}
	// 避免 []int(nil)
	render.ClosestPages.Prev = []int{}
	render.ClosestPages.Next = []int{}
	render.HasPaging = gen.Total != 0
	render.LastPage = int(math.Ceil(float64(gen.Total) / float64(gen.PerPage)))
	// 有时客户端控制的 page 是自由输入的，或者在提交的那一刻是有10页的，但是渲染的时候数据被删了只有9页
	// last page 可能是 0
	if gen.Page > render.LastPage && render.LastPage !=0 {
		gen.Page = render.LastPage
	}
	// LastPage 可能为 0
	render.IsFirstPage = gen.Page == 1

	render.IsLastPage = gen.Page == render.LastPage
	{
		firstPage := 1
		pagesLength := closestPagesLength(*gen, firstPage)
		glist.Run(pagesLength, func(i int) (_break bool) {
			order := i+1
			render.ClosestPages.Prev = append([]int{gen.Page - order}, render.ClosestPages.Prev...)
			return
		})
	}
	{
		pagesLength := closestPagesLength(*gen, render.LastPage)
		glist.Run(pagesLength, func(i int) (_break bool) {
			order := i+1
			render.ClosestPages.Next = append(render.ClosestPages.Next, gen.Page + order)
			return
		})
	}
	{
		prevJumpPage := gen.Page - gen.JumpBatchPageInterval
		firstPage := 1
		if prevJumpPage > firstPage {
			render.JumpBatchPage.HasPrev = true
			render.JumpBatchPage.PrevPage = prevJumpPage
		}
	}
	{
		nextJumpPage := gen.Page + gen.JumpBatchPageInterval
		if nextJumpPage < render.LastPage {
			render.JumpBatchPage.HasNext = true
			render.JumpBatchPage.NextPage = nextJumpPage
		}
	}
	return render
}

func closestPagesLength (gen Gen, firstPageOrLastPage int) (pagesLength int) {
	pagesLength = naturalNumberInterval(gen.Page, firstPageOrLastPage)
	if pagesLength > gen.ClosestPagesLength {
		pagesLength =  gen.ClosestPagesLength
	}
	return
}
func naturalNumberInterval(a int, b int) (interval int) {
	var min, max int
	if a > b {
		min = b
		max = a
	} else {
		min = a
		max = b
	}
	// min:1  max:1
	// min:9 max:9
	if min == max { return 0 }
	// 1 2 3 4 5 6 7 8 9
	// min 1 max 5
	// 5-1-1
	// min 5 max 9
	// 9-5-1
	return max - min - 1
}
