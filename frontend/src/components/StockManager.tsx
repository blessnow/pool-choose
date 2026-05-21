import { useState, useEffect, useCallback, useRef } from 'react';
import {
  Stock,
  StockWithQuote,
  StockSearchResult,
  getStocks,
  createStock,
  updateStock,
  deleteStock,
  searchStocks,
  getStockInfo,
} from '../services/api';
import { StockPreview } from './StockPreview';

interface StockManagerProps {
  token: string;
  onBack: () => void;
}

type FormState = Omit<Stock, 'id' | 'createdAt' | 'updatedAt'>;

const emptyForm: FormState = {
  code: '',
  name: '',
  sector: '',
  industry: '',
  isRecommend: false,
  entryPrice: 0,
  heavyPrice: 0,
  targetPrice: 0,
  coreLogic: '',
};

export function StockManager({ token, onBack }: StockManagerProps) {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [editing, setEditing] = useState<FormState | null>(null);
  const [isNew, setIsNew] = useState(false);
  const [error, setError] = useState('');
  const [enriching, setEnriching] = useState(false); // 拉取行业板块时的 loading

  const load = useCallback(async () => {
    setLoading(true);
    const r = await getStocks(token);
    if (r.ok && r.data) {
      setStocks(r.data.map((s: StockWithQuote) => s.stock));
    }
    setLoading(false);
  }, [token]);

  useEffect(() => {
    load();
  }, [load]);

  const openAdd = () => {
    setIsNew(true);
    setEditing({ ...emptyForm });
    setError('');
  };

  const openEdit = (s: Stock) => {
    setIsNew(false);
    setEditing({
      code: s.code,
      name: s.name,
      sector: s.sector,
      industry: s.industry,
      isRecommend: s.isRecommend,
      entryPrice: s.entryPrice,
      heavyPrice: s.heavyPrice,
      targetPrice: s.targetPrice,
      coreLogic: s.coreLogic,
    });
    setError('');
  };

  const closeForm = () => {
    setEditing(null);
    setError('');
  };

  const onSave = async () => {
    if (!editing) return;
    if (!editing.code.trim() || !editing.name.trim()) {
      setError('股票代码和名称必填');
      return;
    }
    if (isNew && !/^\d{6}$/.test(editing.code.trim())) {
      setError('股票代码必须是 6 位数字');
      return;
    }
    setLoading(true);
    setError('');
    try {
      const result = isNew
        ? await createStock(token, editing)
        : await updateStock(token, editing.code, editing);
      if (result.ok) {
        closeForm();
        await load();
      } else {
        setError(result.error || '保存失败');
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : '网络错误');
    } finally {
      setLoading(false);
    }
  };

  const onDelete = async (s: Stock) => {
    if (!window.confirm(`确认删除 ${s.code} ${s.name}？此操作不可恢复。`)) return;
    setLoading(true);
    const r = await deleteStock(token, s.code);
    if (r.ok) {
      await load();
    } else {
      setError(r.error || '删除失败');
    }
    setLoading(false);
  };

  const filtered = stocks.filter((s) => {
    if (!search) return true;
    const q = search.toLowerCase();
    return (
      s.code.includes(q) ||
      s.name.toLowerCase().includes(q) ||
      (s.sector || '').toLowerCase().includes(q) ||
      (s.industry || '').toLowerCase().includes(q)
    );
  });

  const updateField = <K extends keyof FormState>(key: K, value: FormState[K]) => {
    setEditing((prev) => (prev ? { ...prev, [key]: value } : prev));
  };

  return (
    <div className="min-h-screen bg-[#0a0e17]">
      <div className="max-w-[1400px] mx-auto px-4 py-6">
        {/* 顶栏 */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <button
              onClick={onBack}
              className="px-3 py-1.5 rounded-lg text-sm border border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b] transition-all"
            >
              ← 返回面板
            </button>
            <h1 className="text-xl font-bold text-[#f8fbff]">股票池配置</h1>
            <span className="text-xs text-[#b4c2d6]">共 {stocks.length} 只</span>
          </div>
          <div className="flex items-center gap-2">
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="搜索代码 / 名称 / 板块"
              className="bg-[#111827] border border-[#1e293b] rounded-lg px-3 py-1.5 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] placeholder:text-[#b4c2d6] w-64"
            />
            <button
              onClick={openAdd}
              className="px-4 py-1.5 rounded-lg text-sm font-bold bg-gradient-to-r from-[#f59e0b] to-[#f97316] text-white hover:opacity-90 transition-all"
            >
              + 新增股票
            </button>
          </div>
        </div>

        {error && !editing && (
          <div className="mb-4 px-4 py-2 rounded-lg bg-[rgba(239,68,68,0.1)] border border-[#ef4444] text-sm text-[#ef4444]">
            {error}
          </div>
        )}

        {/* 列表 */}
        <div className="bg-[#111827] border border-[#1e293b] rounded-xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-white/[0.04] text-[#f8fbff]">
                <tr>
                  <th className="px-3 py-2.5 text-left">代码</th>
                  <th className="px-3 py-2.5 text-left">名称</th>
                  <th className="px-3 py-2.5 text-left">板块</th>
                  <th className="px-3 py-2.5 text-left">行业</th>
                  <th className="px-3 py-2.5 text-center">推荐</th>
                  <th className="px-3 py-2.5 text-right font-mono">入场</th>
                  <th className="px-3 py-2.5 text-right font-mono">重仓</th>
                  <th className="px-3 py-2.5 text-right font-mono">止盈</th>
                  <th className="px-3 py-2.5 text-center">操作</th>
                </tr>
              </thead>
              <tbody>
                {loading && stocks.length === 0 ? (
                  <tr>
                    <td colSpan={9} className="px-3 py-8 text-center text-[#b4c2d6]">
                      加载中...
                    </td>
                  </tr>
                ) : filtered.length === 0 ? (
                  <tr>
                    <td colSpan={9} className="px-3 py-8 text-center text-[#b4c2d6]">
                      {search ? '没有匹配的股票' : '暂无股票，点击右上角新增'}
                    </td>
                  </tr>
                ) : (
                  filtered.map((s) => (
                    <tr key={s.code} className="border-t border-[#1e293b] hover:bg-white/[0.02]">
                      <td className="px-3 py-2 font-mono text-[#f8fbff]">{s.code}</td>
                      <td className="px-3 py-2 text-[#f8fbff]">{s.name}</td>
                      <td className="px-3 py-2 text-[#d6e0ee]">{s.sector}</td>
                      <td className="px-3 py-2 text-[#d6e0ee] text-xs">{s.industry}</td>
                      <td className="px-3 py-2 text-center">
                        {s.isRecommend ? (
                          <span className="px-2 py-0.5 rounded-full text-xs bg-[rgba(245,158,11,0.15)] text-[#f59e0b]">推荐</span>
                        ) : (
                          <span className="text-[#475569] text-xs">—</span>
                        )}
                      </td>
                      <td className="px-3 py-2 text-right font-mono text-[#d6e0ee]">{s.entryPrice > 0 ? s.entryPrice.toFixed(2) : '—'}</td>
                      <td className="px-3 py-2 text-right font-mono text-[#d6e0ee]">{s.heavyPrice > 0 ? s.heavyPrice.toFixed(2) : '—'}</td>
                      <td className="px-3 py-2 text-right font-mono text-[#d6e0ee]">{s.targetPrice > 0 ? s.targetPrice.toFixed(2) : '—'}</td>
                      <td className="px-3 py-2">
                        <div className="flex justify-center gap-2">
                          <button
                            onClick={() => openEdit(s)}
                            className="px-2 py-1 rounded text-xs border border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b] transition-all"
                          >
                            编辑
                          </button>
                          <button
                            onClick={() => onDelete(s)}
                            className="px-2 py-1 rounded text-xs border border-[#1e293b] text-[#d6e0ee] hover:border-[#ef4444] hover:text-[#ef4444] transition-all"
                          >
                            删除
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* 编辑/新增 弹窗 */}
      {editing && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
          onClick={closeForm}
        >
          <div
            className="bg-[#111827] border border-[#1e293b] rounded-2xl w-[1200px] max-w-[95vw] max-h-[92vh] overflow-y-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="px-6 py-4 border-b border-[#1e293b] flex justify-between items-center">
              <h2 className="text-lg font-bold text-[#f8fbff]">
                {isNew ? '新增股票' : `编辑 ${editing.code} ${editing.name}`}
              </h2>
              <button
                onClick={closeForm}
                className="text-[#b4c2d6] hover:text-[#f8fbff] text-xl leading-none"
              >
                ×
              </button>
            </div>

            <div className="px-6 py-4 grid grid-cols-[480px_minmax(0,1fr)] gap-6">
              <div className="space-y-3">
              <div className="grid grid-cols-2 gap-3">
                <Field label="股票代码 *" hint={isNew ? '输入拼音/代码/中文' : '不可改'}>
                  {isNew ? (
                    <StockSearchInput
                      token={token}
                      value={editing.code}
                      onChange={(v) => updateField('code', v)}
                      onPick={async (r) => {
                        // 立即回填 code+name，避免用户感到延迟
                        setEditing((prev) => prev ? { ...prev, code: r.code, name: r.name } : prev);
                        setEnriching(true);
                        try {
                          const info = await getStockInfo(token, r.code);
                          if (info.ok && info.data) {
                            setEditing((prev) => {
                              if (!prev || prev.code !== r.code) return prev; // 用户已切走
                              return {
                                ...prev,
                                // 用户已手填的不覆盖
                                sector: prev.sector || info.data!.sector,
                                industry: prev.industry || info.data!.industry,
                                name: prev.name || info.data!.name,
                              };
                            });
                          }
                        } finally {
                          setEnriching(false);
                        }
                      }}
                    />
                  ) : (
                    <input
                      type="text"
                      value={editing.code}
                      disabled
                      className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none disabled:opacity-60 disabled:cursor-not-allowed font-mono"
                    />
                  )}
                </Field>
                <Field label="股票名称 *">
                  <input
                    type="text"
                    value={editing.name}
                    onChange={(e) => updateField('name', e.target.value)}
                    placeholder="如 中国神华"
                    className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b]"
                  />
                </Field>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <Field label="板块" hint={enriching ? '自动获取中...' : '可手改'}>
                  <input
                    type="text"
                    value={editing.sector}
                    onChange={(e) => updateField('sector', e.target.value)}
                    placeholder={enriching ? '获取中...' : '如 白酒 / 纯碱 / 煤炭'}
                    className={`w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] ${enriching ? 'opacity-60' : ''}`}
                  />
                </Field>
                <Field label="行业" hint={enriching ? '自动获取中...' : '可手改'}>
                  <input
                    type="text"
                    value={editing.industry}
                    onChange={(e) => updateField('industry', e.target.value)}
                    placeholder={enriching ? '获取中...' : '如 纯碱 / 化学原料 / 基础化工'}
                    className={`w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] ${enriching ? 'opacity-60' : ''}`}
                  />
                </Field>
              </div>

              <div className="grid grid-cols-3 gap-3">
                <Field label="入场价" hint="0 = 未设定">
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={editing.entryPrice}
                    onChange={(e) => updateField('entryPrice', parseFloat(e.target.value) || 0)}
                    className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] font-mono"
                  />
                </Field>
                <Field label="重仓价" hint="低于入场">
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={editing.heavyPrice}
                    onChange={(e) => updateField('heavyPrice', parseFloat(e.target.value) || 0)}
                    className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] font-mono"
                  />
                </Field>
                <Field label="止盈价" hint="高于入场">
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={editing.targetPrice}
                    onChange={(e) => updateField('targetPrice', parseFloat(e.target.value) || 0)}
                    className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] font-mono"
                  />
                </Field>
              </div>

              <Field label="核心逻辑">
                <textarea
                  value={editing.coreLogic}
                  onChange={(e) => updateField('coreLogic', e.target.value)}
                  rows={5}
                  placeholder="✅ 利好：…  ⚠️ 风险：…  巴菲特视角：…"
                  className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] resize-y"
                />
              </Field>

              <label className="flex items-center gap-2 text-sm text-[#d6e0ee] cursor-pointer select-none">
                <input
                  type="checkbox"
                  checked={editing.isRecommend}
                  onChange={(e) => updateField('isRecommend', e.target.checked)}
                  className="w-4 h-4 accent-[#f59e0b]"
                />
                标记为推荐
              </label>

              {error && (
                <div className="px-3 py-2 rounded-lg bg-[rgba(239,68,68,0.1)] border border-[#ef4444] text-sm text-[#ef4444]">
                  {error}
                </div>
              )}
              </div>

              {/* 右列：实时估值 + K线 / PE Band / PB Band 预览 */}
              <div className="min-w-0">
                <StockPreview token={token} code={editing.code} />
              </div>
            </div>

            <div className="px-6 py-4 border-t border-[#1e293b] flex justify-end gap-2">
              <button
                onClick={closeForm}
                disabled={loading}
                className="px-4 py-2 rounded-lg text-sm border border-[#1e293b] text-[#d6e0ee] hover:border-[#f59e0b] hover:text-[#f59e0b] transition-all disabled:opacity-50"
              >
                取消
              </button>
              <button
                onClick={onSave}
                disabled={loading}
                className="px-4 py-2 rounded-lg text-sm font-bold bg-gradient-to-r from-[#f59e0b] to-[#f97316] text-white hover:opacity-90 transition-all disabled:opacity-50"
              >
                {loading ? '保存中...' : '保存'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function Field({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) {
  return (
    <label className="block">
      <div className="text-xs text-[#b4c2d6] mb-1 flex items-center justify-between">
        <span>{label}</span>
        {hint && <span className="text-[#475569]">{hint}</span>}
      </div>
      {children}
    </label>
  );
}

interface StockSearchInputProps {
  token: string;
  value: string;
  onChange: (v: string) => void;
  onPick: (r: StockSearchResult) => void;
}

// 拼音 / 代码 / 中文模糊补全。300ms 防抖，最多 10 条候选，键盘 ↑↓ Enter Esc 操作
function StockSearchInput({ token, value, onChange, onPick }: StockSearchInputProps) {
  const [results, setResults] = useState<StockSearchResult[]>([]);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [activeIdx, setActiveIdx] = useState(-1);
  const debounceRef = useRef<number | null>(null);
  const reqIdRef = useRef(0);
  const wrapperRef = useRef<HTMLDivElement>(null);

  // 防抖 + 竞态保护
  useEffect(() => {
    if (debounceRef.current) window.clearTimeout(debounceRef.current);
    const q = value.trim();
    if (q.length === 0) {
      setResults([]);
      setOpen(false);
      return;
    }
    debounceRef.current = window.setTimeout(async () => {
      const myId = ++reqIdRef.current;
      setLoading(true);
      try {
        const r = await searchStocks(token, q, 10);
        // 丢弃过时请求
        if (myId !== reqIdRef.current) return;
        if (r.ok && r.data) {
          setResults(r.data);
          setOpen(r.data.length > 0);
          setActiveIdx(-1);
        }
      } finally {
        if (myId === reqIdRef.current) setLoading(false);
      }
    }, 300);
    return () => {
      if (debounceRef.current) window.clearTimeout(debounceRef.current);
    };
  }, [value, token]);

  // 点击外部关闭
  useEffect(() => {
    const onDocClick = (e: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', onDocClick);
    return () => document.removeEventListener('mousedown', onDocClick);
  }, []);

  const pick = (r: StockSearchResult) => {
    onPick(r);
    setOpen(false);
    setActiveIdx(-1);
  };

  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (!open || results.length === 0) return;
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActiveIdx((i) => (i + 1) % results.length);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActiveIdx((i) => (i <= 0 ? results.length - 1 : i - 1));
    } else if (e.key === 'Enter') {
      if (activeIdx >= 0 && activeIdx < results.length) {
        e.preventDefault();
        pick(results[activeIdx]);
      }
    } else if (e.key === 'Escape') {
      setOpen(false);
    }
  };

  return (
    <div ref={wrapperRef} className="relative">
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value.trim())}
        onFocus={() => results.length > 0 && setOpen(true)}
        onKeyDown={onKeyDown}
        placeholder="ZGSH / 600028 / 中国石化"
        className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-3 py-2 text-sm text-[#f8fbff] outline-none focus:border-[#f59e0b] font-mono"
      />
      {open && (
        <div className="absolute z-50 left-0 right-0 mt-1 bg-[#0a0e17] border border-[#1e293b] rounded-lg shadow-xl max-h-72 overflow-y-auto">
          {loading && results.length === 0 ? (
            <div className="px-3 py-2 text-xs text-[#b4c2d6]">搜索中...</div>
          ) : results.length === 0 ? (
            <div className="px-3 py-2 text-xs text-[#b4c2d6]">无匹配股票</div>
          ) : (
            results.map((r, i) => (
              <div
                key={r.code}
                onMouseDown={(e) => {
                  e.preventDefault();
                  pick(r);
                }}
                onMouseEnter={() => setActiveIdx(i)}
                className={`px-3 py-2 cursor-pointer text-sm flex items-center justify-between gap-3 ${
                  i === activeIdx ? 'bg-[rgba(245,158,11,0.12)]' : 'hover:bg-white/[0.04]'
                }`}
              >
                <div className="flex items-center gap-2 min-w-0">
                  <span className="font-mono text-[#f8fbff] text-xs w-14 flex-shrink-0">{r.code}</span>
                  <span className="text-[#f8fbff] truncate">{r.name}</span>
                </div>
                <div className="flex items-center gap-2 flex-shrink-0 text-xs">
                  <span className="text-[#b4c2d6]">{r.pinyin}</span>
                  <span className="px-1.5 py-0.5 rounded bg-white/[0.06] text-[#94a3b8]">{r.market}</span>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
