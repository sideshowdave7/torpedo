package torpedo

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Resource interface {
	Get(values url.Values) (int, string)
	Post(values url.Values) (int, string)
	Put(values url.Values) (int, string)
	Delete(values url.Values) (int, string)
}

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(values url.Values) (int, string) {
	return 405, ""
}

func (PostNotSupported) Post(values url.Values) (int, string) {
	return 405, ""
}

func (PutNotSupported) Put(values url.Values) (int, string) {
	return 405, ""
}

func (DeleteNotSupported) Delete(values url.Values) (int, string) {
	return 405, ""
}

type API struct{}

func (api *API) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func (api *API) requestHandler(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data string
		var code int

		request.ParseForm()
		method := request.Method
		values := request.Form

		switch method {
		case GET:
			code, data = resource.Get(values)
		case POST:
			code, data = resource.Post(values)
		case PUT:
			code, data = resource.Put(values)
		case DELETE:
			code, data = resource.Delete(values)
		default:
			api.Abort(rw, 405)
			return
		}

		rw.WriteHeader(code)
		rw.Write([]byte(data))
	}
}

func (api *API) AddResource(resource Resource, path string) {
	http.HandleFunc(path, api.requestHandler(resource))
}

func (api *API) Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
}
