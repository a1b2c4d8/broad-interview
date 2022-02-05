package main

import (
	. "broad-interview/internal/analyzer"
	. "broad-interview/internal/client"
	. "broad-interview/internal/service"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	var f func(*MbtaService, []string) error

	if len(os.Args) < 2 {
		log.Printf("ERROR: No command provided")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "list-routes":
		f = listRoutes

	case "list-stops":
		f = listStops

	case "examine-routes":
		f = examineRoutes

	case "find-route-path":
		f = findRoutePath

	default:
		log.Printf("ERROR: Unknown command: %s", command)
		os.Exit(1)
	}

	client := NewMbtaApiClient()
	service := NewMbtaService(client)
	err := f(service, args)
	if err != nil {
		log.Printf("ERROR: command '%s' failed: %s", command, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func listRoutes(service *MbtaService, args []string) error {
	routeData, err := service.GetRoutes()
	if err != nil {
		return err
	}

	fmt.Printf("MBTA Routes:\n")

	routes := routeData.Routes
	sort.Slice(routes, func(i, j int) bool { return routes[i].Name < routes[j].Name })
	for _, r := range routes {
		fmt.Printf("  %s\n", r.Name)
	}

	return nil
}

func listStops(service *MbtaService, args []string) error {
	_, stopData, err := service.GetRoutesAndStops()
	if err != nil {
		return err
	}

	sort.Slice(stopData.Stops, func(i, j int) bool {
		return stopData.Stops[i].Name < stopData.Stops[j].Name
	})

	fmt.Printf("MBTA Stops:\n")
	for _, stop := range stopData.Stops {
		fmt.Printf("  %s\n", stop.Name)
	}

	return nil
}

func examineRoutes(service *MbtaService, args []string) error {
	routeData, stopData, err := service.GetRoutesAndStops()
	if err != nil {
		return err
	}

	analyzer := NewAnalyzer(routeData, stopData)
	routeWithMostStops := analyzer.GetRouteWithMostStops()
	routeWithLeastStops := analyzer.GetRouteWithLeastStops()
	connectingStops := analyzer.GetConnectingStops()

	sort.Slice(connectingStops, func(i, j int) bool {
		return connectingStops[i].Stop.Name < connectingStops[j].Stop.Name
	})

	for _, connectingStop := range connectingStops {
		sort.Slice(connectingStop.Routes, func(i, j int) bool {
			return connectingStop.Routes[i].Name < connectingStop.Routes[j].Name
		})
	}

	fmt.Printf("MBTA Route Examination:\n")
	fmt.Printf(
		"  Route with the most stops: %s (%d stops)\n",
		routeWithMostStops.Name, len(routeWithMostStops.StopIds),
	)
	fmt.Printf(
		"  Route with the least stops: %s (%d stops)\n",
		routeWithLeastStops.Name, len(routeWithLeastStops.StopIds),
	)

	fmt.Printf("\nConnecting Stops:")
	for _, connectingStop := range connectingStops {
		fmt.Printf("\n  %s connects the following routes:\n", connectingStop.Stop.Name)

		// Sort connected route ids by route names
		for _, route := range connectingStop.Routes {
			fmt.Printf("    %s\n", route.Name)
		}
	}

	return nil
}

func findRoutePath(service *MbtaService, args []string) error {
	if len(args) != 2 {
		return errors.New("requires two stops as arguments")
	}

	routeData, stopData, err := service.GetRoutesAndStops()
	if err != nil {
		return err
	}

	startStop, found := stopData.StopsByName[StopName(args[0])]
	if !found {
		return errors.New(fmt.Sprintf("unknown stop: %s", args[0]))
	}

	endStop, found := stopData.StopsByName[StopName(args[1])]
	if !found {
		return errors.New(fmt.Sprintf("unknown stop: %s", args[1]))
	}

	analyzer := NewAnalyzer(routeData, stopData)
	routePath := analyzer.GetRoutePath(startStop, endStop)
	if routePath != nil {
		var routeNames []string
		for _, route := range routePath {
			routeNames = append(routeNames, string(route.Name))
		}
		fmt.Printf("%s to %s -> %s\n", startStop.Name, endStop.Name, strings.Join(routeNames, ", "))
	} else {
		return errors.New(fmt.Sprintf("no route path found between stops %s and %s", startStop.Name, endStop.Name))
	}

	return nil
}
