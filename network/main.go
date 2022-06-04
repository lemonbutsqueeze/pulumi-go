package main

import (

	// "github.com/aws/aws-sdk-go-v2/config"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// func authenticateAws() {
// 	cfg, err := config.LoadDefaultConfig(context.TODO())

// 	if err != nil {
// 		log.Fatalf("Failed to load configuration, %v", err)
// 	}
// }

func createVpc(ctx *pulumi.Context) (*ec2.Vpc, error) {
	vpc, err := ec2.NewVpc(ctx, "ApplicationVpc", &ec2.VpcArgs{
		CidrBlock:        pulumi.String("10.0.0.0/16"),
		EnableDnsSupport: pulumi.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	return vpc, err
}

// func createSubnets(ctx *pulumi.Context, vpc *ec2.Vpc) {
// 	ec2.NewSubnet(ctx, "PrivateSubnet", &ec2.SubnetArgs{
// 		VpcId: pulumi.Any(vpc.),
// 		CidrBlock: pulumi.String("10.0.0.0/16"),
// 	})
// }

func main() {
	// authenticateAws()
	println("Hello this got run first")

	pulumi.Run(func(ctx *pulumi.Context) error {

		vpc, err := createVpc(ctx)

		if err != nil {
			return err
		}

		println(vpc)

		return nil
	})
}
