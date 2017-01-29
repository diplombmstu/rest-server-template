package sessions

import (
    "github.com/diplombmstu/rest-server-template/domain"
)

const (
    GetSession = "GetSession"
    CreateSession = "CreateSession"
    DeleteSession = "DeleteSession"
)

func (resource *Resource) generateRoutes(basePath string) {
    if basePath == "" {
        basePath = "/api/sessions"
    }

    resource.routes = &domain.Routes{
        domain.Route{
            Name:           GetSession,
            Method:         "GET",
            Pattern:        "/api/sessions",
            DefaultVersion: "0.0",
            RouteHandlers: domain.RouteHandlers{
                "0.0": resource.HandleGetSession_v0,
            },
            ACLHandler: resource.HandleGetSessionACL,
        },
        domain.Route{
            Name:           CreateSession,
            Method:         "POST",
            Pattern:        "/api/sessions",
            DefaultVersion: "0.0",
            RouteHandlers: domain.RouteHandlers{
                "0.0": resource.HandleCreateSession_v0,
            },
            ACLHandler: resource.HandleCreateSessionACL,
        },
        domain.Route{
            Name:           DeleteSession,
            Method:         "DELETE",
            Pattern:        "/api/sessions",
            DefaultVersion: "0.0",
            RouteHandlers: domain.RouteHandlers{
                "0.0": resource.HandleDeleteSession_v0,
            },
            ACLHandler: resource.HandleDeleteSessionACL,
        },
    }
}
