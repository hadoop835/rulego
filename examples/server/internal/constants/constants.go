package constants

const (
	// DirWorkflows 工作流目录
	DirWorkflows        = "workflows"
	DirWorkflowsRun     = "runs"
	DirWorkflowsRule    = "rules"
	DirWorkflowsSubRule = "subrules"
)
const (
	KeyMsgType         = "msgType"
	KeyChainId         = "chainId"
	KeyNodeId          = "nodeId"
	KeyUsername        = "username"
	KeyClientId        = "clientId"
	KeyVarType         = "varType"
	KeySize            = "size"
	KeyPage            = "page"
	KeyId              = "id"
	KeyKeywords        = "keywords"
	KeyType            = "type"
	KeyDeploy          = "start"
	KeyUndeploy        = "stop"
	KeyWebhookSecret   = "webhookSecret"
	KeyIntegrationType = "integrationType"
	// KeyWorkDir 工作目录
	KeyWorkDir = "workDir"
	// KeyDefaultIntegrationChainId 应用集成规则链ID
	KeyDefaultIntegrationChainId = "$event_bus"
)
const (
	UserSuper = "super"
	UserAdmin = "admin"
)
const (
	RuleChainFileSuffix = ".json"
)

//const (
//	DefaultPoolDef = `
//	{
//	  "ruleChain": {
//		"id": "$default_node_pool",
//		"name": "全局共享节点池"
//	  },
//	  "metadata": {
//		"endpoints": [
//		  {
//			"id": "core_endpoint_http",
//			"type": "endpoint/http",
//			"name": "http:9090",
//			"configuration": {
//			  "allowCors": true,
//			  "server": ":9090"
//			}
//		  }
//		]
//	  }
//	}
//`
//)
