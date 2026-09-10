import { useEffect, useMemo, useState } from 'react';
import { App, Input, Table, Tag } from '@/components/ui';
import type { ColumnsType } from '@/components/ui';
import { SearchOutlined } from '@/utils/icons';
import { BasicModal } from '@/components/common/BasicModal';
import { roleApi } from '@/api/role';
import type { SnowflakeId } from '@/types/api';
import type { RoleRecord, UserRecord } from '@/types/system';
import { isZeroStatus } from '@/utils/number';

interface UserRoleAssignModalProps {
  open: boolean;
  user?: UserRecord;
  onCancel: () => void;
  onSuccess: () => void;
}

function toKey(id?: SnowflakeId) {
  return String(id ?? '');
}

function diffKeys(nextKeys: string[], originKeys: string[]) {
  const nextSet = new Set(nextKeys);
  const originSet = new Set(originKeys);

  return {
    added: nextKeys.filter((key) => !originSet.has(key)),
    removed: originKeys.filter((key) => !nextSet.has(key)),
  };
}

export function UserRoleAssignModal({
  open,
  user,
  onCancel,
  onSuccess,
}: UserRoleAssignModalProps) {
  const { message } = App.useApp();
  const [roles, setRoles] = useState<RoleRecord[]>([]);
  const [originKeys, setOriginKeys] = useState<string[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [keyword, setKeyword] = useState('');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open || !user?.id) {
      return;
    }

    setKeyword('');
    setLoading(true);
    Promise.all([
      roleApi.page({ pageNum: 1, pageSize: 1000, status: 0 }),
      roleApi.getUserRoles(user.id),
    ])
      .then(([rolePage, userRoles]) => {
        const nextOriginKeys = userRoles.map((role) => toKey(role.id));
        setRoles(rolePage.records ?? []);
        setOriginKeys(nextOriginKeys);
        setSelectedKeys(nextOriginKeys);
      })
      .finally(() => setLoading(false));
  }, [open, user?.id]);

  const filteredRoles = useMemo(() => {
    const value = keyword.trim().toLowerCase();
    if (!value) {
      return roles;
    }

    return roles.filter((role) =>
      [role.roleName, role.roleKey]
        .filter(Boolean)
        .some((item) => item.toLowerCase().includes(value)),
    );
  }, [keyword, roles]);

  const roleIdByKey = useMemo(() => {
    const map = new Map<string, SnowflakeId>();
    roles.forEach((role) => map.set(toKey(role.id), role.id));
    return map;
  }, [roles]);

  const changes = useMemo(
    () => diffKeys(selectedKeys, originKeys),
    [originKeys, selectedKeys],
  );

  const columns: ColumnsType<RoleRecord> = [
    {
      title: '角色名称',
      dataIndex: 'roleName',
      width: 180,
    },
    {
      title: '角色标识',
      dataIndex: 'roleKey',
      width: 180,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (value) => (
        <Tag color={isZeroStatus(value) ? 'success' : 'error'}>
          {isZeroStatus(value) ? '正常' : '停用'}
        </Tag>
      ),
    },
  ];

  const handleSave = async () => {
    if (!user?.id) {
      return;
    }

    setSaving(true);
    try {
      const addedRoleIds = changes.added
        .map((key) => roleIdByKey.get(key))
        .filter((id): id is SnowflakeId => Boolean(id));
      const removedRoleIds = changes.removed
        .map((key) => roleIdByKey.get(key))
        .filter((id): id is SnowflakeId => Boolean(id));

      await Promise.all([
        ...addedRoleIds.map((roleId) => roleApi.assignRole(user.id, roleId)),
        ...removedRoleIds.map((roleId) => roleApi.removeRole(user.id, roleId)),
      ]);

      message.success('角色分配已保存');
      onSuccess();
    } finally {
      setSaving(false);
    }
  };

  return (
    <BasicModal
      confirmLoading={saving}
      open={open}
      title="分配角色"
      width={760}
      onCancel={onCancel}
      onOk={() => void handleSave()}
    >
      <div className="assign-modal user-role-assign-modal">
        <Input
          allowClear
          className="assign-search"
          placeholder="搜索角色名称或标识"
          prefix={<SearchOutlined />}
          value={keyword}
          onChange={(event) => setKeyword(event.target.value)}
        />

        <Table<RoleRecord>
          className="assign-fixed-table user-role-table"
          columns={columns}
          dataSource={filteredRoles}
          loading={loading}
          pagination={false}
          rowKey={(record) => toKey(record.id)}
          rowSelection={{
            selectedRowKeys: selectedKeys,
            onChange: (keys) => setSelectedKeys(keys.map(String)),
          }}
          scroll={{ y: 360 }}
          size="small"
        />

        <div className="assign-modal-summary">
          已选 <strong>{selectedKeys.length}</strong> 个角色
          {changes.added.length || changes.removed.length ? (
            <span>
              ，新增 {changes.added.length} 个，移除 {changes.removed.length} 个
            </span>
          ) : (
            <span>，无变更</span>
          )}
        </div>
      </div>
    </BasicModal>
  );
}
