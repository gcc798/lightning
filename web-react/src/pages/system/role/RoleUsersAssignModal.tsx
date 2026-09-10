import { useCallback, useEffect, useMemo, useState } from 'react';
import { App, Button, Input, Space, Table, Typography } from '@/components/ui';
import type { ColumnsType } from '@/components/ui';
import {
  ArrowLeftOutlined,
  ArrowRightOutlined,
  SearchOutlined,
} from '@/utils/icons';
import { BasicModal } from '@/components/common/BasicModal';
import { roleApi } from '@/api/role';
import { userApi } from '@/api/user';
import type { SnowflakeId } from '@/types/api';
import type { RoleRecord, UserRecord } from '@/types/system';

interface RoleUsersAssignModalProps {
  open: boolean;
  role?: RoleRecord;
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

function toUserQuery(keyword: string) {
  const value = keyword.trim();
  if (!value) {
    return {};
  }

  return /^\d+$/.test(value) ? { phonenumber: value } : { username: value };
}

export function RoleUsersAssignModal({
  open,
  role,
  onCancel,
  onSuccess,
}: RoleUsersAssignModalProps) {
  const { message } = App.useApp();
  const [allUsers, setAllUsers] = useState<UserRecord[]>([]);
  const [originKeys, setOriginKeys] = useState<string[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [leftCheckedKeys, setLeftCheckedKeys] = useState<string[]>([]);
  const [rightCheckedKeys, setRightCheckedKeys] = useState<string[]>([]);
  const [knownUsers, setKnownUsers] = useState<Record<string, UserRecord>>({});
  const [keyword, setKeyword] = useState('');
  const [pageNum, setPageNum] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  const mergeUsers = useCallback((users: UserRecord[]) => {
    setKnownUsers((prev) => {
      const next = { ...prev };
      users.forEach((user) => {
        next[toKey(user.id)] = user;
      });
      return next;
    });
  }, []);

  const loadUsers = useCallback(async () => {
    setLoading(true);
    try {
      const page = await userApi.page({
        pageNum,
        pageSize,
        ...toUserQuery(keyword),
      });
      const records = page.records ?? [];
      setAllUsers(records);
      setTotal(page.total ?? 0);
      mergeUsers(records);
    } finally {
      setLoading(false);
    }
  }, [keyword, mergeUsers, pageNum, pageSize]);

  useEffect(() => {
    if (!open || !role?.id) {
      return;
    }

    setKeyword('');
    setPageNum(1);
    setPageSize(10);
    setLeftCheckedKeys([]);
    setRightCheckedKeys([]);
    setKnownUsers({});
    roleApi.getRoleUsers(role.id)
      .then((assignedUsers) => {
        const nextOriginKeys = assignedUsers.map((user) => toKey(user.id));
        setOriginKeys(nextOriginKeys);
        setSelectedKeys(nextOriginKeys);
        mergeUsers(assignedUsers);
      })
      .catch(() => {
        setOriginKeys([]);
        setSelectedKeys([]);
      });
  }, [mergeUsers, open, role?.id]);

  useEffect(() => {
    if (!open || !role?.id) {
      return;
    }
    void loadUsers();
  }, [loadUsers, open, role?.id]);

  const selectedKeySet = useMemo(() => new Set(selectedKeys), [selectedKeys]);
  const leftUsers = useMemo(
    () => allUsers.filter((user) => !selectedKeySet.has(toKey(user.id))),
    [allUsers, selectedKeySet],
  );
  const selectedUsers = useMemo(
    () => selectedKeys.map((key) => knownUsers[key]).filter(Boolean),
    [knownUsers, selectedKeys],
  );
  const changes = useMemo(
    () => diffKeys(selectedKeys, originKeys),
    [originKeys, selectedKeys],
  );

  const userColumns: ColumnsType<UserRecord> = [
    {
      title: '用户名',
      dataIndex: 'userName',
      width: 130,
    },
    {
      title: '昵称',
      dataIndex: 'nickName',
      width: 130,
    },
    {
      title: '手机号',
      dataIndex: 'phonenumber',
      render: (value) => value || '-',
    },
  ];

  const assignedColumns: ColumnsType<UserRecord> = [
    {
      title: '用户名',
      dataIndex: 'userName',
      width: 140,
    },
    {
      title: '昵称',
      dataIndex: 'nickName',
      width: 140,
    },
    {
      title: '手机号',
      dataIndex: 'phonenumber',
      render: (value) => value || '-',
    },
  ];

  const handleAdd = () => {
    setSelectedKeys((prev) => Array.from(new Set([...prev, ...leftCheckedKeys])));
    setLeftCheckedKeys([]);
  };

  const handleRemove = () => {
    const removeSet = new Set(rightCheckedKeys);
    setSelectedKeys((prev) => prev.filter((key) => !removeSet.has(key)));
    setRightCheckedKeys([]);
  };

  const handleSave = async () => {
    if (!role?.id) {
      return;
    }

    setSaving(true);
    try {
      const addedUserIds = changes.added
        .map((key) => knownUsers[key]?.id)
        .filter((id): id is SnowflakeId => Boolean(id));
      const removedUserIds = changes.removed
        .map((key) => knownUsers[key]?.id)
        .filter((id): id is SnowflakeId => Boolean(id));

      await Promise.all([
        addedUserIds.length ? roleApi.assignUsers(role.id, addedUserIds) : Promise.resolve(),
        removedUserIds.length ? roleApi.removeUsers(role.id, removedUserIds) : Promise.resolve(),
      ]);

      message.success('角色用户已保存');
      onSuccess();
    } finally {
      setSaving(false);
    }
  };

  return (
    <BasicModal
      confirmLoading={saving}
      open={open}
      title="分配用户"
      width={960}
      onCancel={onCancel}
      onOk={() => void handleSave()}
    >
      <div className="assign-modal role-users-assign-modal">
        <div className="role-users-transfer">
          <div className="role-users-panel">
            <div className="role-users-panel-header">
              <span>所有用户</span>
              <Input
                allowClear
                className="role-users-search"
                placeholder="搜索用户名或手机号"
                prefix={<SearchOutlined />}
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value);
                  setPageNum(1);
                }}
              />
            </div>
            <Table<UserRecord>
              className="assign-fixed-table role-users-table role-users-table-left"
              columns={userColumns}
              dataSource={leftUsers}
              loading={loading}
              pagination={{
                current: pageNum,
                pageSize,
                total,
                showSizeChanger: true,
                showTotal: (currentTotal) => `共 ${currentTotal} 条`,
                onChange: (nextPage, nextSize) => {
                  setPageNum(nextPage);
                  setPageSize(nextSize);
                },
              }}
              rowKey={(record) => toKey(record.id)}
              rowSelection={{
                selectedRowKeys: leftCheckedKeys,
                onChange: (keys) => setLeftCheckedKeys(keys.map(String)),
              }}
              scroll={{ y: 340 }}
              size="small"
            />
          </div>

          <Space className="role-users-transfer-actions" direction="vertical" size={10}>
            <Button
              disabled={!leftCheckedKeys.length}
              icon={<ArrowRightOutlined />}
              type="primary"
              onClick={handleAdd}
            >
              加入
            </Button>
            <Button
              disabled={!rightCheckedKeys.length}
              icon={<ArrowLeftOutlined />}
              onClick={handleRemove}
            >
              移除
            </Button>
          </Space>

          <div className="role-users-panel">
            <div className="role-users-panel-header">
              <span>已分配用户</span>
              <Typography.Text type="secondary">
                新增 {changes.added.length} / 移除 {changes.removed.length}
              </Typography.Text>
            </div>
            <Table<UserRecord>
              className="assign-fixed-table role-users-table role-users-table-right"
              columns={assignedColumns}
              dataSource={selectedUsers}
              pagination={false}
              rowKey={(record) => toKey(record.id)}
              rowSelection={{
                selectedRowKeys: rightCheckedKeys,
                onChange: (keys) => setRightCheckedKeys(keys.map(String)),
              }}
              scroll={{ y: 410 }}
              size="small"
            />
          </div>
        </div>
      </div>
    </BasicModal>
  );
}
