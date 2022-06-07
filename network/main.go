package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateVpc(ctx *pulumi.Context) (*ec2.Vpc, error) {
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
	pulumi.Run(func(ctx *pulumi.Context) error {

		vpc, err := CreateVpc(ctx)

		if err != nil {
			return err
		}

		println(vpc)

		return nil
	})
}
