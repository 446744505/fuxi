package main

import (
	"fuxi/robot/internal"
)

func main() {
	r := internal.NewRobot()
	r.Start()
	r.Wait()
}

