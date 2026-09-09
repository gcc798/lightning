import { useEffect, useMemo, useState } from 'react';
import type { Key } from 'react';
import { App, Input, Radio, Space, Spin, Table, Tree } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { DataNode } from 'antd/es/tree';
import { apiPermissionApi } from '@/api/apiPermission';
import type { SnowflakeId } from '@/types/api';
import type { ApiPermissionRecord } from '@/types/system';
import { BasicModal } from './BasicModal';

interface ApiPermissionAssignModalProps {
  open: boolean;
  targetType: 'role' | 'user';
  targetId?: SnowflakeId;
  title: string;
  onCancel: () => void;
  onSuccess: () => void;
}

interface PermissionTreeNode extends DataNode {
  raw: ApiPermissionRecord;
  children?: PermissionTreeNode[];
}

function flatten(nodes: ApiPermissionRecord[]): ApiPermissionRecord[] {
  return nodes.flatMap((node) => [{ ...node, children: undefined }, ...(node.children ? flatten(node.children) : [])]);
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

function toPermissionTreeNodes(nodes: ApiPermissionRecord[], derivedKeys: Key[]): PermissionTreeNode[] {
  return nodes.map((node) => ({
    key: node.id,
    title: node.name,
    raw: node,
    disableCheckbox: derivedKeys.includes(node.id),
    children: node.children ? toPermissionTreeNodes(node.children, derivedKeys) : undefined,
  }));
}

export function ApiPermissionAssignModal({
  open,
  targetType,
  targetId,
  title,
  onCancel,
  onSuccess,
}: ApiPermissionAssignModalProps) {
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [tree, setTree] = useState<ApiPermissionRecord[]>([]);
  const [checkedKeys, setCheckedKeys] = useState<Key[]>([]);
  const [derivedKeys, setDerivedKeys] = useState<Key[]>([]);
  const [keyword, setKeyword] = useState('');
  const [viewMode, setViewMode] = useState<'tree' | 'flat'>('tree');

  useEffect(() => {
    if (!open || !targetId) {
      setCheckedKeys([]);
      setDerivedKeys([]);
      return;
    }

    void (async () => {
      setLoading(true);
      try {
        const [permissionTree, assignedIds] = await Promise.all([
          apiPermissionApi.tree(),
          targetType === 'role'
            ? apiPermissionApi.getRolePermissions(targetId)
            : apiPermissionApi.getUserPermissions(targetId),
        ]);
        setTree(permissionTree);
        if (targetType === 'role') {
          const grants = assignedIds as Awaited<ReturnType<typeof apiPermissionApi.getRolePermissions>>;
          setDerivedKeys(grants.derivedIds);
          setCheckedKeys([...grants.manualIds, ...grants.derivedIds]);
        } else {
          setDerivedKeys([]);
          setCheckedKeys(assignedIds as SnowflakeId[]);
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [open, targetId, targetType]);

  const flatData = useMemo(() => flatten(tree), [tree]);
  const dataSource = useMemo(() => {
    const source = viewMode === 'tree' ? tree : flatData;
    return filterTree(source, keyword);
  }, [flatData, keyword, tree, viewMode]);
  const treeData = useMemo(() => toPermissionTreeNodes(dataSource, derivedKeys), [dataSource, derivedKeys]);

  const columns: ColumnsType<ApiPermissionRecord> = [
    { title: '名称', dataIndex: 'name', width: 200 },
    { title: '权限标识', dataIndex: 'code', width: 260, render: (value) => <span className="table-code">{value}</span> },
    { title: '动作', dataIndex: 'action', width: 90 },
    { title: '路径', dataIndex: 'path', render: (value) => <span className="table-code">{value || '-'}</span> },
  ];

  const handleSubmit = async () => {
    if (!targetId) {
      return;
    }
    setLoading(true);
    try {
      const permissionIds = checkedKeys.filter(
        (key): key is SnowflakeId =>
          (typeof key === 'string' || typeof key === 'number') && !derivedKeys.includes(key),
      );
      if (targetType === 'role') {
        await apiPermissionApi.assignRolePermissions(targetId, permissionIds);
      } else {
        await apiPermissionApi.assignUserPermissions(targetId, permissionIds);
      }
      message.success('API 权限授权成功');
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title={title}
      width={860}
      confirmLoading={loading}
      onCancel={onCancel}
      onOk={() => void handleSubmit()}
    >
      <Spin spinning={loading}>
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          <Space wrap>
            <Input.Search
              allowClear
              placeholder="本地搜索名称、标识、模块、路径"
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
          </Space>
          {viewMode === 'tree' ? (
            <Tree
              checkable
              checkedKeys={checkedKeys}
              defaultExpandAll
              selectable={false}
              treeData={treeData}
              titleRender={(node: PermissionTreeNode) => (
                <span>
                  {node.raw.name} <span className="table-code">{node.raw.code}</span>
                </span>
              )}
              onCheck={(keys) => {
                const selected = Array.isArray(keys) ? keys : keys.checked;
                setCheckedKeys([...new Set([...selected, ...derivedKeys])]);
              }}
            />
          ) : (
            <Table<ApiPermissionRecord>
              columns={columns}
              dataSource={dataSource}
              pagination={false}
              rowKey="id"
              rowSelection={{
                selectedRowKeys: checkedKeys,
                getCheckboxProps: (record) => ({ disabled: derivedKeys.includes(record.id) }),
                onChange: (keys) => setCheckedKeys([...new Set([...keys, ...derivedKeys])]),
              }}
              scroll={{ x: 900, y: 420 }}
            />
          )}
        </Space>
      </Spin>
    </BasicModal>
  );
}
