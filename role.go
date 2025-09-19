package main

import (
	"fmt"

    "github.com/pulumi/pulumi/sdk/go/pulumi"
    "github.com/pulumi/pulumi-aws/sdk/go/aws/iam"
)

func CreateIamEksRole(ctx *pulumi.Context, roleName string) (*iam.Role, error) {
    iamRole, err := iam.NewRole(ctx, roleName, &iam.RoleArgs{
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
            return iamRole, err
        }

        return iamRole, nil
}

func CreateIamEc2Role(ctx *pulumi.Context, roleName string) (*iam.Role, error) {
    iamRole, err := iam.NewRole(ctx, roleName, &iam.RoleArgs{
            AssumeRolePolicy: pulumi.String(fmt.Sprintf(`{
                 "Version": "2012-10-17",
                 "Statement": [
                     {
                         "Effect": "Allow",
                         "Principal": {
                              "Service": "ec2.amazonaws.com"
                         },
                         "Action": "sts:AssumeRole"
                     }
                 ]
            }`)),
        })

        if err != nil {
            return iamRole, err
        }

        return iamRole, nil
}