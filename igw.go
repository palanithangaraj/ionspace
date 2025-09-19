package main

import (
    "github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
)

// Create AWS InternetGateway
func CreateInternetGateway(ctx *pulumi.Context, igwName string, vpcId pulumi.IDOutput) (*ec2.InternetGateway, error) {
        igw, err := ec2.NewInternetGateway(ctx, igwName, &ec2.InternetGatewayArgs{
    		VpcId: vpcId, // Replace 'vpc.ID()' with the actual ID of your VPC
//     		Tags: pulumi.StringMap{
//     			"Name": pulumi.String("ionspace-gateway"),
//     		},
    	})
    	if err != nil {
    		return igw, err
    	}

    	// Export the Internet Gateway ID
    	ctx.Export("internetGatewayId", igw.ID())
        return igw, nil
}
