import { useRef, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Popconfirm, Tag } from '@/components/ui';
import { DeleteOutlined, EditOutlined, KeyOutlined, PlusOutlined } from '@/utils/icons';
import type { BasicTableRef } from '@/components/common/BasicTable';
import { ApiPermissionAssignModal } from '@/components/common/ApiPermissionAssignModal';
import { BasicTable } from '@/components/common/BasicTable';
import { PermissionGate } from '@/components/common/PermissionGate';
import { TableAction } from '@/components/common/TableAction';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { UserRecord } from '@/types/system';
import { userApi } from '@/api/user';
import { isZeroStatus } from '@/utils/number';
import { UserRoleAssignModal } from './UserRoleAssignModal';
import { UserModal } from './UserModal';

const searchSchemas: FormSchema[] = [
  {
    name: 'username',
    label: '用户名',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    name: 'phonenumber',
    label: '手机号',
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

export default function UserPage() {
  const { message } = App.useApp();
  const tableRef = useRef<BasicTableRef<UserRecord>>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [apiPermissionModalOpen, setApiPermissionModalOpen] = useState(false);
  const [roleAssignModalOpen, setRoleAssignModalOpen] = useState(false);
  const [currentUserId, setCurrentUserId] = useState<SnowflakeId>();
  const [currentUser, setCurrentUser] = useState<UserRecord>();

  const columns: ColumnsType<UserRecord> = [
    { title: '用户名', dataIndex: 'userName', width: 140 },
    { title: '昵称', dataIndex: 'nickName', width: 140 },
    { title: '邮箱', dataIndex: 'email', width: 220 },
    { title: '手机号', dataIndex: 'phonenumber', width: 160 },
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
              permission: 'user.update',
              onClick: () => {
                setCurrentUserId(record.id);
                setModalOpen(true);
              },
            },
            {
              key: 'apiPermission',
              label: 'API权限',
              permission: 'api_permission.assign',
              onClick: () => {
                setCurrentUserId(record.id);
                setApiPermissionModalOpen(true);
              },
            },
            {
              key: 'assignRole',
              label: '分配角色',
              permission: 'role.assign',
              onClick: () => {
                setCurrentUser(record);
                setCurrentUserId(record.id);
                setRoleAssignModalOpen(true);
              },
            },
            {
              key: 'resetPassword',
              label: '重置密码',
              permission: 'user.update',
              onClick: async () => {
                await userApi.resetPassword(record.id, '123456');
                message.success('密码已重置为 123456');
              },
            },
            {
              key: 'delete',
              label: '删除',
              permission: 'user.delete',
              danger: true,
              confirmTitle: '确定删除该用户吗？',
              onClick: async () => {
                await userApi.delete(record.id);
                message.success('删除成功');
                tableRef.current?.reload();
              },
            },
          ]}
        />
      ),
    },
  ];

  const handleBatchDelete = async () => {
    const rows = tableRef.current?.getSelectedRows() ?? [];
    if (!rows.length) {
      message.warning('请先选择要删除的用户');
      return;
    }

    await userApi.batchDelete(rows.map((item) => item.id));
    message.success('批量删除成功');
    tableRef.current?.reload();
  };

  return (
    <>
      <BasicTable<UserRecord>
        ref={tableRef}
        columns={columns}
        fetchData={userApi.page}
        rowKey="id"
        searchSchemas={searchSchemas}
        scroll={{ x: 1200 }}
        bulkActions={(
          <PermissionGate permission="user.delete">
            <Popconfirm title="确定批量删除选中的用户吗？" onConfirm={() => void handleBatchDelete()}>
              <Button danger icon={<DeleteOutlined />}>
                删除所选
              </Button>
            </Popconfirm>
          </PermissionGate>
        )}
        toolbar={
          <PermissionGate permission="user.create">
            <Button
              icon={<PlusOutlined />}
              type="primary"
              onClick={() => {
                setCurrentUserId(undefined);
                setCurrentUser(undefined);
                setModalOpen(true);
              }}
            >
              新增
            </Button>
          </PermissionGate>
        }
      />

      <UserModal
        open={modalOpen}
        userId={currentUserId}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          setModalOpen(false);
          tableRef.current?.reload();
        }}
      />

      <UserRoleAssignModal
        open={roleAssignModalOpen}
        user={currentUser}
        onCancel={() => setRoleAssignModalOpen(false)}
        onSuccess={() => {
          setRoleAssignModalOpen(false);
          tableRef.current?.reload();
        }}
      />

      <ApiPermissionAssignModal
        open={apiPermissionModalOpen}
        targetId={currentUserId}
        targetType="user"
        title="分配用户 API 权限"
        onCancel={() => setApiPermissionModalOpen(false)}
        onSuccess={() => {
          setApiPermissionModalOpen(false);
          tableRef.current?.reload();
        }}
      />
    </>
  );
}
