package main

import "os"

func main() {

	if len(os.Args) > 1 {
		generateDataSource(os.Args[1])
	} else {

		knownGoodPaths := [...]string{
			"/users/{user-id}",
			"/groups/{group-id}",
			"/sites/{site-id}",
		}

		for _, path := range knownGoodPaths {
			generateDataSource(path)
		}

	}

}
