package main

import "fuxi/providee"

func main() {
	p := providee.NewProvidee(1, "gs")
	p.Start()
	p.Wait()
}