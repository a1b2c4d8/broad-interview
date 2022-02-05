package client

type RoutesResponse struct {
	Routes []*RouteResponseData `json:"data"`
}

type RouteResponseData struct {
	Id              string                   `json:"id"`
	RouteAttributes *RouteResponseAttributes `json:"attributes"`
}

type RouteResponseAttributes struct {
	Name string `json:"long_name"`
}
