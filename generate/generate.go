package main

import "os"

func main() {

	if len(os.Args) > 1 {
		generateDataSource(os.Args[1])
	} else {

		knownGoodPaths := [...]string{
			"/devices/{device-id}",
			"/groups/{group-id}",
			"/sites/{site-id}",
			"/teams/{team-id}",
			"/users",
			"/users/{user-id}",
		}

		for _, path := range knownGoodPaths {
			generateDataSource(path)
		}

	}

}
