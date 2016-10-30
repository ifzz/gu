package main

import (
	. "github.com/influx6/gu/design"
	. "github.com/influx6/gu/dom"
	_ "github.com/influx6/gu/examples/dashboard/apps"
	"honnef.co/go/js/dom"
)

func main() {

	DashboardApp := New(&DOMRender{
		Document: dom.GetWindow().Document(),
	}).Init()

}
