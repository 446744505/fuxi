package main

import "fuxi/providee"

func main() {
	prov := providee.NewProvidee(1, "gs")
	prov.Start()
}