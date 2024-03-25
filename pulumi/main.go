package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a custom VPC network
		network, err := compute.NewNetwork(ctx, "custom-network", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(false), // We'll define custom subnetwork
		})
		if err != nil {
			return err
		}

		// Create a subnetwork for the GKE cluster
		subnetwork, err := compute.NewSubnetwork(ctx, "custom-subnetwork", &compute.SubnetworkArgs{
			IpCidrRange: pulumi.String("10.2.0.0/16"),   // Example CIDR, adjust as needed
			Network:     network.ID(),                   // Reference to the created custom network
			Region:      pulumi.String("us-central1"), // Example region, adjust as needed
		})
		if err != nil {
			return err
		}

		// Create a GKE cluster attached to the specified subnetwork
		_, err = container.NewCluster(ctx, "gke-cluster", &container.ClusterArgs{
			InitialNodeCount: pulumi.Int(1),
			Network:          network.ID(),
			Subnetwork:       subnetwork.ID(),
			Location:         pulumi.String("us-central1"), // Match the region with the subnetwork
		})
		if err != nil {
			return err
		}

		return nil
	})
}
