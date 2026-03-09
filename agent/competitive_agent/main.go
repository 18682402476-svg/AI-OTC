package main

import (
	"agent/common"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	// SystemPrompt 定义了竞技竞争型智能体的人设
	SystemPrompt = `你是一名注重竞争、关注手续费的竞技型交易智能体。

你的目标不只是增加总资产价值，更是最大化自己成为当日收益排行榜第一的概率。

你完全清楚以下规则：
- 每一笔场外交易都会向买方收取 1% 的手续费
- 每日收取手续费的 90% 会奖励给排名第一的智能体

你会主动寻找机会：
- 鼓励其他智能体交易（例如通过挂出极具吸引力的价格）
- 挂出大概率会被成交的场外订单
- 努力提升平台总交易量，从而增加总奖励池

你不会考虑撤单。

你可以接受不利交易，只要这能帮助你锁定榜首位置、赢得每日手续费奖励，并且坚持挂单，直到完成交易。

你会从策略层面思考其他智能体的行为，努力超越对手，而不只是公平交易。

只输出一个符合以下固定格式的JSON动作：
{
  "action": "动作类型",
  "params": {
    "mon_amount": 100,
    "usdc_expected": 95,
    "order_id": "仅RE_ORDER、BUY_MON_FROM_ORDER、BUY_USDC_FROM_ORDER时需传入"
  },
  "reasoning": "为什么做出这个决定，请用中文简述理由"
}

动作类型说明：
HOLD：不操作（params为空，无需传入order_id）
SELL_MON_FOR_USDC：挂单卖MON（params需包含mon_amount、usdc_expected，无需传入order_id，单次交易数量限制在 50-250 MON）
SELL_USDC_FOR_MON：挂单卖USDC买MON（params需包含usdc_amount、mon_expected，无需传入order_id，单次交易数量限制在 50-250 MON）
BUY_MON_FROM_ORDER：吃MON卖单（params需包含order_id）
BUY_USDC_FROM_ORDER：吃USDC卖单（params需包含order_id）
RE_ORDER：撤销挂单（params需包含order_id）`
)

func main() {
	// 初始化客户端
	mcpClient := common.NewMCPClient("") // 初始不传 Token
	llmClient := common.NewOpenAIStandardClient()

	log.Println("竞技竞争型 Agent 启动中...")

	// 1. 注册 Agent
	addr := os.Getenv("COMP_WALLET_ADDRESS")
	if addr == "" {
		addr = "0x000000000000000000000000000000000000b003"
	}
	name := os.Getenv("COMP_AGENT_NAME")
	if name == "" {
		// name = fmt.Sprintf("CompetitiveAgent_%d", time.Now().Unix())
		// 尝试从本地凭证文件加载已有的名字，避免每次生成新名字导致无法复用凭证
		if creds, err := common.LoadCredentials(); err == nil {
			name = creds.AgentName
			log.Printf("复用本地凭证中的 Agent 名称: %s", name)
		} else {
			name = fmt.Sprintf("CompetitiveAgent_%d", time.Now().Unix())
			log.Printf("未找到本地凭证，生成新 Agent 名称: %s", name)
		}
	}
	token, err := mcpClient.RegisterAgent(name, addr)
	if err != nil {
		log.Fatalf("Agent 注册失败: %v", err)
	}
	log.Printf("Agent 注册成功，Token: %s", token)

	common.StartRatesListener("MON/USDC")
	common.StartOrdersListener(token)

	// 模拟心跳循环
	for {
		log.Println("开始分析排行榜竞争策略与市场机会...")

		// 2. 获取上下文数据
		usdc := os.Getenv("USDC_TOKEN_ADDRESS")
		if usdc == "" {
			usdc = "0x1E85f6e91e5370E91D74196d249ce703E0993fb7"
		}
		data, err := mcpClient.FetchAllMarketData("MON/USDC", usdc)
		if err != nil {
			log.Printf("获取市场数据失败: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// 3. 构建 Prompt
		userPrompt := fmt.Sprintf(`请基于以下参数进行竞技化策略决策：
1. 当前 MON/USDC 价格: %s
2. 24 小时价格变动记录: %s
3. 当前平台所有挂单（分析对手盘）: %s
4. 我目前的挂单状态: %s
5. 自己的可用余额: %s`, 
			data.CurrentPrice, data.History24h, data.ActiveOrders, data.MyOrders, data.MyBalance)

		// 打印完整的 Prompt 内容
		log.Printf("================ 发送给 LLM 的 Prompt ================")
		log.Printf("%s", userPrompt)
		log.Printf("====================================================")

		// 4. 调用 LLM 获取决策
		action, err := llmClient.GetDecision(SystemPrompt, userPrompt)
		if err != nil {
			log.Printf("获取 LLM 决策失败: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// 5. 打印决策结果
		log.Printf("================ 决策结果 ================")
		log.Printf("动作: %s", action.Action)
		log.Printf("理由: %s", action.Reasoning)
		if len(action.Params) > 0 {
			log.Printf("参数: %+v", action.Params)
		}
		log.Printf("==========================================")

		// 解析当前汇率
		var currentRate float64
		var rateResp struct {
			Data struct {
				LatestPrice float64 `json:"latestPrice"`
			} `json:"data"`
		}
		if err := common.ParseJSON(data.CurrentPrice, &rateResp); err == nil {
			currentRate = rateResp.Data.LatestPrice
		}

		// 6. 执行真实决策
		execResult, err := mcpClient.ExecuteAction(action, currentRate)
		if err != nil {
			log.Printf("执行决策失败: %v", err)
		} else {
			log.Printf("执行成功: %s", execResult)
		}

		// 循环间隔
		time.Sleep(1 * time.Minute)
	}
}
