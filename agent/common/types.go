package common

// AgentAction 定义了智能体决策后的 JSON 输出结构
type AgentAction struct {
	Action    string                 `json:"action"`    // 动作类型: HOLD, SELL_MON_FOR_USDC, SELL_USDC_FOR_MON, BUY_MON_FROM_ORDER, BUY_USDC_FROM_ORDER, RE_ORDER
	Params    map[string]interface{} `json:"params"`    // 动作所需的参数
	Reasoning string                 `json:"reasoning"` // 决策理由 (中文)
}

// MarketDataContext 聚合了来自 MCP 的所有上下文数据
type MarketDataContext struct {
	CurrentPrice string `json:"current_price"` // 当前实时汇率
	History24h   string `json:"history_24h"`   // 24 小时价格变动记录
	ActiveOrders string `json:"active_orders"` // 当前平台的所有订单
	MyOrders     string `json:"my_orders"`     // 当前 Agent 自己的挂单
	MyBalance    string `json:"my_balance"`    // 当前 Agent 的账户余额
}
