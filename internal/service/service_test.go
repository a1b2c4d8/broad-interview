package service_test

import (
	. "broad-interview/internal/analyzer"
	. "broad-interview/internal/client"
	. "broad-interview/internal/service"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("Service", func() {
	var service *MbtaService
	var routeData *RouteData
	var stopData *StopData
	var analyzer *Analyzer

	BeforeEach(func() {
		// NOTE: I would usually mock a responses from the MbtaClient and never rely upon an external resource.
		// However, given time restraints, I'm using data fetched from the MBTA API.
		client := NewMbtaApiClient()
		service = NewMbtaService(client)
		var err error
		routeData, stopData, err = service.GetRoutesAndStops()
		Expect(err).ShouldNot(HaveOccurred())
		analyzer = NewAnalyzer(routeData, stopData)
	})

	Context("route listing", func() {
		routeCount := 8

		It(fmt.Sprintf("should contain %d routes", routeCount), func() {
			Expect(routeData.Routes).Should(HaveLen(routeCount))
		})
	})

	Context("stop listing", func() {
		stopCount := 117

		It(fmt.Sprintf("should contain %d stops", stopCount), func() {
			Expect(stopData.Stops).Should(HaveLen(stopCount))
		})
	})

	Context("get route with most stops", func() {
		routeName := "Green Line B"
		stopCount := 23

		It(fmt.Sprintf("should return %s", routeName), func() {
			route := analyzer.GetRouteWithMostStops()
			Expect(route.Name).Should(BeEquivalentTo(routeName))
			Expect(route.StopIds).Should(HaveLen(stopCount))
		})
	})

	Context("get route with least stops", func() {
		routeName := "Mattapan Trolley"
		stopCount := 8

		It(fmt.Sprintf("should return %s", routeName), func() {
			route := analyzer.GetRouteWithLeastStops()
			Expect(route.Name).Should(BeEquivalentTo(routeName))
			Expect(route.StopIds).Should(HaveLen(stopCount))
		})
	})

	Context("get connecting stops", func() {
		stopCount := 12

		It(fmt.Sprintf("should return %d stops", stopCount), func() {
			connectingStops := analyzer.GetConnectingStops()
			Expect(connectingStops).Should(HaveLen(stopCount))
		})

		stopName := "Arlington"
		It(fmt.Sprintf("should return %s", stopName), func() {
			connectingStops := analyzer.GetConnectingStops()
			sort.Slice(connectingStops, func(i, j int) bool {
				return connectingStops[i].Stop.Name < connectingStops[j].Stop.Name
			})

			connectingStop := connectingStops[0]
			Expect(connectingStop.Stop.Name).Should(BeEquivalentTo(stopName))

			var routeNames []string
			for _, route := range connectingStop.Routes {
				routeNames = append(routeNames, string(route.Name))
			}
			Expect(routeNames).Should(ConsistOf(
				"Green Line B",
				"Green Line C",
				"Green Line D",
				"Green Line E",
			))
		})
	})

	Context("get route path", func() {
		It(fmt.Sprintf("should return the path from Porter to Wonderland"), func() {
			routePath := analyzer.GetRoutePath(stopData.StopsByName["Porter"], stopData.StopsByName["Wonderland"])

			var routeNames []string
			for _, route := range routePath {
				routeNames = append(routeNames, string(route.Name))
			}
			Expect(routeNames).Should(BeEquivalentTo([]string{
				"Red Line",
				"Green Line B",
				"Blue Line",
			}))
		})
	})
})
