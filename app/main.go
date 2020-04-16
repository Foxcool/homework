package main

import (
	"fmt"
	"os"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var version string

func main() {
	router := routing.New()

	router.Get("/health", func(c *routing.Context) error {
		fmt.Fprintf(c, `{"status": "OK"}`)
		return nil
	})

	router.Get("/", func(c *routing.Context) error {
		fmt.Fprintf(c, `{"Hello world from": "%s"}`, os.Getenv("HOSTNAME"))
		return nil
	})

	router.Get("/version", func(c *routing.Context) error {
		fmt.Fprintf(c, `{"version": "%s"}`, version)
		return nil
	})

	panic(fasthttp.ListenAndServe(":8000", router.HandleRequest))
}
