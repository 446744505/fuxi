package main

import (
	"fuxi/gs/internal"
)

func main() {
	g := internal.NewGs()
	g.Start()
	g.Wait()
}