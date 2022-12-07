package main

import "fuxi/map/internal"

func main() {
	m := internal.NewMap()
	m.Start()
	m.Wait()
}
