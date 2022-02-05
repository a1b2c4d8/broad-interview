# Broad Institute Interview Project

This project is the result of answering the questions for the Broad Institute take home assignment. The code for this project is written in [Go](https://go.dev/).

## Building and running the command

This is the easiest way to build and run the command locally is to use the Go Docker image.
All commands and tests can be run with `bin/main.sh` executable. This will build the Docker image
if needed and run the command in a Docker container based on that image.

### Question 1

Run the command `bin/main.sh list-routes` to list the MBTA routes.

This command relies upon the MBTA API `https://api-v3.mbta.com/routes?filter[type]=0,1` to pre-filter desired routes instead of calling `https://api-v3.mbta.com/routes` and filtering the results.
This is done in order to fetch the minimal amount of data needed from the API server in order to speed up the transmission of the data.
Additionally, if there was a large amount of data, I would have made two separate calls for types `0` and `1` concurrently, then collated the results in order to speed up request time.

### Question 2

Run the command `bin/main.sh examine-routes` in order to:

1. List the subway route with the most stops and the count of its stops.
2. List the subway route with the least stops and the count of its stops.
3. A list of the stops that connect multiple routes, and the routes they connect.

### Question 3

Run the command `bin/main.sh find-route-path STOP1 STOP2` in order to find the routes needed to travel from `STOP1` to `STOP2`.
To get a list of stops, run `bin/main.sh list-stops`.

This command starts with route the starting stop is on and the routes the ending stop is on.
If they share a common route, this is the route chosen. If no common route is found, then the routes that connect to the starting routes are examined recursively.
This recursion continues until one of the ending routes is found. A set of previously used routes is retained to prevent the recursion
from backtracking, which could result in infinite recursion.

## Running Tests

Run the command `bin/main.sh --test` to run the Ginkgo tests.