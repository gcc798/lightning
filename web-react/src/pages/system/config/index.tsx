import { useRef, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Popconfirm, Tag } from '@/components/ui';
import { DeleteOutlined, PlusOutlined } from '@/utils/icons';
import { BasicTable, type BasicTableRef } from '@/components/common/BasicTable';
import { JsonViewerModal } from '@/components/common/JsonViewerModal';
import { TableAction } from '@/components/common/TableAction';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { ConfigRecord } from '@/types/system';
import { configApi } from '@/api/config';
import { ConfigModal } from './ConfigModal';

const searchSchemas: FormSchema[] = [
  {
    name: 'name',
    label: '配置名称',
    component: 'Input',
    colProps: { span: 12 },
  },
  {
    name: 'code',
    label: '配置编码',
    component: 'Input',
    colProps: { span: 12 },
  },
];

function previewJson(data: unknown) {
  const text =
    typeof data === 'string' ? data : JSON.stringify(data ?? {}, null, 2);
  const singleLine = text.replace(/\s+/g, ' ').trim();
  return singleLine.length > 72 ? `${singleLine.slice(0, 72)}...` : singleLine;
}

export default function ConfigPage() {
  const { message } = App.useApp();
  const tableRef = useRef<BasicTableRef<ConfigRecord>>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [currentConfigId, setCurrentConfigId] = useState<SnowflakeId>();
  const [viewerOpen, setViewerOpen] = useState(false);
  const [viewerData, setViewerData] = useState<unknown>(null);

  const columns: ColumnsType<ConfigRecord> = [
    { title: '名称', dataIndex: 'name', width: 180 },
    {
      title: '编码',
      dataIndex: 'code',
      width: 180,
      render: (value) => <Tag color="blue">{value}</Tag>,
    },
    {
      title: '内容',
      dataIndex: 'data',
      width: 340,
      render: (value) => (
        <span
          className="json-preview-inline table-code"
          role="button"
          tabIndex={0}
          onClick={() => {
            setViewerData(value);
            setViewerOpen(true);
          }}
          onKeyDown={(event) => {
            if (event.key === 'Enter' || event.key === ' ') {
              event.preventDefault();
              setViewerData(value);
              setViewerOpen(true);
            }
          }}
        >
          {previewJson(value)}
        </span>
      ),
    },
    { title: '备注', dataIndex: 'remark', width: 220 },
    { title: '创建时间', dataIndex: 'createdTime', width: 180 },
    {
      title: '操作',
      key: 'action',
      width: 180,
      fixed: 'right',
      render: (_, record) => (
        <TableAction
          actions={[
            {
              key: 'edit',
              label: '编辑',
              onClick: () => {
                setCurrentConfigId(record.id);
                setModalOpen(true);
              },
            },
            {
              key: 'delete',
              label: '删除',
              danger: true,
              confirmTitle: '确定删除该配置吗？',
              onClick: async () => {
                await configApi.delete(record.id);
                message.success('删除成功');
                tableRef.current?.reload();
              },
            },
          ]}
        />
      ),
    },
  ];

  return (
    <>
      <BasicTable<ConfigRecord>
        ref={tableRef}
        columns={columns}
        fetchData={configApi.page}
        rowKey="id"
        searchSchemas={searchSchemas}
        scroll={{ x: 1300 }}
        bulkActions={(
          <Popconfirm
            title="确定批量删除选中的配置吗？"
            onConfirm={async () => {
              const rows = tableRef.current?.getSelectedRows() ?? [];
              await configApi.batchDelete(rows.map((row) => row.id));
              message.success('批量删除成功');
              tableRef.current?.reload();
            }}
          >
            <Button danger icon={<DeleteOutlined />}>
              删除所选
            </Button>
          </Popconfirm>
        )}
        toolbar={
          <Button
            icon={<PlusOutlined />}
            type="primary"
            onClick={() => {
              setCurrentConfigId(undefined);
              setModalOpen(true);
            }}
          >
            新增
          </Button>
        }
      />

      <ConfigModal
        open={modalOpen}
        configId={currentConfigId}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          setModalOpen(false);
          tableRef.current?.reload();
        }}
      />

      <JsonViewerModal
        open={viewerOpen}
        title="查看配置数据"
        data={viewerData}
        onCancel={() => setViewerOpen(false)}
      />
    </>
  );
}
