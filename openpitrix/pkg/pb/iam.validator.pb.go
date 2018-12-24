// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: iam.proto

package pb

import regexp "regexp"
import fmt "fmt"
import go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/golang/protobuf/proto"
import math "math"
import _ "github.com/golang/protobuf/ptypes/timestamp"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import _ "github.com/mwitkow/go-proto-validators"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_RequestHeader_UserId = regexp.MustCompile("^[a-z]+$")

func (this *RequestHeader) Validate() error {
	if !_regex_RequestHeader_UserId.MatchString(this.UserId) {
		return go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.UserId))
	}
	return nil
}

var _regex_ResponseHeader_UserId = regexp.MustCompile("^[a-z]+$")
var _regex_ResponseHeader_OwnerPath = regexp.MustCompile("^[a-z]+$")
var _regex_ResponseHeader_AccessPath = regexp.MustCompile("^[a-z]+$")

func (this *ResponseHeader) Validate() error {
	if !_regex_ResponseHeader_UserId.MatchString(this.UserId) {
		return go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.UserId))
	}
	if !_regex_ResponseHeader_OwnerPath.MatchString(this.OwnerPath) {
		return go_proto_validators.FieldError("OwnerPath", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.OwnerPath))
	}
	if !_regex_ResponseHeader_AccessPath.MatchString(this.AccessPath) {
		return go_proto_validators.FieldError("AccessPath", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.AccessPath))
	}
	return nil
}
func (this *Range) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *Bool) Validate() error {
	return nil
}
func (this *String) Validate() error {
	return nil
}

var _regex_Action_ActionId = regexp.MustCompile("^[a-z]+$")
var _regex_Action_ActionName = regexp.MustCompile("^[a-z]+$")
var _regex_Action_Method = regexp.MustCompile("^[a-z]+$")

func (this *Action) Validate() error {
	if !_regex_Action_ActionId.MatchString(this.ActionId) {
		return go_proto_validators.FieldError("ActionId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.ActionId))
	}
	if !_regex_Action_ActionName.MatchString(this.ActionName) {
		return go_proto_validators.FieldError("ActionName", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.ActionName))
	}
	if !_regex_Action_Method.MatchString(this.Method) {
		return go_proto_validators.FieldError("Method", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.Method))
	}
	return nil
}
func (this *ActionList) Validate() error {
	for _, item := range this.Value {
		if item != nil {
			if err := go_proto_validators.CallValidatorIfExists(item); err != nil {
				return go_proto_validators.FieldError("Value", err)
			}
		}
	}
	return nil
}

var _regex_Role_RoleId = regexp.MustCompile("^[a-z]+$")

func (this *Role) Validate() error {
	if !_regex_Role_RoleId.MatchString(this.RoleId) {
		return go_proto_validators.FieldError("RoleId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.RoleId))
	}
	if this.CreateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.CreateTime); err != nil {
			return go_proto_validators.FieldError("CreateTime", err)
		}
	}
	if this.UpdateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.UpdateTime); err != nil {
			return go_proto_validators.FieldError("UpdateTime", err)
		}
	}
	return nil
}
func (this *CreateRoleRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *CreateRoleResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *DeleteRoleRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *DeleteRoleResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *ModifyRoleRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *ModifyRoleResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *GetRoleRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *GetRoleResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *DescribeRolesRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Range != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Range); err != nil {
			return go_proto_validators.FieldError("Range", err)
		}
	}
	return nil
}
func (this *DescribeRolesResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	for _, item := range this.Value {
		if item != nil {
			if err := go_proto_validators.CallValidatorIfExists(item); err != nil {
				return go_proto_validators.FieldError("Value", err)
			}
		}
	}
	return nil
}

var _regex_RoleModuleBinding_BindingId = regexp.MustCompile("^[a-z]+$")
var _regex_RoleModuleBinding_RoleId = regexp.MustCompile("^[a-z]+$")
var _regex_RoleModuleBinding_ModuleId = regexp.MustCompile("^[a-z]+$")

func (this *RoleModuleBinding) Validate() error {
	if !_regex_RoleModuleBinding_BindingId.MatchString(this.BindingId) {
		return go_proto_validators.FieldError("BindingId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.BindingId))
	}
	if !_regex_RoleModuleBinding_RoleId.MatchString(this.RoleId) {
		return go_proto_validators.FieldError("RoleId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.RoleId))
	}
	if !_regex_RoleModuleBinding_ModuleId.MatchString(this.ModuleId) {
		return go_proto_validators.FieldError("ModuleId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.ModuleId))
	}
	if this.CreateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.CreateTime); err != nil {
			return go_proto_validators.FieldError("CreateTime", err)
		}
	}
	if this.UpdateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.UpdateTime); err != nil {
			return go_proto_validators.FieldError("UpdateTime", err)
		}
	}
	return nil
}
func (this *ModifyRoleModuleBindingsRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Binding != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Binding); err != nil {
			return go_proto_validators.FieldError("Binding", err)
		}
	}
	return nil
}
func (this *ModifyRoleModuleBindingsResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}

var _regex_Group_GroupId = regexp.MustCompile("^[a-z]+$")
var _regex_Group_GroupName = regexp.MustCompile("^[a-z]+$")
var _regex_Group_ParentGroupId = regexp.MustCompile("^[a-z]+$")

