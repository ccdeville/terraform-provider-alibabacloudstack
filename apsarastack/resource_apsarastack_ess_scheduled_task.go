package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

func resourceApsaraStackEssScheduledTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEssScheduledTaskCreate,
		Read:   resourceApsaraStackEssScheduledTaskRead,
		Update: resourceApsaraStackEssScheduledTaskUpdate,
		Delete: resourceApsaraStackEssScheduledTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scheduled_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"launch_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scheduled_task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"launch_expiration_time": {
				Type:         schema.TypeInt,
				Default:      600,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 21600),
			},
			"recurrence_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Daily", "Weekly", "Monthly"}, false),
			},
			"recurrence_value": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"recurrence_end_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"task_enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackEssScheduledTaskCreate(d *schema.ResourceData, meta interface{}) error {

	request := buildApsaraStackEssScheduledTaskArgs(d)
	client := meta.(*connectivity.ApsaraStackClient)
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScheduledTask(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ess_scheduled_task", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.CreateScheduledTaskResponse)
	d.SetId(response.ScheduledTaskId)

	return resourceApsaraStackEssScheduledTaskRead(d, meta)
}

func resourceApsaraStackEssScheduledTaskRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)

	client := meta.(*connectivity.ApsaraStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScheduledTask(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scheduled_action", object.ScheduledAction)
	d.Set("launch_time", object.LaunchTime)
	d.Set("scheduled_task_name", object.ScheduledTaskName)
	d.Set("description", object.Description)
	d.Set("launch_expiration_time", object.LaunchExpirationTime)
	d.Set("recurrence_type", object.RecurrenceType)
	d.Set("recurrence_value", object.RecurrenceValue)
	d.Set("recurrence_end_time", object.RecurrenceEndTime)
	d.Set("task_enabled", object.TaskEnabled)

	return nil
}

func resourceApsaraStackEssScheduledTaskUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)

	request := ess.CreateModifyScheduledTaskRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.ScheduledTaskId = d.Id()
	request.LaunchExpirationTime = requests.NewInteger(d.Get("launch_expiration_time").(int))

	if d.HasChange("scheduled_task_name") {
		request.ScheduledTaskName = d.Get("scheduled_task_name").(string)
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("scheduled_action") {
		request.ScheduledAction = d.Get("scheduled_action").(string)
	}

	if d.HasChange("launch_time") {
		request.LaunchTime = d.Get("launch_time").(string)
	}

	if d.HasChange("recurrence_type") || d.HasChange("recurrence_value") || d.HasChange("recurrence_end_time") {
		request.RecurrenceType = d.Get("recurrence_type").(string)
		request.RecurrenceValue = d.Get("recurrence_value").(string)
		request.RecurrenceEndTime = d.Get("recurrence_end_time").(string)
	}

	if d.HasChange("task_enabled") {
		request.TaskEnabled = requests.NewBoolean(d.Get("task_enabled").(bool))
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScheduledTask(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceApsaraStackEssScheduledTaskRead(d, meta)
}

func resourceApsaraStackEssScheduledTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	essService := EssService{client}

	request := ess.CreateDeleteScheduledTaskRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.ScheduledTaskId = d.Id()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScheduledTask(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScheduledTaskId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(essService.WaitForEssScheduledTask(d.Id(), Deleted, DefaultTimeout))
}

func buildApsaraStackEssScheduledTaskArgs(d *schema.ResourceData) *ess.CreateScheduledTaskRequest {
	request := ess.CreateCreateScheduledTaskRequest()

	request.ScheduledAction = d.Get("scheduled_action").(string)
	request.LaunchTime = d.Get("launch_time").(string)

	if v, ok := d.GetOk("task_enabled"); ok {
		request.TaskEnabled = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("scheduled_task_name"); ok && v.(string) != "" {
		request.ScheduledTaskName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("recurrence_type"); ok && v.(string) != "" {
		request.RecurrenceType = v.(string)
	}

	if v, ok := d.GetOk("recurrence_value"); ok && v.(string) != "" {
		request.RecurrenceValue = v.(string)
	}

	if v, ok := d.GetOk("recurrence_end_time"); ok && v.(string) != "" {
		request.RecurrenceEndTime = v.(string)
	}

	if v, ok := d.GetOk("launch_expiration_time"); ok && v.(int) != 0 {
		request.LaunchExpirationTime = requests.NewInteger(v.(int))
	}

	return request
}
