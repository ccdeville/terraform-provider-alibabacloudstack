package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApsaraStackEssNotification_basic(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.NotificationConfigurationModel
	resourceId := "apsarastack_ess_notification.default"

	basicMap := map[string]string{
		"notification_types.#": "2",
		"scaling_group_id":     CHECKSET,
		"notification_arn":     CHECKSET,
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssNotification(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEssNotification_update_notification_types(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_types.#": "4",
					}),
				),
			},
			{
				Config: testAccEssNotification_update_scaling_group_id(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccEssNotification_update_notification_arn(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckEssNotificationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	essService := EssService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_ess_notification" {
			continue
		}
		if _, err := essService.DescribeEssNotification(rs.Primary.ID); err != nil {
			if IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("ess notification %s still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccEssNotification(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	data "apsarastack_regions" "default" {
		current = true
	}

	data "apsarastack_account" "default" {
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
	}

	resource "apsarastack_mns_queue" "default"{
		name="${var.name}"
	}
	
	resource "apsarastack_ess_notification" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR"]
		notification_arn = "acs:ess:${data.apsarastack_regions.default.regions.0.id}:${data.apsarastack_account.default.id}:queue/${apsarastack_mns_queue.default.name}"
	}
	`, common, rand)
}

func testAccEssNotification_update_notification_types(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	data "apsarastack_regions" "default" {
		current = true
	}

	data "apsarastack_account" "default" {
	}
	
	resource "apsarastack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
	}

	resource "apsarastack_mns_queue" "default"{
		name="${var.name}"
	}
	
	resource "apsarastack_ess_notification" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR","AUTOSCALING:SCALE_IN_SUCCESS","AUTOSCALING:SCALE_IN_ERROR"]
		notification_arn = "acs:ess:${data.apsarastack_regions.default.regions.0.id}:${data.apsarastack_account.default.id}:queue/${apsarastack_mns_queue.default.name}"
	}
	`, common, rand)
}

func testAccEssNotification_update_scaling_group_id(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	variable "newname" {
		default = "tf-testAccEssNotification_new-%d"
	}

	data "apsarastack_regions" "default" {
		current = true
	}

	data "apsarastack_account" "default" {
	}
	
	resource "apsarastack_ess_scaling_group" "default1" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.newname}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
	}

	resource "apsarastack_mns_queue" "default"{
		name="${var.name}"
	}
	
	resource "apsarastack_ess_notification" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default1.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR","AUTOSCALING:SCALE_IN_SUCCESS","AUTOSCALING:SCALE_IN_ERROR"]
		notification_arn = "acs:ess:${data.apsarastack_regions.default.regions.0.id}:${data.apsarastack_account.default.id}:queue/${apsarastack_mns_queue.default.name}"
	}
	`, common, rand, rand)
}

func testAccEssNotification_update_notification_arn(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	variable "newname" {
		default = "tf-testAccEssNotification-new-%d"
	}

	data "apsarastack_regions" "default" {
		current = true
	}

	data "apsarastack_account" "default" {
	}
	
	resource "apsarastack_ess_scaling_group" "default1" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.newname}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${apsarastack_vswitch.default.id}"]
	}

	resource "apsarastack_mns_queue" "default1"{
		name="${var.newname}"
	}
	
	resource "apsarastack_ess_notification" "default" {
		scaling_group_id = "${apsarastack_ess_scaling_group.default1.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR","AUTOSCALING:SCALE_IN_SUCCESS","AUTOSCALING:SCALE_IN_ERROR"]
		notification_arn = "acs:ess:${data.apsarastack_regions.default.regions.0.id}:${data.apsarastack_account.default.id}:queue/${apsarastack_mns_queue.default1.name}"
	}
	`, common, rand, rand)
}
