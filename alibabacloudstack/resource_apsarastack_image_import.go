package alibabacloudstack

import (
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackImageImport() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackImageImportCreate,
		Read:   resourceAlibabacloudStackImageImportRead,
		Update: resourceAlibabacloudStackImageImportUpdate,
		Delete: resourceAlibabacloudStackImageImportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "x86_64",
				ValidateFunc: validation.StringInSlice([]string{
					"x86_64",
					"i386",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Auto",
				ValidateFunc: validation.StringInSlice([]string{
					"Auto",
					"Aliyun",
					"BYOL",
				}, false),
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Ubuntu",
				ValidateFunc: validation.StringInSlice([]string{
					"CentOS",
					"Ubuntu",
					"SUSE",
					"OpenSUSE",
					"Debian",
					"CoreOS",
					"Windows Server 2003",
					"Windows Server 2008",
					"Windows Server 2012",
					"Windows 7",
					"Customized Linux",
					"Others Linux",
				}, false),
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "linux",
				ValidateFunc: validation.StringInSlice([]string{
					"windows",
					"linux",
				}, false),
			},
			"disk_device_mapping": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_image_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
							ForceNew: true,
						},
						"format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"RAW",
								"VHD",
								"qcow2",
							}, false),
						},
						"oss_bucket": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"oss_object": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlibabacloudStackImageImportCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}

	request := ecs.CreateImportImageRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.Architecture = d.Get("architecture").(string)
	request.Description = d.Get("description").(string)
	request.ImageName = d.Get("image_name").(string)
	request.LicenseType = d.Get("license_type").(string)
	request.OSType = d.Get("os_type").(string)
	request.Platform = d.Get("platform").(string)

	diskDeviceMappings := d.Get("disk_device_mapping").([]interface{})
	if diskDeviceMappings != nil && len(diskDeviceMappings) > 0 {
		mappings := make([]ecs.ImportImageDiskDeviceMapping, 0, len(diskDeviceMappings))
		for _, diskDeviceMapping := range diskDeviceMappings {
			mapping := diskDeviceMapping.(map[string]interface{})
			size := strconv.Itoa(mapping["disk_image_size"].(int))
			diskmapping := ecs.ImportImageDiskDeviceMapping{
				DiskImageSize: size,
				Format:        mapping["format"].(string),
				OSSBucket:     mapping["oss_bucket"].(string),
				OSSObject:     mapping["oss_object"].(string),
			}
			mappings = append(mappings, diskmapping)
		}
		request.DiskDeviceMapping = &mappings
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ImportImage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_import_image", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.ImportImageResponse)
	d.SetId(resp.ImageId)
	stateConf := BuildStateConfByTimes([]string{"Waiting"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFunc(d.Id(), []string{"CreateFailed", "UnAvailable"}), 200)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlibabacloudStackImageImportRead(d, meta)
}

func resourceAlibabacloudStackImageImportRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}

	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("image_name", object.ImageName)
	d.Set("description", object.Description)
	d.Set("architecture", object.Architecture)
	d.Set("os_type", object.OSType)
	d.Set("platform", object.Platform)
	d.Set("disk_device_mapping", FlattenImageImportDiskDeviceMappings(object.DiskDeviceMappings.DiskDeviceMapping))

	return nil
}

func resourceAlibabacloudStackImageImportUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlibabacloudStackImageImportRead(d, meta)
}

func resourceAlibabacloudStackImageImportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	return ecsService.deleteImage(d)
}

func FlattenImageImportDiskDeviceMappings(list []ecs.DiskDeviceMapping) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		size, _ := strconv.Atoi(i.Size)
		l := map[string]interface{}{
			"disk_image_size": size,
			"format":          i.Format,
			"oss_bucket":      i.ImportOSSBucket,
			"oss_object":      i.ImportOSSObject,
		}
		result = append(result, l)
	}
	return result
}
