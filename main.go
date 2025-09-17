package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
    "github.com/pulumi/pulumi/sdk/go/pulumi"
    "github.com/pulumi/pulumi-aws/sdk/go/aws/iam"
    "github.com/pulumi/pulumi-aws/sdk/go/aws/eks"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
        getEnv(ctx, "vpc:cidr-block", "unknown")
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
        eksRole, err := iam.NewRole(ctx, "eksClusterRole", &iam.RoleArgs{
            AssumeRolePolicy: pulumi.String(fmt.Sprintf(`{
                 "Version": "2012-10-17",
                 "Statement": [
                     {
                         "Effect": "Allow",
                         "Principal": {
                              "Service": "eks.amazonaws.com"
                         },
                         "Action": "sts:AssumeRole"
                     }
                 ]
            }`)),
        })
        if err != nil {
            return err
        }

        // Attach the AmazonEKSClusterPolicy to the role
        _, err = iam.NewRolePolicyAttachment(ctx, "eksClusterPolicyAttachment", &iam.RolePolicyAttachmentArgs{
            Role:      eksRole.Name,
            PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
        })
        if err != nil {
            return err
        }

        // Optionally, attach the AmazonEKSVPCResourceController policy if needed for specific functionalities
        _, err = iam.NewRolePolicyAttachment(ctx, "eksVpcResourceControllerPolicyAttachment", &iam.RolePolicyAttachmentArgs{
            Role:      eksRole.Name,
            PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"),
        })
        if err != nil {
            return err
        }

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

		return nil
	})
}

// getEnv searches for the requested key in the pulumi context and provides either the value of the key or the fallback.
func getEnv(ctx *pulumi.Context, key string, fallback string) string {
	if value, ok := ctx.GetConfig(key); ok {
	    fmt.Println(value)
		return value
	}
	return fallback
}

func toPulumiStringArray(arr []string) pulumi.StringArray {
    result := make(pulumi.StringArray, len(arr))
    for i, v := range arr {
        result[i] = pulumi.String(v)
    }
    return result
}