func (this *Group) Validate() error {
	if !_regex_Group_GroupId.MatchString(this.GroupId) {
		return go_proto_validators.FieldError("GroupId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.GroupId))
	}
	if !_regex_Group_GroupName.MatchString(this.GroupName) {
		return go_proto_validators.FieldError("GroupName", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.GroupName))
	}
	if !_regex_Group_ParentGroupId.MatchString(this.ParentGroupId) {
		return go_proto_validators.FieldError("ParentGroupId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.ParentGroupId))
	}
	if this.CreateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.CreateTime); err != nil {
			return go_proto_validators.FieldError("CreateTime", err)
		}
	}
	if this.UpdateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.UpdateTime); err != nil {
			return go_proto_validators.FieldError("UpdateTime", err)
		}
	}
	return nil
}
func (this *CreateGroupRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *CreateGroupResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *DeleteGroupsRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *DeleteGroupsResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *ModifyGroupRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *ModifyGroupResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *GetGroupRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *GetGroupResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *DescribeGroupsRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Range != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Range); err != nil {
			return go_proto_validators.FieldError("Range", err)
		}
	}
	return nil
}
func (this *DescribeGroupsResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	for _, item := range this.Value {
		if item != nil {
			if err := go_proto_validators.CallValidatorIfExists(item); err != nil {
				return go_proto_validators.FieldError("Value", err)
			}
		}
	}
	return nil
}

var _regex_User_UserId = regexp.MustCompile("^[a-z]+$")
var _regex_User_GroupId = regexp.MustCompile("^[a-z]+$")
var _regex_User_RoleId = regexp.MustCompile("^[a-z]+$")
var _regex_User_UserName = regexp.MustCompile("^[a-z]+$")

func (this *User) Validate() error {
	if !_regex_User_UserId.MatchString(this.UserId) {
		return go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.UserId))
	}
	if !_regex_User_GroupId.MatchString(this.GroupId) {
		return go_proto_validators.FieldError("GroupId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.GroupId))
	}
	if !_regex_User_RoleId.MatchString(this.RoleId) {
		return go_proto_validators.FieldError("RoleId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.RoleId))
	}
	if !_regex_User_UserName.MatchString(this.UserName) {
		return go_proto_validators.FieldError("UserName", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.UserName))
	}
	if this.CreateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.CreateTime); err != nil {
			return go_proto_validators.FieldError("CreateTime", err)
		}
	}
	if this.StatusTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.StatusTime); err != nil {
			return go_proto_validators.FieldError("StatusTime", err)
		}
	}
	if this.UpdateTime != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.UpdateTime); err != nil {
			return go_proto_validators.FieldError("UpdateTime", err)
		}
	}
	return nil
}

var _regex_UserPassword_UserId = regexp.MustCompile("^[a-z]+$")
var _regex_UserPassword_Password = regexp.MustCompile("^[a-z]+$")

func (this *UserPassword) Validate() error {
	if !_regex_UserPassword_UserId.MatchString(this.UserId) {
		return go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.UserId))
	}
	if !_regex_UserPassword_Password.MatchString(this.Password) {
		return go_proto_validators.FieldError("Password", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.Password))
	}
	return nil
}
func (this *CreateUserRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *CreateUserResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *DeleteUsersRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *DeleteUsersResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *ModifyUserRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *ModifyUserResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *GetUserRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}
func (this *GetUserResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Value != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Value); err != nil {
			return go_proto_validators.FieldError("Value", err)
		}
	}
	return nil
}
func (this *DescribeUsersRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if this.Range != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Range); err != nil {
			return go_proto_validators.FieldError("Range", err)
		}
	}
	return nil
}
func (this *DescribeUsersResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	for _, item := range this.Value {
		if item != nil {
			if err := go_proto_validators.CallValidatorIfExists(item); err != nil {
				return go_proto_validators.FieldError("Value", err)
			}
		}
	}
	return nil
}

var _regex_CanDoActionRequest_Method = regexp.MustCompile("^[a-z]+$")

func (this *CanDoActionRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if !_regex_CanDoActionRequest_Method.MatchString(this.Method) {
		return go_proto_validators.FieldError("Method", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.Method))
	}
	return nil
}
func (this *CanDoActionResponse) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	return nil
}

var _regex_GetOwnerPathRequest_Method = regexp.MustCompile("^[a-z]+$")

func (this *GetOwnerPathRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if !_regex_GetOwnerPathRequest_Method.MatchString(this.Method) {
		return go_proto_validators.FieldError("Method", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.Method))
	}
	return nil
}

var _regex_GetAccessPathRequest_Method = regexp.MustCompile("^[a-z]+$")

func (this *GetAccessPathRequest) Validate() error {
	if nil == this.Head {
		return go_proto_validators.FieldError("Head", fmt.Errorf("message must exist"))
	}
	if this.Head != nil {
		if err := go_proto_validators.CallValidatorIfExists(this.Head); err != nil {
			return go_proto_validators.FieldError("Head", err)
		}
	}
	if !_regex_GetAccessPathRequest_Method.MatchString(this.Method) {
		return go_proto_validators.FieldError("Method", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]+$"`, this.Method))
	}
	return nil
}
