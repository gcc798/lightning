import {
  ApiOutlined,
  ArrowDownOutlined,
  ArrowUpOutlined,
  CloudOutlined,
  DashboardOutlined,
  EyeOutlined,
  FieldTimeOutlined,
  ThunderboltOutlined,
  UserSwitchOutlined,
} from '@/utils/icons';
import { Card, Col, Progress, Row, Space, Tag, Typography } from '@/components/ui';

const visitSeries = [
  { time: '00:00', value: 128 },
  { time: '03:00', value: 86 },
  { time: '06:00', value: 144 },
  { time: '09:00', value: 386 },
  { time: '12:00', value: 468 },
  { time: '15:00', value: 532 },
  { time: '18:00', value: 491 },
  { time: '21:00', value: 354 },
];

const metricCards = [
  {
    title: '今日访问',
    value: '2,489',
    delta: '+18.4%',
    trend: 'up',
    icon: <EyeOutlined />,
    tone: 'blue',
  },
  {
    title: '在线用户',
    value: '136',
    delta: '+9.2%',
    trend: 'up',
    icon: <UserSwitchOutlined />,
    tone: 'green',
  },
  {
    title: '接口均耗时',
    value: '86ms',
    delta: '-12.6%',
    trend: 'down',
    icon: <ThunderboltOutlined />,
    tone: 'orange',
  },
  {
    title: 'API 成功率',
    value: '99.2%',
    delta: '+0.8%',
    trend: 'up',
    icon: <ApiOutlined />,
    tone: 'cyan',
  },
];

const weatherItems = [
  { label: '体感', value: '24°C' },
  { label: '湿度', value: '62%' },
  { label: '风速', value: '东南风 3级' },
];

const channelStats = [
  { name: '后台直达', value: 46, color: '#087ea4' },
  { name: '移动端入口', value: 31, color: '#16a34a' },
  { name: '小程序', value: 23, color: '#f59e0b' },
];

const apiHealth = [
  { name: '用户与角色', value: 98 },
  { name: '菜单与权限', value: 96 },
  { name: '日志查询', value: 91 },
  { name: '附件服务', value: 94 },
];

function buildLinePath(values: number[], width: number, height: number) {
  const min = Math.min(...values);
  const max = Math.max(...values);
  const range = max - min || 1;

  return values
    .map((value, index) => {
      const x = (index / (values.length - 1)) * width;
      const y = height - ((value - min) / range) * height;
      return `${index === 0 ? 'M' : 'L'} ${x.toFixed(2)} ${y.toFixed(2)}`;
    })
    .join(' ');
}

function buildAreaPath(values: number[], width: number, height: number) {
  const linePath = buildLinePath(values, width, height);
  return `${linePath} L ${width} ${height} L 0 ${height} Z`;
}

