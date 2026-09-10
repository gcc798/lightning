import { useRef } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Modal, Tag } from '@/components/ui';
import { DeleteOutlined } from '@/utils/icons';
import { BasicTable, type BasicTableRef } from '@/components/common/BasicTable';
import { PermissionGate } from '@/components/common/PermissionGate';
import type { FormSchema } from '@/types/form';
import type { LoginLogRecord } from '@/types/system';
import { loginLogApi } from '@/api/loginlog';
import { isZeroStatus } from '@/utils/number';

const searchSchemas: FormSchema[] = [
  {
    name: 'userName',
    label: '用户名',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    name: 'ipaddr',
    label: 'IP 地址',
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

export default function LoginLogPage() {
  const { message } = App.useApp();
  const tableRef = useRef<BasicTableRef<LoginLogRecord>>(null);

  const columns: ColumnsType<LoginLogRecord> = [
    { title: '用户名', dataIndex: 'userName', width: 140 },
    { title: 'IP 地址', dataIndex: 'ipaddr', width: 160 },
    { title: '登录地点', dataIndex: 'loginLocation', width: 160 },
    { title: '浏览器', dataIndex: 'browser', width: 140 },
    { title: '操作系统', dataIndex: 'os', width: 140 },
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
    { title: '提示消息', dataIndex: 'msg', width: 240 },
    { title: '登录时间', dataIndex: 'loginTime', width: 180 },
  ];

  return (
    <BasicTable<LoginLogRecord>
      ref={tableRef}
      columns={columns}
      fetchData={loginLogApi.page}
      rowKey="id"
      searchSchemas={searchSchemas}
      selectable={false}
      scroll={{ x: 1200 }}
      toolbar={
        <PermissionGate permission="login_log.delete">
          <Button
            danger
            icon={<DeleteOutlined />}
            onClick={() => {
              Modal.confirm({
                title: '确认清空',
                content: '确定要清空所有登录日志吗？此操作不可恢复。',
                onOk: async () => {
                  await loginLogApi.clean();
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
  );
}
