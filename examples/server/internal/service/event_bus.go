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

package service

import (
	"examples/server/config"
	"examples/server/internal/constants"
	"github.com/rulego/rulego/api/types"
)

// EventBusEngine 全局核心规则引擎服务
var EventBusEngine *EventBusService

type EventBusService struct {
	config   config.Config
	eventBus types.RuleEngine
}

func InitCoreEngineService(c config.Config) error {
	//sharedNodeCtx, ok := node_pool.DefaultNodePool.Get("core_endpoint_http")
	//if !ok {
	//	return fmt.Errorf("找不到默认http endpoint节点")
	//}
	//node := sharedNodeCtx.GetNode()
	//restEndpoint, ok := node.(*rest.Rest)

	//ruleConfig := rulego.NewConfig(types.WithDefaultPool(), types.WithLogger(logger.Logger), types.WithNetPool(node_pool.DefaultNodePool))
	//s, err := service.NewRuleEngineService(c, ruleConfig, "super")
	//if err != nil {
	//	return err
	//}
	////初始化rulego
	//s.InitRuleGo(logger.Logger, c.DataDir, "super")
	if c.EventBusChainId == "" {
		c.EventBusChainId = constants.KeyDefaultIntegrationChainId
	}
	s, _ := UserRuleEngineServiceImpl.Get(constants.UserSuper)
	engine, _ := s.GetEngine(c.EventBusChainId)
	//if !ok {
	//	return fmt.Errorf("找不到集成流程")
	//}

	coreS := &EventBusService{
		config:   c,
		eventBus: engine,
	}
	//coreS.registerFunc()
	//coreS.loadServeFiles()
	//coreS.startWebsocketServe()
	EventBusEngine = coreS
	return nil
}

//func (s *EventBusService) registerFunc() {
//	var flowsFunc = &funcs.Flows{}
//	var eventFunc = &funcs.Event{}
//
//	action.Functions.Register("$components", flowsFunc.Components)
//	action.Functions.Register("$chains/delete", flowsFunc.DeleteDsl)
//	action.Functions.Register("$chains", flowsFunc.GetList)
//	action.Functions.Register("$chains/execute", flowsFunc.Execute)
//	action.Functions.Register("$chains/get", flowsFunc.Get)
//	action.Functions.Register("$chains/saveBase", flowsFunc.SaveBaseInfo)
//	action.Functions.Register("$chains/saveConfig", flowsFunc.SaveConfiguration)
//	action.Functions.Register("$chains/save", flowsFunc.Save)
//	action.Functions.Register("$chains/postMsg", flowsFunc.PostMsg)
//	action.Functions.Register("$nodes/shared", flowsFunc.NodePoolList)
//	action.Functions.Register("$logs/debugData", eventFunc.GetDebugData)
//	action.Functions.Register("$logs/runs", eventFunc.GetRuns)
//	action.Functions.Register("$logs/runs/delete", eventFunc.DeleteRuns)
//}

//func (s *EventBusService) loadServeFiles() {
//	loadServeFiles(s.config, s.rest)
//}
//func (s *EventBusService) startWebsocketServe() {
//	ws := NewWebsocketServe(s.config, s.rest)
//	_ = ws.Start()
//}
