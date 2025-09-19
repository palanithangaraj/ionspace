package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
    "github.com/pulumi/pulumi/sdk/go/pulumi"
    "github.com/pulumi/pulumi-aws/sdk/go/aws/eks"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Prepare the tags that are used for each individual resource so they can be found
		// using the Resource Groups service in the AWS Console
		tags := pulumi.Map{
            "version": pulumi.String(getEnv(ctx, "tags:version", "unknown")),
            "author": pulumi.String(getEnv(ctx, "tags:author", "unknown")),
            "team": pulumi.String(getEnv(ctx, "tags:team", "unknown")),
            "feature": pulumi.String(getEnv(ctx, "tags:feature", "unknown")),
            "region": pulumi.String(getEnv(ctx, "aws:region", "unknown")),
		}

		// Create a VPC for the EKS cluster
		cidrBlock := pulumi.String(getEnv(ctx, "vpc:cidr-block", "unknown"))
		vpcArgs := &ec2.VpcArgs{
			CidrBlock: cidrBlock,
			Tags:      tags,
		}

		vpcName := getEnv(ctx, "vpc:name", "unknown")
        vpc, _ := CreateVpc(ctx, vpcName, vpcArgs)

        // Create an Internet Gateway
    	igw, err := ec2.NewInternetGateway(ctx, "ionspaceInternetGateway", &ec2.InternetGatewayArgs{
    		VpcId: vpc.ID(), // Replace 'vpc.ID()' with the actual ID of your VPC
//     		Tags: pulumi.StringMap{
//     			"Name": pulumi.String("ionspace-gateway"),
//     		},
    	})
    	if err != nil {
    		return err
    	}

    	// Export the Internet Gateway ID
    	ctx.Export("internetGatewayId", igw.ID())

    	routeTable, err := ec2.NewRouteTable(ctx, "myRouteTable", &ec2.RouteTableArgs{
         VpcId: vpc.ID(),
         Routes: ec2.RouteTableRouteArray{
             &ec2.RouteTableRouteArgs{
                 CidrBlock: pulumi.String("0.0.0.0/0"), // Default route for all outbound traffic
                 GatewayId: igw.ID(),                    // Route to the Internet Gateway
             },
         },
//          Tags: pulumi.StringMap{
//              "Name": pulumi.String("my-public-route-table"),
//          },
        })
        if err != nil {
             return err
        }

        ctx.Export("routeTableId", routeTable.ID())

		// Create the required number of subnets
        subnets := pulumi.Map{
            "subnet_ids": pulumi.StringArray{},
        }

        subnetZones := strings.Split(getEnv(ctx, "vpc:subnet-zones", "unknown"), ",")
        subnetIPs := strings.Split(getEnv(ctx, "vpc:subnet-ips", "unknown"), ",")

        for idx, availabilityZone := range subnetZones {
             subnetArgs := &ec2.SubnetArgs{
                 Tags:             tags,
                 VpcId:            vpc.ID(),
                 CidrBlock:         pulumi.String(subnetIPs[idx]),
                 AvailabilityZone: pulumi.StringPtr(availabilityZone),
             }

             subnet, err := ec2.NewSubnet(ctx, fmt.Sprintf("%s-subnet-%d", vpcName, idx), subnetArgs)
             if err != nil {
                 fmt.Println(err.Error())
                 return err
             }

             subnets["subnet_ids"] = append(subnets["subnet_ids"].(pulumi.StringArray), subnet.ID())
        }

        ctx.Export("SUBNET-IDS", subnets["subnet_ids"])
        // Define the EKS Cluster IAM Role
        eksRole, err := CreateIamRole(ctx, "eksClusterRole")
        AttachPolicyToIamRole(ctx, eksRole.Name, "eksClusterPolicyAttachment", "AmazonEKSClusterPolicy")
        // Optionally, attach the AmazonEKSVPCResourceController policy if needed for specific functionalities
        AttachPolicyToIamRole(ctx, eksRole.Name, "eksVpcResourceControllerPolicyAttachment", "AmazonEKSVPCResourceController")
        // Export the role ARN
        ctx.Export("eksClusterRoleArn", eksRole.Arn)

        // Create an EKS cluster
        clusterName := getEnv(ctx, "eks:cluster-name", "unknown")
        enabledClusterLogTypes := strings.Split(getEnv(ctx, "eks:cluster-log-types", "unknown"), ",")

        clusterArgs := &eks.ClusterArgs{
             Name:                   pulumi.String(clusterName),
             Version:                pulumi.String(getEnv(ctx, "eks:k8s-version", "unknown")),
             RoleArn:                pulumi.StringInput(eksRole.Arn),
             Tags:                   tags,
             VpcConfig: &eks.ClusterVpcConfigArgs{
             				SubnetIds: subnets["subnet_ids"].(pulumi.StringArray),
             			},
             EnabledClusterLogTypes: toPulumiStringArray(enabledClusterLogTypes),
        }

        cluster, err := eks.NewCluster(ctx, clusterName, clusterArgs)
        if err != nil {
             fmt.Println(err.Error())
             return err
        }

        ctx.Export("CLUSTER-ID", cluster.ID())

        // Create an IAM Role for the Node Group
        nodeGroupRole, err := CreateIamRole(ctx, "nodeGroupRole")

    		// Attach policies to the Node Group Role
            AttachPolicyToIamRole(ctx, nodeGroupRole.Name, "nodeGroupRolePolicyAttachment1", "AmazonEKSWorkerNodePolicy")
            AttachPolicyToIamRole(ctx, nodeGroupRole.Name, "nodeGroupRolePolicyAttachment2", "AmazonEKS_CNI_Policy")
            AttachPolicyToIamRole(ctx, nodeGroupRole.Name, "nodeGroupRolePolicyAttachment3", "AmazonEC2ContainerRegistryReadOnly")

    		// Create a new Node Group and associate it with the EKS cluster
            _, err = eks.NewNodeGroup(ctx, "ionspaceNodeGroup", &eks.NodeGroupArgs{
    			ClusterName:       cluster.Name,
    			//InstanceTypes:  instanceTypes.(pulumi.StringArray), // Choose your desired instance type
                ScalingConfig: &eks.NodeGroupScalingConfigArgs{
				// Sets the initial desired number of nodes to 2.
				DesiredSize: pulumi.Int(2),
				// Sets the minimum number of nodes to 1.
				MinSize:     pulumi.Int(1),
				// Sets the maximum number of nodes to 5.
				MaxSize:     pulumi.Int(5),
			},
// 			InstanceTypes: pulumi.StringArrayInput{
//                 				pulumi.StringPtr("t3.medium"), // Adjust instance type
//                 			},
    			NodeRoleArn:   nodeGroupRole.Arn,
    			SubnetIds:     subnets["subnet_ids"].(pulumi.StringArray), // Use subnets from your EKS cluster's VPC
    		})
    		if err != nil {
    			return err
    		}

		return nil
	})
}
