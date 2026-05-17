import React, { useState } from 'react';

interface FilterBarProps {
  filters: FilterState;
  onFilterChange: (filters: FilterState) => void;
  onSortChange: (sort: SortState) => void;
  onSearch: (query: string) => void;
}

export interface FilterState {
  tab: 'all' | 'holdings' | 'status' | 'sector' | 'recommend';
  status?: string;
  sector?: string;
}

export interface SortState {
  field: string;
  direction: 'asc' | 'desc';
}

const STATUS_OPTIONS = ['到达入场', '到达重仓', '低于入场5%', '已触止盈'];
const SECTOR_OPTIONS = ['化工', '金属/材料', '农业/食品', '工程/建筑', '能源', '军工/民爆', '稀有金属'];
const SORT_OPTIONS = [
  { value: 'recommend', label: '推荐优先' },
  { value: 'change', label: '股价涨跌幅' },
  { value: 'entryDistance', label: '距离入场价' },
  { value: 'targetDistance', label: '距止盈价' },
];

const TIME_OPTIONS = ['今日', '本周', '本月', '近一年'];

export function FilterBar({ filters, onFilterChange, onSortChange, onSearch }: FilterBarProps) {
  const [activeSort, setActiveSort] = useState('recommend');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');
  const [activeTime, setActiveTime] = useState('今日');
  const [searchQuery, setSearchQuery] = useState('');

  const handleTabClick = (tab: typeof filters.tab) => {
    onFilterChange({ tab });
  };

  const handleStatusClick = (status: string) => {
    const newStatus = filters.status === status ? undefined : status;
    onFilterChange({ tab: 'status', status: newStatus });
  };

  const handleSectorClick = (sector: string) => {
    const newSector = filters.sector === sector ? undefined : sector;
    onFilterChange({ tab: 'sector', sector: newSector });
  };

  const handleSortClick = (sort: string) => {
    setActiveSort(sort);
    onSortChange({ field: sort, direction: sortDirection });
  };

  const toggleSortDirection = () => {
    const newDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    setSortDirection(newDirection);
    onSortChange({ field: activeSort, direction: newDirection });
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const query = e.target.value;
    setSearchQuery(query);
    onSearch(query);
  };

  return (
    <div className="flex gap-4 items-start mb-2.5">
      <div className="flex-1 min-w-0 flex gap-3 items-center flex-wrap">
        {/* 搜索框 */}
        <input
          type="text"
          placeholder="搜索股票名称/代码..."
          value={searchQuery}
          onChange={handleSearch}
          className="bg-[#111827] border border-[#1e293b] rounded-lg px-3.5 py-2 text-[#f8fbff] text-sm w-[200px] outline-none focus:border-[#f59e0b] placeholder:text-[#b4c2d6]"
        />

        {/* 主标签 */}
        <div className="flex gap-2 flex-wrap items-center">
          {[
            { key: 'all', label: '全部' },
            { key: 'holdings', label: '我的持仓' },
            { key: 'status', label: '当前状态' },
            { key: 'sector', label: '行业分类' },
            { key: 'recommend', label: '推荐优先' },
          ].map(({ key, label }) => (
            <button
              key={key}
              onClick={() => handleTabClick(key as typeof filters.tab)}
              className={`px-3.5 py-2 rounded-lg text-sm cursor-pointer transition-all border ${
                filters.tab === key
                  ? 'bg-[rgba(245,158,11,0.14)] border-[rgba(245,158,11,0.42)] text-[#fbbf24]'
                  : 'bg-white/[0.02] border-[#1e293b] text-[#b4c2d6] hover:border-[rgba(251,191,36,0.38)] hover:text-[#f8fbff]'
              }`}
            >
              {label}
            </button>
          ))}
        </div>

        {/* 状态筛选 */}
        {filters.tab === 'status' && (
          <div className="flex gap-2 flex-wrap">
            {STATUS_OPTIONS.map((status) => (
              <button
                key={status}
                onClick={() => handleStatusClick(status)}
                className={`px-3.5 py-2 rounded-lg text-xs cursor-pointer transition-all border ${
                  filters.status === status
                    ? 'bg-[rgba(245,158,11,0.15)] border-[#f59e0b] text-[#f59e0b]'
                    : 'bg-[#111827] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
                }`}
              >
                {status}
              </button>
            ))}
          </div>
        )}

        {/* 行业筛选 */}
        {filters.tab === 'sector' && (
          <div className="flex gap-2 flex-wrap">
            {SECTOR_OPTIONS.map((sector) => (
              <button
                key={sector}
                onClick={() => handleSectorClick(sector)}
                className={`px-3.5 py-2 rounded-lg text-xs cursor-pointer transition-all border ${
                  filters.sector === sector
                    ? 'bg-[rgba(245,158,11,0.15)] border-[#f59e0b] text-[#f59e0b]'
                    : 'bg-[#111827] border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b]'
                }`}
              >
                {sector}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* 排序工具 */}
      <div className="flex items-center gap-2.5 flex-shrink-0">
        <details className="relative">
          <summary className="list-none flex items-center gap-2.5 min-h-9 px-3.5 py-2 bg-white/[0.02] border border-[#1e293b] rounded-lg text-sm text-[#b4c2d6] cursor-pointer transition-all hover:border-[rgba(251,191,36,0.38)] hover:text-[#f8fbff]">
            {SORT_OPTIONS.find(o => o.value === activeSort)?.label || '排序'}
            <span className="w-2 h-2 border-r border-b border-current" style={{ transform: 'rotate(45deg)' }} />
          </summary>
          <div className="absolute top-full right-0 mt-2.5 w-[min(360px,calc(100vw-24px))] bg-[rgba(15,23,42,0.96)] border border-[rgba(148,163,184,0.18)] rounded-2xl shadow-2xl p-3.5 backdrop-blur-sm z-30">
            {SORT_OPTIONS.map((option) => (
              <button
                key={option.value}
                onClick={() => handleSortClick(option.value)}
                className={`w-full text-left px-3 py-2 rounded-full text-xs font-semibold transition-all border mb-2 ${
                  activeSort === option.value
                    ? 'bg-[rgba(245,158,11,0.14)] border-[rgba(245,158,11,0.42)] text-[#f59e0b]'
                    : 'bg-white/[0.03] border-[rgba(148,163,184,0.16)] text-[#d6e0ee] hover:border-[rgba(245,158,11,0.38)] hover:text-[#f8fbff]'
                }`}
              >
                {option.label}
              </button>
            ))}
          </div>
        </details>

        <button
          onClick={toggleSortDirection}
          className="min-w-9 h-9 flex items-center justify-center bg-white/[0.02] border border-[#1e293b] rounded-lg text-[#b4c2d6] cursor-pointer transition-all hover:border-[rgba(251,191,36,0.38)] hover:text-[#f8fbff]"
          title={sortDirection === 'asc' ? '当前正序，点击切到倒序' : '当前倒序，点击切到正序'}
        >
          {sortDirection === 'asc' ? '↑' : '↓'}
        </button>
      </div>
    </div>
  );
}