export default function DashboardPage() {
  const chartWidth = 720;
  const chartHeight = 260;
  const values = visitSeries.map((item) => item.value);
  const linePath = buildLinePath(values, chartWidth, chartHeight);
  const areaPath = buildAreaPath(values, chartWidth, chartHeight);
  const peak = visitSeries.reduce((current, item) => (item.value > current.value ? item : current), visitSeries[0]);

  return (
    <div className="dashboard-page">
      <section className="dashboard-hero">
        <div>
          <Tag className="dashboard-eyebrow" color="processing">
            Operations Center
          </Tag>
          <Typography.Title className="dashboard-title" level={2}>
            今日运营概览
          </Typography.Title>
          <Typography.Paragraph className="dashboard-subtitle">
            聚合访问趋势、接口健康度、天气和关键服务状态。当前为前端内置演示数据，后续可直接接入统计接口。
          </Typography.Paragraph>
        </div>
        <div className="dashboard-hero-status">
          <span className="dashboard-status-dot" />
          <span>系统运行正常</span>
        </div>
      </section>

      <Row gutter={[16, 16]}>
        {metricCards.map((item) => (
          <Col key={item.title} lg={6} md={12} sm={24} xs={24}>
            <Card className={`metric-card metric-card-${item.tone}`} variant="borderless">
              <div className="metric-card-top">
                <span className="metric-icon">{item.icon}</span>
                <Tag className="metric-delta" color={item.trend === 'up' ? 'success' : 'gold'}>
                  {item.trend === 'up' ? <ArrowUpOutlined /> : <ArrowDownOutlined />} {item.delta}
                </Tag>
              </div>
              <div className="metric-value">{item.value}</div>
              <div className="metric-title">{item.title}</div>
            </Card>
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]}>
        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
          <Card className="page-card dashboard-chart-card" variant="borderless">
            <div className="dashboard-card-head">
              <div>
                <Typography.Title level={4}>访问量趋势</Typography.Title>
                <Typography.Text type="secondary">按 3 小时聚合，峰值 {peak.time} / {peak.value} 次</Typography.Text>
              </div>
              <Tag color="blue">今日</Tag>
            </div>

            <div className="visit-chart">
              <svg viewBox={`0 0 ${chartWidth} ${chartHeight}`} preserveAspectRatio="none">
                <defs>
                  <linearGradient id="visitArea" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stopColor="var(--app-accent)" stopOpacity="0.26" />
                    <stop offset="100%" stopColor="var(--app-accent)" stopOpacity="0.02" />
                  </linearGradient>
                </defs>
                <path className="visit-chart-grid" d="M 0 65 L 720 65 M 0 130 L 720 130 M 0 195 L 720 195" />
                <path className="visit-chart-area" d={areaPath} />
                <path className="visit-chart-line" d={linePath} />
                {visitSeries.map((item, index) => {
                  const x = (index / (visitSeries.length - 1)) * chartWidth;
                  const min = Math.min(...values);
                  const max = Math.max(...values);
                  const y = chartHeight - ((item.value - min) / (max - min || 1)) * chartHeight;
                  return <circle key={item.time} className="visit-chart-dot" cx={x} cy={y} r="5" />;
                })}
              </svg>
              <div className="visit-chart-axis">
                {visitSeries.map((item) => (
                  <span key={item.time}>{item.time}</span>
                ))}
              </div>
            </div>
          </Card>
        </Col>

        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
          <Card className="weather-card" variant="borderless">
            <div className="weather-glow" />
            <div className="weather-top">
              <div>
                <Typography.Text className="weather-city">上海 · 今日</Typography.Text>
                <div className="weather-temp">23°C</div>
                <Typography.Text className="weather-desc">多云，适合外出巡检</Typography.Text>
              </div>
              <CloudOutlined className="weather-icon" />
            </div>
            <div className="weather-grid">
              {weatherItems.map((item) => (
                <div key={item.label}>
                  <span>{item.label}</span>
                  <strong>{item.value}</strong>
                </div>
              ))}
            </div>
          </Card>
        </Col>
      </Row>

      <Row gutter={[12, 12]}>
        <Col xl={7} lg={12} md={24} sm={24} xs={24}>
          <Card className="page-card dashboard-panel" variant="borderless">
            <div className="dashboard-card-head">
              <Typography.Title level={4}>访问来源</Typography.Title>
              <DashboardOutlined />
            </div>
            <Space direction="vertical" size={14} style={{ width: '100%' }}>
              {channelStats.map((item) => (
                <div key={item.name} className="channel-row">
                  <div className="channel-label">
                    <span style={{ background: item.color }} />
                    {item.name}
                  </div>
                  <Progress percent={item.value} showInfo={false} strokeColor={item.color} />
                  <strong>{item.value}%</strong>
                </div>
              ))}
            </Space>
          </Card>
        </Col>

        <Col xl={9} lg={12} md={24} sm={24} xs={24}>
          <Card className="page-card dashboard-panel" variant="borderless">
            <div className="dashboard-card-head">
              <Typography.Title level={4}>接口健康度</Typography.Title>
              <ApiOutlined />
            </div>
            <Space direction="vertical" size={12} style={{ width: '100%' }}>
              {apiHealth.map((item) => (
                <div key={item.name} className="health-row">
                  <span>{item.name}</span>
                  <Progress percent={item.value} size="small" strokeColor="var(--app-accent)" />
                </div>
              ))}
            </Space>
          </Card>
        </Col>

        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
          <Card className="page-card dashboard-panel dashboard-summary-panel" variant="borderless">
            <div className="dashboard-card-head">
              <Typography.Title level={4}>运行摘要</Typography.Title>
              <FieldTimeOutlined />
            </div>
            <div className="summary-list">
              <div>
                <span>最近同步</span>
                <strong>2 分钟前</strong>
              </div>
              <div>
                <span>慢请求</span>
                <strong>3 条</strong>
              </div>
              <div>
                <span>安全事件</span>
                <strong>0 条</strong>
              </div>
              <div>
                <span>存储状态</span>
                <strong>正常</strong>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
