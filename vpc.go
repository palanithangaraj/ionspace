package main

import (
	"fmt"

    "github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
    "github.com/pulumi/pulumi/sdk/go/pulumi"
)

func CreateVpc(ctx *pulumi.Context, vpcName string, vpcArgs *ec2.VpcArgs) (*ec2.Vpc, error) {
		vpc, err := ec2.NewVpc(ctx, vpcName, vpcArgs)
		if err != nil {
			fmt.Println(err.Error())
			return vpc, err
		}

		// Export IDs of the created resources to the Pulumi stack
		ctx.Export("VPC-ID", vpc.ID())
		return vpc, nil
}
