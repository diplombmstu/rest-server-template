package sessions

import (
    "github.com/diplombmstu/rest-server-template/domain"
    "gopkg.in/mgo.v2/bson"
    "time"
)

type IUserRepositoryFactory interface {
    New() IUserRepository
}

type IUserRepository interface {
    CreateUser(user domain.IUser) error
    GetUsers() domain.IUsers
    FilterUsers(field string, query string, lastID string, limit int, sort string) domain.IUsers
    CountUsers(field string, query string) int
    DeleteUsers(ids []string) error
    DeleteAllUsers() error
    GetUserById(id string) (domain.IUser, error)
    GetUserByUsername(username string) (domain.IUser, error)
    UserExistsByUsername(username string) bool
    UserExistsByEmail(email string) bool
    UpdateUser(id string, inUser domain.IUser) (domain.IUser, error)
    DeleteUser(id string) error
}

type IRevokedTokenRepositoryFactory interface {
    New(db domain.IDatabase) IRevokedTokenRepository
}

type IRevokedTokenRepository interface {
    CreateRevokedToken(token IRevokedToken) error
    DeleteExpiredTokens() error
    IsTokenRevoked(id string) bool
}

// User collection name
const RevokedTokenCollections string = "revoked_tokens"

func NewRevokedTokenRepositoryFactory() IRevokedTokenRepositoryFactory {
    return &RevokedTokenRepositoryFactory{}
}

type RevokedTokenRepositoryFactory struct{}

func (factory *RevokedTokenRepositoryFactory) New(db domain.IDatabase) IRevokedTokenRepository {
    return &RevokedTokenRepository{db}
}

type RevokedTokenRepository struct {
    DB domain.IDatabase
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) CreateRevokedToken(token IRevokedToken) error {
    t := token.(*RevokedToken)
    t.RevokedDate = time.Now()
    return repo.DB.Insert(RevokedTokenCollections, t)
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) DeleteExpiredTokens() error {
    return repo.DB.RemoveAll(RevokedTokenCollections, domain.Query{
        "exp": domain.Query{
            "$lt": time.Now(),
        },
    })
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) IsTokenRevoked(id string) bool {
    if !bson.IsObjectIdHex(id) {
        return false
    }
    return repo.DB.Exists(RevokedTokenCollections, domain.Query{
        "_id": bson.ObjectIdHex(id),
    })
}
