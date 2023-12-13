package main

func main() {
	knownGoodPaths := [...]string {
		"/users/{user-id}",
		"/groups/{group-id}",
	}

	for _, path := range knownGoodPaths {
		generateDataSource(path)
	}

}
