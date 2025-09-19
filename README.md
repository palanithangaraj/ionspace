# Ionspace Platfom
## AWS EKS Platform 

 This Pulumi bootstraps a minimal Go project that provisions an AWS EKS Platform.
 It uses the Pulumi AWS provider to create the bucket and exports its name as an output.

 ## When to Use This Template

 Use this template when you want to:
 - Get started quickly with Pulumi and Go for AWS Platforms
 - Provision a simple EKS with minimal configuration
 - Learn how Pulumi programs are structured in Go

 ## Providers & Resources

- **AWS Provider** (`github.com/pulumi/pulumi-aws/sdk/v6`)
- **Resource Created**: `https://github.com/palanithangaraj/ionspace`

 ## Outputs

- `EKS Platform`: The unique name/ID of the created resources

 ## Prerequisites

- Go 1.20 or later
- An AWS account with credentials configured in your environment (for example, via the AWS CLI)
- Pulumi CLI installed and authenticated

 ## Getting Started
### Uses Strategy Pattern

 1. Create a new creates a platform using this POC:
   - Modify the pulumi.dev.yaml and AWS config and create the EKS platform on the fly. 

 2. Follow the interactive prompts to set:
    - **Project name** (used in `Pulumi.yaml`)
    - **Project description**
    - **AWS region** (defaults to `us-east-1`)

 3. Preview your changes:
    ```bash
    pulumi preview
    ```

 4. Deploy your stack as per Environment:
    - **Dev**
        ```bash
        pulumi up --stack dev 
        ```
    - **Staging**
      ```bash
      pulumi up --stack stg
      ```
    - **Production**
      ```bash
      pulumi up --stack prd
      ```
 5. Verify the output:
    ```bash
    pulumi stack output bucketName
    ```

 ## Project Layout

    .
    ├── Pulumi.yaml     Pulumi project settings and metadata
    ├── go.mod          Go module definition (updated on `pulumi new`)
    └── main.go         Pulumi program that creates the S3 bucket

 ## Configuration

 The following configuration keys are available:

 - `aws:region` (string): The AWS region where resources will be created. Default: `us-east-1`.

 You can view or update configuration values with:
```bash
pulumi config get aws:region
pulumi config set aws:region us-west-2
```

 ## Next Steps

 - Customize the platform, such as versioning, encryption, tags, etc.
 - Add additional AWS resources (e.g., IAM roles, DynamoDB tables) by importing and using other Pulumi AWS SDK packages.
 - Explore Pulumi’s Go SDK patterns, such as component resources and functions.
 - Integrate infrastructure as code into your CI/CD pipelines.

 ## Getting Help

 If you run into issues or have questions:
 - Check the Pulumi documentation: https://www.pulumi.com/docs/
 - Browse AWS provider reference: https://www.pulumi.com/registry/packages/aws/api-docs/
 - Ask in the Pulumi Community Slack: https://slack.pulumi.com/
 - Open an issue on the template’s repository (if available)

 Enjoy building cloud infrastructure with Go and Pulumi!