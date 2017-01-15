// Package redom contains the DOM RenderTarget for rendering gu markup into a
// browser DOM.
package redom

import (
	"honnef.co/go/js/dom"

	gjs "github.com/gopherjs/gopherjs/js"
	"github.com/gu-io/gu"
	"github.com/gu-io/gu/dispatch"
	"github.com/gu-io/gu/js"
	"github.com/gu-io/gu/trees"
)

// AddStylesheet adds an external stylesheet to the document into the document
// DOM using the provide URL.
func AddStylesheet(url string) {
	link := gjs.Global.Get("document").Call("createElement", "link")
	link.Set("href", url)
	link.Set("rel", "stylesheet")

	gjs.Global.Get("document").Get("head").Call("appendChild", link)
}

// DOMRenderer defines an implementation for gu.gu.ResourceRenderere
// and handles rendering of a giving group of resources into the live DOM body root.
type DOMRenderer struct {
	Document dom.Document
}

// Render renders the giving set of resources into the provided body and header
// of the DOM.
func (dm *DOMRenderer) Render(rs ...*gu.ResourceDefinition) {
	head := dm.Document.QuerySelector("head")
	body := dm.Document.QuerySelector("body")

	// clear all children of head and body if the belong to us.
	for _, item := range head.QuerySelectorAll("[data-gen*='gu']") {
		if !item.HasAttribute("gu-resource-root") {
			item.ParentNode().RemoveChild(item)
		}
	}

	for _, item := range body.QuerySelectorAll("[data-gen*='gu']") {
		if !item.HasAttribute("gu-resource-root") {
			item.ParentNode().RemoveChild(item)
		}
	}

	// Render the important resources and normal links first.
	for _, res := range rs {

		// Retrieve base resources which must be rendered for each component and
		// add them according to the head and body.
		toHead, toBody := res.Resources()

		for _, item := range toHead {
			js.Patch(js.CreateFragment(item.HTML()), head.Underlying(), false)
		}

		for _, item := range toBody {
			js.Patch(js.CreateFragment(item.HTML()), body.Underlying(), false)
		}

		for _, item := range res.Links {
			markup := item.Render()

			markup.EachEvent(func(ev *trees.Event, root *trees.Markup) {
				dm.BindEvent(ev, head.Underlying())
			})

			js.Patch(js.CreateFragment(markup.HTML()), head.Underlying(), false)
		}
	}

	// Render all basic views which are not to be deffered.
	for _, res := range rs {
		for _, render := range res.Renderables {
			dm.RenderUpdate(render.View, render.Targets, false)
		}
	}

	// Render all defered views.
	for _, res := range rs {
		for _, render := range res.DRenderables {
			dm.RenderUpdate(render.View, render.Targets, false)
		}
	}

	// Render the defered links.
	for _, res := range rs {
		for _, item := range res.DeferLinks {
			markup := item.Render()

			markup.EachEvent(func(ev *trees.Event, root *trees.Markup) {
				dm.BindEvent(ev, body.Underlying())
			})

			js.Patch(js.CreateFragment(markup.HTML()), body.Underlying(), false)
		}
	}
}

// RenderUpdate handles rendering calls for individual renderers which may have
// determined targets within the body.
func (dm *DOMRenderer) RenderUpdate(rv gu.Renderable, targets string, update bool) {
	body := dm.Document.QuerySelector("body")

	if targets == "" {
		markup := rv.Render()

		if !update {
			markup.EachEvent(func(ev *trees.Event, root *trees.Markup) {
				dm.BindEvent(ev, body.Underlying())
			})
		}

		if kvr, ok := rv.(gu.RenderView); ok {
			js.Patch(js.CreateFragment(markup.HTML()), body.Underlying(), !kvr.RenderedBefore())

			if cvs, ok := rv.(gu.ViewHooks); ok {
				mounted, _, _ := cvs.Hooks()
				if !mounted.Used() {
					mounted.Publish()
				}
			}

			return
		}

		js.Patch(js.CreateFragment(markup.HTML()), body.Underlying(), false)

		if cvs, ok := rv.(gu.ViewHooks); ok {
			mounted, _, _ := cvs.Hooks()
			if !mounted.Used() {
				mounted.Publish()
			}
		}

		return
	}

	kernels := body.QuerySelectorAll(targets)

	for _, targetDOM := range kernels {
		markup := rv.Render()

		if !update {
			markup.EachEvent(func(ev *trees.Event, root *trees.Markup) {
				dm.BindEvent(ev, targetDOM.Underlying())
			})
		}

		if kvr, ok := rv.(gu.RenderView); ok {
			js.Patch(js.CreateFragment(markup.HTML()), body.Underlying(), !kvr.RenderedBefore())
			continue
		}

		js.Patch(js.CreateFragment(markup.HTML()), targetDOM.Underlying(), false)
	}

	if cvs, ok := rv.(gu.ViewHooks); ok {
		mounted, _, _ := cvs.Hooks()

		if !mounted.Used() {
			mounted.Publish()
		}
	}
}

// BindNodeEvent binds the removal events and addition events to fire the mount and
// unmount handlers for a rendered source.
// NOTE: Currently only handles registering for removals has resources RenderUpdate
// handles insertion(mount) calls.
// TODO: Implement removal calls.
func (dm *DOMRenderer) BindNodeEvent(source *trees.Event, root *gjs.Object) {

}

// BindEvent connects the event with the provided event object and root.
func (dm *DOMRenderer) BindEvent(source *trees.Event, root *gjs.Object) {
	source.Link = func(ev *gjs.Object) { dm.TriggerBindEvent(ev, root, source) }

	root.Call("addEventListener", source.Type, source.Link, true)

	source.Handle.AddEnd(func() {
		root.Call("removeEventListener", source.Type, source.Link, true)
	})
}

// TriggerBindEvent connects the giving event with the provided dom target.
func (dm *DOMRenderer) TriggerBindEvent(event *gjs.Object, root *gjs.Object, source *trees.Event) {
	target := event.Get("target")

	children := root.Call("querySelectorAll", source.Target())
	if children == nil || children == gjs.Undefined {
		return
	}

	kids := js.DOMObjectToList(children)
	var match bool

	for _, item := range kids {
		if item != target {
			continue
		}

		match = true
		break
	}

	// if we match then run the listeners registered.
	if match {
		dispatch.Dispatch(trees.EventBroadcast{
			EventID: source.EventID,
			Event:   trees.NewWrapperEvent(event),
		})
	}
}
