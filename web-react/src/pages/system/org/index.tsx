import { useEffect, useMemo, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Card, Form, Popconfirm, Space, Table, Tag } from '@/components/ui';
import { DeleteOutlined, MenuFoldOutlined, MenuUnfoldOutlined, PlusOutlined } from '@/utils/icons';
import { BasicForm } from '@/components/common/BasicForm';
import { PermissionGate } from '@/components/common/PermissionGate';
import { TableAction } from '@/components/common/TableAction';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { OrgRecord } from '@/types/system';
import { orgApi } from '@/api/org';
import { isZeroStatus, toOptionalNumber } from '@/utils/number';
import { OrgModal } from './OrgModal';

function collectKeys(nodes: OrgRecord[]): SnowflakeId[] {
  const keys: SnowflakeId[] = [];
  const walk = (items: OrgRecord[]) => {
    items.forEach((item) => {
      keys.push(item.id);
      if (item.children?.length) {
        walk(item.children);
      }
    });
  };
  walk(nodes);
  return keys;
}

const searchSchemas: FormSchema[] = [
  {
    name: 'orgName',
    label: '组织名称',
    component: 'Input',
  },
  {
    name: 'orgCode',
    label: '组织编码',
    component: 'Input',
  },
  {
    name: 'status',
    label: '状态',
    component: 'Select',
    props: {
      allowClear: true,
      options: [
        { label: '正常', value: 0 },
        { label: '停用', value: 1 },
      ],
    },
  },
];

function filterTree(
  nodes: OrgRecord[],
  keywordName: string,
  keywordCode: string,
  status?: number,
): OrgRecord[] {
  const normalizedName = keywordName.trim().toLowerCase();
  const normalizedCode = keywordCode.trim().toLowerCase();

  return nodes
    .map((node) => {
      const children = node.children
        ? filterTree(node.children, keywordName, keywordCode, status)
        : [];

      const selfMatch =
        (!normalizedName ||
          node.orgName.toLowerCase().includes(normalizedName)) &&
        (!normalizedCode ||
          (node.orgCode ?? '').toLowerCase().includes(normalizedCode)) &&
        (status === undefined || toOptionalNumber(node.status) === status);

      if (selfMatch || children.length) {
        return {
          ...node,
          children,
        };
      }

      return null;
    })
    .filter(Boolean) as OrgRecord[];
}

export default function OrgPage() {
  const { message } = App.useApp();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [fullTree, setFullTree] = useState<OrgRecord[]>([]);
  const [tableData, setTableData] = useState<OrgRecord[]>([]);
  const [expandedRowKeys, setExpandedRowKeys] = useState<React.Key[]>([]);
  const [expandAll, setExpandAll] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [currentOrgId, setCurrentOrgId] = useState<SnowflakeId>();
  const [parentId, setParentId] = useState<SnowflakeId>();

  const loadOrgTree = async () => {
    setLoading(true);
    try {
      const data = await orgApi.tree();
      setFullTree(data);
      const values = form.getFieldsValue();
      const filtered = filterTree(
        data,
        String(values.orgName ?? ''),
        String(values.orgCode ?? ''),
        toOptionalNumber(values.status),
      );
      setTableData(filtered);
      setExpandedRowKeys(expandAll ? collectKeys(filtered) : []);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadOrgTree();
  }, []);

  const columns: ColumnsType<OrgRecord> = useMemo(
    () => [
      { title: '组织名称', dataIndex: 'orgName', width: 240 },
      { title: '组织编码', dataIndex: 'orgCode', width: 180 },
      { title: '显示顺序', dataIndex: 'sort', width: 120 },
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
        width: 240,
        fixed: 'right',
        render: (_, record) => (
          <TableAction
            actions={[
              {
                key: 'edit',
                label: '编辑',
                permission: 'org.update',
                onClick: () => {
                  setCurrentOrgId(record.id);
                  setParentId(undefined);
                  setModalOpen(true);
                },
              },
              {
                key: 'add',
                label: '新增下级',
                permission: 'org.create',
                onClick: () => {
                  setCurrentOrgId(undefined);
                  setParentId(record.id);
                  setModalOpen(true);
                },
              },
              {
                key: 'delete',
                label: '删除',
                permission: 'org.delete',
                danger: true,
                confirmTitle: '确定删除该组织吗？',
                onClick: async () => {
                  await orgApi.delete(record.id);
                  message.success('删除成功');
                  await loadOrgTree();
                },
              },
            ]}
          />
        ),
      },
    ],
    [message],
  );

  const handleSearch = (values: Record<string, unknown> = form.getFieldsValue()) => {
    const status = toOptionalNumber(values.status);
    const filtered = filterTree(
      fullTree,
      typeof values.orgName === 'string' ? values.orgName : '',
      typeof values.orgCode === 'string' ? values.orgCode : '',
      status,
    );
    setTableData(filtered);
    setExpandedRowKeys(expandAll ? collectKeys(filtered) : []);
  };

  const handleReset = () => {
    form.resetFields();
    setTableData(fullTree);
    setExpandedRowKeys(expandAll ? collectKeys(fullTree) : []);
  };

  return (
    <>
      <div className="page-search">
        <BasicForm
          form={form}
          schemas={searchSchemas}
          layout="vertical"
          onReset={handleReset}
          onSubmit={handleSearch}
          resetText="重置"
          submitText="查询"
          variant="search"
        />
      </div>

      <Card className="page-card" variant="borderless">
        <div className="page-toolbar">
          <Space>
            <PermissionGate permission="org.create">
              <Button
                icon={<PlusOutlined />}
                type="primary"
                onClick={() => {
                  setCurrentOrgId(undefined);
                  setParentId(undefined);
                  setModalOpen(true);
                }}
              >
                新增
              </Button>
            </PermissionGate>
            <Button
              icon={expandAll ? <MenuFoldOutlined /> : <MenuUnfoldOutlined />}
              onClick={() => {
                const next = !expandAll;
                setExpandAll(next);
                setExpandedRowKeys(next ? collectKeys(tableData) : []);
              }}
            >
              {expandAll ? '折叠' : '展开'}
            </Button>
          </Space>
        </div>

        <Table<OrgRecord>
          columns={columns}
          dataSource={tableData}
          loading={loading}
          pagination={false}
          rowKey="id"
          expandable={{
            expandedRowKeys,
            onExpandedRowsChange: (keys) => setExpandedRowKeys([...keys]),
          }}
          scroll={{ x: 1200 }}
        />
      </Card>

      <OrgModal
        open={modalOpen}
        orgId={currentOrgId}
        parentId={parentId}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          setModalOpen(false);
          void loadOrgTree();
        }}
      />
    </>
  );
}
