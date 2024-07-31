package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_ec2_subnet_cidr_reservations", name="Subnet CIDR Reservations")
func dataSourceSubnetIpv4CIDRReservations() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSubnetIpv4CIDRReservationsRead,

		Schema: map[string]*schema.Schema{
			names.AttrSubnetID: {
				Type:     schema.TypeString,
				Required: true,
			},
			"reservations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrCIDRBlock: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrDescription: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrOwnerID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_cidr_reservation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSubnetIpv4CIDRReservationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	subnetID := d.Get(names.AttrSubnetID).(string)

	input := &ec2.GetSubnetCidrReservationsInput{
		SubnetId: aws.String(subnetID),
	}

	output, err := conn.GetSubnetCidrReservations(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Subnet CIDR Reservations: %s", err)
	}

	if len(output.SubnetIpv4CidrReservations) == 0 {
		d.SetId(subnetID)
		reservations := make([]map[string]interface{}, 0)
		if err := d.Set("reservations", reservations); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting reservations: %s", err)
		}
		return diags
	}

	// Set a unique ID for the data source based on the subnet ID
	d.SetId(subnetID)

	// Prepare the list of reservations
	reservations := make([]map[string]interface{}, len(output.SubnetIpv4CidrReservations))
	for i, reservation := range output.SubnetIpv4CidrReservations {
		reservations[i] = map[string]interface{}{
			names.AttrCIDRBlock:          aws.ToString(reservation.Cidr),
			names.AttrDescription:        aws.ToString(reservation.Description),
			names.AttrOwnerID:            aws.ToString(reservation.OwnerId),
			"reservation_type":           reservation.ReservationType,
			"subnet_cidr_reservation_id": aws.ToString(reservation.SubnetCidrReservationId),
		}
	}

	// Set the reservations in the schema
	if err := d.Set("reservations", reservations); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting reservations: %s", err)
	}

	return diags
}
