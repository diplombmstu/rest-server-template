package events

import "github.com/diplombmstu/rest-server-template/domain"

const (
    Subscribe = "Subscribe"
)

func (resource *Resource) generateRoutes(basePath string) {
    resource.routes = &domain.Routes{
        domain.Route{
            Name:           Subscribe,
            Method:         "POST",
            Pattern:        "/events/subscribe",
            DefaultVersion: "0.0",
            RouteHandlers: domain.RouteHandlers{
                "0.0": resource.HandleSubscribe_v0,
            },
            ACLHandler: resource.HandleSubscribeACL,
        },
    }
}