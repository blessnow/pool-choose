import { useState, useEffect, useCallback } from 'react';
import { getStocks, getQuotes, StockWithQuote, Quote } from '../services/api';

interface UseStocksOptions {
  token: string;
  refreshInterval?: number;
}

interface UseStocksReturn {
  stocks: StockWithQuote[];
  quotes: Record<string, Quote>;
  loading: boolean;
  error: string | null;
  refresh: () => void;
}

export function useStocks(options: UseStocksOptions): UseStocksReturn {
  const { token, refreshInterval = 30000 } = options;
  const [stocks, setStocks] = useState<StockWithQuote[]>([]);
  const [quotes, setQuotes] = useState<Record<string, Quote>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    if (!token) return;

    try {
      setLoading(true);
      setError(null);

      const result = await getStocks(token);
      if (result.ok && result.data) {
        setStocks(result.data);

        // 提取股票代码获取行情
        const codes = result.data.map(s => s.stock.code);
        const quotesResult = await getQuotes(token, codes);
        if (quotesResult.ok && quotesResult.data) {
          setQuotes(quotesResult.data);
        }
      } else {
        setError(result.error || '获取数据失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '网络错误');
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    fetchData();

    // 设置定时刷新
    if (refreshInterval > 0) {
      const interval = setInterval(fetchData, refreshInterval);
      return () => clearInterval(interval);
    }
  }, [fetchData, refreshInterval]);

  return {
    stocks,
    quotes,
    loading,
    error,
    refresh: fetchData,
  };
}