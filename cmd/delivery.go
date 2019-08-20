package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	deliveryID string
)

func init() {
	RootCmd.AddCommand(DeliveryCmd)
	DeliveryCmd.AddCommand(GetDestinationsCmd)
	DeliveryCmd.AddCommand(DeleteDestinationCmd)
	DeliveryCmd.AddCommand(CreateSkeletonCmd)

	GetDestinationsCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "ID of the team for the deliveries (required)")
	GetDestinationsCmd.MarkFlagRequired("team-id")

	DeleteDestinationCmd.Flags().StringVarP(&deliveryID, "delivery-id", "d", "", "ID of the destination to be deleted (required)")
	DeleteDestinationCmd.MarkFlagRequired("delivery-id")

}

// DeliveryCmd - Container for holding delivery root and secondary commands
var DeliveryCmd = &cobra.Command{
	Use:   "delivery",
	Short: "Delivery resource",
	Long:  `Delivery resource - access data relating to deliveries and their associations`,
}

// GetDestinationsCmd - Container for holding destinations root and secondary commands
var GetDestinationsCmd = &cobra.Command{
	Use:   "get-destinations",
	Short: "Get Delivery Destinations",
	Long:  `Get the data for destintations in a team`,
	Run: func(cmd *cobra.Command, args []string) {
		ps, e := ion.GetDeliveryDestinations(teamID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint(ps)
	},
}

// DeleteDestinationCmd - Container for holding delete destination root and secondary commands
var DeleteDestinationCmd = &cobra.Command{
	Use:   "delete-destination",
	Short: "Delete Delivery Destination",
	Long:  `Delete a single delivery destination`,
	Run: func(cmd *cobra.Command, args []string) {
		e := ion.DeleteDeliveryDestination(deliveryID, viper.GetString(secretKey))
		if e != nil {
			fmt.Println(e.Error())
		}

		PPrint("delivery marked as deleted")
	},
}

// CreateSkeletonCmd - Container for holding create skeleton root
var CreateSkeletonCmd = &cobra.Command{
	Use:   "create-skeleton",
	Short: "Prints a skeleton JSON",
	Long:  `Prints a skeleton JSON to edit and use with create-destination`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("{")
		fmt.Println("  \"data\":{")
		fmt.Println("    \"id\":\"\",")
		fmt.Println("    \"team_id\":\"\",")
		fmt.Println("    \"location\":\"\",")
		fmt.Println("    \"region\":\"\",")
		fmt.Println("    \"name\":\"\",")
		fmt.Println("    \"type\":\"s3\",")
		fmt.Println("    \"access_key\":\"\",")
		fmt.Println("    \"secret_key\":\"\"")
		fmt.Println("  }")
		fmt.Println("}")
	},
}
