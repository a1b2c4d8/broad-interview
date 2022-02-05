package service

import "sort"

type StopId string
type StopName string

type Stop struct {
	Id       StopId
	Name     StopName
	RouteIds map[RouteId]struct{}
}

func NewStop(id StopId, name StopName) *Stop {
	return &Stop{id, name, make(map[RouteId]struct{})}
}

type StopData struct {
	Stops                []*Stop
	StopsById            map[StopId]*Stop
	StopsByName          map[StopName]*Stop
	ConnectingSidsToRids map[StopId]map[RouteId]struct{}
	RidsToConnectingSids map[RouteId]map[StopId]struct{}
}

func NewStopData(routes []*Route, stopsById map[StopId]*Stop) *StopData {
	stops := make([]*Stop, len(stopsById))
	stopsIndex := 0
	stopsByName := make(map[StopName]*Stop)
	connectingSidsToRids := make(map[StopId]map[RouteId]struct{})
	ridsToConnectingSids := make(map[RouteId]map[StopId]struct{})
	processedStopIds := make(map[StopId]struct{})

	for i, route := range routes {
		for sid := range route.StopIds {
			stop := stopsById[sid]
			stop.RouteIds[route.Id] = struct{}{}

			_, processed := processedStopIds[sid]
			if !processed {
				stops[stopsIndex] = stop
				stopsIndex++
				stopsByName[stop.Name] = stop
				processedStopIds[sid] = struct{}{}
			}

			findConnectingSidsToRids(
				sid,
				route,
				routes[(i+1):],
				connectingSidsToRids,
				ridsToConnectingSids,
			)
		}
	}

	return &StopData{
		stops,
		stopsById,
		stopsByName,
		connectingSidsToRids,
		ridsToConnectingSids,
	}
}

func findConnectingSidsToRids(
	sid StopId,
	route *Route,
	otherRoutes []*Route,
	connectingSidsToRids map[StopId]map[RouteId]struct{},
	ridsToConnectingSids map[RouteId]map[StopId]struct{},
) {
	for _, otherRoute := range otherRoutes {
		_, connects := otherRoute.StopIds[sid]
		if connects {
			addConnectingStop(sid, route.Id, otherRoute.Id, connectingSidsToRids, ridsToConnectingSids)
		}
	}
}

func addConnectingStop(
	sid StopId,
	rid1 RouteId,
	rid2 RouteId,
	connectingSidsToRids map[StopId]map[RouteId]struct{},
	ridsToConnectingSids map[RouteId]map[StopId]struct{},
) {
	ridSet, exists := connectingSidsToRids[sid]
	if !exists {
		ridSet = make(map[RouteId]struct{})
		connectingSidsToRids[sid] = ridSet
	}
	ridSet[rid1] = struct{}{}
	ridSet[rid2] = struct{}{}

	sidSet, exists := ridsToConnectingSids[rid1]
	if !exists {
		sidSet = make(map[StopId]struct{})
		ridsToConnectingSids[rid1] = sidSet
	}
	sidSet[sid] = struct{}{}

	sidSet, exists = ridsToConnectingSids[rid2]
	if !exists {
		sidSet = make(map[StopId]struct{})
		ridsToConnectingSids[rid2] = sidSet
	}
	sidSet[sid] = struct{}{}
}

func SortStopIds(stopIds map[StopId]struct{}) []StopId {
	sorted := make([]StopId, len(stopIds))
	index := 0

	for k := range stopIds {
		sorted[index] = k
		index++
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted
}
