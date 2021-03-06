App, View, Components, Drivers and the Connection
=================================================

The `App`, `View`, `Component` and `Drivers` are concepts which are internal to Gu, but are important in the understanding of how the applications developed with Gu will be developed.

Apps are centralized registry for views in Gu. In Gu, Views are likened to a page with it's content and the `App` structure is responsible to manage and render the appropriate view for the appropriate route provided. It also handles the underlying details of view updates and encapsulates that away from the user.

Views as said are the pages in Gu, they themselves define the markup which will be the mount points for components. Views were created in this manager to allow a fine tuned control on what gets rendered and where individual components will be rendered without having to manage markup declared in html files. Using this approach views fully encapsulate their details and provide a testable bed for the current state of the application.

In Gu, though it's possible to render multiple views based on the provided route but it's more suited to have one view encapsulate the displayed components for a given app.

See [Drivers](./drivers.md), for more on components. See [Components](./components.md), for more on components.

Relations of Apps and Drivers
-----------------------------

Drivers define the rendering target for an app. Gu was built with the forsight that there will be the need to provide the capability to target different rendering systems (eg. browser, mobile and deskop webviews), hence the design of the way Gu works, was created to suite this.

Drivers provide an encapsulated means by which such rendering details can easily be used to render the outputs of any giving apps. By providing structures that satisfy the `Driver` interface, practically any rendering platform that can translate the html generated by the App, can be used for rendering. This in itself provides a powerful functionality.

Relations of Views and Components
---------------------------------

Views are the central system to manage components, their `App` has no business in the managed of a component nor it's lifecycles, it centrally only manages the views and expects each view to handle the rendering calls and requirements of the it's components. This allows us to create specific views which individually represent a given page of a larger app more effectively, has only the components for that page ever exists for that view.

Example
-------

Example of a View registered to a App

```go

router := router.NewRouter(_, memorycache.New("greeter"))
app := gu.App("GreeterApp", router)

index := app.View(elems.Parse(`
		<div class="greeter-view view wrapper">
			<h1 class="view-header">Greeter App</h1>

			<div class="greeter-app" id="greeter-app-component">
			</div>
		</div>
	`, elems.CSS(`
			&{
				color: #fff;
				width: 100%;
				padding: 10px;
				min-height: 100%;
				margin: 0px auto;
				background: rgba(0,0,0,0.4);
			}

			& .greeter-app {
				width: 90%;
				height: auto;
				margin: 30px auto;
				padding-top: 100px;
				text-align: center;
			}
	`, nil), "/*", 0),
})

index.Component(components.NewGreeter(), gu.AnyOrder, "/*", "#greeter-app-component")

```
