package amp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/prometheusservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfamp "github.com/hashicorp/terraform-provider-aws/internal/service/amp"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccAMPLoggingConfiguration_basic(t *testing.T) {
	resourceName := "aws_prometheus_logging_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(prometheusservice.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, prometheusservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoggingConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingConfig_basic(defaultLoggingConfiguration()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "log_group_arn", defaultLoggingConfiguration()),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingConfig_basic(anotherLoggingConfiguration()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "log_group_arn", anotherLoggingConfiguration()),
				),
			},
			{
				Config: testAccLoggingConfig_basic(defaultLoggingConfiguration()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "log_group_arn", defaultLoggingConfiguration()),
				),
			},
		},
	})
}

func TestLoggingConfiguration_disappears(t *testing.T) {
	resourceName := "aws_prometheus_logging_configuration.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(prometheusservice.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, prometheusservice.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoggingConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingConfig_basic(defaultLoggingConfiguration()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingConfigurationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfamp.ResourceLoggingConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckLoggingConfigurationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Prometheus Logging Configuration ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AMPConn

		_, err := tfamp.FindLogGroupByID(context.Background(), conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckLoggingConfigurationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).AMPConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_prometheus_logging_configuration" {
			continue
		}

		_, err := tfamp.FindLogGroupByID(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("prometheus Logging Configuration %s still exists", rs.Primary.ID)
	}

	return nil
}

func defaultLoggingConfiguration() string {
	return `arn:aws:logs:us-west-2:447597502929:log-group:testVL1:*`
}

func anotherLoggingConfiguration() string {
	return `arn:aws:logs:us-west-2:447597502929:log-group:testVL2:*`
}

func testAccLoggingConfig_basic(logGroupArn string) string {
	return fmt.Sprintf(`
resource "aws_prometheus_workspace" "test" {
}
resource "aws_prometheus_logging_configuration" "test" {
  workspace_id = aws_prometheus_workspace.test.id
  log_group_arn = %[1]q
}
`, logGroupArn)
}
