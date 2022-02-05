package analyzer

import (
	. "broad-interview/internal/service"
	. "broad-interview/internal/util"
	"sort"
)

type ConnectingStop struct {
	Stop   *Stop
	Routes []*Route
}

type Analyzer struct {
	routeData *RouteData
	stopData  *StopData
}

func NewAnalyzer(routeData *RouteData, stopData *StopData) *Analyzer {
	return &Analyzer{routeData, stopData}
}

func (a *Analyzer) GetRouteWithMostStops() *Route {
	routes := a.routeData.Routes
	sort.Slice(routes, func(i, j int) bool { return len(routes[i].StopIds) > len(routes[j].StopIds) })
	return routes[0]
}

func (a *Analyzer) GetRouteWithLeastStops() *Route {
	routes := a.routeData.Routes
	sort.Slice(routes, func(i, j int) bool { return len(routes[i].StopIds) < len(routes[j].StopIds) })
	return routes[0]
}

func (a *Analyzer) GetConnectingStops() []*ConnectingStop {
	var connectingStops []*ConnectingStop
	for sid, rids := range a.stopData.ConnectingSidsToRids {
		var routes []*Route

		for rid := range rids {
			routes = append(routes, a.routeData.RoutesById[rid])
		}

		connectingStop := &ConnectingStop{a.stopData.StopsById[sid], routes}
		connectingStops = append(connectingStops, connectingStop)
	}
	return connectingStops
}

func (a *Analyzer) GetRoutePath(
	startStop *Stop,
	endStop *Stop,
) []*Route {
	allRoutePaths := a.findRoutePath(
		NewRouteIdSet(),
		NewRouteIdSetFrom(startStop.RouteIds),
		NewRouteIdSetFrom(endStop.RouteIds),
	)

	if allRoutePaths != nil {
		shortestRoutePath := allRoutePaths[0]
		for _, routePath := range allRoutePaths[1:] {
			shortestRoutePath = a.getShortestRoutePath(shortestRoutePath, routePath)
		}

		var routes []*Route
		for _, rid := range shortestRoutePath {
			routes = append(routes, a.routeData.RoutesById[rid])
		}
		return routes
	}

	return nil
}

func (a *Analyzer) findRoutePath(
	usedRouteIds *RouteIdSet,
	startRoutes *RouteIdSet,
	endRoutes *RouteIdSet,
) [][]RouteId {
	// If all starting routes have been used, stop going down this path... this can't
	// lead to any of the end routes.
	if usedRouteIds.ContainsAll(startRoutes) {
		return nil
	}

	// Sort start routes for deterministic results.
	sortedStartRoutes := startRoutes.Sorted()

	// See if any of the start routes are in the end routes. That ends a path.
	for _, rid := range sortedStartRoutes {
		_, found := endRoutes.Set[rid]
		// Make sure we don't back track... otherwise we can get in an endless loop.
		_, used := usedRouteIds.Set[rid]
		if found && !used {
			// We found a start route that is in the end route set, so return
			// this as the last (or only) route in th path.
			return [][]RouteId{{rid}}
		}
	}

	// If no common route is found, recurse with the route(s) that connect to the starting route(s).
	var allRoutePaths [][]RouteId

	usedRouteIds.AddAll(startRoutes)
	for _, rid := range sortedStartRoutes {
		// Get the stop(s) on this starting route that have connections to other routes.
		connectingStopIds := a.stopData.RidsToConnectingSids[rid]

		// Search the connecting routes for paths to the end route(s).
		connectingRouteIds := NewRouteIdSet()
		for sid := range connectingStopIds {
			connectingRouteIds.AddAllFromMap(a.stopData.ConnectingSidsToRids[sid])
		}
		connectingRouteIds.DeleteAll(usedRouteIds)

		routePaths := a.findRoutePath(usedRouteIds, connectingRouteIds, endRoutes)
		if routePaths != nil {
			// Prepend this starting route to the determined paths.
			for i, routePath := range routePaths {
				routePaths[i] = append([]RouteId{rid}, routePath...)
				allRoutePaths = append(allRoutePaths, routePaths[i])
			}
		}
	}

	return allRoutePaths
}

func (a *Analyzer) getShortestRoutePath(x []RouteId, y []RouteId) []RouteId {
	if len(x) < len(y) {
		return x
	}

	if len(x) > len(y) {
		return y
	}

	// On a length tie, compare entries for deterministic results
	for i, rid := range x {
		if rid < y[i] {
			return x
		}

		if rid > y[i] {
			return y
		}
	}

	// On a tie, just return the first path
	return x
}
