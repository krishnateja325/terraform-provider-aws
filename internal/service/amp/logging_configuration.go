package amp

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/prometheusservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"log"
)

func ResourceLoggingConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoggingConfigurationCreate,
		ReadContext:   resourceLoggingConfigurationRead,
		UpdateContext: resourceLoggingConfigurationUpdate,
		DeleteContext: resourceLoggingConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"log_group_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoggingConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AMPConn

	workspaceID := d.Get("workspace_id").(string)
	logGrpArn := d.Get("log_group_arn").(string)
	input := &prometheusservice.CreateLoggingConfigurationInput{
		WorkspaceId: aws.String(workspaceID),
		LogGroupArn: aws.String(logGrpArn),
	}

	log.Printf("[DEBUG] Creating Prometheus Logging Configuration: %s", input)
	_, err := conn.CreateLoggingConfigurationWithContext(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating  Prometheus Logging Configuration:(%s): %w", workspaceID, err))
	}

	d.SetId(workspaceID)

	if _, err := waitLoggingConfigurataionCreated(ctx, conn, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for Prometheus Logging Configuration: (%s) create: %w", d.Id(), err))
	}

	return resourceLoggingConfigurationRead(ctx, d, meta)
}

func resourceLoggingConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AMPConn

	input := &prometheusservice.UpdateLoggingConfigurationInput{
		WorkspaceId: aws.String(d.Get("workspace_id").(string)),
		LogGroupArn: aws.String(d.Get("log_group_arn").(string)),
	}

	log.Printf("[DEBUG] Updating Prometheus Logging Configuration: %s", input)
	_, err := conn.UpdateLoggingConfigurationWithContext(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error Updating Prometheus Logging Configuration (%s): %w", d.Id(), err))
	}

	if _, err := waitLoggingConfigurationUpdated(ctx, conn, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for Prometheus Logging Configuration (%s) update: %w", d.Id(), err))
	}

	return resourceLoggingConfigurationRead(ctx, d, meta)
}

func resourceLoggingConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AMPConn

	logConfig, err := FindLogGroupByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Prometheus Logging Configuration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Prometheus Logging Configuration (%s): %w", d.Id(), err))
	}

	d.Set("log_group_arn", logConfig.LogGroupArn)
	d.Set("workspace_id", d.Id())

	return nil
}

func resourceLoggingConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Logging Configuration %s", d.Id())
	conn := meta.(*conns.AWSClient).AMPConn

	_, err := conn.DeleteLoggingConfigurationWithContext(ctx, &prometheusservice.DeleteLoggingConfigurationInput{
		WorkspaceId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, prometheusservice.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Prometheus Logging Configuration (%s): %w", d.Id(), err))
	}

	if _, err := waitLoggingConfigurationDeleted(ctx, conn, d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for Prometheus Logging Configuration (%s) to be deleted: %w", d.Id(), err))
	}

	return nil
}
