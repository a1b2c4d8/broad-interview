package service

import (
	. "broad-interview/internal/client"
)

type MbtaService struct {
	client *MbtaClient
}

func NewMbtaService(client *MbtaClient) *MbtaService {
	return &MbtaService{client}
}

func (s *MbtaService) GetRoutes() (*RouteData, error) {
	resp, err := s.client.GetRoutes()
	if err != nil {
		return nil, err
	}

	routes := make([]*Route, len(resp.Routes))
	for i, routeData := range resp.Routes {
		routes[i] = s.routeRespToRoute(routeData)
	}

	return NewRouteData(routes), nil
}

func (s *MbtaService) GetRoutesAndStops() (*RouteData, *StopData, error) {
	routeData, err := s.GetRoutes()
	if err != nil {
		return nil, nil, err
	}

	stopsById := make(map[StopId]*Stop)
	for _, route := range routeData.Routes {
		var routeStops []*Stop
		routeStops, err = s.getStopsForRoute(route.Id)
		if err != nil {
			return nil, nil, err
		}

		for _, stop := range routeStops {
			route.StopIds[stop.Id] = struct{}{}
			stopsById[stop.Id] = stop
		}
	}

	return routeData, NewStopData(routeData.Routes, stopsById), nil
}

func (s *MbtaService) getStopsForRoute(routeId RouteId) ([]*Stop, error) {
	var resp *StopsResponse
	resp, err := s.client.GetStopsForRoutes(string(routeId))
	if err != nil {
		return nil, err
	}

	stops := make([]*Stop, len(resp.Stops))
	for i, stopData := range resp.Stops {
		stops[i] = s.stopRespToStop(stopData)
	}

	return stops, nil
}

func (s *MbtaService) routeRespToRoute(data *RouteResponseData) *Route {
	return NewRoute(RouteId(data.Id), RouteName(data.RouteAttributes.Name))
}

func (s *MbtaService) stopRespToStop(data *StopResponseData) *Stop {
	return NewStop(StopId(data.Id), StopName(data.StopAttributes.Name))
}
