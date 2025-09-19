package main

import (
    "fmt"

    "github.com/pulumi/pulumi/sdk/go/pulumi"
    "github.com/pulumi/pulumi-aws/sdk/go/aws/eks"
)

func CreateCluster(ctx *pulumi.Context, clusterName string, clusterArgs *eks.ClusterArgs) (*eks.Cluster, error) {
    cluster, err := eks.NewCluster(ctx, clusterName, clusterArgs)
    if err != nil {
        fmt.Println(err.Error())
        return cluster, err
    }

    ctx.Export("CLUSTER-ID", cluster.ID())
    return cluster, nil
}



