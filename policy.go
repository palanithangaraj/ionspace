package main

import (
    "github.com/pulumi/pulumi/sdk/go/pulumi"
    "github.com/pulumi/pulumi-aws/sdk/go/aws/iam"
)

// Attach the AmazonEKSClusterPolicy to the role
func AttachPolicyToIamRole(ctx *pulumi.Context, roleName pulumi.Input, attachmentName string, policyName string) error {
        _, err := iam.NewRolePolicyAttachment(ctx, attachmentName, &iam.RolePolicyAttachmentArgs{
            Role:      roleName,
            PolicyArn: pulumi.String("arn:aws:iam::aws:policy/"+policyName),
        })
        if err != nil {
            return err
        }

    return nil
}

