package main

import (
	_ "github.com/bblfsh/php-driver/driver/impl"
	"github.com/bblfsh/php-driver/driver/normalizer"

	"gopkg.in/bblfsh/sdk.v2/sdk/driver"
)

func main() {
	driver.Run(normalizer.Transforms)
}
