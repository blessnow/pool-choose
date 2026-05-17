import React from 'react';

interface SummaryCardProps {
  num: number;
  label: string;
  color?: 'gold' | 'green' | 'red' | 'blue' | 'purple';
  onClick?: () => void;
  active?: boolean;
}

export function SummaryCard({ num, label, color = 'gold', onClick, active }: SummaryCardProps) {
  const colorClasses = {
    gold: 'text-[#fbbf24]',
    green: 'text-[#22c55e]',
    red: 'text-[#ef4444]',
    blue: 'text-[#3b82f6]',
    purple: 'text-[#a78bfa]',
  };

  return (
    <div
      className={`bg-[#111827] border rounded-xl px-4 py-3 text-center cursor-pointer transition-all duration-200 hover:-translate-y-0.5 ${
        active
          ? 'border-[#f59e0b] bg-[rgba(245,158,11,0.1)]'
          : 'border-[#1e293b] hover:border-[#f59e0b]'
      }`}
      onClick={onClick}
    >
      <div className={`font-mono text-2xl font-bold ${colorClasses[color]}`}>
        {num}
      </div>
      <div className={`text-xs mt-0.5 ${active ? 'text-[#f59e0b]' : 'text-[#d6e0ee]'}`}>{label}</div>
    </div>
  );
}

interface SummaryGridProps {
  total: number;
  recommend: number;
  entryReached: number;
  heavyReached: number;
  belowEntry5: number;
  profitTriggered: number;
  onCardClick?: (type: string) => void;
  activeCard?: string;
}

export function SummaryGrid({
  total,
  recommend,
  entryReached,
  heavyReached,
  belowEntry5,
  profitTriggered,
  onCardClick,
  activeCard,
}: SummaryGridProps) {
  return (
    <div className="grid grid-cols-[repeat(auto-fit,minmax(160px,1fr))] gap-3 mb-0">
      <SummaryCard num={total} label="股票池总数" color="blue" onClick={() => onCardClick?.('total')} active={activeCard === 'total'} />
      <SummaryCard num={recommend} label="推荐" color="gold" onClick={() => onCardClick?.('recommend')} active={activeCard === 'recommend'} />
      <SummaryCard num={entryReached} label="到达入场价" color="green" onClick={() => onCardClick?.('entryReached')} active={activeCard === 'entryReached'} />
      <SummaryCard num={heavyReached} label="到达重仓价" color="blue" onClick={() => onCardClick?.('heavyReached')} active={activeCard === 'heavyReached'} />
      <SummaryCard num={belowEntry5} label="低于入场 5%" color="red" onClick={() => onCardClick?.('belowEntry5')} active={activeCard === 'belowEntry5'} />
      <SummaryCard num={profitTriggered} label="触发止盈" color="purple" onClick={() => onCardClick?.('profitTriggered')} active={activeCard === 'profitTriggered'} />
    </div>
  );
}