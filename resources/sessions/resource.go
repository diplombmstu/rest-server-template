package sessions

import (
    "net/http"
    "github.com/diplombmstu/rest-server-template/domain"
    "crypto/rsa"
)

type PostCreateSessionHookPayload struct {
    TokenString string
}

type PostDeleteSessionHookPayload struct {
    Claims ITokenClaims
}

type ControllerHooks struct {
    PostCreateSessionHook func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostCreateSessionHookPayload) error
    PostDeleteSessionHook func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostDeleteSessionHookPayload) error
}

type Options struct {
    BasePath                      string
    TokenAuthority                ITokenAuthority
    PrivateSigningKey             *rsa.PrivateKey
    PublicSigningKey              *rsa.PublicKey
    Database                      domain.IDatabase
    Renderer                      domain.IRenderer
    RevokedTokenRepositoryFactory IRevokedTokenRepositoryFactory
    UserRepositoryFactory         IUserRepositoryFactory
    ControllerHooks               *ControllerHooks
}

func NewResource(ctx domain.IContext, options *Options) *Resource {

    database := options.Database
    if database == nil {
        panic("resources.sessions.Options.Database is required")
    }

    renderer := options.Renderer
    if renderer == nil {
        panic("resources.sessions.Options.Renderer is required")
    }

    userRepositoryFactory := options.UserRepositoryFactory
    if userRepositoryFactory == nil {
        panic("resources.sessions.options.UserRepositoryFactory is required")
    }

    tokenAuthority := options.TokenAuthority
    if tokenAuthority == nil {
        if options.PrivateSigningKey == nil {
            panic("resources.sessions.options.PrivateSigningKey is required")
        }
        if options.PublicSigningKey == nil {
            panic("resources.sessions.options.PublicSigningKey is required")
        }
        // init default RevokedTokenRepository
        tokenAuthority = NewTokenAuthority(&TokenAuthorityOptions{
            PrivateSigningKey: options.PrivateSigningKey,
            PublicSigningKey:  options.PublicSigningKey,
        })
    }

    revokedTokenRepositoryFactory := options.RevokedTokenRepositoryFactory
    if revokedTokenRepositoryFactory == nil {
        // init default RevokedTokenRepositoryFactory
        revokedTokenRepositoryFactory = NewRevokedTokenRepositoryFactory()
    }

    controllerHooks := options.ControllerHooks
    if controllerHooks == nil {
        controllerHooks = &ControllerHooks{nil, nil}
    }

    resource := &Resource{
        ctx,
        options,
        nil,
        database,
        renderer,
        tokenAuthority,
        revokedTokenRepositoryFactory,
        userRepositoryFactory,
        controllerHooks,
    }

    resource.generateRoutes(options.BasePath)

    return resource
}

// SessionsResource implements IResource
type Resource struct {
    ctx                           domain.IContext
    options                       *Options
    routes                        *domain.Routes
    Database                      domain.IDatabase
    Renderer                      domain.IRenderer
    TokenAuthority                ITokenAuthority
    RevokedTokenRepositoryFactory IRevokedTokenRepositoryFactory
    UserRepositoryFactory         IUserRepositoryFactory
    ControllerHooks               *ControllerHooks
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

func (resource *Resource) NewAuthenticator() *Authenticator {
    return NewAuthenticator(resource)
}

func (resource *Resource) RevokedTokenRepository(req *http.Request) IRevokedTokenRepository {
    return resource.RevokedTokenRepositoryFactory.New(resource.Database)
}

func (resource *Resource) UserRepository(req *http.Request) IUserRepository {
    return resource.UserRepositoryFactory.New()
}
