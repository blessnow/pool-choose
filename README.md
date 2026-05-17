# 周期股投资管理系统

一个用于管理周期股投资的Web应用，包含实时行情、K线图表、持仓记录等功能。

## 功能特性

- 📊 **宏观经济数据展示** - CPI/PPI数据、周期定位分析
- 📈 **股票池管理** - 68只周期股的筛选、分类、状态跟踪
- 💹 **实时行情** - 股价、涨跌幅、PE/PB/ROE等财务指标
- 📉 **K线图表** - 日线/周线/月线图表，技术指标(MACD/RSI/KDJ)
- 💰 **持仓记录** - 入场价、止盈价、重仓位管理
- 🔍 **筛选过滤** - 按状态、行业、推荐优先级筛选

## 技术栈

- **Backend**: Go 1.21+, Gin, GORM, SQLite
- **Frontend**: React 18, TypeScript, Vite, TailwindCSS, ECharts
- **Data Source**: 新浪财经API (免费实时行情)
- **Deployment**: Docker容器化部署

## 快速开始

### 本地开发

1. **启动后端** (需要安装Go)
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

2. **启动前端**
```bash
cd frontend
npm install
npm run dev
```

3. 访问 http://localhost:5173

### Docker部署

```bash
docker-compose up -d
```

访问 http://localhost

默认密码: `dayuchi`

## API文档

### 认证
- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出

### 股票
- `GET /api/stocks` - 获取股票列表
- `GET /api/stocks/:code` - 获取单只股票详情
- `POST /api/stocks` - 创建股票
- `PUT /api/stocks/:code` - 更新股票
- `DELETE /api/stocks/:code` - 删除股票

### 行情
- `GET /api/quotes?codes=000707,300470` - 批量获取实时行情
- `GET /api/company-summaries?codes=000707` - 获取公司财务数据
- `GET /api/chart/:code?period=daily` - 获取K线数据

### 周期分析
- `GET /api/cycle-insight` - 获取周期分析数据
- `PUT /api/cycle-insight` - 更新周期分析数据

### 持仓记录
- `GET /api/positions` - 获取持仓记录
- `POST /api/positions` - 创建持仓记录
- `PUT /api/positions/:id` - 更新持仓记录
- `DELETE /api/positions/:id` - 删除持仓记录

## 数据来源

- **实时行情**: 新浪财经API `https://hq.sinajs.cn/list={股票代码}`
- **K线数据**: 新浪财经历史数据接口
- **宏观经济**: 国家统计局官网

## 项目结构

```
yuchi/
├── frontend/          # React前端
│   ├── src/
│   │   ├── components/   # UI组件
│   │   ├── services/     # API服务
│   │   └── hooks/        # 自定义hooks
│   └── package.json
├── backend/           # Go后端
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── handlers/     # HTTP handlers
│   │   ├── services/     # 业务逻辑
│   │   ├── models/       # 数据模型
│   │   └── repository/   # 数据访问
│   └── go.mod
└── docker-compose.yml
```

## 注意事项

- 新浪财经API有频率限制，建议适当缓存
- 本系统仅供个人投资研究使用，不构成任何投资建议
- 投资有风险，入市需谨慎

## License

MIT