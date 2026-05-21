import { useState, useEffect, useCallback } from 'react';
import ReactECharts from 'echarts-for-react';
import { Stock, Quote, CompanySummary, getChart, getPositions, getValuationBand, Position, ValuationBandPoint } from '../services/api';

interface StockDetailProps {
  stock: Stock;
  quote?: Quote;
  summary?: CompanySummary;
  token: string;
}

type ChartMode = 'kline' | 'pb_band' | 'pe_band';
type ChartPeriod = 'daily' | 'weekly' | 'monthly';
type RangePreset = '1M' | '3M' | '6M' | '1Y' | '3Y' | 'ALL';

interface CandlestickParams {
  dataIndex: number;
}

function calcMA(period: number, data: { close: number }[]): (number | null)[] {
  const result: (number | null)[] = [];
  let sum = 0;
  for (let i = 0; i < data.length; i++) {
    sum += data[i].close;
    if (i >= period) sum -= data[i - period].close;
    result.push(i + 1 >= period ? +(sum / period).toFixed(2) : null);
  }
  return result;
}

function rangeToZoom(preset: RangePreset, total: number, period: ChartPeriod): [number, number] {
  if (preset === 'ALL' || total === 0) return [0, 100];
  const monthsMap: Record<Exclude<RangePreset, 'ALL'>, number> = { '1M': 1, '3M': 3, '6M': 6, '1Y': 12, '3Y': 36 };
  const bars = period === 'daily' ? 21 : period === 'weekly' ? 4.3 : 1;
  const show = Math.min(total, Math.ceil(monthsMap[preset] * bars));
  return [Math.max(0, ((total - show) / total) * 100), 100];
}

