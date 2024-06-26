package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSnapshotPolicyCreate,
		Read:   resourceAlibabacloudStackSnapshotPolicyRead,
		Update: resourceAlibabacloudStackSnapshotPolicyUpdate,
		Delete: resourceAlibabacloudStackSnapshotPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"repeat_weekdays": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": tagsSchema(),
			"time_points": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlibabacloudStackSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateCreateAutoSnapshotPolicyRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AutoSnapshotPolicyName = d.Get("name").(string)
	request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateAutoSnapshotPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_snapshot_policy", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.CreateAutoSnapshotPolicyResponse)
	d.SetId(response.AutoSnapshotPolicyId)

	ecsService := EcsService{client}
	if err := ecsService.WaitForSnapshotPolicy(d.Id(), SnapshotPolicyNormal, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlibabacloudStackSnapshotPolicyRead(d, meta)
}

func resourceAlibabacloudStackSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeSnapshotPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.AutoSnapshotPolicyName)
	weekdays, err := convertJsonStringToList(object.RepeatWeekdays)
	if err != nil {
		return WrapError(err)
	}
	d.Set("repeat_weekdays", weekdays)
	d.Set("retention_days", object.RetentionDays)
	timePoints, err := convertJsonStringToList(object.TimePoints)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", ecsService.tagsToMap(object.Tags.Tag))
	d.Set("time_points", timePoints)

	return nil
}

func resourceAlibabacloudStackSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	ecsService := EcsService{client}
	if d.HasChange("tags") {
		if err := ecsService.SetResourceTagsNew(d, "auto_snapshot_policy"); err != nil {
			return WrapError(err)
		}
	}

	request := ecs.CreateModifyAutoSnapshotPolicyExRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AutoSnapshotPolicyId = d.Id()
	if d.HasChange("name") {
		request.AutoSnapshotPolicyName = d.Get("name").(string)
	}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if d.HasChange("repeat_weekdays") {
		request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	}
	if d.HasChange("retention_days") {
		request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	}
	if d.HasChange("time_points") {
		request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyAutoSnapshotPolicyEx(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlibabacloudStackSnapshotPolicyRead(d, meta)
}

func resourceAlibabacloudStackSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteAutoSnapshotPolicyRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AutoSnapshotPolicyId = d.Id()
	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteAutoSnapshotPolicy(request)
		})
		if err != nil {
			if IsExpectedErrors(err, SnapshotPolicyInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	return WrapError(ecsService.WaitForSnapshotPolicy(d.Id(), Deleted, DefaultTimeout))
}
