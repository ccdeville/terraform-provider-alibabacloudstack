package apsarastack

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type QuickbiPublicService struct {
	client *connectivity.ApsaraStackClient
}

func (s *QuickbiPublicService) DescribeQuickBiUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewQuickbiClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "QueryUserInfoByUserId"
	request := map[string]interface{}{
		"UserId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("QuickBI:User", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, ApsaraStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiPublicService) QueryUserInfoByUserId(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewQuickbiClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "QueryUserInfoByUserId"
	request := map[string]interface{}{
		"UserId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("QuickBI:User", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, ApsaraStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiPublicService) DescribeQuickBiUserGroup(id string) (object map[string]interface{}, err error) {
	//var response map[string]interface{}
	//conn, err := s.client.NewQuickbiClient()
	//if err != nil {
	//	return nil, WrapError(err)
	//}
	//action := "ListByUserGroupId"
	//request := map[string]interface{}{
	//	"UserGroupIds": id,
	//}
	//runtime := util.RuntimeOptions{}
	//runtime.SetAutoretry(true)
	//wait := incrementalWait(3*time.Second, 3*time.Second)
	//err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	//	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &runtime)
	//	if err != nil {
	//		if NeedRetry(err) {
	//			wait()
	//			return resource.RetryableError(err)
	//		}
	//		return resource.NonRetryableError(err)
	//	}
	//	return nil
	//})
	//addDebug(action, response, request)
	//if err != nil {
	//	if IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
	//		return object, WrapErrorf(Error(GetNotFoundMessage("QuickBI:User", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	//	}
	//	return object, WrapErrorf(err, DefaultErrorMsg, id, action, ApsaraStackSdkGoERROR)
	//}
	//v, err := jsonpath.Get("$.Result", response)
	//if err != nil {
	//	return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	//}
	//object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiPublicService) DescribeQuickBiWorkspace(id string) (object map[string]interface{}, err error) {
	//var response map[string]interface{}
	//conn, err := s.client.NewQuickbiClient()
	//if err != nil {
	//	return nil, WrapError(err)
	//}
	//action := "QueryWorkspaceUserList"
	//request := map[string]interface{}{
	//	"UserId": id,
	//}
	//runtime := util.RuntimeOptions{}
	//runtime.SetAutoretry(true)
	//wait := incrementalWait(3*time.Second, 3*time.Second)
	//err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	//	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &runtime)
	//	if err != nil {
	//		if NeedRetry(err) {
	//			wait()
	//			return resource.RetryableError(err)
	//		}
	//		return resource.NonRetryableError(err)
	//	}
	//	return nil
	//})
	//addDebug(action, response, request)
	//if err != nil {
	//	if IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
	//		return object, WrapErrorf(Error(GetNotFoundMessage("QuickBI:User", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	//	}
	//	return object, WrapErrorf(err, DefaultErrorMsg, id, action, ApsaraStackSdkGoERROR)
	//}
	//v, err := jsonpath.Get("$.Result", response)
	//if err != nil {
	//	return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	//}
	//object = v.(map[string]interface{})
	return object, nil
}
