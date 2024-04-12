package main

import (
	"banners/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	// temporary line
	fmt.Println(cfg)
}
