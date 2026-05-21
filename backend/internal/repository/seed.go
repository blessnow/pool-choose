package repository

import (
	"log"

	"github.com/yuchi/cycle-stock/internal/models"
)

// seedStocks 如果数据库中没有股票数据，自动导入默认的68只股票
func seedStocks() {
	var count int64
	if err := DB.Model(&models.Stock{}).Count(&count).Error; err != nil {
		log.Printf("[seed] 检查股票数量失败: %v", err)
		return
	}
	if count > 0 {
		return
	}

	log.Printf("[seed] 数据库为空，导入默认 %d 只股票", len(defaultStocks))
	if err := DB.Create(&defaultStocks).Error; err != nil {
		log.Printf("[seed] 导入失败: %v", err)
		return
	}
	log.Printf("[seed] 导入完成")
}

var defaultStocks = []models.Stock{
	{Code: "000707", Name: "双环科技", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 6.15, HeavyPrice: 6.0, TargetPrice: 10.0, CoreLogic: "核心逻辑：纯碱和氯化铵联合生产。纯碱受益于光伏玻璃和浮法玻璃需求。✅ 利好：低于入场价2.3%，PB接近1倍，下行空间有限。⚠️ 风险：纯碱产能过剩（天然碱法冲击），ROE仅5%。巴菲特视角：PB接近1倍有资产保护，但ROE过低。安全边际靠\"不会更差\"。"},
	{Code: "300470", Name: "中密股份", Sector: "机械密封", Industry: "通用机械 / 工业品 / 工业", IsRecommend: true, EntryPrice: 34.0, HeavyPrice: 30.0, TargetPrice: 60.0, CoreLogic: "核心逻辑：国内机械密封件龙头，进口替代逻辑。✅ 利好：高于入场价6.6%。国产替代空间大。⚠️ 风险：PE25倍不算便宜。"},
	{Code: "000819", Name: "岳阳兴长", Sector: "石化", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 14.0, HeavyPrice: 0, TargetPrice: 30.0, CoreLogic: "核心逻辑：中国石化体系内石化企业，主营石油化工产品、化工新材料与能源化工。✅ 利好：石化主业清晰，具备央企背景和产业协同。⚠️ 风险：当前交易参数尚未配置，需补充入场价、止盈价后再纳入完整策略判断。"},
	{Code: "002408", Name: "齐翔腾达", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 5.5, HeavyPrice: 4.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：碳四产业链龙头，向碳三延伸（丙烷脱氢+环氧丙烷）。受美伊冲突影响，油价高位推升产品价格，顺酐、丙烯酸等主产品价格大幅上涨。✅ 利好：产品价格全线上涨，装置开工率90%+，碳三产业链即将贡献增量。山东国资背景。⚠️ 风险：2025年1-9月归母净利-1.46亿，尚未扭亏，需Q1业绩验证。巴菲特视角：已高于入场价14.5%，但距止盈目标10元仍有空间。化工涨价周期启动则弹性极大，但亏损是硬伤。"},
	{Code: "688707", Name: "振华新材", Sector: "锂电材料", Industry: "电气部件与设备 / 工业品 / 工业", IsRecommend: true, EntryPrice: 13.0, HeavyPrice: 8.0, TargetPrice: 20.0, CoreLogic: "核心逻辑：三元正极材料企业，受益于锂电池/新能源车产业链。碳酸锂价格触底回升预期。✅ 利好：价格12.92几乎精准落在入场价13附近，提供极佳建仓窗口。锂价企稳回升则业绩弹性巨大。⚠️ 风险：当前亏损，正极材料行业竞争激烈，产能过剩。科创板流动性偏弱。巴菲特视角：亏损企业不符合传统价值投资，但属于\"困境反转\"逻辑。重仓价8元才是真正的\"捡便宜\"价位。"},
	{Code: "002100", Name: "天康生物", Sector: "农业", Industry: "农牧渔产品 / 主要消费 / 农林牧渔", IsRecommend: true, EntryPrice: 7.0, HeavyPrice: 6.0, TargetPrice: 10.0, CoreLogic: "核心逻辑：生猪养殖+动物疫苗双轮驱动。猪周期底部回升阶段。✅ 利好：当前高于入场价9%，猪周期景气上行中。疫苗业务穿越猪周期。止盈2高达20元暗示高弹性期待。⚠️ 风险：猪价波动剧烈，饲料成本挤压利润。巴菲特视角：养殖周期性极强不适合永久持有，但周期底部安全边际充足。"},
	{Code: "601118", Name: "海南橡胶", Sector: "农业", Industry: "农牧渔产品 / 主要消费 / 农林牧渔", IsRecommend: true, EntryPrice: 6.0, HeavyPrice: 4.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：中国天然橡胶龙头，全产业链布局。353万亩胶林覆盖海南17个市县。✅ 利好：高于入场价9.2%。东南亚主产区减产，价格有支撑。海南自贸港+碳中和多重概念。⚠️ 风险：2025年1-9月归母净利-2.75亿持续亏损。巴菲特视角：持续亏损巴菲特通常不碰，但国内唯一天然橡胶全产业链上市公司具稀缺性。"},
	{Code: "300034", Name: "钢研高纳", Sector: "高温合金", Industry: "航天航空 / 工业品 / 工业", IsRecommend: true, EntryPrice: 20.0, HeavyPrice: 10.0, TargetPrice: 30.0, CoreLogic: "核心逻辑：国内高温合金龙头，航空发动机核心材料供应商。✅ 利好：低于入场价20约10.7%，仍在\"折扣\"区间。军工需求稳定增长，高温合金国产化率提升。⚠️ 风险：PE较高~35倍，年初至今下跌20.4%。军工板块热度降温。巴菲特视角：35倍PE不是传统价值标的，但低于入场价可考虑左侧布局。"},
	{Code: "603977", Name: "国泰集团", Sector: "化工", Industry: "化学制品 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 13.0, HeavyPrice: 11.5, TargetPrice: 18.0, CoreLogic: "核心逻辑：国内草酸和DMC龙头，DMC是锂电池电解液关键原料。✅ 利好：高于入场价21.5%。PE仅12倍，股息率3.5%，ROE 15%——非常符合价值投资标准。⚠️ 风险：距止盈1仅13.9%空间收窄。巴菲特视角：12倍PE+3.5%股息+15%ROE，本批中最接近巴菲特审美的标的之一。"},
	{Code: "002783", Name: "凯龙股份", Sector: "民爆", Industry: "化学制品 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 9.0, HeavyPrice: 7.5, TargetPrice: 15.0, CoreLogic: "核心逻辑：民爆行业龙头之一，拓展氢能源。✅ 利好：略高于入场价1.7%，几乎在入场价附近。民爆行业格局稳定。⚠️ 风险：氢能源尚处早期，民爆增长弹性有限。巴菲特视角：民爆行业有准入壁垒（许可证），竞争格局较好。安全边际一般。"},
	{Code: "000928", Name: "中钢国际", Sector: "工程", Industry: "建筑与工程 / 工业服务 / 工业", IsRecommend: true, EntryPrice: 6.75, HeavyPrice: 5.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：中国宝武旗下冶金工程龙头，业务遍布40+国家。✅ 利好：PE仅8倍，PB<1（破净），股息率4.2%——典型深度价值股。中字头+央企改革。⚠️ 风险：工程企业应收账款大，现金流风险。钢铁下游需求承压。巴菲特视角：8倍PE+破净+4.2%股息率——被市场严重忽视的价值洼地。安全边际最高的标的之一。"},
	{Code: "002136", Name: "安纳达", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 12.0, HeavyPrice: 10.0, TargetPrice: 20.0, CoreLogic: "核心逻辑：钛白粉企业，受益于涂料、塑料下游需求回暖及出口增长。✅ 利好：略高于入场价2.8%。钛白粉价格受益于房地产竣工端回暖。⚠️ 风险：钛白粉产能过剩，龙佰集团等对中小企业构成挤压。巴菲特视角：典型大宗商品缺乏定价权。体量较小竞争力有限。但化工涨价周期启动弹性也大。"},
	{Code: "600500", Name: "中化国际", Sector: "化工", Industry: "原材料 / 基础化工 / 化学制品", IsRecommend: true, EntryPrice: 4.0, HeavyPrice: 3.8, TargetPrice: 7.0, CoreLogic: "核心逻辑：中化集团旗下精细化工平台。央企整合预期。✅ 利好：PE~10，PB仅0.7（严重破净），股息率3.8%。央企背景+化工涨价+市值管理考核。⚠️ 风险：业务分散，聚焦度不足。巴菲特视角：0.7倍PB央企平台，\"以7折买入1块钱资产\"。3.8%股息+破净修复，安全边际数一数二。"},
	{Code: "000731", Name: "四川美丰", Sector: "化工", Industry: "农用化工 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 7.2, HeavyPrice: 5.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：氮肥龙头之一，受益于尿素价格上涨（油价传导）。✅ 利好：精准踩线入场价。尿素价格受油价推升走强。PE12倍合理。⚠️ 风险：化肥行业受政策调控影响大，出口限制。巴菲特视角：化肥是刚需，12倍PE+3%股息+涨价周期，风险收益比合理。"},
	{Code: "601618", Name: "中国中冶", Sector: "工程", Industry: "建筑与工程 / 工业服务 / 工业", IsRecommend: true, EntryPrice: 3.15, HeavyPrice: 3.14, TargetPrice: 5.5, CoreLogic: "核心逻辑：中国基建国家队成员，极低估值+高股息。✅ 利好：PE仅5倍，PB仅0.5（极度破净），股息率高达5.5%！估值最低、股息最高。央企市值管理。⚠️ 风险：低于入场价5.1%。房地产拖累，应收账款周期长。巴菲特视角：5倍PE、0.5倍PB、5.5%股息——安全边际的教科书案例。强烈推荐关注。"},
	{Code: "603227", Name: "雪峰科技", Sector: "军工/民爆", Industry: "化学制品 / 基础化工 / 原材料", IsRecommend: true, EntryPrice: 7.5, HeavyPrice: 7.5, TargetPrice: 12.0, CoreLogic: "核心逻辑：民爆+军工企业，新疆地区龙头。✅ 利好：高于入场价34.1%，浮盈丰厚。新疆基建+军工双催化。⚠️ 风险：PE~25不便宜，军工板块回调压力。巴菲特视角：准入壁垒好但25倍PE不符合价值审美。已有34%浮盈，建议逐步减仓锁定利润。"},
	{Code: "301058", Name: "中粮科工", Sector: "工程", Industry: "建筑与工程 / 工业服务 / 工业", IsRecommend: true, EntryPrice: 10.0, HeavyPrice: 10.0, TargetPrice: 18.0, CoreLogic: "核心逻辑：中粮集团旗下粮油食品工程设计龙头。受益于粮食安全政策。✅ 利好：央企背景+粮食安全国策。止盈目标18，向上空间仍然明确。⚠️ 风险：当前低于入场价1%，年初至今跌9.9%。估值PE25偏高。巴菲特视角：粮食安全是长期逻辑，但25倍PE缺乏安全边际。"},
	{Code: "600866", Name: "星湖科技", Sector: "化工", Industry: "味精+氨基酸+生物发酵", IsRecommend: true, EntryPrice: 5.8, HeavyPrice: 4.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：味精及氨基酸生产企业，生物发酵领域布局。受益于饲料氨基酸需求增长。✅ 利好：氨基酸市场需求稳定增长，生物发酵技术积累深厚。⚠️ 风险：味精行业增长有限，原材料价格波动。巴菲特视角：消费品属性+技术壁垒，估值合理时具备安全边际。"},
	{Code: "600459", Name: "贵研铂业", Sector: "贵金属", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 18.5, HeavyPrice: 12.0, TargetPrice: 25.0, CoreLogic: "核心逻辑：中国铂族金属深加工龙头，受益于贵金属价格上涨和新能源催化剂需求。✅ 利好：高于入场价4.2%。贵金属受避险情绪驱动上涨。央企+五矿背景。⚠️ 风险：铂族金属价格波动大，受全球经济影响。巴菲特视角：贵金属加工企业有一定护城河，但需跟踪金属价格走势。"},
	{Code: "002237", Name: "恒邦股份", Sector: "黄金", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 12.0, HeavyPrice: 9.0, TargetPrice: 18.0, CoreLogic: "核心逻辑：国内黄金冶炼龙头之一，受益于金价上行周期。✅ 利好：高于入场价22.9%。国际金价持续走高，黄金股弹性释放中。⚠️ 风险：黄金价格高位波动风险。巴菲特视角：黄金股在避险需求旺盛时表现好。当前浮盈可观。"},
	{Code: "000554", Name: "泰山石油", Sector: "能源", Industry: "石油天然气 / 能源 / 化工", IsRecommend: false, EntryPrice: 6.5, HeavyPrice: 6.0, TargetPrice: 20.0, CoreLogic: "核心逻辑：中石化旗下成品油零售企业，受益于油价上行和央企改革。✅ 利好：高于入场价31.5%。中石化注资预期+油价走高。止盈目标20元空间巨大。⚠️ 风险：PE较高，更多依赖概念炒作而非基本面。巴菲特视角：估值较高不符合价值投资，但央企改革注入预期提供想象空间。"},
	{Code: "600328", Name: "中盐化工", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 8.0, HeavyPrice: 7.5, TargetPrice: 0, CoreLogic: "核心逻辑：中盐集团旗下盐化工平台，纯碱+PVC双主业。✅ 利好：高于入场价7%。央企背景+PB约1倍+3%股息。⚠️ 风险：纯碱和PVC均面临产能过剩。巴菲特视角：央企+低估值+合理股息，安全边际可接受。"},
	{Code: "600963", Name: "岳阳林纸", Sector: "林业/造纸", Industry: "纸类与林业产品 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 4.0, HeavyPrice: 3.5, TargetPrice: 0, CoreLogic: "核心逻辑：造纸+碳汇林业双概念。受益于碳中和政策。✅ 利好：高于入场价18%。碳汇概念+央企背景。⚠️ 风险：造纸行业景气度低迷。碳交易市场尚未成熟。巴菲特视角：碳汇是长期概念，但当前业绩支撑不足。"},
	{Code: "600348", Name: "华阳股份", Sector: "煤炭", Industry: "煤炭 / 能源 / 采掘", IsRecommend: false, EntryPrice: 6.4, HeavyPrice: 6.0, TargetPrice: 12.0, CoreLogic: "核心逻辑：传统煤炭+新能源储能转型。飞轮储能和钠离子电池提供成长性。✅ 利好：高于入场价40.6%！PE仅8倍，股息率5%，PB<1。煤价高位+新能源概念加持。⚠️ 风险：煤炭长期需求下行。新能源业务尚处早期。巴菲特视角：8倍PE+5%股息+破净——典型的深度价值+成长转型标的。浮盈已非常丰厚。"},
	{Code: "603639", Name: "湖南海利", Sector: "农药", Industry: "农用化工 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 6.0, HeavyPrice: 5.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：农药制剂企业，受益于农药涨价周期。✅ 利好：高于入场价123%！止盈1已触达。涨幅惊人。⚠️ 风险：涨幅过大，回调风险极高。巴菲特视角：已大幅超越止盈目标，强烈建议锁定利润。"},
	{Code: "600299", Name: "安迪苏", Sector: "饲料添加剂", Industry: "主要消费 / 农牧渔产品 / 化工", IsRecommend: false, EntryPrice: 9.0, HeavyPrice: 8.0, TargetPrice: 15.0, CoreLogic: "核心逻辑：全球蛋氨酸龙头，动物营养刚需品。中化集团旗下。✅ 利好：高于入场价48%。蛋氨酸价格触底回升。全球寡头垄断格局。⚠️ 风险：距止盈15元仅12.6%。蛋氨酸新增产能压力。巴菲特视角：寡头垄断格局巴菲特最喜欢，但当前估值偏高。接近止盈区间。"},
	{Code: "000923", Name: "河钢资源", Sector: "矿业", Industry: "黑色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 13.5, HeavyPrice: 13.5, TargetPrice: 25.0, CoreLogic: "核心逻辑：海外矿产资源标的，受益于大宗商品上行周期。✅ 利好：高于入场价33.4%。铁矿石和铜价走强。⚠️ 风险：海外政治风险，汇率波动。"},
	{Code: "600583", Name: "海油工程", Sector: "油服", Industry: "能源设备与服务 / 能源 / 采掘", IsRecommend: false, EntryPrice: 5.0, HeavyPrice: 4.5, TargetPrice: 7.5, CoreLogic: "核心逻辑：中海油旗下海洋油服龙头。油价高位带动资本开支增加。✅ 利好：高于入场价38%。油价突破100美元利好油服。⚠️ 风险：距止盈仅8.7%。海洋工程项目周期长。"},
	{Code: "601600", Name: "中国铝业", Sector: "铝业", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 6.5, HeavyPrice: 6.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：中国铝业龙头，止盈1和止盈2均已触达。✅ 利好：高于入场价79.7%！铝价高位。⚠️ 风险：已远超止盈目标。巴菲特视角：止盈已触达，建议逐步兑现利润。"},
	{Code: "600509", Name: "天富能源", Sector: "电力", Industry: "电力公用事业 / 公用事业 / 电力", IsRecommend: false, EntryPrice: 5.5, HeavyPrice: 5.5, TargetPrice: 10.0, CoreLogic: "核心逻辑：新疆热电联产+煤化工企业。受益于电力需求增长和煤化工景气。✅ 利好：高于入场价48.2%。⚠️ 风险：新疆区域性企业，流动性一般。"},
	{Code: "600378", Name: "昊华科技", Sector: "化工", Industry: "原材料 / 基础化工 / 化学制品", IsRecommend: false, EntryPrice: 23.0, HeavyPrice: 23.0, TargetPrice: 43.0, CoreLogic: "核心逻辑：央企氟化工+电子气体平台。受益于半导体和新能源。✅ 利好：高于入场价32.7%。电子气体国产替代逻辑强。⚠️ 风险：估值不算便宜。"},
	{Code: "000629", Name: "钒钛股份", Sector: "钒钛", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 2.4, HeavyPrice: 2.2, TargetPrice: 5.0, CoreLogic: "核心逻辑：国内钒钛资源龙头。钒电池储能带来新需求。✅ 利好：高于入场价37.1%。钒电池储能前景广阔。⚠️ 风险：钒电池商业化进度不确定。"},
	{Code: "601061", Name: "中信金属", Sector: "稀有金属", Industry: "工业 / 工业服务 / 工业贸易经销商", IsRecommend: false, EntryPrice: 7.0, HeavyPrice: 7.0, TargetPrice: 14.0, CoreLogic: "核心逻辑：中信集团旗下稀有金属平台。铌铁全球龙头。✅ 利好：高于入场价65.1%。稀有金属战略地位提升。⚠️ 风险：距止盈14元仅21.1%。"},
	{Code: "002140", Name: "东华科技", Sector: "化工工程", Industry: "建筑与工程 / 工业服务 / 工业", IsRecommend: false, EntryPrice: 7.5, HeavyPrice: 7.5, TargetPrice: 15.0, CoreLogic: "核心逻辑：化工工程设计龙头，受益于化工行业新一轮资本开支周期。✅ 利好：高于入场价70%。化工复苏带动工程设计需求。⚠️ 风险：距止盈仅17.6%，空间收窄。"},
	{Code: "600230", Name: "沧州大化", Sector: "化工", Industry: "化学制品 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 10.0, HeavyPrice: 9.0, TargetPrice: 20.0, CoreLogic: "核心逻辑：TDI国内重要生产商。受益于聚氨酯需求回升。✅ 利好：高于入场价89.8%！距止盈20元仅5.4%。⚠️ 风险：即将触及止盈线。巴菲特视角：接近止盈，应考虑兑现利润。"},
	{Code: "600282", Name: "南钢股份", Sector: "钢铁", Industry: "黑色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 4.0, HeavyPrice: 3.3, TargetPrice: 5.5, CoreLogic: "核心逻辑：特钢龙头，止盈1已触达。✅ 利好：高于入场价37.5%。PB仅0.7+5%股息。⚠️ 风险：止盈1已达，钢铁行业产能过剩。止盈1已触达。"},
	{Code: "002258", Name: "利尔化学", Sector: "农药", Industry: "农用化工 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 8.0, HeavyPrice: 8.0, TargetPrice: 15.0, CoreLogic: "核心逻辑：国内除草剂龙头，草铵膦全球重要供应商。✅ 利好：高于入场价81.9%。距止盈仅3.1%。⚠️ 风险：即将触及止盈。接近止盈，建议减仓。"},
	{Code: "000830", Name: "鲁西化工", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 9.0, HeavyPrice: 9.0, TargetPrice: 18.0, CoreLogic: "核心逻辑：化工新材料龙头，产品线丰富。受益于化工周期反转。✅ 利好：高于入场价76.1%。PE仅10倍。化工涨价弹性大。⚠️ 风险：距止盈13.6%，空间不大。"},
	{Code: "002430", Name: "杭氧股份", Sector: "工业气体", Industry: "原材料 / 基础化工 / 化学原料", IsRecommend: false, EntryPrice: 16.0, HeavyPrice: 15.0, TargetPrice: 35.0, CoreLogic: "核心逻辑：国内空分设备龙头，工业气体运营转型。✅ 利好：高于入场价72.8%。工业气体长期增长确定性高。⚠️ 风险：估值不算便宜。"},
	{Code: "600486", Name: "扬农化工", Sector: "农药", Industry: "农用化工 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 45.0, HeavyPrice: 45.0, TargetPrice: 90.0, CoreLogic: "核心逻辑：先正达旗下农药龙头，菊酯全球最大生产商。✅ 利好：高于入场价66.9%。农药涨价周期+全球领先地位。⚠️ 风险：距止盈19.8%。"},
	{Code: "601958", Name: "金钼股份", Sector: "钼业", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 10.0, HeavyPrice: 6.0, TargetPrice: 20.0, CoreLogic: "核心逻辑：国内钼业绝对龙头，全球重要钼生产商。✅ 利好：高于入场价86.8%。钼价高位运行。ROE约20%。⚠️ 风险：距止盈仅7.1%。"},
	{Code: "000960", Name: "锡业股份", Sector: "锡业", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 14.0, HeavyPrice: 8.0, TargetPrice: 25.0, CoreLogic: "核心逻辑：全球锡业龙头，止盈1和止盈2均已触达。✅ 利好：高于入场价120%！⚠️ 双止盈均已触达，建议兑现利润。"},
	{Code: "000962", Name: "东方钽业", Sector: "稀有金属", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 8.0, HeavyPrice: 12.0, TargetPrice: 20.0, CoreLogic: "核心逻辑：钽铌深加工龙头，止盈1和2均已大幅超越。✅ 利好：高于入场价425.9%！惊人涨幅。⚠️ 远超止盈，极高回调风险，强烈建议兑现。"},
	{Code: "600426", Name: "华鲁恒升", Sector: "化工", Industry: "基础化工 / 原材料 / 化学原料", IsRecommend: false, EntryPrice: 20.0, HeavyPrice: 15.0, TargetPrice: 40.0, CoreLogic: "核心逻辑：煤化工行业标杆企业，成本控制能力全行业领先。✅ 利好：高于入场价78.2%。化工龙头中ROE最稳定的之一。⚠️ 风险：距止盈12.3%。巴菲特视角：行业龙头+稳定ROE+合理估值，最接近巴菲特标准。"},
	{Code: "600379", Name: "宝光股份", Sector: "电气", Industry: "电气部件与设备 / 工业品 / 工业", IsRecommend: false, EntryPrice: 7.0, HeavyPrice: 6.0, TargetPrice: 20.0, CoreLogic: "核心逻辑：真空灭弧室龙头，受益于电网投资。✅ 利好：高于入场价67.4%。电网投资持续增长。⚠️ 风险：市值较小，波动性大。"},
	{Code: "600456", Name: "宝钛股份", Sector: "钛合金", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 17.0, HeavyPrice: 16.0, TargetPrice: 34.0, CoreLogic: "核心逻辑：国内钛材加工绝对龙头，军工+民用双市场。✅ 利好：高于入场价95.8%。距止盈仅2.2%！⚠️ 风险：几乎触及止盈线。即将止盈，建议兑现。"},
	{Code: "600961", Name: "株冶集团", Sector: "锌冶炼", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 7.5, HeavyPrice: 7.5, TargetPrice: 13.0, CoreLogic: "核心逻辑：锌冶炼龙头，止盈1和2均已触达。✅ 利好：高于入场价129.9%。⚠️ 双止盈均已触达，建议兑现。"},
	{Code: "300527", Name: "ST应急", Sector: "应急装备", Industry: "航天航空 / 工业品 / 工业", IsRecommend: false, EntryPrice: 7.5, HeavyPrice: 4.0, TargetPrice: 12.0, CoreLogic: "核心逻辑：应急救援装备企业，ST股有摘帽预期。✅ 利好：高于入场价7.3%。摘帽预期+应急产业政策支持。⚠️ 风险：ST股风险极高。亏损企业。巴菲特视角：ST股完全不符合价值投资框架，谨慎。"},
	{Code: "600301", Name: "华锡有色", Sector: "锡锌", Industry: "原材料 / 基础材料 / 有色金属", IsRecommend: false, EntryPrice: 14.0, HeavyPrice: 14.0, TargetPrice: 25.0, CoreLogic: "核心逻辑：锡锌冶炼企业，止盈1和2均已触达。✅ 利好：高于入场价235%！ROE约35%极高。⚠️ 远超止盈，强烈建议兑现。"},
	{Code: "688750", Name: "金天钛业", Sector: "钛合金", Industry: "原材料 / 基础材料 / 有色金属", IsRecommend: false, EntryPrice: 8.0, HeavyPrice: 8.0, TargetPrice: 0, CoreLogic: "核心逻辑：军用钛合金材料企业。✅ 利好：高于入场价119.4%。军工钛合金国产替代。⚠️ 风险：未设止盈目标，涨幅巨大。"},
	{Code: "600549", Name: "厦门钨业", Sector: "钨钼稀土", Industry: "原材料 / 基础材料 / 有色金属", IsRecommend: false, EntryPrice: 0, HeavyPrice: 12.0, TargetPrice: 0, CoreLogic: "核心逻辑：钨钼稀土龙头+锂电正极材料。多元化布局。✅ 利好：无入场价设定但远高于重仓价12。止盈2已触达。⚠️ 止盈2已触达。"},
	{Code: "600392", Name: "盛和资源", Sector: "稀土", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 0, HeavyPrice: 7.0, TargetPrice: 0, CoreLogic: "核心逻辑：稀土采选冶炼企业。受益于稀土战略资源地位提升。✅ 利好：远高于重仓价7。止盈2已触达。稀土战略意义重大。⚠️ 止盈2已触达。"},
	{Code: "600111", Name: "北方稀土", Sector: "稀土", Industry: "有色金属 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 0, HeavyPrice: 10.0, TargetPrice: 0, CoreLogic: "核心逻辑：全球最大稀土企业。✅ 利好：远高于重仓价10。止盈2已触达。稀土龙头地位无可撼动。⚠️ 止盈2已触达，涨幅巨大。"},
	{Code: "600737", Name: "中粮糖业", Sector: "食品", Industry: "食品 / 食品饮料与烟草 / 主要消费", IsRecommend: false, EntryPrice: 0, HeavyPrice: 6.0, TargetPrice: 0, CoreLogic: "核心逻辑：中粮旗下食糖龙头。✅ 利好：远高于重仓价6。糖价高位运行。央企背景。⚠️ 风险：糖价周期性波动。"},
	{Code: "002057", Name: "中钢天源", Sector: "磁材料", Industry: "原材料 / 基础材料 / 合成金属", IsRecommend: false, EntryPrice: 7.0, HeavyPrice: 6.0, TargetPrice: 15.0, CoreLogic: "核心逻辑：中钢集团旗下磁性材料企业。✅ 利好：高于入场价31.4%。磁性材料受益于新能源。⚠️ 风险：估值偏高。"},
	{Code: "601568", Name: "北元化工", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 4.0, HeavyPrice: 3.5, TargetPrice: 0, CoreLogic: "核心逻辑：PVC龙头之一，煤-电-化一体化循环经济。✅ 利好：高于入场价6.3%。PE8倍+PB0.8破净+4%股息率。⚠️ 风险：PVC产能过剩严重。巴菲特视角：破净+4%股息+低PE，深度价值。"},
	{Code: "000738", Name: "航发控制", Sector: "军工", Industry: "航天航空 / 工业品 / 工业", IsRecommend: false, EntryPrice: 15.0, HeavyPrice: 0, TargetPrice: 0, CoreLogic: "核心逻辑：航空发动机控制系统垄断供应商。✅ 利好：高于入场价42.6%。军工核心赛道。⚠️ 风险：PE50倍极贵。巴菲特视角：50倍PE完全不符合价值投资，但军工垄断有护城河。"},
	{Code: "002092", Name: "中泰化学", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 4.0, HeavyPrice: 0, TargetPrice: 0, CoreLogic: "核心逻辑：新疆PVC+烧碱龙头。资源禀赋优势明显。✅ 利好：高于入场价66.3%。PB约1倍。⚠️ 风险：PVC行业产能过剩。"},
	{Code: "000657", Name: "中钨高新", Sector: "钨业", Industry: "基础材料 / 原材料 / 有色金属", IsRecommend: false, EntryPrice: 6.0, HeavyPrice: 0, TargetPrice: 18.0, CoreLogic: "核心逻辑：中国五矿旗下钨业龙头。✅ 利好：高于入场价669.5%！止盈1已大幅超越。⚠️ 涨幅惊人，远超止盈，极高风险。"},
	{Code: "002556", Name: "辉隆股份", Sector: "农资", Industry: "工业 / 工业服务 / 工业贸易经销商", IsRecommend: false, EntryPrice: 5.3, HeavyPrice: 4.3, TargetPrice: 8.0, CoreLogic: "核心逻辑：安徽农资流通龙头，覆盖化肥农药种子。✅ 利好：高于入场价15.1%。农资行业刚需属性。⚠️ 风险：农资流通行业利润率低。"},
	{Code: "600298", Name: "安琪酵母", Sector: "食品", Industry: "食品 / 食品饮料与烟草 / 主要消费", IsRecommend: false, EntryPrice: 36.0, HeavyPrice: 0, TargetPrice: 55.0, CoreLogic: "核心逻辑：全球酵母龙头，消费品属性提供穿越周期能力。✅ 利好：高于入场价13.1%。全球酵母寡头垄断格局。⚠️ 风险：估值较高。巴菲特视角：消费品+寡头垄断，巴菲特最爱的类型。但当前PE偏高。"},
	{Code: "603970", Name: "中农立华", Sector: "农资", Industry: "原材料 / 基础化工 / 农用化工", IsRecommend: false, EntryPrice: 11.0, HeavyPrice: 0, TargetPrice: 18.0, CoreLogic: "核心逻辑：中国农资旗下农药流通龙头。✅ 利好：高于入场价15.6%。央企背景+农药涨价周期。⚠️ 风险：流通企业利润率低。"},
	{Code: "600618", Name: "氯碱化工", Sector: "化工", Industry: "化学原料 / 基础化工 / 原材料", IsRecommend: false, EntryPrice: 9.0, HeavyPrice: 0, TargetPrice: 18.0, CoreLogic: "核心逻辑：上海国资PVC+烧碱企业。✅ 利好：高于入场价37.6%。PE10倍+PB约1+3.5%股息。⚠️ 风险：PVC需求疲软。巴菲特视角：低PE+破净边缘+高股息，价值属性明显。"},
	{Code: "600819", Name: "耀皮玻璃", Sector: "建材", Industry: "非金属材料与制品 / 基础材料 / 原材料", IsRecommend: false, EntryPrice: 5.0, HeavyPrice: 0, TargetPrice: 8.0, CoreLogic: "核心逻辑：节能玻璃企业。受益于建筑节能政策。✅ 利好：高于入场价41.8%。PB0.8破净+4%股息。⚠️ 风险：距止盈仅12.8%。巴菲特视角：破净+高股息，深度价值。"},
	{Code: "688295", Name: "中复神鹰", Sector: "碳纤维", Industry: "原材料 / 基础化工 / 合成纤维", IsRecommend: false, EntryPrice: 18.0, HeavyPrice: 17.0, TargetPrice: 36.0, CoreLogic: "核心逻辑：国内碳纤维龙头，止盈1已触达。✅ 利好：高于入场价188.9%。碳纤维国产替代+军工+风电。⚠️ 止盈已触达，涨幅巨大。"},
	{Code: "600409", Name: "三友化工", Sector: "化工", Industry: "基础化工 / 原材料 / 合成纤维", IsRecommend: false, EntryPrice: 6.0, HeavyPrice: 5.5, TargetPrice: 0, CoreLogic: "核心逻辑：纯碱+粘胶纤维联合生产。✅ 利好：高于入场价18.2%。PB0.8破净+4%股息。⚠️ 风险：纯碱和粘胶纤维均产能过剩。巴菲特视角：破净+高股息，但行业景气度不佳。"},
	{Code: "600195", Name: "中牧股份", Sector: "畜牧", Industry: "主要消费 / 农牧渔产品 / 农林牧渔", IsRecommend: false, EntryPrice: 7.4, HeavyPrice: 7.4, TargetPrice: 20.0, CoreLogic: "核心逻辑：央企兽用疫苗龙头。受益于猪周期+动保行业。✅ 利好：新入仓（3月27日），央企+猪周期+动保。⚠️ 风险：价格数据暂未加载。巴菲特视角：需等待价格确认后再评估安全边际。"},
	{Code: "002096", Name: "易普力", Sector: "化工", Industry: "民爆器材+工程爆破服务", IsRecommend: false, EntryPrice: 10.85, HeavyPrice: 9.5, TargetPrice: 18.0, CoreLogic: "核心逻辑：民爆行业龙头企业之一，主营民爆器材生产和工程爆破服务。✅ 利好：民爆行业准入壁垒高，格局稳定。基建投资持续拉动爆破服务需求。⚠️ 风险：民爆增长弹性有限，受基建投资节奏影响。巴菲特视角：民爆牌照构成天然护城河，行业竞争格局较好。"},
}
