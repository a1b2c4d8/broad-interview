package service

type RouteId string
type RouteName string

type Route struct {
	Id      RouteId
	Name    RouteName
	StopIds map[StopId]struct{}
}

func NewRoute(id RouteId, name RouteName) *Route {
	return &Route{id, name, make(map[StopId]struct{})}
}

func (r *Route) AddStop(stop *Stop) {
	r.StopIds[stop.Id] = struct{}{}
}

type RouteData struct {
	Routes       []*Route
	RoutesById   map[RouteId]*Route
	RoutesByName map[RouteName]*Route
}

func NewRouteData(routes []*Route) *RouteData {
	routesById := make(map[RouteId]*Route)
	routesByName := make(map[RouteName]*Route)
	for _, route := range routes {
		routesById[route.Id] = route
		routesByName[route.Name] = route
	}
	return &RouteData{routes, routesById, routesByName}
}
