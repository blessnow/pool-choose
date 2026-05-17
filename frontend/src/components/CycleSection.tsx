import React, { useState } from 'react';
import { CycleInsight, MacroCard, CycleBar, Source } from '../services/api';

interface CycleSectionProps {
  insight: CycleInsight | null;
}

export function CycleSection({ insight }: CycleSectionProps) {
  const [collapsed, setCollapsed] = useState(false);

  if (!insight) return null;

  return (
    <div className={`bg-[#111827] border border-[#1e293b] rounded-2xl px-4 py-4 mb-5 transition-all ${collapsed ? 'py-4' : ''}`}>
      <div
        className="flex items-center justify-between gap-3 cursor-pointer"
        onClick={() => setCollapsed(!collapsed)}
      >
        <h2 className="text-base font-bold text-[#f8fbff] flex items-center gap-2.5 leading-tight">
          <span>📊</span>
          {insight.title}
        </h2>
        <button className="flex items-center justify-center w-2.5 h-2.5 p-0 border-none bg-transparent text-[#b4c2d6] cursor-pointer flex-shrink-0 transition-all">
          <span className={`block w-1.5 h-1.5 border-r border-b border-current transition-transform ${collapsed ? '' : 'rotate-180'}`}
            style={{ transform: collapsed ? 'rotate(45deg)' : 'rotate(225deg)' }}
          />
        </button>
      </div>

      <div className={`overflow-hidden transition-all ${collapsed ? 'max-h-0 opacity-0 mt-0' : 'max-h-[900px] opacity-100 mt-4'}`}>
        {/* 数据层 */}
        <div className="text-sm text-[#d6e0ee] leading-relaxed mb-4">
          <strong className="text-[#f8fbff]">数据层：</strong>
          {insight.dataLayer}
        </div>

        {/* 观点层 */}
        <div className="text-sm text-[#d6e0ee] leading-relaxed mb-4">
          <strong className="text-[#f8fbff]">观点层：</strong>
          {insight.opinionLayer}
        </div>

        {/* 宏观卡片 */}
        <div className="grid grid-cols-4 gap-3 mb-4">
          {insight.macroCards.map((card, i) => (
            <div key={i} className="p-3 rounded-xl border border-[#1e293b] bg-gradient-to-b from-white/[0.03] to-white/[0.015]">
              <div className="text-xs text-[#b4c2d6] mb-2">{card.label}</div>
              <div className="text-lg font-extrabold text-[#f8fbff]">{card.value}</div>
              <div className="text-xs text-[#d6e0ee] mt-1.5">{card.detail}</div>
            </div>
          ))}
        </div>

        {/* 周期条 */}
        {insight.bars.map((bar, i) => (
          <div key={i} className="mb-4">
            <div className="text-sm text-[#b4c2d6] mb-2">{bar.label}</div>
            <div className="h-2 rounded bg-[#1e293b] overflow-hidden">
              <div
                className="h-full rounded transition-all duration-1000"
                style={{ width: `${bar.value}%`, background: bar.gradient }}
              />
            </div>
            <div className="flex justify-between text-xs text-[#b4c2d6] mt-1">
              {bar.labels.map((label, j) => (
                <span key={j}>{label}</span>
              ))}
            </div>
          </div>
        ))}

        {/* 结论 */}
        <div className="bg-[rgba(245,158,11,0.05)] border-l-[3px] border-[#f59e0b] px-4 py-3.5 mt-4 rounded-r-lg text-sm leading-relaxed text-[#d6e0ee]">
          <p className="mb-2"><strong className="text-[#f8fbff]">结论：</strong>{insight.conclusion}</p>
          <p className="mb-2"><strong className="text-[#f8fbff]">关注点：</strong>{insight.focus}</p>
          <p><strong className="text-[#f8fbff]">风险：</strong>{insight.risk}</p>
        </div>

        {/* 来源 */}
        <div className="mt-4 text-sm text-[#b4c2d6]">
          <strong>来源：</strong>
          {insight.sources.map((source, i) => (
            <span key={i}>
              <a href={source.url} target="_blank" rel="noopener noreferrer" className="text-[#ffd36b] hover:text-[#ffe39a]">
                {source.label}
              </a>
              {i < insight.sources.length - 1 && ' · '}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}