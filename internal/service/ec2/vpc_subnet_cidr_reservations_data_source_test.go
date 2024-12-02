package ec2_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSubnetIpv4CIDRReservationsDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetIpv4CIDRReservationsDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_ec2_subnet_ipv4_cidr_reservations.test", "reservations.#", "1"),
					resource.TestCheckResourceAttr("data.aws_ec2_subnet_ipv4_cidr_reservations.test", "reservations.0.cidr_block", "10.0.0.12/32"),
				),
			},
		},
	})
}

func testAccSubnetIpv4CIDRReservationsDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id     = aws_vpc.test.id
  cidr_block = "10.0.0.0/24"

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_subnet_cidr_reservation" "reservation" {
    subnet_id = aws_subnet.test.id
    cidr_block = "10.0.0.12/32"
    reservation_type = "explicit"
    description = "test reservation"
}
data "aws_ec2_subnet_ipv4_cidr_reservations" "test" {
  subnet_id = aws_subnet.test.id
  depends_on = ["aws_ec2_subnet_cidr_reservation.reservation"]
}
`, rName)
}
