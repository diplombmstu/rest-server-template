package events

import "net/http"

func (resource *Resource) HandleSubscribe_v0(w http.ResponseWriter, req *http.Request) {
    resource.Broker.ServeHTTP(w, req, resource.ctx.GetCurrentUserCtx(req).GetID())
}
