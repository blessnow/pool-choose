const API_BASE = import.meta.env.VITE_API_BASE || '/api';

interface ApiResponse<T> {
  ok: boolean;
  data?: T;
  error?: string;
  token?: string;  // 登录接口直接返回token
}

// 认证
export async function login(password: string): Promise<ApiResponse<{ token: string }>> {
  try {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password }),
    });
    if (!res.ok) {
      return { ok: false, error: `HTTP错误: ${res.status}` };
    }
    const data = await res.json();
    // API返回 {ok: true, token: "..."}，需要转换格式
    if (data.ok && data.token) {
      return { ok: true, data: { token: data.token } };
    }
    return data;
  } catch (err) {
    console.error('Login error:', err);
    return { ok: false, error: err instanceof Error ? err.message : '网络连接失败' };
  }
}

export async function logout(token: string): Promise<ApiResponse<null>> {
  const res = await fetch(`${API_BASE}/auth/logout`, {
    method: 'POST',
    headers: { Authorization: token },
  });
  return res.json();
}

// 股票
export interface Stock {
  id: number;
  code: string;
  name: string;
  sector: string;
  industry: string;
  isRecommend: boolean;
  entryPrice: number;
  heavyPrice: number;
  targetPrice: number;
  coreLogic: string;
  createdAt: string;
  updatedAt: string;
}

export interface Quote {
  stockCode: string;
  open: number;
  close: number;
  high: number;
  low: number;
  volume: number;
  amount: number;
  change: number;
  quote_time: string;
  prev_close: number;
}

export interface CompanySummary {
  code: string;
  syncedAt: string;
  totalShares: number;
  netAssets: number;
  annualProfit: number;
  latestDividendPer10: number;
  marketCap: number;
  pe: number;
  pb: number;
  roe: number;
  dividendYield: number;
  sector: string;
}

export interface StockWithQuote {
  stock: Stock;
  quote?: Quote;
}

export async function getStocks(token: string): Promise<ApiResponse<StockWithQuote[]>> {
  const res = await fetch(`${API_BASE}/stocks`, {
    headers: { Authorization: token },
  });
  return res.json();
}

export async function getStock(token: string, code: string): Promise<ApiResponse<Stock>> {
  const res = await fetch(`${API_BASE}/stocks/${code}`, {
    headers: { Authorization: token },
  });
  return res.json();
}

export async function createStock(token: string, stock: Partial<Stock>): Promise<ApiResponse<Stock>> {
  const res = await fetch(`${API_BASE}/stocks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: token },
    body: JSON.stringify(stock),
  });
  return res.json();
}

export async function updateStock(token: string, code: string, stock: Partial<Stock>): Promise<ApiResponse<Stock>> {
  const res = await fetch(`${API_BASE}/stocks/${code}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: token },
    body: JSON.stringify(stock),
  });
  return res.json();
}

export async function deleteStock(token: string, code: string): Promise<ApiResponse<null>> {
  const res = await fetch(`${API_BASE}/stocks/${code}`, {
    method: 'DELETE',
    headers: { Authorization: token },
  });
  return res.json();
}

// 行情
export async function getQuotes(token: string, codes: string[]): Promise<ApiResponse<Record<string, Quote>>> {
  const res = await fetch(`${API_BASE}/quotes?codes=${codes.join(',')}`, {
    headers: { Authorization: token },
  });
  return res.json();
}

export async function getCompanySummaries(token: string, codes: string[]): Promise<ApiResponse<Record<string, CompanySummary>>> {
  const res = await fetch(`${API_BASE}/company-summaries?codes=${codes.join(',')}`, {
    headers: { Authorization: token },
  });
  return res.json();
}

export async function getChart(token: string, code: string, period: string = 'daily'): Promise<ApiResponse<any[]>> {
  const res = await fetch(`${API_BASE}/chart/${code}?period=${period}`, {
    headers: { Authorization: token },
  });
  return res.json();
}

