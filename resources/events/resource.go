package events

import (
    "github.com/diplombmstu/rest-server-template/domain"
    "net/http"
    "github.com/diplombmstu/rest-server-template/sse-service/sse_server"
)

type Options struct {
    BasePath string
    Renderer domain.IRenderer
    Broker   *sse_server.SseBroker
}

func NewResource(ctx domain.IContext, options *Options) *Resource {
    renderer := options.Renderer
    if renderer == nil {
        panic("users.Options.Renderer is required")
    }

    broker := options.Broker
    if broker == nil {
        panic("Broker is required")
    }

    res := &Resource{
        ctx,
        options,
        nil,
        renderer,
        broker,
    }

    res.generateRoutes(options.BasePath)

    return res
}

// UsersResource implements IResource
type Resource struct {
    ctx      domain.IContext
    options  *Options
    routes   *domain.Routes
    Renderer domain.IRenderer
    Broker   *sse_server.SseBroker
}

func (resource *Resource) Context() domain.IContext {
    return resource.ctx
}

func (resource *Resource) Routes() *domain.Routes {
    return resource.routes
}

func (resource *Resource) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
    resource.Renderer.Render(w, req, status, v)
}
