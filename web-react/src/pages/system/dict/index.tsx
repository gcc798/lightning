import { useRef, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Popconfirm, Tag } from '@/components/ui';
import { DeleteOutlined, PlusOutlined } from '@/utils/icons';
import { BasicTable, type BasicTableRef } from '@/components/common/BasicTable';
import { PermissionGate } from '@/components/common/PermissionGate';
import { TableAction } from '@/components/common/TableAction';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { DictRecord } from '@/types/system';
import { dictApi } from '@/api/dict';
import { isZeroStatus } from '@/utils/number';
import { DictModal } from './DictModal';
import { DictSubItemsModal } from './DictSubItemsModal';

const searchSchemas: FormSchema[] = [
  {
    name: 'dictType',
    label: '字典类型',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    name: 'dictLabel',
    label: '字典标签',
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
        { label: '正常', value: 0 },
        { label: '停用', value: 1 },
      ],
    },
  },
];

export default function DictPage() {
  const { message } = App.useApp();
  const tableRef = useRef<BasicTableRef<DictRecord>>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [subItemsOpen, setSubItemsOpen] = useState(false);
  const [currentDictId, setCurrentDictId] = useState<SnowflakeId>();
  const [currentParentId, setCurrentParentId] = useState<SnowflakeId>();

  const columns: ColumnsType<DictRecord> = [
    { title: '字典类型', dataIndex: 'dictType', width: 160 },
    { title: '字典标签', dataIndex: 'dictLabel', width: 160 },
    { title: '字典键值', dataIndex: 'dictValue', width: 160 },
    { title: '排序', dataIndex: 'sort', width: 100 },
    {
      title: '默认',
      dataIndex: 'isDefault',
      width: 100,
      render: (value) => <Tag color={value ? 'success' : 'default'}>{value ? '是' : '否'}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (value) => (
        <Tag color={isZeroStatus(value) ? 'success' : 'error'}>
          {isZeroStatus(value) ? '正常' : '停用'}
        </Tag>
      ),
    },
    { title: '创建时间', dataIndex: 'createdTime', width: 180 },
    {
      title: '操作',
      key: 'action',
      width: 220,
      fixed: 'right',
      render: (_, record) => (
        <TableAction
          actions={[
            {
              key: 'edit',
              label: '编辑',
              permission: 'dict.update',
              onClick: () => {
                setCurrentDictId(record.id);
                setModalOpen(true);
              },
            },
            {
              key: 'subItems',
              label: '编辑子项',
              permission: 'dict.update',
              onClick: () => {
                setCurrentParentId(record.id);
                setSubItemsOpen(true);
              },
            },
            {
              key: 'delete',
              label: '删除',
              permission: 'dict.delete',
              danger: true,
              confirmTitle: '确定删除该字典吗？',
              onClick: async () => {
                await dictApi.delete(record.id);
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
      <BasicTable<DictRecord>
        ref={tableRef}
        columns={columns}
        fetchData={dictApi.page}
        rowKey="id"
        searchSchemas={searchSchemas}
        scroll={{ x: 1300 }}
        bulkActions={(
          <PermissionGate permission="dict.delete">
            <Popconfirm
              title="确定批量删除选中的字典吗？"
              onConfirm={async () => {
                const rows = tableRef.current?.getSelectedRows() ?? [];
                await dictApi.batchDelete(rows.map((row) => row.id));
                message.success('批量删除成功');
                tableRef.current?.reload();
              }}
            >
              <Button danger icon={<DeleteOutlined />}>
                删除所选
              </Button>
            </Popconfirm>
          </PermissionGate>
        )}
        toolbar={
          <PermissionGate permission="dict.create">
            <Button
              icon={<PlusOutlined />}
              type="primary"
              onClick={() => {
                setCurrentDictId(undefined);
                setModalOpen(true);
              }}
            >
              新增
            </Button>
          </PermissionGate>
        }
      />

      <DictModal
        open={modalOpen}
        dictId={currentDictId}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          setModalOpen(false);
          tableRef.current?.reload();
        }}
      />

      <DictSubItemsModal
        open={subItemsOpen}
        parentId={currentParentId}
        onCancel={() => setSubItemsOpen(false)}
      />
    </>
  );
}
