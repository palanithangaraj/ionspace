package main

import (
    "github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
)

// Create AWS Route Table
func CreateRouteTable(ctx *pulumi.Context, rtbName string, vpcId pulumi.IDOutput, igwId pulumi.IDOutput) (*ec2.RouteTable, error) {
        rtb, err := ec2.NewRouteTable(ctx, rtbName, &ec2.RouteTableArgs{
             VpcId: vpcId,
             Routes: ec2.RouteTableRouteArray{
                 &ec2.RouteTableRouteArgs{
                     CidrBlock: pulumi.String("0.0.0.0/0"), // Default route for all outbound traffic
                     GatewayId: igwId,                    // Route to the Internet Gateway
                 },
             },
    	})
    	if err != nil {
    		return rtb, err
    	}

        ctx.Export("routeTableId", rtb.ID())
        return rtb, nil
}
