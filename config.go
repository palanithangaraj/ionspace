package main

import (
	"fmt"

    "github.com/pulumi/pulumi/sdk/go/pulumi"
)

// getEnv searches for the requested key in the pulumi context and provides either the value of the key or the fallback.
func getEnv(ctx *pulumi.Context, key string, fallback string) string {
	if value, ok := ctx.GetConfig(key); ok {
	    fmt.Println(value)
		return value
	}
	return fallback
}
