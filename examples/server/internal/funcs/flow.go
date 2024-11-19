/*
 * Copyright 2024 The RuleGo Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package funcs

import (
	"examples/server/config"
	"examples/server/internal/constants"
	"examples/server/internal/service"
	"fmt"
	"github.com/rulego/rulego/utils/str"
	"strconv"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/builtin/processor"
	"github.com/rulego/rulego/components/action"
	"github.com/rulego/rulego/endpoint"
	"github.com/rulego/rulego/node_pool"
	"github.com/rulego/rulego/utils/json"
)

type Flows struct {
}

// Components 创建获取规则引擎节点组件列表路由
func (f *Flows) Components(ctx types.RuleContext, msg types.RuleMsg) {
	nodePool, _ := node_pool.DefaultNodePool.GetAllDef()
	//响应endpoint和节点组件配置表单列表
	list, err := json.Marshal(map[string]interface{}{
		//endpoint组件
		"endpoints": endpoint.Registry.GetComponentForms().Values(),
		//节点组件
		"nodes": rulego.Registry.GetComponentForms().Values(),
		//组件配置内置选项
		"builtins": map[string]interface{}{
			// functions节点组件
			"functions": map[string]interface{}{
				//函数名选项
				"functionName": action.Functions.Names(),
			},
			//endpoints内置路由选项
			"endpoints": map[string]interface{}{
				//in 处理器列表
				"inProcessors": processor.InBuiltins.Names(),
				//in 处理器列表
				"outProcessors": processor.OutBuiltins.Names(),
			},
			//共享节点池
			"nodePool": nodePool,
		},
	})
	if err != nil {
		ctx.TellFailure(msg, err)
	} else {
		msg.Data = string(list)
		ctx.TellSuccess(msg)
	}
}

// GetList 创建流程列表
func (f *Flows) GetList(ctx types.RuleContext, msg types.RuleMsg) {
	username := getUsername(ctx)
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		if list, err := json.Marshal(s.List()); err == nil {
			msg.Data = string(list)
			ctx.TellSuccess(msg)
		} else {
			ctx.TellFailure(msg, err)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}
}

// GetDsl 创建获取指定规则链路由
func (f *Flows) GetDsl(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	username := getUsernameAndChange(ctx, chainId)
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		if def, err := s.GetFromDb(chainId); err == nil {
			msg.Data = string(def)
			ctx.TellSuccess(msg)
		} else {
			ctx.TellFailure(msg, err)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}
}

// SaveDsl 创建保存/更新指定规则链路由
func (f *Flows) SaveDsl(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	nodeId := msg.Metadata.GetValue(constants.KeyNodeId)
	username := getUsernameAndChange(ctx, chainId)
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		if err := s.Save(chainId, nodeId, []byte(msg.Data)); err == nil {
			ctx.TellSuccess(msg)
		} else {
			ctx.TellFailure(msg, err)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}
}

// DeleteDsl 创建删除指定规则链路由
func (f *Flows) DeleteDsl(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	username := getUsername(ctx)
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		if err := s.Delete(chainId); err == nil {
			ctx.TellSuccess(msg)
		} else {
			ctx.TellFailure(msg, err)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}
}

// SaveBaseInfo 保存规则链扩展信息
func (f *Flows) SaveBaseInfo(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	username := getUsernameAndChange(ctx, chainId)
	var req types.RuleChainBaseInfo
	if err := json.Unmarshal([]byte(msg.Data), &req); err != nil {
		ctx.TellFailure(msg, err)
	} else {
		if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
			if err := s.SaveBaseInfo(chainId, req); err != nil {
				ctx.TellFailure(msg, err)
			} else {
				ctx.TellSuccess(msg)
			}
		} else {
			ctx.TellFailure(msg, errUserNotFound(username))
		}
	}
}

// SaveConfiguration 保存规则链配置
func (f *Flows) SaveConfiguration(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	varType := msg.Metadata.GetValue(constants.KeyVarType)
	username := getUsernameAndChange(ctx, chainId)
	var req interface{}
	if err := json.Unmarshal([]byte(msg.Data), &req); err != nil {
		ctx.TellFailure(msg, err)
	} else {
		if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
			if err := s.SaveConfiguration(chainId, varType, req); err != nil {
				ctx.TellFailure(msg, err)
			} else {
				ctx.TellSuccess(msg)
			}
		} else {
			ctx.TellFailure(msg, errUserNotFound(username))
		}
	}
}

// Execute 处理请求，并转发到规则引擎，同步等待规则链执行结果返回给调用方
func (f *Flows) Execute(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	msgType := msg.Metadata.GetValue(constants.KeyMsgType)
	//获取消息类型
	msg.Type = msgType
	username := getUsername(ctx)
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		if err := s.ExecuteAndWait(chainId, msg); err == nil {
			ctx.TellSuccess(msg)
		} else {
			ctx.TellFailure(msg, err)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}
}

// PostMsg 处理请求，并转发到规则引擎
func (f *Flows) PostMsg(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	msgType := msg.Metadata.GetValue(constants.KeyMsgType)
	//获取消息类型
	msg.Type = msgType
	username := getUsername(ctx)
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		if err := s.Execute(chainId, msg); err == nil {
			ctx.TellSuccess(msg)
		} else {
			ctx.TellFailure(msg, err)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}
}
func (f *Flows) NodePoolList(ctx types.RuleContext, msg types.RuleMsg) {
	username := getUsername(ctx)
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		var result = map[string][]*types.RuleNode{}
		var err error
		if s.GetRuleConfig().NetPool != nil {
			result, err = s.GetRuleConfig().NetPool.GetAllDef()
		}
		if err != nil {
			ctx.TellFailure(msg, err)
		}
		if v, err := json.Marshal(result); err == nil {
			msg.Data = string(v)
			ctx.TellSuccess(msg)
		} else {
			ctx.TellFailure(msg, err)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}

}

func errUserNotFound(username string) error {
	return fmt.Errorf("user %s not found", username)
}

func getUsername(ctx types.RuleContext) string {
	usernameValue := ctx.GetContext().Value(constants.KeyUsername)
	u := str.ToString(usernameValue)
	if u == "" {
		return constants.UserAdmin
	}
	return u
}
func getUsernameAndChange(ctx types.RuleContext, chainId string) string {
	u := getUsername(ctx)
	//使用系统超级管理员，修改主规则链
	if chainId == config.C.EventBusChainId {
		return constants.UserSuper
	}
	return u
}

type Event struct {
}

func (f *Event) GetDebugData(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	nodeId := msg.Metadata.GetValue(constants.KeyNodeId)
	username := getUsername(ctx)
	var current = 1
	var pageSize = 20
	currentStr := msg.Metadata.GetValue(constants.KeyPage)
	if i, err := strconv.Atoi(currentStr); err == nil {
		current = i
	}
	pageSizeStr := msg.Metadata.GetValue(constants.KeySize)
	if i, err := strconv.Atoi(pageSizeStr); err == nil {
		pageSize = i
	}
	if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
		page := s.DebugData().GetToPage(chainId, nodeId, pageSize, current)
		if v, err := json.Marshal(page); err != nil {
			ctx.TellFailure(msg, err)
		} else {
			msg.Data = string(v)
			ctx.TellSuccess(msg)
		}
	} else {
		ctx.TellFailure(msg, errUserNotFound(username))
	}
}
func (f *Event) DeleteRuns(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	id := msg.Metadata.GetValue(constants.KeyId)
	username := getUsername(ctx)

	if err := service.EventServiceImpl.Delete(username, chainId, id); err != nil {
		ctx.TellFailure(msg, err)
	} else {
		ctx.TellSuccess(msg)
	}
}
func (f *Event) GetRuns(ctx types.RuleContext, msg types.RuleMsg) {
	chainId := msg.Metadata.GetValue(constants.KeyChainId)
	id := msg.Metadata.GetValue(constants.KeyId)
	username := getUsername(ctx)
	var result interface{}
	if id == "" {
		var current = 1
		var pageSize = 20
		currentStr := msg.Metadata.GetValue(constants.KeyPage)
		if i, err := strconv.Atoi(currentStr); err == nil {
			current = i
		}
		pageSizeStr := msg.Metadata.GetValue(constants.KeySize)
		if i, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = i
		}
		if v, total, err := service.EventServiceImpl.List(username, chainId, current, pageSize); err != nil {
			ctx.TellFailure(msg, err)
		} else {
			result = map[string]interface{}{
				"total": total,
				"data":  v,
			}
		}
	} else {
		if v, err := service.EventServiceImpl.Get(username, chainId, id); err != nil {
			ctx.TellFailure(msg, err)
		} else {
			result = v
		}
	}

	if v, err := json.Marshal(result); err != nil {
		ctx.TellFailure(msg, err)
	} else {
		msg.Data = string(v)
		ctx.TellSuccess(msg)
	}
}
