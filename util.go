package main

import (
    "github.com/pulumi/pulumi/sdk/go/pulumi"
)


func toPulumiStringArray(arr []string) pulumi.StringArray {
    result := make(pulumi.StringArray, len(arr))
    for i, v := range arr {
        result[i] = pulumi.String(v)
    }
    return result
}
