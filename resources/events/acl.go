package events

import (
    "github.com/diplombmstu/rest-server-template/domain"
    "net/http"
    "github.com/diplombmstu/rest-server-template/resources/users"
)

func (resource *Resource) HandleSubscribeACL(req *http.Request, user domain.IUser) (bool, string) {
    if user == nil {
        return false, "Anonymous access is denied"
    }

    if !user.HasRole(users.RoleAdmin) {
        return false, "User has not enough rights"
    }

    return true, ""
}