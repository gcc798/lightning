import { useEffect, useMemo, useState } from 'react';
import { App, Form } from '@/components/ui';
import { apiPermissionApi } from '@/api/apiPermission';
import { BasicForm } from '@/components/common/BasicForm';
import { BasicModal } from '@/components/common/BasicModal';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { ApiPermissionRecord } from '@/types/system';

interface ApiPermissionModalProps {
  open: boolean;
  permissionId?: SnowflakeId;
  parentId?: SnowflakeId;
  onCancel: () => void;
  onSuccess: () => void;
}

function toTreeSelect(nodes: ApiPermissionRecord[]): Record<string, unknown>[] {
  return nodes.map((node) => ({
    title: `${node.name} (${node.code})`,
    value: node.id,
    key: node.id,
    children: node.children ? toTreeSelect(node.children) : undefined,
  }));
}

function flatten(nodes: ApiPermissionRecord[]): ApiPermissionRecord[] {
  return nodes.flatMap((node) => [node, ...(node.children ? flatten(node.children) : [])]);
}

export function ApiPermissionModal({
  open,
  permissionId,
  parentId,
  onCancel,
  onSuccess,
}: ApiPermissionModalProps) {
  const [form] = Form.useForm<Partial<ApiPermissionRecord>>();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [tree, setTree] = useState<ApiPermissionRecord[]>([]);

  const schemas = useMemo<FormSchema[]>(
    () => [
      {
        name: 'parentId',
        label: '上级权限',
        component: 'TreeSelect',
        initialValue: parentId ?? 0,
        props: {
          treeData: [{ title: '顶级模块', value: 0, key: 0 }, ...toTreeSelect(tree)],
          treeDefaultExpandAll: true,
          showSearch: true,
          treeNodeFilterProp: 'title',
          placeholder: '请选择上级权限',
        },
      },
      {
        name: 'nodeType',
        label: '节点类型',
        component: 'RadioGroup',
        initialValue: 0,
        props: {
          optionType: 'button',
          buttonStyle: 'solid',
          options: [
            { label: '模块', value: 0 },
            { label: '分组', value: 1 },
            { label: '权限', value: 2 },
          ],
        },
      },
      {
        name: 'name',
        label: '名称',
        component: 'Input',
        rules: [{ required: true, message: '请输入名称' }],
      },
      {
        name: 'module',
        label: '模块',
        component: 'Input',
        rules: [{ required: true, message: '请输入模块' }],
        helpMessage: '例如 user / role / autodoor',
      },
      {
        name: 'code',
        label: '权限标识',
        component: 'Input',
        rules: [{ required: true, message: '请输入权限标识' }],
        helpMessage: '例如 user.* / user.create / autodoor.device.read',
      },
      {
        name: 'action',
        label: '动作',
        component: 'Select',
        initialValue: '*',
        props: {
          options: [
            { label: '*', value: '*' },
            { label: 'read', value: 'read' },
            { label: 'write', value: 'write' },
            { label: 'control', value: 'control' },
          ],
        },
      },
      {
        name: 'method',
        label: '请求方法',
        component: 'Select',
        props: {
          allowClear: true,
          options: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map((item) => ({ label: item, value: item })),
        },
      },
      {
        name: 'path',
        label: '接口路径',
        component: 'Input',
      },
      {
        name: 'sort',
        label: '排序',
        component: 'InputNumber',
        initialValue: 0,
        props: { min: 0 },
      },
      {
        name: 'status',
        label: '状态',
        component: 'RadioGroup',
        initialValue: 0,
        props: {
          options: [
            { label: '正常', value: 0 },
            { label: '停用', value: 1 },
          ],
        },
      },
      {
        name: 'remark',
        label: '备注',
        component: 'TextArea',
        props: { rows: 3 },
      },
    ],
    [parentId, tree],
  );

  useEffect(() => {
    if (!open) {
      form.resetFields();
      return;
    }

    void (async () => {
      const data = await apiPermissionApi.tree();
      setTree(data);
      if (!permissionId) {
        form.resetFields();
        form.setFieldsValue({ parentId: parentId ?? 0, nodeType: 0, action: '*', sort: 0, status: 0 });
        return;
      }
      const current = flatten(data).find((item) => item.id === permissionId);
      if (current) {
        form.setFieldsValue(current);
      }
    })();
  }, [form, open, parentId, permissionId]);

  const handleSubmit = async () => {
    const values = await form.validateFields();
    setLoading(true);
    try {
      if (permissionId) {
        await apiPermissionApi.update(permissionId, values);
        message.success('API 权限更新成功');
      } else {
        await apiPermissionApi.create(values);
        message.success('API 权限创建成功');
      }
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title={permissionId ? '编辑 API 权限' : '新增 API 权限'}
      width={760}
      confirmLoading={loading}
      onCancel={onCancel}
      onOk={() => void handleSubmit()}
    >
      <BasicForm form={form} layout="vertical" schemas={schemas} showActionButtons={false} variant="modal" />
    </BasicModal>
  );
}
