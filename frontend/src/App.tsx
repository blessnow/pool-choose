import { useState, useEffect } from 'react';
import { login, getStocks, getCompanySummaries, getCycleInsight, Stock, Quote, CompanySummary, CycleInsight, StockWithQuote } from './services/api';
import { Header } from './components/Header';
import { SummaryGrid } from './components/SummaryGrid';
import { CycleSection } from './components/CycleSection';
import { FilterBar, FilterState, SortState } from './components/FilterBar';
import { StockCard } from './components/StockCard';
import './index.css';

function App() {
  const [token, setToken] = useState<string>('');
  const [password, setPassword] = useState('');
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [error, setError] = useState('');

  const [stocks, setStocks] = useState<Stock[]>([]);
  const [quotes, setQuotes] = useState<Record<string, Quote>>({});
  const [summaries, setSummaries] = useState<Record<string, CompanySummary>>({});
  const [cycleInsight, setCycleInsight] = useState<CycleInsight | null>(null);
  const [loading, setLoading] = useState(false);
  const [dataLoaded, setDataLoaded] = useState(false);

  const [filters, setFilters] = useState<FilterState>({ tab: 'all' });
  const [sort, setSort] = useState<SortState>({ field: 'recommend', direction: 'asc' });
  const [searchQuery, setSearchQuery] = useState('');

  // 登录
  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const result = await login(password);
      if (result.ok && result.data) {
        setToken(result.data.token);
        setIsLoggedIn(true);
        localStorage.setItem('token', result.data.token);
      } else {
        setError(result.error || '登录失败');
      }
    } catch {
      setError('网络错误');
    } finally {
      setLoading(false);
    }
  };

  // 加载数据
  const loadData = async (authToken: string) => {
    setLoading(true);
    setDataLoaded(false);
    try {
      // 并行获取所有数据
      const [stocksResult, cycleResult] = await Promise.all([
        getStocks(authToken),
        getCycleInsight(authToken)
      ]);

      if (stocksResult.ok && stocksResult.data) {
        const stockList = stocksResult.data.map((s: StockWithQuote) => s.stock);
        const quotesMap: Record<string, Quote> = {};
        stocksResult.data.forEach((s: StockWithQuote) => {
          if (s.quote) {
            quotesMap[s.stock.code] = s.quote;
          }
        });

        setStocks(stockList);
        setQuotes(quotesMap);

        // 获取公司财务数据
        const codes = stockList.map(s => s.code);
        const summariesResult = await getCompanySummaries(authToken, codes);
        if (summariesResult.ok && summariesResult.data) {
          setSummaries(summariesResult.data);
        }
      }

      if (cycleResult.ok && cycleResult.data) {
        setCycleInsight(cycleResult.data);
      }

      setDataLoaded(true);
    } catch (err) {
      console.error('加载数据失败:', err);
      // 如果加载失败，可能是token过期，清除登录状态
      if (err instanceof Error && err.message.includes('401')) {
        setIsLoggedIn(false);
        setToken('');
        localStorage.removeItem('token');
      }
    } finally {
      setLoading(false);
    }
  };

  // 当token变化时加载数据
  useEffect(() => {
    if (token && isLoggedIn) {
      loadData(token);
    }
  }, [token, isLoggedIn]);

  // 从localStorage恢复token
  useEffect(() => {
    const savedToken = localStorage.getItem('token');
    if (savedToken) {
      setToken(savedToken);
      setIsLoggedIn(true);
    }
  }, []);

  // 筛选和排序
  const getFilteredStocks = () => {
    let filtered = [...stocks];

    // 搜索
    if (searchQuery) {
      filtered = filtered.filter(s =>
        s.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        s.code.includes(searchQuery)
      );
    }

    // Tab筛选
    if (filters.tab === 'recommend') {
      filtered = filtered.filter(s => s.isRecommend);
    }

    // 状态筛选
    if (filters.status) {
      filtered = filtered.filter(s => {
        const quote = quotes[s.code];
        const price = quote?.close || 0;
        switch (filters.status) {
          case '到达入场':
            return price <= s.entryPrice;
          case '到达重仓':
            return price <= s.heavyPrice;
          case '低于入场5%':
            return price < s.entryPrice * 0.95;
          case '已触止盈':
            return price >= s.targetPrice;
          default:
            return true;
        }
      });
    }

    // 行业筛选
    if (filters.sector) {
      filtered = filtered.filter(s => s.sector === filters.sector);
    }

    // 排序
    filtered.sort((a, b) => {
      const quoteA = quotes[a.code];
      const quoteB = quotes[b.code];
      const priceA = quoteA?.close || 0;
      const priceB = quoteB?.close || 0;

      let cmp = 0;
      switch (sort.field) {
        case 'recommend':
          cmp = (b.isRecommend ? 1 : 0) - (a.isRecommend ? 1 : 0);
          break;
        case 'change':
          cmp = (quoteB?.change || 0) - (quoteA?.change || 0);
          break;
        case 'entryDistance':
          const distA = a.entryPrice > 0 ? (priceA - a.entryPrice) / a.entryPrice : 0;
          const distB = b.entryPrice > 0 ? (priceB - b.entryPrice) / b.entryPrice : 0;
          cmp = distB - distA;
          break;
        case 'targetDistance':
          const targetA = a.targetPrice > 0 ? (priceA - a.targetPrice) / a.targetPrice : 0;
          const targetB = b.targetPrice > 0 ? (priceB - b.targetPrice) / b.targetPrice : 0;
          cmp = targetB - targetA;
          break;
        default:
          cmp = 0;
      }

      return sort.direction === 'asc' ? cmp : -cmp;
    });

    return filtered;
  };

  // 计算统计数据
  const getStats = () => {
    const filteredStocks = stocks;
    let recommend = 0;
    let entryReached = 0;
    let heavyReached = 0;
    let belowEntry5 = 0;
    let profitTriggered = 0;

    filteredStocks.forEach(s => {
      if (s.isRecommend) recommend++;
      const quote = quotes[s.code];
      const price = quote?.close || 0;
      if (price <= s.entryPrice) entryReached++;
      if (price <= s.heavyPrice) heavyReached++;
      if (price < s.entryPrice * 0.95) belowEntry5++;
      if (price >= s.targetPrice) profitTriggered++;
    });

    return {
      total: stocks.length,
      recommend,
      entryReached,
      heavyReached,
      belowEntry5,
      profitTriggered,
    };
  };

  // 登录页面
  if (!isLoggedIn) {
    return (
      <div className="min-h-screen bg-[#0a0e17] flex items-center justify-center">
        <form onSubmit={handleLogin} className="bg-[#111827] border border-[#1e293b] rounded-xl p-8 w-[400px]">
          <h1 className="text-2xl font-bold text-[#f8fbff] mb-6 text-center">
            🐊 鳄鱼·周期股
          </h1>
          <input
            type="password"
            placeholder="请输入访问密码"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full bg-[#0a0e17] border border-[#1e293b] rounded-lg px-4 py-3 text-[#f8fbff] outline-none focus:border-[#f59e0b] placeholder:text-[#b4c2d6]"
          />
          {error && <p className="text-[#ef4444] text-sm mt-2">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full mt-4 bg-gradient-to-r from-[#f59e0b] to-[#f97316] text-white font-bold py-3 rounded-lg cursor-pointer transition-all hover:opacity-90 disabled:opacity-50"
          >
            {loading ? '登录中...' : '进入面板'}
          </button>
        </form>
      </div>
    );
  }

  const filteredStocks = getFilteredStocks();
  const stats = getStats();

  // 如果数据未加载完成，显示加载状态
  if (!dataLoaded && loading) {
    return (
      <div className="min-h-screen bg-[#0a0e17] flex items-center justify-center">
        <div className="text-[#b4c2d6]">加载中...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#0a0e17]">
      <div className="max-w-[1200px] mx-auto px-4 py-6 relative z-10">
        {/* Hero背景 */}
        <div className="relative -mt-6 -mx-4 mb-0 pt-[67px]">
          <div className="absolute inset-0 bg-gradient-to-b from-[rgba(56,92,138,0.12)] via-[rgba(36,58,88,0.08)] to-transparent opacity-50" />
        </div>

        <Header title="鳄鱼·周期股" subtitle="好公司　低位买　买资源　信周期　耐心拿　狂涨卖" />

        {/* 统计卡片 */}
        <div className="mb-2">
          <SummaryGrid
            {...stats}
            activeCard={
              filters.tab === 'recommend' ? 'recommend' :
              filters.tab === 'status' && filters.status === '到达入场' ? 'entryReached' :
              filters.tab === 'status' && filters.status === '到达重仓' ? 'heavyReached' :
              filters.tab === 'status' && filters.status === '低于入场5%' ? 'belowEntry5' :
              filters.tab === 'status' && filters.status === '已触止盈' ? 'profitTriggered' :
              filters.tab === 'all' ? 'total' : undefined
            }
            onCardClick={(type: string) => {
              switch (type) {
                case 'recommend':
                  setFilters({ tab: 'recommend' });
                  break;
                case 'entryReached':
                  setFilters({ tab: 'status', status: '到达入场' });
                  break;
                case 'heavyReached':
                  setFilters({ tab: 'status', status: '到达重仓' });
                  break;
                case 'belowEntry5':
                  setFilters({ tab: 'status', status: '低于入场5%' });
                  break;
                case 'profitTriggered':
                  setFilters({ tab: 'status', status: '已触止盈' });
                  break;
                default:
                  setFilters({ tab: 'all' });
              }
            }}
          />
        </div>

        {/* 周期分析 */}
        <div className="mt-2.5">
          <CycleSection insight={cycleInsight} />
        </div>

        {/* 筛选栏 */}
        <div className="mt-5">
          <FilterBar
            filters={filters}
            onFilterChange={setFilters}
            onSortChange={setSort}
            onSearch={setSearchQuery}
          />
        </div>

        {/* 股票列表 */}
        <div className="flex flex-col gap-2.5 mt-2.5">
          {loading && !dataLoaded ? (
            <div className="text-center py-8 text-[#b4c2d6]">加载中...</div>
          ) : filteredStocks.length === 0 && dataLoaded ? (
            <div className="text-center py-8 text-[#b4c2d6]">暂无股票数据</div>
          ) : (
            filteredStocks.map((stock) => (
              <StockCard
                key={stock.code}
                stock={stock}
                quote={quotes[stock.code]}
                summary={summaries[stock.code]}
                token={token}
              />
            ))
          )}
        </div>

        {/* 底部声明 */}
        <div className="text-center text-[#b4c2d6] text-xs py-8 mt-8 leading-relaxed">
          本系统仅供个人投资研究使用，不构成任何投资建议。投资有风险，入市需谨慎。
        </div>
      </div>
    </div>
  );
}

export default App;