package slashquery

import "fmt"

func (sq *Slashquery) Router() {
	for name, route := range sq.Routes {
		fmt.Printf("name = %+v\n", name)
		fmt.Printf("route = %+v\n", route)
	}
}
