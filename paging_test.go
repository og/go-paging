package gpaging_test

import (
	gpaging "github.com/og/go-paging"
	gis "github.com/og/x/test"
	"testing"
)

func TestGenCheck(t *testing.T) {
	is := gis.New(t)
	func(){
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "og/go-paging: gpaging.CreateData(gen) gen.Page can not less 0")
		}()
		// page < 0
		gpaging.CreateData(&gpaging.Gen{
			Page:    -1,
		})
	}()
	{
		gen := gpaging.Gen{
			Page:    0,
			PerPage: 10,
		}
		waringMsg := ""
		gpaging.WarningLog = func(msg string) {
			waringMsg = msg
		}
		render := gpaging.CreateData(&gen)
		_=render
		is.Eql(gen.Page, 1)
		is.Eql(waringMsg, "og/go-paging: gpaging.CreateData(gen) gen.Page can not be 0, gpaging will set page 1, but you need check your code")
	}
	func(){
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "og/go-paging: gpaging.CreateData(gen) gen.Total can not less 0")
		}()
		// total < 0
		gpaging.CreateData(&gpaging.Gen{
			Page:    1,
			Total: -1,
			PerPage: 10,
		})
	}()
	func(){
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "og/go-paging: gpaging.CreateData(gen) gen.PerPage can not less 0")
		}()
		// total < 0
		gpaging.CreateData(&gpaging.Gen{
			Page:    1,
			Total: 1,
			PerPage: -1,
		})
	}()
	func(){
		defer func() {
			r := recover()
			is.Eql(r.(error).Error(), "og/go-paging: gpaging.CreateData(gen) gen.PerPage can not be 0")
		}()
		// total < 0
		gpaging.CreateData(&gpaging.Gen{
			Page:    1,
			Total: 1,
			PerPage: 0,
		})
	}()
}

func TestHasPaging(t *testing.T) {
	is := gis.New(t)
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page:    1,
			PerPage: 10,
			Total:   0,
		})
		is.Eql(render.HasPaging, false)
	}
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page:    1,
			PerPage: 10,
			Total:   1,
		})
		is.Eql(render.HasPaging, true)
	}
}

func TestLastPage(t *testing.T) {
	is := gis.New(t)
	_= is
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 1,
			PerPage: 10,
			Total: 99,
		})
		is.Eql(render.LastPage, 10)
	}
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 1,
			PerPage: 10,
			Total: 100,
		})
		is.Eql(render.LastPage, 10)
	}
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 1,
			PerPage: 10,
			Total: 109,
		})
		is.Eql(render.LastPage, 11)
	}
}

func TestIsFirstPage(t *testing.T) {
	is := gis.New(t)
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 1,
			PerPage: 10,
			Total: 99,
		})
		is.Eql(render.IsFirstPage, true)
	}
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 2,
			PerPage: 10,
			Total: 99,
		})
		is.Eql(render.IsFirstPage, false)
	}
}


func TestIsLastPage(t *testing.T) {
	is := gis.New(t)
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 10,
			PerPage: 10,
			Total: 100,
		})
		is.Eql(render.IsLastPage, true)
	}
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 9,
			PerPage: 10,
			Total: 100,
		})
		is.Eql(render.IsLastPage, false)
	}
}
func TestClosestPages(t *testing.T) {
	is := gis.New(t)
	_=is

	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page:    10,
			PerPage: 10,
			Total:   200,
			ClosestPagesLength: 3,
		})
		is.Eql(render.ClosestPages.Prev, []int{7,8,9})
		is.Eql(render.ClosestPages.Next, []int{11,12,13})
	}
}

func TestCreateData_jumpBatchPages(t *testing.T) {
	is := gis.New(t)
	{
		render := gpaging.CreateData(&gpaging.Gen{
			Page: 10,
			PerPage: 10,
			Total:200,
			JumpBatchPageInterval: 3,
		})
		is.Eql(render.JumpBatchPage, gpaging.RenderJumpBatchPage{
			HasPrev: true,
			PrevPage: 7,
			HasNext: true,
			NextPage: 13,
		})
	}
}

type TestData struct {
	Gen gpaging.Gen
	Render gpaging.Render
}
func TestCreateDataTestData(t *testing.T) {
	testDataList := []TestData{
		{
			gpaging.Gen{1,10,101,3,3,},
			gpaging.Render{
				HasPaging:true,
				LastPage:11,
				IsFirstPage:true,
				IsLastPage:false,
				ClosestPages:gpaging.RenderClosestPages{
					[]int{},
					[]int{ 2, 3, 4, },
				},
				JumpBatchPage:gpaging.RenderJumpBatchPage{ false, 0, true, 4, },
			},
		},
		{
			gpaging.Gen{2,10,101,3,3,},
			gpaging.Render{
				HasPaging:true,
				LastPage:11,
				IsFirstPage:false,
				IsLastPage:false,
				ClosestPages:gpaging.RenderClosestPages{
					[]int{},
					[]int{ 3, 4, 5,},
				},
				JumpBatchPage:gpaging.RenderJumpBatchPage{ false, 0, true, 5, },
			},
		},

	}
	is := gis.New(t)
	for _, test := range testDataList {
		render := gpaging.CreateData(&test.Gen)
		is.Eql(render, test.Render)
	}
}