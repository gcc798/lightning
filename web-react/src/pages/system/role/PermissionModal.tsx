import { useEffect, useMemo, useState } from 'react';
import type { Key } from 'react';
import { App, Space, Spin, Tag, Tree } from '@/components/ui';
import type { SnowflakeId } from '@/types/api';
import { BasicModal } from '@/components/common/BasicModal';
import { menuApi } from '@/api/menu';
import { roleApi } from '@/api/role';
import type { MenuRecord } from '@/types/menu';

interface PermissionModalProps {
  open: boolean;
  roleId?: SnowflakeId;
  onCancel: () => void;
  onSuccess: () => void;
}

export function PermissionModal({
  open,
  roleId,
  onCancel,
  onSuccess,
}: PermissionModalProps) {
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [treeData, setTreeData] = useState<MenuRecord[]>([]);
  const [checkedKeys, setCheckedKeys] = useState<Key[]>([]);

  const derivedCount = useMemo(() => {
    const selected = new Set(checkedKeys);
    const permissions = new Set<SnowflakeId>();
    const walk = (nodes: MenuRecord[]) => nodes.forEach((menu) => {
      if (selected.has(menu.id)) {
        menu.apiPermissionIds?.forEach((id) => permissions.add(id));
      }
      if (menu.children) walk(menu.children);
    });
    walk(treeData);
    return permissions.size;
  }, [checkedKeys, treeData]);

  useEffect(() => {
    if (!open || !roleId) {
      setCheckedKeys([]);
      return;
    }

    void (async () => {
      setLoading(true);
      try {
        const [menus, roleMenus] = await Promise.all([
          menuApi.getMenuTree(),
          roleApi.getMenus(roleId),
        ]);
        setTreeData(menus);
        setCheckedKeys(roleMenus);
      } finally {
        setLoading(false);
      }
    })();
  }, [open, roleId]);

  const handleSubmit = async () => {
    if (!roleId) {
      return;
    }

    setLoading(true);
    try {
      await roleApi.assignMenus(
        roleId,
        checkedKeys.filter(
          (key): key is SnowflakeId => typeof key === 'string' || typeof key === 'number',
        ),
      );
      message.success('角色菜单分配成功');
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title="分配菜单权限"
      width={700}
      confirmLoading={loading}
      onCancel={onCancel}
      onOk={() => void handleSubmit()}
    >
      <Spin spinning={loading}>
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          <Tag color="blue">将派生 {derivedCount} 个 API 权限</Tag>
          <Tree
            checkable
            checkedKeys={checkedKeys}
            defaultExpandAll
            fieldNames={{ title: 'menuName', key: 'id', children: 'children' }}
            selectable={false}
            treeData={treeData}
            onCheck={(keys) => setCheckedKeys(Array.isArray(keys) ? keys : keys.checked)}
          />
        </Space>
      </Spin>
    </BasicModal>
  );
}
