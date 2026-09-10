import { useEffect, useMemo, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Card, Input, Radio, Space, Table, Tag } from '@/components/ui';
import { DeleteOutlined, DownOutlined, PlusOutlined, UpOutlined } from '@/utils/icons';
import { apiPermissionApi } from '@/api/apiPermission';
import { PermissionGate } from '@/components/common/PermissionGate';
import { TableAction } from '@/components/common/TableAction';
import type { SnowflakeId } from '@/types/api';
import type { ApiPermissionRecord } from '@/types/system';
import { isNumericValue, isZeroStatus, toNumberValue } from '@/utils/number';
import { ApiPermissionModal } from './ApiPermissionModal';

function collectKeys(nodes: ApiPermissionRecord[]): SnowflakeId[] {
  return nodes.flatMap((node) => [node.id, ...(node.children ? collectKeys(node.children) : [])]);
}

function filterTree(nodes: ApiPermissionRecord[], keyword: string): ApiPermissionRecord[] {
  const value = keyword.trim().toLowerCase();
  if (!value) {
    return nodes;
  }
  return nodes
    .map((node) => {
      const children = node.children ? filterTree(node.children, value) : [];
      const matched = [node.name, node.code, node.module, node.path].some((item) =>
        (item ?? '').toLowerCase().includes(value),
      );
      return matched || children.length ? { ...node, children } : null;
    })
    .filter(Boolean) as ApiPermissionRecord[];
}

export default function ApiPermissionPage() {
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [tree, setTree] = useState<ApiPermissionRecord[]>([]);
  const [expandedRowKeys, setExpandedRowKeys] = useState<React.Key[]>([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [currentId, setCurrentId] = useState<SnowflakeId>();
  const [parentId, setParentId] = useState<SnowflakeId>();
  const [keyword, setKeyword] = useState('');
  const [viewMode, setViewMode] = useState<'tree' | 'flat'>('tree');

  const loadTree = async () => {
    setLoading(true);
    try {
      const data = await apiPermissionApi.tree();
      setTree(data);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadTree();
  }, []);

  const flatData = useMemo(() => {
    const flatten = (nodes: ApiPermissionRecord[]): ApiPermissionRecord[] =>
      nodes.flatMap((node) => [{ ...node, children: undefined }, ...(node.children ? flatten(node.children) : [])]);
    return flatten(tree);
  }, [tree]);

  const dataSource = useMemo(() => {
    const source = viewMode === 'tree' ? tree : flatData;
    return filterTree(source, keyword);
  }, [flatData, keyword, tree, viewMode]);

  const columns: ColumnsType<ApiPermissionRecord> = [
    { title: '名称', dataIndex: 'name', width: 220 },
    { title: '权限标识', dataIndex: 'code', width: 260, render: (value) => <span className="table-code">{value}</span> },
    { title: '模块', dataIndex: 'module', width: 140 },
    {
      title: '类型',
      dataIndex: 'nodeType',
      width: 100,
      render: (value) => {
        const nodeType = toNumberValue(value);
        const label = nodeType === 0 ? '模块' : nodeType === 1 ? '分组' : '权限';
        const color = nodeType === 0 ? 'blue' : nodeType === 1 ? 'cyan' : 'green';
        return <Tag color={color}>{label}</Tag>;
      },
    },
    { title: '动作', dataIndex: 'action', width: 100, render: (value) => <Tag>{value}</Tag> },
    { title: '方法', dataIndex: 'method', width: 100, render: (value) => value || '-' },
    { title: '路径', dataIndex: 'path', width: 280, render: (value) => <span className="table-code">{value || '-'}</span> },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (value) => <Tag color={isZeroStatus(value) ? 'success' : 'error'}>{isZeroStatus(value) ? '正常' : '停用'}</Tag>,
    },
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
              permission: 'api_permission.update',
              onClick: () => {
                setCurrentId(record.id);
                setParentId(undefined);
                setModalOpen(true);
              },
            },
            {
              key: 'addChild',
              label: '新增子级',
              permission: 'api_permission.create',
              hidden: isNumericValue(record.nodeType, 2),
              onClick: () => {
                setCurrentId(undefined);
                setParentId(record.id);
                setModalOpen(true);
              },
            },
            {
              key: 'delete',
              label: '删除',
              permission: 'api_permission.delete',
              danger: true,
              confirmTitle: '确定删除该 API 权限吗？',
              onClick: async () => {
                await apiPermissionApi.delete(record.id);
                message.success('删除成功');
                await loadTree();
              },
            },
          ]}
        />
      ),
    },
  ];

  return (
    <>
      <Card className="page-card" variant="borderless">
        <div className="page-toolbar">
          <Space wrap>
            <PermissionGate permission="api_permission.create">
              <Button
                icon={<PlusOutlined />}
                type="primary"
                onClick={() => {
                  setCurrentId(undefined);
                  setParentId(undefined);
                  setModalOpen(true);
                }}
              >
                新增
              </Button>
            </PermissionGate>
            <Input.Search
              allowClear
              placeholder="搜索名称、标识、模块、路径"
              style={{ width: 320 }}
              onChange={(event) => setKeyword(event.target.value)}
            />
            <Radio.Group
              optionType="button"
              value={viewMode}
              options={[
                { label: '树形', value: 'tree' },
                { label: '平铺', value: 'flat' },
              ]}
              onChange={(event) => setViewMode(event.target.value)}
            />
            {viewMode === 'tree' ? (
              <>
                <Button icon={<DownOutlined />} onClick={() => setExpandedRowKeys(collectKeys(tree))}>
                  展开全部
                </Button>
                <Button icon={<UpOutlined />} onClick={() => setExpandedRowKeys([])}>
                  折叠全部
                </Button>
              </>
            ) : null}
          </Space>
        </div>
        <Table<ApiPermissionRecord>
          columns={columns}
          dataSource={dataSource}
          expandable={viewMode === 'tree' ? { expandedRowKeys, onExpandedRowsChange: (keys) => setExpandedRowKeys([...keys]) } : undefined}
          loading={loading}
          pagination={false}
          rowKey="id"
          scroll={{ x: 1500 }}
        />
      </Card>
      <ApiPermissionModal
        open={modalOpen}
        parentId={parentId}
        permissionId={currentId}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          setModalOpen(false);
          void loadTree();
        }}
      />
    </>
  );
}
