import { useRef, useState } from 'react';
import type { ColumnsType } from 'antd/es/table';
import { App, Button, Popconfirm, Space, Tag } from 'antd';
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { BasicTable, type BasicTableRef } from '@/components/common/BasicTable';
import { ApiPermissionAssignModal } from '@/components/common/ApiPermissionAssignModal';
import { PermissionGate } from '@/components/common/PermissionGate';
import { TableAction } from '@/components/common/TableAction';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { RoleRecord } from '@/types/system';
import { roleApi } from '@/api/role';
import { isZeroStatus } from '@/utils/number';
import { PermissionModal } from './PermissionModal';
import { RoleUsersAssignModal } from './RoleUsersAssignModal';
import { RoleModal } from './RoleModal';

const searchSchemas: FormSchema[] = [
  {
    name: 'roleName',
    label: '角色名称',
    component: 'Input',
    colProps: { span: 12 },
  },
  {
    name: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 12 },
    props: {
      allowClear: true,
      options: [
        { label: '正常', value: 0 },
        { label: '停用', value: 1 },
      ],
    },
  },
];

export default function RolePage() {
  const { message } = App.useApp();
  const tableRef = useRef<BasicTableRef<RoleRecord>>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [permissionModalOpen, setPermissionModalOpen] = useState(false);
  const [apiPermissionModalOpen, setApiPermissionModalOpen] = useState(false);
  const [usersAssignModalOpen, setUsersAssignModalOpen] = useState(false);
  const [currentRoleId, setCurrentRoleId] = useState<SnowflakeId>();
  const [currentRole, setCurrentRole] = useState<RoleRecord>();

  const columns: ColumnsType<RoleRecord> = [
    { title: '角色名称', dataIndex: 'roleName', width: 160 },
    { title: '角色标识', dataIndex: 'roleKey', width: 180 },
    { title: '排序', dataIndex: 'sort', width: 100 },
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
      width: 168,
      fixed: 'right',
      render: (_, record) => (
        <TableAction
          actions={[
            {
              key: 'edit',
              label: '编辑',
              permission: 'role.update',
              onClick: () => {
                setCurrentRoleId(record.id);
                setModalOpen(true);
              },
            },
            {
              key: 'permission',
              label: '菜单权限',
              permission: 'role.update',
              onClick: () => {
                setCurrentRoleId(record.id);
                setPermissionModalOpen(true);
              },
            },
            {
              key: 'apiPermission',
              label: 'API权限（高级）',
              permission: 'api_permission.assign',
              onClick: () => {
                setCurrentRoleId(record.id);
                setApiPermissionModalOpen(true);
              },
            },
            {
              key: 'assignUsers',
              label: '分配用户',
              permission: 'role.assign',
              onClick: () => {
                setCurrentRole(record);
                setCurrentRoleId(record.id);
                setUsersAssignModalOpen(true);
              },
            },
            {
              key: 'delete',
              label: '删除',
              permission: 'role.delete',
              danger: true,
              confirmTitle: '确定删除该角色吗？',
              onClick: async () => {
                await roleApi.delete(record.id);
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
      <BasicTable<RoleRecord>
        ref={tableRef}
        columns={columns}
        fetchData={roleApi.page}
        rowKey="id"
        searchSchemas={searchSchemas}
        scroll={{ x: 'max-content' }}
        toolbar={
          <Space>
            <PermissionGate permission="role.create">
              <Button
                icon={<PlusOutlined />}
                type="primary"
                onClick={() => {
                  setCurrentRoleId(undefined);
                  setCurrentRole(undefined);
                  setModalOpen(true);
                }}
              >
                新增
              </Button>
            </PermissionGate>
            <PermissionGate permission="role.delete">
              <Popconfirm
                title="确定删除选中的角色吗？"
                onConfirm={async () => {
                  const rows = tableRef.current?.getSelectedRows() ?? [];
                  if (!rows.length) {
                    message.warning('请先选择角色');
                    return;
                  }
                  await Promise.all(rows.map((row) => roleApi.delete(row.id)));
                  message.success('批量删除成功');
                  tableRef.current?.reload();
                }}
              >
                <Button danger icon={<DeleteOutlined />}>
                  批量删除
                </Button>
              </Popconfirm>
            </PermissionGate>
          </Space>
        }
      />

      <RoleModal
        open={modalOpen}
        roleId={currentRoleId}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          setModalOpen(false);
          tableRef.current?.reload();
        }}
      />

      <PermissionModal
        open={permissionModalOpen}
        roleId={currentRoleId}
        onCancel={() => setPermissionModalOpen(false)}
        onSuccess={() => {
          setPermissionModalOpen(false);
          tableRef.current?.reload();
        }}
      />

      <ApiPermissionAssignModal
        open={apiPermissionModalOpen}
        targetId={currentRoleId}
        targetType="role"
        title="高级 API 授权"
        onCancel={() => setApiPermissionModalOpen(false)}
        onSuccess={() => {
          setApiPermissionModalOpen(false);
          tableRef.current?.reload();
        }}
      />

      <RoleUsersAssignModal
        open={usersAssignModalOpen}
        role={currentRole}
        onCancel={() => setUsersAssignModalOpen(false)}
        onSuccess={() => {
          setUsersAssignModalOpen(false);
          tableRef.current?.reload();
        }}
      />
    </>
  );
}
