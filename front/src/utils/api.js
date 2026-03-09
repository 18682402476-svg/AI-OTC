// API请求封装

const API_BASE_URL = 'http://************/api';
const AGENT1_TOKEN = '2b******************************52';

// 通用请求方法
const request = async (url, options = {}) => {
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    const data = await response.json();

    if (!data.success) {
      throw new Error(data.error || '请求失败');
    }

    return data.data;
  } catch (error) {
    console.error('API请求错误:', error);
    throw error;
  }
};

// 市场概览接口
export const marketApi = {
  // 获取当前汇率
  getCurrentRate: async (tradingPair = 'MON/USDC') => {
    return request(`/rate/current?tradingPair=${encodeURIComponent(tradingPair)}`);
  },

  // 获取市场数据
  getMarketData: async () => {
    return request('/rate/market-data');
  },
};

// Agent相关接口
export const agentApi = {
  // 获取Agent排行榜
  getAgentRanking: async (limit = 10) => {
    return request(`/agent/ranking?limit=${limit}`);
  },

  // 获取Agent详情
  getAgentDetail: async (agentId) => {
    return request(`/agent/detail/${agentId}`);
  },
};

// 交易市场接口
export const tradingApi = {
  // 获取活跃订单列表
  getActiveOrders: async () => {
    return request('/order/active', {
      headers: {
        'Authorization': `Bearer ${AGENT1_TOKEN}`
      }
    });
  },

  // 获取买入订单列表
  getBuyOrders: async () => {
    return request('/order/buy');
  },

  // 获取卖出订单列表
  getSellOrders: async () => {
    return request('/order/sell');
  },
};

// 交易记录接口
export const transactionApi = {
  // 获取最新交易记录
  getLatestTransactions: async (limit = 10) => {
    return request(`/api/transaction/latest?limit=${limit}`);
  },

  // 获取Agent的AI推理分析记录
  getAIReasoning: async (agentId, startTime, endTime, limit = 50) => {
    let url = `/transaction/ai-reasoning?agentId=${agentId}&limit=${limit}`;
    if (startTime) url += `&startTime=${encodeURIComponent(startTime)}`;
    if (endTime) url += `&endTime=${encodeURIComponent(endTime)}`;
    return request(url);
  },
};

export default {
  market: marketApi,
  agent: agentApi,
  trading: tradingApi,
  transaction: transactionApi,
};