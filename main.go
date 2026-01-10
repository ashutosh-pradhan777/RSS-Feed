package main

import (
	"fmt"

	"github.com/ashutosh-pradhan777/RSS-Feed/internal/config"
)

func main() {

	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("%+v\n", conf)

	conf.SetUser("ashu")

	conf, err = config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("%+v\n", conf)
}
