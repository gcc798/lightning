import { useRef, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Descriptions, Modal, Tag } from '@/components/ui';
import { DeleteOutlined } from '@/utils/icons';
import { BasicTable, type BasicTableRef } from '@/components/common/BasicTable';
import { PermissionGate } from '@/components/common/PermissionGate';
import type { FormSchema } from '@/types/form';
import type { OperLogRecord } from '@/types/system';
import { operLogApi } from '@/api/operlog';
import { isZeroStatus } from '@/utils/number';

const searchSchemas: FormSchema[] = [
  {
    name: 'title',
    label: '系统模块',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    name: 'operName',
    label: '用户',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    name: 'businessType',
    label: '业务类型',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    name: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 8 },
    props: {
      allowClear: true,
      options: [
        { label: '成功', value: 0 },
        { label: '失败', value: 1 },
      ],
    },
  },
];

function getMethodColor(method?: string) {
  const map: Record<string, string> = {
    GET: 'blue',
    POST: 'green',
    PUT: 'orange',
    DELETE: 'red',
    PATCH: 'purple',
  };
  return method ? map[method] ?? 'default' : 'default';
}

export default function OperLogPage() {
  const { message } = App.useApp();
  const tableRef = useRef<BasicTableRef<OperLogRecord>>(null);
  const [detailOpen, setDetailOpen] = useState(false);
  const [currentRecord, setCurrentRecord] = useState<OperLogRecord>();

  const columns: ColumnsType<OperLogRecord> = [
    { title: '系统模块', dataIndex: 'title', width: 140 },
    { title: '业务类型', dataIndex: 'businessType', width: 120 },
    {
      title: '请求方式',
      dataIndex: 'requestMethod',
      width: 110,
      render: (value) => <Tag color={getMethodColor(value)}>{value}</Tag>,
    },
    { title: 'URL', dataIndex: 'operUrl', width: 260 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (value) => (
        <Tag color={isZeroStatus(value) ? 'success' : 'error'}>
          {isZeroStatus(value) ? '成功' : '失败'}
        </Tag>
      ),
    },
    { title: '耗时(ms)', dataIndex: 'costTime', width: 110 },
    { title: '终端类型', dataIndex: 'deviceType', width: 120 },
    { title: 'IP', dataIndex: 'operIp', width: 140 },
    { title: '地点', dataIndex: 'operLocation', width: 160 },
    { title: '用户', dataIndex: 'operName', width: 160 },
    { title: '时间', dataIndex: 'operTime', width: 180 },
    {
      title: '操作',
      key: 'action',
      width: 100,
      fixed: 'right',
      render: (_, record) => (
        <Button
          size="small"
          type="link"
          onClick={() => {
            setCurrentRecord(record);
            setDetailOpen(true);
          }}
        >
          详情
        </Button>
      ),
    },
  ];

  return (
    <>
      <BasicTable<OperLogRecord>
        ref={tableRef}
        columns={columns}
        fetchData={operLogApi.page}
        rowKey="id"
        searchSchemas={searchSchemas}
      selectable={false}
      scroll={{ x: 1800 }}
      toolbar={
        <PermissionGate permission="oper_log.delete">
          <Button
            danger
            icon={<DeleteOutlined />}
            onClick={() => {
              Modal.confirm({
                title: '确认清空',
                content: '确定要清空所有操作日志吗？此操作不可恢复。',
                onOk: async () => {
                  await operLogApi.clean();
                  message.success('清空成功');
                  tableRef.current?.reload();
                },
              });
            }}
          >
            清空日志
          </Button>
        </PermissionGate>
      }
    />

      <Modal
        footer={null}
        open={detailOpen}
        title="操作日志详情"
        width={900}
        onCancel={() => setDetailOpen(false)}
      >
        <Descriptions bordered column={2}>
          <Descriptions.Item label="系统模块">{currentRecord?.title}</Descriptions.Item>
          <Descriptions.Item label="业务类型">{currentRecord?.businessType}</Descriptions.Item>
          <Descriptions.Item label="请求方式">{currentRecord?.requestMethod}</Descriptions.Item>
          <Descriptions.Item label="终端类型">{currentRecord?.deviceType}</Descriptions.Item>
          <Descriptions.Item label="用户">{currentRecord?.operName}</Descriptions.Item>
          <Descriptions.Item label="IP">{currentRecord?.operIp}</Descriptions.Item>
          <Descriptions.Item label="地点">{currentRecord?.operLocation || '-'}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={isZeroStatus(currentRecord?.status) ? 'success' : 'error'}>
              {isZeroStatus(currentRecord?.status) ? '成功' : '失败'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="时间" span={2}>
            {currentRecord?.operTime}
          </Descriptions.Item>
          <Descriptions.Item label="耗时">{currentRecord?.costTime}ms</Descriptions.Item>
          <Descriptions.Item label="请求 URL" span={2}>
            {currentRecord?.operUrl}
          </Descriptions.Item>
          <Descriptions.Item label="调用方法" span={2}>
            {currentRecord?.method}
          </Descriptions.Item>
          <Descriptions.Item label="请求参数" span={2}>
            <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{currentRecord?.operParam || '-'}</pre>
          </Descriptions.Item>
          <Descriptions.Item label="返回结果" span={2}>
            <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{currentRecord?.jsonResult || '-'}</pre>
          </Descriptions.Item>
          <Descriptions.Item label="错误信息" span={2}>
            <span style={{ color: '#ff4d4f' }}>{currentRecord?.errorMsg || '-'}</span>
          </Descriptions.Item>
          <Descriptions.Item label="User-Agent" span={2}>
            {currentRecord?.userAgent || '-'}
          </Descriptions.Item>
        </Descriptions>
      </Modal>
    </>
  );
}