// 估值数据
export interface ValuationData {
  current: {
    code: string;
    name: string;
    pe_ttm: number;
    pe: number;
    pb: number;
    roe: number;
    ps: number;
    close: number;
    high: number;
    low: number;
    open: number;
  };
  financial: Array<{
    date: string;
    roe: number;
    bps: number;
  }>;
  kline: Array<{
    date: string;
    open: number;
    close: number;
    high: number;
    low: number;
    volume: number;
  }>;
}

export async function getValuation(token: string, code: string): Promise<ApiResponse<ValuationData>> {
  const res = await fetch(`${API_BASE}/valuation/${code}`, {
    headers: { Authorization: token },
  });
  return res.json();
}

// 估值带数据（PE/PB Band）
export interface ValuationBandPoint {
  date: string;
  rawValue?: number;
  price: number;
  value: number;
  percentile?: number;
  tracks: {
    p90?: number;
    p70?: number;
    p50?: number;
    p30?: number;
    p10?: number;
  };
}

export interface ValuationBandData {
  metric: string;
  years: number;
  points: ValuationBandPoint[];
  updatedAt?: string;
}

export interface ValuationBandResponse {
  band: ValuationBandData;
  updatedAt: string;
  latestDate: string;
}

export async function getValuationBand(
  token: string,
  code: string,
  metric: 'pe_ttm' | 'pb' = 'pe_ttm',
  years: number = 5
): Promise<ApiResponse<ValuationBandResponse>> {
  const res = await fetch(`${API_BASE}/valuation-band?code=${code}&metric=${metric}&years=${years}`, {
    headers: { Authorization: token },
  });
  return res.json();
}

// 周期分析
export interface MacroCard {
  label: string;
  value: string;
  detail: string;
}

export interface CycleBar {
  label: string;
  value: number;
  gradient: string;
  labels: string[];
}

export interface Source {
  label: string;
  url: string;
}

export interface CycleInsight {
  title: string;
  updatedAt: string;
  dataLayer: string;
  opinionLayer: string;
  macroCards: MacroCard[];
  bars: CycleBar[];
  conclusion: string;
  focus: string;
  risk: string;
  sources: Source[];
}

export async function getCycleInsight(token: string): Promise<ApiResponse<CycleInsight>> {
  const res = await fetch(`${API_BASE}/cycle-insight`, {
    headers: { Authorization: token },
  });
  return res.json();
}

export async function updateCycleInsight(token: string, insight: Partial<CycleInsight>): Promise<ApiResponse<CycleInsight>> {
  const res = await fetch(`${API_BASE}/cycle-insight`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: token },
    body: JSON.stringify(insight),
  });
  return res.json();
}

// 持仓记录
export interface Position {
  id: number;
  stockCode: string;
  type: string;
  price: number;
  quantity: number;
  date: string;
  note: string;
}

export async function getPositions(token: string, code?: string): Promise<ApiResponse<Position[]>> {
  const url = code ? `${API_BASE}/positions?code=${code}` : `${API_BASE}/positions`;
  const res = await fetch(url, {
    headers: { Authorization: token },
  });
  return res.json();
}

export async function createPosition(token: string, position: Partial<Position>): Promise<ApiResponse<Position>> {
  const res = await fetch(`${API_BASE}/positions`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: token },
    body: JSON.stringify(position),
  });
  return res.json();
}

export async function updatePosition(token: string, id: number, position: Partial<Position>): Promise<ApiResponse<Position>> {
  const res = await fetch(`${API_BASE}/positions/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: token },
    body: JSON.stringify(position),
  });
  return res.json();
}

export async function deletePosition(token: string, id: number): Promise<ApiResponse<null>> {
  const res = await fetch(`${API_BASE}/positions/${id}`, {
    method: 'DELETE',
    headers: { Authorization: token },
  });
  return res.json();
}