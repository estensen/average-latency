package main

import "fmt"

func querySite(site string) {
	fmt.Println("Querying ", site)
}

func main() {
	sites := []string{"vg.no", "nrk.no", "adressa.no"}

	for _, site := range sites {
		querySite(site)
	}
}
