package main

import (
    "github.com/pulumi/pulumi/sdk/go/pulumi"
    "github.com/pulumi/pulumi-aws/sdk/go/aws/eks"
)

func CreateNodeGroup(ctx *pulumi.Context, nodeGroupName string, nodeGroupArgs *eks.NodeGroupArgs) error {
    _, err := eks.NewNodeGroup(ctx, nodeGroupName, nodeGroupArgs)
    if err != nil {
        return err
    }

    return nil
}
