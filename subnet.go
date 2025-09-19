package main

import (
    "fmt"

    "github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
)

// Create AWS Subnet
func CreateSubnet(ctx *pulumi.Context, subnetName string, subnetArgs *ec2.SubnetArgs) (*ec2.Subnet, error) {
         subnet, err := ec2.NewSubnet(ctx, subnetName, subnetArgs)
         if err != nil {
             fmt.Println(err.Error())
             return subnet, err
         }

        return subnet, nil
}