export function StockDetail({ stock, quote, summary, token }: StockDetailProps) {
  const [chartMode, setChartMode] = useState<ChartMode>('kline');
  const [chartPeriod, setChartPeriod] = useState<ChartPeriod>('daily');
  const [rangePreset, setRangePreset] = useState<RangePreset>('1Y');
  const [chartData, setChartData] = useState<{day?: string; date?: string; open: number; close: number; high: number; low: number; volume: number}[]>([]);
  const [bandData, setBandData] = useState<ValuationBandPoint[]>([]);
  const [positions, setPositions] = useState<Position[]>([]);
  const [showPositions, setShowPositions] = useState(false);
  const [expanded] = useState(true);
  const [loading, setLoading] = useState(false);

  const loadChartData = useCallback(async () => {
    setLoading(true);
    try {
      if (chartMode === 'kline') {
        const result = await getChart(token, stock.code, chartPeriod);
        if (result.ok && result.data) {
          setChartData(result.data);
        }
      } else {
        // PE/PB Band模式
        const metric = chartMode === 'pb_band' ? 'pb' : 'pe_ttm';
        const result = await getValuationBand(token, stock.code, metric, 5);
        if (result.ok && result.data) {
          setBandData(result.data.band.points);
        }
      }
    } finally {
      setLoading(false);
    }
  }, [token, stock.code, chartPeriod, chartMode]);

  const loadPositions = useCallback(async () => {
    const result = await getPositions(token, stock.code);
    if (result.ok && result.data) {
      setPositions(result.data);
    }
  }, [token, stock.code]);

  useEffect(() => {
    if (token) {
      loadChartData();
      loadPositions();
    }
  }, [token, loadChartData, loadPositions]);

  const currentPrice = quote?.close || 0;
  const heavyDistance = stock.heavyPrice > 0 ? ((currentPrice - stock.heavyPrice) / stock.heavyPrice) * 100 : 0;
  const targetPercent = stock.targetPrice > 0 && currentPrice > 0 ? ((stock.targetPrice - currentPrice) / currentPrice) * 100 : 0;
  const entryDistance = stock.entryPrice > 0 ? ((currentPrice - stock.entryPrice) / stock.entryPrice) * 100 : null;
  // 关键价位条：现价相对 [重仓价, 止盈价] 的位置；若范围无效则隐藏现价标记
  const priceBarSpan = stock.targetPrice - stock.heavyPrice;
  const currentLeftPct = priceBarSpan > 0
    ? Math.max(0, Math.min(100, ((currentPrice - stock.heavyPrice) / priceBarSpan) * 100))
    : null;

  const pe = summary?.pe || 0;
  const pb = summary?.pb || 0;
  const roe = summary?.roe || 0;
  const dividendYield = summary?.dividendYield || 0;
  const marketCap = summary?.marketCap || 0;

  // K线图配置
  const [zoomStart, zoomEnd] = rangeToZoom(rangePreset, chartData.length, chartPeriod);
  const xAxisCategories = chartData.map(d => d.day || d.date);
  const ma5 = calcMA(5, chartData);
  const ma10 = calcMA(10, chartData);
  const ma20 = calcMA(20, chartData);
  const ma60 = calcMA(60, chartData);

  const klineChartOption = {
    animation: false,
    legend: {
      show: true,
      top: 0,
      data: ['MA5', 'MA10', 'MA20', 'MA60'],
      textStyle: { color: '#b4c2d6', fontSize: 11 },
      itemGap: 12,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
      backgroundColor: 'rgba(10,14,23,0.94)',
      borderColor: 'rgba(245,158,11,0.28)',
      textStyle: { color: '#f8fbff' },
    },
    axisPointer: {
      link: [{ xAxisIndex: 'all' }],
    },
    grid: [
      { left: '10%', right: '8%', top: '8%', height: '55%' },
      { left: '10%', right: '8%', top: '68%', height: '15%' },
    ],
    xAxis: [
      {
        type: 'category',
        data: xAxisCategories,
        boundaryGap: true,
        axisLine: { lineStyle: { color: '#1e293b' } },
        axisTick: { show: false },
        axisLabel: { color: '#b4c2d6', fontSize: 10 },
      },
      {
        type: 'category',
        gridIndex: 1,
        data: xAxisCategories,
        boundaryGap: true,
        axisLine: { lineStyle: { color: '#1e293b' } },
        axisTick: { show: false },
        axisLabel: { show: false },
      },
    ],
    yAxis: [
      {
        scale: true,
        axisLine: { lineStyle: { color: '#1e293b' } },
        axisTick: { show: false },
        axisLabel: { color: '#b4c2d6', fontSize: 10 },
        splitLine: { lineStyle: { color: '#1e293b' } },
      },
      {
        scale: true,
        gridIndex: 1,
        splitNumber: 2,
        axisLine: { lineStyle: { color: '#1e293b' } },
        axisTick: { show: false },
        axisLabel: { show: false },
        splitLine: { lineStyle: { color: '#1e293b' } },
      },
    ],
    dataZoom: [
      {
        type: 'inside',
        xAxisIndex: [0, 1],
        start: zoomStart,
        end: zoomEnd,
      },
      {
        show: true,
        xAxisIndex: [0, 1],
        type: 'slider',
        top: '90%',
        height: 18,
        start: zoomStart,
        end: zoomEnd,
        textStyle: { color: '#b4c2d6' },
        borderColor: '#1e293b',
      },
    ],
    series: [
      {
        name: 'K线',
        type: 'candlestick',
        data: chartData.map(d => [d.open, d.close, d.low, d.high]),
        itemStyle: {
          color: '#ef4444',
          color0: '#22c55e',
          borderColor: '#ef4444',
          borderColor0: '#22c55e',
        },
      },
      { name: 'MA5', type: 'line', data: ma5, smooth: true, symbol: 'none', lineStyle: { color: '#ffd36b', width: 1 } },
      { name: 'MA10', type: 'line', data: ma10, smooth: true, symbol: 'none', lineStyle: { color: '#a78bfa', width: 1 } },
      { name: 'MA20', type: 'line', data: ma20, smooth: true, symbol: 'none', lineStyle: { color: '#22c55e', width: 1 } },
      { name: 'MA60', type: 'line', data: ma60, smooth: true, symbol: 'none', lineStyle: { color: '#ef4444', width: 1 } },
      {
        name: '成交量',
        type: 'bar',
        xAxisIndex: 1,
        yAxisIndex: 1,
        data: chartData.map(d => d.volume),
        itemStyle: {
          color: (params: CandlestickParams) => {
            const d = chartData[params.dataIndex];
            return d && d.close >= d.open ? '#ef4444' : '#22c55e';
          },
        },
      },
    ],
  };

  // PE/PB Band 图配置：单 Y 轴（股价），P10-P90 反推股价作为色块堆叠背景，实际股价线叠加在上方
  const metricLabel = chartMode === 'pe_band' ? 'PE' : 'PB';
  const priceP10 = bandData.map(d => d.priceTracks?.p10 ?? null);
  const priceP30 = bandData.map(d => d.priceTracks?.p30 ?? null);
  const priceP50 = bandData.map(d => d.priceTracks?.p50 ?? null);
  const priceP70 = bandData.map(d => d.priceTracks?.p70 ?? null);
  const priceP90 = bandData.map(d => d.priceTracks?.p90 ?? null);
  const diff = (a: number | null, b: number | null) =>
    a === null || b === null ? null : +(a - b).toFixed(4);
  const diffP30 = priceP30.map((v, i) => diff(v, priceP10[i]));
  const diffP50 = priceP50.map((v, i) => diff(v, priceP30[i]));
  const diffP70 = priceP70.map((v, i) => diff(v, priceP50[i]));
  const diffP90 = priceP90.map((v, i) => diff(v, priceP70[i]));
  const hasPriceTracks = priceP50.some(v => v !== null);

  const [bandZoomStart, bandZoomEnd] = rangeToZoom(rangePreset, bandData.length, 'daily');

  const bandChartOption = bandData.length > 0 ? {
    animation: false,
    legend: {
      show: true,
      top: 0,
      data: ['股价', '低估区(<P30)', '合理偏低(P30-P50)', '合理偏高(P50-P70)', '高估区(>P70)'],
      textStyle: { color: '#b4c2d6', fontSize: 11 },
      itemGap: 10,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
      backgroundColor: 'rgba(10,14,23,0.94)',
      borderColor: 'rgba(245,158,11,0.28)',
      textStyle: { color: '#f8fbff' },
      formatter: (params: { axisValue: string }[]) => {
        const date = params[0].axisValue;
        const point = bandData.find(p => p.date === date);
        if (!point) return date;
        let result = `<div style="font-weight:bold;margin-bottom:4px">${date}</div>`;
        result += `<div>股价: ${point.price.toFixed(2)}</div>`;
        if (point.rawValue != null) {
          result += `<div>${metricLabel}: ${point.rawValue.toFixed(2)}</div>`;
        }
        if (point.percentile != null) {
          result += `<div>分位: ${point.percentile.toFixed(1)}%</div>`;
        }
        if (point.priceTracks?.p50 != null) {
          result += `<div style="margin-top:4px;color:#94a3b8">公允价(P50): ${point.priceTracks.p50.toFixed(2)}</div>`;
        }
        return result;
      },
    },
    grid: { left: '10%', right: '8%', top: '12%', bottom: '15%' },
    xAxis: {
      type: 'category',
      data: bandData.map(d => d.date),
      boundaryGap: false,
      axisLine: { lineStyle: { color: '#1e293b' } },
      axisTick: { show: false },
      axisLabel: { color: '#b4c2d6', fontSize: 10 },
    },
    yAxis: {
      type: 'value',
      scale: true,
      name: '股价',
      nameTextStyle: { color: '#b4c2d6' },
      axisLine: { lineStyle: { color: '#1e293b' } },
      axisTick: { show: false },
      axisLabel: { color: '#b4c2d6', fontSize: 10 },
      splitLine: { lineStyle: { color: '#1e293b' } },
    },
    dataZoom: [
      { type: 'inside', start: bandZoomStart, end: bandZoomEnd },
      { type: 'slider', top: '92%', height: 16, start: bandZoomStart, end: bandZoomEnd, textStyle: { color: '#b4c2d6' }, borderColor: '#1e293b' },
    ],
    series: hasPriceTracks ? [
      // 底线 P10：不显示线，stack 起点
      {
        name: 'P10基线',
        type: 'line',
        data: priceP10,
        stack: 'band',
        symbol: 'none',
        lineStyle: { opacity: 0 },
        areaStyle: { color: 'transparent' },
        showInLegend: false,
        silent: true,
      },
      // P10–P30：低估区
      {
        name: '低估区(<P30)',
        type: 'line',
        data: diffP30,
        stack: 'band',
        symbol: 'none',
        lineStyle: { opacity: 0 },
        areaStyle: { color: 'rgba(34, 197, 94, 0.22)' },
      },
      // P30–P50：合理偏低
      {
        name: '合理偏低(P30-P50)',
        type: 'line',
        data: diffP50,
        stack: 'band',
        symbol: 'none',
        lineStyle: { opacity: 0 },
        areaStyle: { color: 'rgba(134, 239, 172, 0.14)' },
      },
      // P50–P70：合理偏高
      {
        name: '合理偏高(P50-P70)',
        type: 'line',
        data: diffP70,
        stack: 'band',
        symbol: 'none',
        lineStyle: { opacity: 0 },
        areaStyle: { color: 'rgba(251, 191, 36, 0.16)' },
      },
      // P70–P90：高估区
      {
        name: '高估区(>P70)',
        type: 'line',
        data: diffP90,
        stack: 'band',
        symbol: 'none',
        lineStyle: { opacity: 0 },
        areaStyle: { color: 'rgba(248, 113, 113, 0.22)' },
      },
      // P50 中轨虚线参考
      {
        name: 'P50中轨',
        type: 'line',
        data: priceP50,
        symbol: 'none',
        lineStyle: { color: 'rgba(148, 163, 184, 0.55)', type: 'dashed', width: 1 },
        showInLegend: false,
      },
      // 股价线
      {
        name: '股价',
        type: 'line',
        data: bandData.map(d => d.price),
        symbol: 'none',
        lineStyle: { color: '#f8fbff', width: 2 },
        z: 5,
      },
    ] : [
      // 后端没有 priceTracks（缺失财报数据）时，退回到纯股价 + PE/PB 双轴
      {
        name: '股价',
        type: 'line',
        data: bandData.map(d => d.price),
        symbol: 'none',
        lineStyle: { color: '#f8fbff', width: 2 },
      },
    ],
  } : null;

  const chartOption = chartMode === 'kline' ? klineChartOption : bandChartOption;

  return (
    <div className="bg-[rgba(0,0,0,0.2)] max-h-0 overflow-hidden transition-all" style={{ maxHeight: expanded ? '3200px' : '0' }}>
      <div className="px-5 py-4">
        {/* 详情网格 */}
        <div className="grid grid-cols-[repeat(auto-fit,minmax(160px,1fr))] gap-2.5 mb-4">
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">入场价</div>
            <div className="font-mono text-base font-semibold">{stock.entryPrice > 0 ? stock.entryPrice.toFixed(2) : '--'}</div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">今日涨跌</div>
            <div className={`font-mono text-base font-semibold ${(quote?.change ?? 0) > 0 ? 'text-[#ef4444]' : 'text-[#22c55e]'}`}>
              {(quote?.change ?? 0) > 0 ? '+' : ''}{((quote?.change ?? 0) * 100).toFixed(1)}%
            </div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">距入场</div>
            <div className={`font-mono text-base font-semibold ${
              entryDistance === null ? '' : entryDistance < 0 ? 'text-[#ef4444]' : 'text-[#22c55e]'
            }`}>
              {entryDistance === null ? '--' : `${entryDistance > 0 ? '+' : ''}${entryDistance.toFixed(1)}%`}
            </div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">重仓位</div>
            <div className="font-mono text-base font-semibold">{stock.heavyPrice > 0 ? stock.heavyPrice.toFixed(2) : '--'}</div>
            {stock.heavyPrice > 0 && (
              <div className={`font-mono text-xs font-bold ${heavyDistance > 0 ? 'text-[#22c55e]' : 'text-[#ef4444]'}`}>
                {heavyDistance > 0 ? '+' : ''}{heavyDistance.toFixed(1)}%
              </div>
            )}
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">止盈</div>
            <div className="font-mono text-base font-semibold">{stock.targetPrice > 0 ? `${targetPercent.toFixed(1)}%` : '--'}</div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">PE(TTM)</div>
            <div className="font-mono text-base font-semibold">{pe.toFixed(2)}</div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">PB</div>
            <div className="font-mono text-base font-semibold">{pb.toFixed(2)}</div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">股息率</div>
            <div className="font-mono text-base font-semibold">{dividendYield.toFixed(2)}%</div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">ROE</div>
            <div className="font-mono text-base font-semibold">{roe.toFixed(2)}%</div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">市值</div>
            <div className="font-mono text-base font-semibold">{(marketCap / 1e8).toFixed(2)}亿</div>
          </div>
          <div className="bg-white/[0.03] rounded-lg px-3 py-2.5">
            <div className="text-xs text-[#d6e0ee] uppercase tracking-wider mb-1">所属行业</div>
            <div className="text-sm">{stock.industry || summary?.sector || stock.sector}</div>
          </div>
        </div>

        {/* 关键价位定位 */}
        {priceBarSpan > 0 ? (
          <div className="relative h-7 bg-[#1e293b] rounded my-4">
            <div className="absolute top-[-4px] w-[3px] h-9 rounded z-[2]" style={{ left: '0%', background: '#22c55e' }}>
              <span className="absolute top-[-18px] left-1/2 -translate-x-1/2 text-xs font-mono whitespace-nowrap">重仓 {stock.heavyPrice.toFixed(2)}</span>
            </div>
            {stock.entryPrice > 0 && stock.entryPrice > stock.heavyPrice && stock.entryPrice < stock.targetPrice && (
              <div
                className="absolute top-[-4px] w-[3px] h-9 rounded z-[2]"
                style={{ left: `${((stock.entryPrice - stock.heavyPrice) / priceBarSpan) * 100}%`, background: '#f59e0b' }}
              >
                <span className="absolute top-[-18px] left-1/2 -translate-x-1/2 text-xs font-mono whitespace-nowrap">入场 {stock.entryPrice.toFixed(2)}</span>
              </div>
            )}
            {currentLeftPct !== null && currentPrice > 0 && (
              <div className="absolute top-[-4px] w-[3px] h-9 rounded z-[2]" style={{ left: `${currentLeftPct}%`, background: '#f8fbff' }}>
                <span className="absolute top-[-18px] left-1/2 -translate-x-1/2 text-xs font-mono whitespace-nowrap">现价 {currentPrice.toFixed(2)}</span>
              </div>
            )}
            <div className="absolute top-[-4px] w-[3px] h-9 rounded z-[2]" style={{ left: '100%', background: '#a78bfa' }}>
              <span className="absolute top-[-18px] left-1/2 -translate-x-1/2 text-xs font-mono whitespace-nowrap">止盈 {stock.targetPrice.toFixed(2)}</span>
            </div>
          </div>
        ) : (
          <div className="my-4 text-xs text-[#b4c2d6]">
            关键价位定位不可用（缺少重仓价或止盈价）
          </div>
        )}

        {/* 核心逻辑 */}
        <div className="text-sm leading-relaxed text-[#d6e0ee] border-t border-[#1e293b] pt-3 mt-2.5">
          <strong className="text-[#f8fbff]">核心逻辑：</strong>
          {stock.coreLogic}
        </div>

        {/* 图表面板 */}
        <div className="mt-4 p-4 border border-[#1e293b] rounded-xl bg-white/[0.02]">
          <div className="flex justify-between items-start gap-3 mb-3 flex-wrap">
            <div className="text-sm font-bold text-[#f8fbff]">{stock.name} 走势</div>
            <div className="flex gap-2 flex-wrap items-center mt-7">
              {['K线', 'PB Band', 'PE Band'].map((mode) => (
                <button
                  key={mode}
                  onClick={(e) => {
                    e.stopPropagation();
                    setChartMode(mode === 'K线' ? 'kline' : mode === 'PB Band' ? 'pb_band' : 'pe_band');
                  }}
                  className={`px-2.5 py-1.5 rounded-lg text-xs font-mono cursor-pointer transition-all border ${
                    (mode === 'K线' && chartMode === 'kline') ||
                    (mode === 'PB Band' && chartMode === 'pb_band') ||
                    (mode === 'PE Band' && chartMode === 'pe_band')
                      ? 'text-[#f59e0b] border-[#f59e0b] bg-[rgba(245,158,11,0.1)]'
                      : 'bg-white/[0.03] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
                  }`}
                >
                  {mode}
                </button>
              ))}
              {chartMode === 'kline' && ['日线', '周线', '月线'].map((period) => (
                <button
                  key={period}
                  onClick={(e) => {
                    e.stopPropagation();
                    setChartPeriod(period === '日线' ? 'daily' : period === '周线' ? 'weekly' : 'monthly');
                  }}
                  className={`px-2.5 py-1.5 rounded-lg text-xs font-mono cursor-pointer transition-all border ${
                    (period === '日线' && chartPeriod === 'daily') ||
                    (period === '周线' && chartPeriod === 'weekly') ||
                    (period === '月线' && chartPeriod === 'monthly')
                      ? 'text-[#f59e0b] border-[#f59e0b] bg-[rgba(245,158,11,0.1)]'
                      : 'bg-white/[0.03] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
                  }`}
                >
                  {period}
                </button>
              ))}
              <span className="mx-1 text-[#1e293b]">|</span>
              {(['1M', '3M', '6M', '1Y', '3Y', 'ALL'] as RangePreset[]).map((preset) => (
                <button
                  key={preset}
                  onClick={(e) => {
                    e.stopPropagation();
                    setRangePreset(preset);
                  }}
                  className={`px-2 py-1.5 rounded-lg text-xs font-mono cursor-pointer transition-all border ${
                    rangePreset === preset
                      ? 'text-[#f59e0b] border-[#f59e0b] bg-[rgba(245,158,11,0.1)]'
                      : 'bg-white/[0.03] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
                  }`}
                >
                  {preset === 'ALL' ? '全部' : preset}
                </button>
              ))}
            </div>
          </div>

          {loading ? (
            <div className="h-[400px] flex items-center justify-center text-[#b4c2d6]">
              加载中...
            </div>
          ) : chartOption ? (
            <ReactECharts
              option={chartOption}
              style={{ height: '400px' }}
              opts={{ renderer: 'canvas' }}
            />
          ) : (
            <div className="h-[400px] flex items-center justify-center text-[#b4c2d6]">
              暂无数据
            </div>
          )}

          <div className="flex justify-between gap-2.5 flex-wrap mt-2.5 text-xs text-[#d6e0ee]">
            <span>最新价：{currentPrice.toFixed(2)}</span>
            <span>行情时间：{quote?.quote_time || '--'}</span>
          </div>
        </div>

        {/* 持仓记录 */}
        <div className="mt-4">
          <button
            onClick={(e) => {
              e.stopPropagation();
              setShowPositions(!showPositions);
            }}
            className="px-3.5 py-2.5 rounded-lg text-sm cursor-pointer transition-all border bg-white/[0.02] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]"
          >
            持仓记录
          </button>

          {showPositions && (
            <div className="mt-3 p-4 border border-[#1e293b] rounded-xl bg-white/[0.02]">
              {positions.length === 0 ? (
                <div className="text-sm text-[#b4c2d6]">暂无持仓记录</div>
              ) : (
                <table className="w-full text-sm">
                  <thead>
                    <tr className="text-[#f8fbff] bg-white/[0.04]">
                      <th className="px-2.5 py-2 text-center">类型</th>
                      <th className="px-2.5 py-2 text-center">价格</th>
                      <th className="px-2.5 py-2 text-center">数量</th>
                      <th className="px-2.5 py-2 text-center">日期</th>
                      <th className="px-2.5 py-2 text-center">备注</th>
                    </tr>
                  </thead>
                  <tbody>
                    {positions.map((pos) => (
                      <tr key={pos.id} className="text-[#d6e0ee]">
                        <td className="px-2.5 py-2 text-center">{pos.type}</td>
                        <td className="px-2.5 py-2 text-center font-mono">{pos.price.toFixed(2)}</td>
                        <td className="px-2.5 py-2 text-center">{pos.quantity}</td>
                        <td className="px-2.5 py-2 text-center">{pos.date}</td>
                        <td className="px-2.5 py-2 text-center">{pos.note}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}