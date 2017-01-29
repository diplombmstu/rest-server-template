package users

import (
    "github.com/diplombmstu/rest-server-template/domain"
    "net/http"
    "github.com/diplombmstu/rest-server-template/resources/sessions"
)

type PostCreateUserHookPayload struct {
    User domain.IUser
}

type PostConfirmUserHookPayload struct {
    User domain.IUser
}

type ControllerHooks struct {
    PostCreateUserHook  func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostCreateUserHookPayload) error
    PostConfirmUserHook func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostConfirmUserHookPayload) error
}

type Options struct {
    BasePath              string
    Database              domain.IDatabase
    Renderer              domain.IRenderer
    UserRepositoryFactory sessions.IUserRepositoryFactory
    ControllerHooks       *ControllerHooks
}

func NewResource(ctx domain.IContext, options *Options) *Resource {

    database := options.Database
    if database == nil {
        panic("users.Options.Database is required")
    }
    renderer := options.Renderer
    if renderer == nil {
        panic("users.Options.Renderer is required")
    }

    userRepositoryFactory := options.UserRepositoryFactory
    if userRepositoryFactory == nil {
        // init default UserRepositoryFactory
        userRepositoryFactory = NewUserRepositoryFactory(database)
    }

    controllerHooks := options.ControllerHooks
    if controllerHooks == nil {
        controllerHooks = &ControllerHooks{nil, nil}
    }

    u := &Resource{ctx, options, nil,
        database,
        renderer,
        userRepositoryFactory,
        controllerHooks,
    }
    u.generateRoutes(options.BasePath)
    return u
}

// UsersResource implements IResource
type Resource struct {
    ctx                   domain.IContext
    options               *Options
    routes                *domain.Routes
    Database              domain.IDatabase
    Renderer              domain.IRenderer
    UserRepositoryFactory sessions.IUserRepositoryFactory
    ControllerHooks       *ControllerHooks
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

func (resource *Resource) UserRepository(req *http.Request) sessions.IUserRepository {
    return resource.UserRepositoryFactory.New()
}
