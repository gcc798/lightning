import { useEffect, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Card, Space, Table, Tag } from '@/components/ui';
import { DeleteOutlined, DownOutlined, EyeOutlined, PlusOutlined, UpOutlined } from '@/utils/icons';
import { PermissionGate } from '@/components/common/PermissionGate';
import { TableAction } from '@/components/common/TableAction';
import type { SnowflakeId } from '@/types/api';
import { getMenuIconNode } from '@/utils/icons';
import type { MenuRecord } from '@/types/menu';
import { isNumericValue, isZeroStatus, toNumberValue } from '@/utils/number';
import { menuApi } from '@/api/menu';
import { MenuModal } from './MenuModal';

function collectKeys(nodes: MenuRecord[]): SnowflakeId[] {
  const keys: SnowflakeId[] = [];
  const walk = (items: MenuRecord[]) => {
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

export default function MenuPage() {
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [menuTree, setMenuTree] = useState<MenuRecord[]>([]);
  const [expandedRowKeys, setExpandedRowKeys] = useState<React.Key[]>([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [currentMenuId, setCurrentMenuId] = useState<SnowflakeId>();
  const [parentId, setParentId] = useState<SnowflakeId>();

  const loadMenuTree = async () => {
    setLoading(true);
    try {
      const data = await menuApi.getMenuTree();
      setMenuTree(data);
      setExpandedRowKeys([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadMenuTree();
  }, []);

  const columns: ColumnsType<MenuRecord> = [
    {
      title: '菜单名称',
      dataIndex: 'menuName',
      width: 220,
    },
    {
      title: '图标',
      dataIndex: 'icon',
      width: 90,
      render: (value) => getMenuIconNode(value),
    },
    {
      title: '类型',
      dataIndex: 'menuType',
      width: 100,
      render: (value) => {
        const menuType = toNumberValue(value);
        const color = menuType === 0 ? 'blue' : menuType === 1 ? 'green' : 'orange';
        const text = menuType === 0 ? '目录' : menuType === 1 ? '菜单' : '按钮';
        return <Tag color={color}>{text}</Tag>;
      },
    },
    {
      title: '路由地址',
      dataIndex: 'path',
      width: 180,
      render: (value) => <span className="table-code">{value || '-'}</span>,
    },
    {
      title: '组件路径',
      dataIndex: 'component',
      width: 260,
      render: (value) => <span className="table-code">{value || '-'}</span>,
    },
    {
      title: '权限标识',
      dataIndex: 'perms',
      width: 180,
      render: (value) => <span className="table-code">{value || '-'}</span>,
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
    {
      title: '操作',
      key: 'action',
      width: 260,
      fixed: 'right',
      render: (_, record) => (
        <TableAction
          actions={[
            {
              key: 'view',
              label: '查看',
              onClick: () => {
                setCurrentMenuId(record.id);
                setParentId(undefined);
                setModalOpen(true);
              },
            },
            {
              key: 'edit',
              label: '编辑',
              permission: 'menu.update',
              onClick: () => {
                setCurrentMenuId(record.id);
                setParentId(undefined);
                setModalOpen(true);
              },
            },
            {
              key: 'addChild',
              label: '新增子菜单',
              permission: 'menu.create',
              hidden: isNumericValue(record.menuType, 2),
              onClick: () => {
                setCurrentMenuId(undefined);
                setParentId(record.id);
                setModalOpen(true);
              },
            },
            {
              key: 'delete',
              label: '删除',
              permission: 'menu.delete',
              danger: true,
              confirmTitle: '确定删除该菜单吗？',
              onClick: async () => {
                await menuApi.delete(record.id);
                message.success('删除成功');
                await loadMenuTree();
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
            <PermissionGate permission="menu.create">
              <Button
                icon={<PlusOutlined />}
                type="primary"
                onClick={() => {
                  setCurrentMenuId(undefined);
                  setParentId(undefined);
                  setModalOpen(true);
                }}
              >
                新增
              </Button>
            </PermissionGate>
            <Button
              icon={<DownOutlined />}
              onClick={() => setExpandedRowKeys(collectKeys(menuTree))}
            >
              展开全部
            </Button>
            <Button icon={<UpOutlined />} onClick={() => setExpandedRowKeys([])}>
              折叠全部
            </Button>
          </Space>
        </div>

        <Table<MenuRecord>
          columns={columns}
          dataSource={menuTree}
          expandable={{
            expandedRowKeys,
            onExpandedRowsChange: (keys) => setExpandedRowKeys([...keys]),
          }}
          loading={loading}
          pagination={false}
          rowKey="id"
          scroll={{ x: 1400 }}
        />
      </Card>

      <MenuModal
        open={modalOpen}
        menuId={currentMenuId}
        parentId={parentId}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          setModalOpen(false);
          void loadMenuTree();
        }}
      />
    </>
  );
}
