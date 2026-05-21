import { useState, useEffect, useCallback } from 'react';
import ReactECharts from 'echarts-for-react';
import {
  getChart,
  getValuation,
  getValuationBand,
  ValuationBandPoint,
  ValuationData,
} from '../services/api';

interface StockPreviewProps {
  token: string;
  code: string; // 6 位股票代码；非 6 位时显示空态
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

const isValidCode = (c: string) => /^\d{6}$/.test(c);

export function StockPreview({ token, code }: StockPreviewProps) {
  const [chartMode, setChartMode] = useState<ChartMode>('kline');
  const [chartPeriod, setChartPeriod] = useState<ChartPeriod>('daily');
  const [rangePreset, setRangePreset] = useState<RangePreset>('1Y');
  const [chartData, setChartData] = useState<{day?: string; date?: string; open: number; close: number; high: number; low: number; volume: number}[]>([]);
  const [bandData, setBandData] = useState<ValuationBandPoint[]>([]);
  const [valuation, setValuation] = useState<ValuationData['current'] | null>(null);
  const [loading, setLoading] = useState(false);
  const [valLoading, setValLoading] = useState(false);

  // 当 code 切换时重置图表状态、并拉取估值卡数据
  useEffect(() => {
    setChartData([]);
    setBandData([]);
    setValuation(null);
    if (!isValidCode(code)) return;
    let cancelled = false;
    setValLoading(true);
    getValuation(token, code)
      .then((r) => {
        if (!cancelled && r.ok && r.data) setValuation(r.data.current);
      })
      .finally(() => {
        if (!cancelled) setValLoading(false);
      });
    return () => { cancelled = true; };
  }, [token, code]);

  // 切换模式/周期时按需拉数
  const loadChartData = useCallback(async () => {
    if (!isValidCode(code)) return;
    setLoading(true);
    try {
      if (chartMode === 'kline') {
        const r = await getChart(token, code, chartPeriod);
        if (r.ok && r.data) setChartData(r.data);
      } else {
        const metric = chartMode === 'pb_band' ? 'pb' : 'pe_ttm';
        const r = await getValuationBand(token, code, metric, 5);
        if (r.ok && r.data) setBandData(r.data.band.points);
      }
    } finally {
      setLoading(false);
    }
  }, [token, code, chartMode, chartPeriod]);

  useEffect(() => {
    loadChartData();
  }, [loadChartData]);

  if (!isValidCode(code)) {
    return (
      <div className="h-full flex items-center justify-center text-sm text-[#b4c2d6] border border-dashed border-[#1e293b] rounded-xl">
        输入或选择一只股票后预览图表
      </div>
    );
  }

  // ===== K 线 option =====
  const [zoomStart, zoomEnd] = rangeToZoom(rangePreset, chartData.length, chartPeriod);
  const xAxisCategories = chartData.map(d => d.day || d.date);
  const ma5 = calcMA(5, chartData);
  const ma10 = calcMA(10, chartData);
  const ma20 = calcMA(20, chartData);
  const ma60 = calcMA(60, chartData);

  const klineChartOption = {
    animation: false,
    legend: {
      show: true, top: 0,
      data: ['MA5', 'MA10', 'MA20', 'MA60'],
      textStyle: { color: '#b4c2d6', fontSize: 11 },
      itemGap: 12,
    },
    tooltip: {
      trigger: 'axis', axisPointer: { type: 'cross' },
      backgroundColor: 'rgba(10,14,23,0.94)', borderColor: 'rgba(245,158,11,0.28)',
      textStyle: { color: '#f8fbff' },
    },
    axisPointer: { link: [{ xAxisIndex: 'all' }] },
    grid: [
      { left: '10%', right: '8%', top: '8%', height: '55%' },
      { left: '10%', right: '8%', top: '68%', height: '15%' },
    ],
    xAxis: [
      { type: 'category', data: xAxisCategories, boundaryGap: true, axisLine: { lineStyle: { color: '#1e293b' } }, axisTick: { show: false }, axisLabel: { color: '#b4c2d6', fontSize: 10 } },
      { type: 'category', gridIndex: 1, data: xAxisCategories, boundaryGap: true, axisLine: { lineStyle: { color: '#1e293b' } }, axisTick: { show: false }, axisLabel: { show: false } },
    ],
    yAxis: [
      { scale: true, axisLine: { lineStyle: { color: '#1e293b' } }, axisTick: { show: false }, axisLabel: { color: '#b4c2d6', fontSize: 10 }, splitLine: { lineStyle: { color: '#1e293b' } } },
      { scale: true, gridIndex: 1, splitNumber: 2, axisLine: { lineStyle: { color: '#1e293b' } }, axisTick: { show: false }, axisLabel: { show: false }, splitLine: { lineStyle: { color: '#1e293b' } } },
    ],
    dataZoom: [
      { type: 'inside', xAxisIndex: [0, 1], start: zoomStart, end: zoomEnd },
      { show: true, xAxisIndex: [0, 1], type: 'slider', top: '90%', height: 18, start: zoomStart, end: zoomEnd, textStyle: { color: '#b4c2d6' }, borderColor: '#1e293b' },
    ],
    series: [
      { name: 'K线', type: 'candlestick', data: chartData.map(d => [d.open, d.close, d.low, d.high]), itemStyle: { color: '#ef4444', color0: '#22c55e', borderColor: '#ef4444', borderColor0: '#22c55e' } },
      { name: 'MA5', type: 'line', data: ma5, smooth: true, symbol: 'none', lineStyle: { color: '#ffd36b', width: 1 } },
      { name: 'MA10', type: 'line', data: ma10, smooth: true, symbol: 'none', lineStyle: { color: '#a78bfa', width: 1 } },
      { name: 'MA20', type: 'line', data: ma20, smooth: true, symbol: 'none', lineStyle: { color: '#22c55e', width: 1 } },
      { name: 'MA60', type: 'line', data: ma60, smooth: true, symbol: 'none', lineStyle: { color: '#ef4444', width: 1 } },
      {
        name: '成交量', type: 'bar', xAxisIndex: 1, yAxisIndex: 1,
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

  // ===== PE/PB Band option =====
  const metricLabel = chartMode === 'pe_band' ? 'PE' : 'PB';
  const priceP10 = bandData.map(d => d.priceTracks?.p10 ?? null);
  const priceP30 = bandData.map(d => d.priceTracks?.p30 ?? null);
  const priceP50 = bandData.map(d => d.priceTracks?.p50 ?? null);
  const priceP70 = bandData.map(d => d.priceTracks?.p70 ?? null);
  const priceP90 = bandData.map(d => d.priceTracks?.p90 ?? null);
  const diff = (a: number | null, b: number | null) => a === null || b === null ? null : +(a - b).toFixed(4);
  const diffP30 = priceP30.map((v, i) => diff(v, priceP10[i]));
  const diffP50 = priceP50.map((v, i) => diff(v, priceP30[i]));
  const diffP70 = priceP70.map((v, i) => diff(v, priceP50[i]));
  const diffP90 = priceP90.map((v, i) => diff(v, priceP70[i]));
  const hasPriceTracks = priceP50.some(v => v !== null);

  const [bandZoomStart, bandZoomEnd] = rangeToZoom(rangePreset, bandData.length, 'daily');

  const bandChartOption = bandData.length > 0 ? {
    animation: false,
    legend: {
      show: true, top: 0,
      data: ['股价', '低估区(<P30)', '合理偏低(P30-P50)', '合理偏高(P50-P70)', '高估区(>P70)'],
      textStyle: { color: '#b4c2d6', fontSize: 11 }, itemGap: 10,
    },
    tooltip: {
      trigger: 'axis', axisPointer: { type: 'cross' },
      backgroundColor: 'rgba(10,14,23,0.94)', borderColor: 'rgba(245,158,11,0.28)',
      textStyle: { color: '#f8fbff' },
      formatter: (params: { axisValue: string }[]) => {
        const date = params[0].axisValue;
        const point = bandData.find(p => p.date === date);
        if (!point) return date;
        let result = `<div style="font-weight:bold;margin-bottom:4px">${date}</div>`;
        result += `<div>股价: ${point.price.toFixed(2)}</div>`;
        if (point.rawValue != null) result += `<div>${metricLabel}: ${point.rawValue.toFixed(2)}</div>`;
        if (point.percentile != null) result += `<div>分位: ${point.percentile.toFixed(1)}%</div>`;
        if (point.priceTracks?.p50 != null) result += `<div style="margin-top:4px;color:#94a3b8">公允价(P50): ${point.priceTracks.p50.toFixed(2)}</div>`;
        return result;
      },
    },
    grid: { left: '10%', right: '8%', top: '12%', bottom: '15%' },
    xAxis: { type: 'category', data: bandData.map(d => d.date), boundaryGap: false, axisLine: { lineStyle: { color: '#1e293b' } }, axisTick: { show: false }, axisLabel: { color: '#b4c2d6', fontSize: 10 } },
    yAxis: { type: 'value', scale: true, name: '股价', nameTextStyle: { color: '#b4c2d6' }, axisLine: { lineStyle: { color: '#1e293b' } }, axisTick: { show: false }, axisLabel: { color: '#b4c2d6', fontSize: 10 }, splitLine: { lineStyle: { color: '#1e293b' } } },
    dataZoom: [
      { type: 'inside', start: bandZoomStart, end: bandZoomEnd },
      { type: 'slider', top: '92%', height: 16, start: bandZoomStart, end: bandZoomEnd, textStyle: { color: '#b4c2d6' }, borderColor: '#1e293b' },
    ],
    series: hasPriceTracks ? [
      { name: 'P10基线', type: 'line', data: priceP10, stack: 'band', symbol: 'none', lineStyle: { opacity: 0 }, areaStyle: { color: 'transparent' }, showInLegend: false, silent: true },
      { name: '低估区(<P30)', type: 'line', data: diffP30, stack: 'band', symbol: 'none', lineStyle: { opacity: 0 }, areaStyle: { color: 'rgba(34, 197, 94, 0.22)' } },
      { name: '合理偏低(P30-P50)', type: 'line', data: diffP50, stack: 'band', symbol: 'none', lineStyle: { opacity: 0 }, areaStyle: { color: 'rgba(134, 239, 172, 0.14)' } },
      { name: '合理偏高(P50-P70)', type: 'line', data: diffP70, stack: 'band', symbol: 'none', lineStyle: { opacity: 0 }, areaStyle: { color: 'rgba(251, 191, 36, 0.16)' } },
      { name: '高估区(>P70)', type: 'line', data: diffP90, stack: 'band', symbol: 'none', lineStyle: { opacity: 0 }, areaStyle: { color: 'rgba(248, 113, 113, 0.22)' } },
      { name: 'P50中轨', type: 'line', data: priceP50, symbol: 'none', lineStyle: { color: 'rgba(148, 163, 184, 0.55)', type: 'dashed', width: 1 }, showInLegend: false },
      { name: '股价', type: 'line', data: bandData.map(d => d.price), symbol: 'none', lineStyle: { color: '#f8fbff', width: 2 }, z: 5 },
    ] : [
      { name: '股价', type: 'line', data: bandData.map(d => d.price), symbol: 'none', lineStyle: { color: '#f8fbff', width: 2 } },
    ],
  } : null;

  const chartOption = chartMode === 'kline' ? klineChartOption : bandChartOption;

  const fmt = (n: number | undefined | null, suffix = '') => (n === undefined || n === null ? '--' : n.toFixed(2) + suffix);

  return (
    <div className="flex flex-col gap-3">
      {/* 实时估值卡 */}
      <div className="grid grid-cols-4 gap-2">
        <MetricCard label="现价" value={valLoading ? '...' : fmt(valuation?.close)} />
        <MetricCard label="PE TTM" value={valLoading ? '...' : fmt(valuation?.pe_ttm)} />
        <MetricCard label="PB" value={valLoading ? '...' : fmt(valuation?.pb)} />
        <MetricCard label="ROE" value={valLoading ? '...' : fmt(valuation?.roe, '%')} />
      </div>

      {/* 图表面板 */}
      <div className="p-3 border border-[#1e293b] rounded-xl bg-white/[0.02]">
        <div className="flex flex-wrap gap-2 items-center mb-3">
          {(['K线', 'PB Band', 'PE Band'] as const).map((m) => {
            const mode: ChartMode = m === 'K线' ? 'kline' : m === 'PB Band' ? 'pb_band' : 'pe_band';
            return (
              <button
                key={m}
                onClick={() => setChartMode(mode)}
                className={`px-2.5 py-1.5 rounded-lg text-xs font-mono border transition-all ${
                  chartMode === mode
                    ? 'text-[#f59e0b] border-[#f59e0b] bg-[rgba(245,158,11,0.1)]'
                    : 'bg-white/[0.03] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
                }`}
              >
                {m}
              </button>
            );
          })}
          {chartMode === 'kline' && (['日线', '周线', '月线'] as const).map((p) => {
            const period: ChartPeriod = p === '日线' ? 'daily' : p === '周线' ? 'weekly' : 'monthly';
            return (
              <button
                key={p}
                onClick={() => setChartPeriod(period)}
                className={`px-2.5 py-1.5 rounded-lg text-xs font-mono border transition-all ${
                  chartPeriod === period
                    ? 'text-[#f59e0b] border-[#f59e0b] bg-[rgba(245,158,11,0.1)]'
                    : 'bg-white/[0.03] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
                }`}
              >
                {p}
              </button>
            );
          })}
          <span className="mx-1 text-[#1e293b]">|</span>
          {(['1M', '3M', '6M', '1Y', '3Y', 'ALL'] as RangePreset[]).map((preset) => (
            <button
              key={preset}
              onClick={() => setRangePreset(preset)}
              className={`px-2 py-1.5 rounded-lg text-xs font-mono border transition-all ${
                rangePreset === preset
                  ? 'text-[#f59e0b] border-[#f59e0b] bg-[rgba(245,158,11,0.1)]'
                  : 'bg-white/[0.03] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
              }`}
            >
              {preset === 'ALL' ? '全部' : preset}
            </button>
          ))}
        </div>
        {loading ? (
          <div className="h-[360px] flex items-center justify-center text-[#b4c2d6] text-sm">加载中...</div>
        ) : chartOption ? (
          <ReactECharts option={chartOption} style={{ height: '360px' }} opts={{ renderer: 'canvas' }} />
        ) : (
          <div className="h-[360px] flex items-center justify-center text-[#b4c2d6] text-sm">暂无数据</div>
        )}
      </div>
    </div>
  );
}

function MetricCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="bg-white/[0.03] rounded-lg px-2.5 py-2 border border-[#1e293b]">
      <div className="text-[10px] text-[#b4c2d6] uppercase tracking-wider mb-0.5">{label}</div>
      <div className="font-mono text-sm font-semibold text-[#f8fbff]">{value}</div>
    </div>
  );
}
