import React, { useState } from 'react';
import { Stock, Quote, CompanySummary } from '../services/api';
import { StockDetail } from './StockDetail';

interface StockCardProps {
  stock: Stock;
  quote?: Quote;
  summary?: CompanySummary;
  token: string;
}

const SECTOR_COLORS: Record<string, string> = {
  '化工': 'bg-[rgba(245,158,11,0.15)] text-[#f59e0b]',
  '金属/材料': 'bg-[rgba(59,130,246,0.15)] text-[#3b82f6]',
  '农业/食品': 'bg-[rgba(34,197,94,0.15)] text-[#22c55e]',
  '能源': 'bg-[rgba(167,139,250,0.15)] text-[#a78bfa]',
  '工程/建筑': 'bg-[rgba(249,115,22,0.15)] text-[#f97316]',
  '军工/民爆': 'bg-[rgba(239,68,68,0.15)] text-[#ef4444]',
  '稀有金属': 'bg-[rgba(236,72,153,0.15)] text-[#ec4899]',
};

export function StockCard({ stock, quote, summary, token }: StockCardProps) {
  const [expanded, setExpanded] = useState(false);

  const currentPrice = quote?.close || 0;
  const changePercent = (quote?.change || 0) * 100; // 转换为百分比
  const entryDistance = stock.entryPrice > 0 ? ((currentPrice - stock.entryPrice) / stock.entryPrice) * 100 : 0;
  const targetDistance = stock.targetPrice > 0 ? ((currentPrice - stock.targetPrice) / stock.targetPrice) * 100 : 0;

  const changeClass = changePercent > 0 ? 'text-[#ef4444]' : changePercent < 0 ? 'text-[#22c55e]' : 'text-[#f8fbff]';
  const entryClass = entryDistance < 0 ? 'text-[#ef4444]' : 'text-[#22c55e]';
  const targetClass = targetDistance < 0 ? 'text-[#d6e0ee]' : 'text-[#a78bfa]';

  // 状态标签
  const statusTags: string[] = [];
  if (stock.heavyPrice > 0 && currentPrice <= stock.heavyPrice) {
    statusTags.push('到达重仓价');
  }
  if (stock.entryPrice > 0 && currentPrice <= stock.entryPrice) {
    statusTags.push('到达入场价');
  }
  if (stock.entryPrice > 0 && currentPrice < stock.entryPrice * 0.95) {
    statusTags.push('低于入场价 5%');
  }
  if (stock.targetPrice > 0 && currentPrice >= stock.targetPrice) {
    statusTags.push('触发止盈');
  }

  const sectorColor = SECTOR_COLORS[stock.sector] || 'bg-[rgba(148,163,184,0.1)] text-[#d6e0ee]';

  return (
    <div
      className={`bg-[#111827] border border-[#1e293b] rounded-xl overflow-hidden transition-all cursor-pointer ${
        expanded ? 'border-[#1e293b]' : ''
      }`}
      onClick={() => setExpanded(!expanded)}
    >
      {/* 股票头部 */}
      <div className="grid grid-cols-[2.2fr_1fr_1fr_1fr_10px] items-center px-4 py-2.5 gap-2.5">
        {/* 名称和代码 */}
        <div className="font-bold text-base flex items-center gap-1.5 flex-wrap">
          <span>{stock.name}</span>
          <span className="font-mono text-xs text-[#b4c2d6]">{stock.code}</span>
          <span className={`text-xs px-1.5 py-0.5 rounded font-normal ${sectorColor}`}>
            {stock.sector}
          </span>
          {stock.isRecommend && (
            <span className="text-xs px-1.5 py-0.5 rounded bg-[rgba(251,191,36,0.15)] text-[#fbbf24] font-normal">
              推荐
            </span>
          )}
        </div>

        {/* 当前价 */}
        <div className="flex flex-col items-start gap-1 font-mono text-sm">
          <span className="text-xs text-[#b4c2d6]">当前价</span>
          <div className="flex items-baseline gap-2">
            <span className="text-lg font-bold">{currentPrice.toFixed(2)}</span>
            <span className={`text-xs font-bold ${changeClass}`}>
              {changePercent > 0 ? '+' : ''}{changePercent.toFixed(1)}%
            </span>
          </div>
        </div>

        {/* 入场价 */}
        <div className="flex flex-col items-start gap-1 font-mono text-sm">
          <span className="text-xs text-[#b4c2d6]">入场价</span>
          <div className="flex items-baseline gap-2">
            <span className="text-lg font-bold">{stock.entryPrice.toFixed(1)}</span>
            <span className={`text-xs font-bold ${entryClass}`}>
              {entryDistance > 0 ? '+' : ''}{entryDistance.toFixed(1)}%
            </span>
          </div>
        </div>

        {/* 止盈价 */}
        <div className="flex flex-col items-start gap-1 font-mono text-sm">
          <span className="text-xs text-[#b4c2d6]">止盈价</span>
          <div className="flex items-baseline gap-2">
            <span className="text-lg font-bold">{stock.targetPrice.toFixed(1)}</span>
            <span className={`text-xs font-bold ${targetClass}`}>
              {targetDistance > 0 ? '+' : ''}{targetDistance.toFixed(1)}%
            </span>
          </div>
        </div>

        {/* 展开图标 */}
        <div className="flex items-center justify-center w-5 h-5 text-[#b4c2d6] justify-self-end">
          <span
            className="block w-1.5 h-1.5 border-r border-b border-current transition-transform"
            style={{ transform: expanded ? 'rotate(225deg) translateY(1px)' : 'rotate(45deg) translateY(-1px)' }}
          />
        </div>
      </div>

      {/* 状态标签 */}
      {statusTags.length > 0 && (
        <div className="flex gap-2 px-4 pb-2">
          {statusTags.map((tag) => (
            <span
              key={tag}
              className={`text-xs px-2 py-1 rounded-full font-bold ${
                tag === '到达重仓价' ? 'bg-[rgba(59,130,246,0.14)] text-[#3b82f6] border border-[rgba(59,130,246,0.25)]' :
                tag === '到达入场价' ? 'bg-[rgba(239,68,68,0.14)] text-[#ef4444] border border-[rgba(239,68,68,0.25)]' :
                tag === '低于入场价 5%' ? 'bg-[rgba(245,158,11,0.14)] text-[#f59e0b] border border-[rgba(245,158,11,0.25)]' :
                'bg-[rgba(167,139,250,0.14)] text-[#a78bfa] border border-[rgba(167,139,250,0.25)]'
              }`}
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* 详情面板 */}
      {expanded && (
        <StockDetail
          stock={stock}
          quote={quote}
          summary={summary}
          token={token}
        />
      )}
    </div>
  );
}