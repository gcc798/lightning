import { useEffect, useMemo, useState } from 'react';
import { App, Form } from '@/components/ui';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { UserFormData } from '@/types/system';
import { orgApi } from '@/api/org';
import { userApi } from '@/api/user';
import { BasicForm } from '@/components/common/BasicForm';
import { BasicModal } from '@/components/common/BasicModal';

interface UserModalProps {
  open: boolean;
  userId?: SnowflakeId;
  onCancel: () => void;
  onSuccess: () => void;
}

export function UserModal({
  open,
  userId,
  onCancel,
  onSuccess,
}: UserModalProps) {
  const [form] = Form.useForm<UserFormData>();
  const { message } = App.useApp();
  const isEdit = Boolean(userId);
  const [loading, setLoading] = useState(false);
  const [orgTree, setOrgTree] = useState<Record<string, unknown>[]>([]);

  const schemas = useMemo<FormSchema[]>(
    () => [
      {
        name: 'orgId',
        label: '所属组织',
        component: 'TreeSelect',
        props: {
          treeData: orgTree,
          treeDefaultExpandAll: true,
          fieldNames: { label: 'title', value: 'value', children: 'children' },
          placeholder: '请选择所属组织',
        },
      },
      {
        name: 'userName',
        label: '用户名',
        component: 'Input',
        rules: [
          { required: true, message: '请输入用户名' },
          { min: 3, message: '用户名至少 3 位' },
        ],
        props: {
          disabled: isEdit,
        },
      },
      {
        name: 'nickName',
        label: '昵称',
        component: 'Input',
        rules: [{ required: true, message: '请输入昵称' }],
      },
      {
        name: 'password',
        label: '密码',
        component: 'Password',
        hidden: isEdit,
        rules: isEdit
          ? undefined
          : [
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少 6 位' },
            ],
        props: { autoComplete: 'new-password' },
      },
      {
        name: 'email',
        label: '邮箱',
        component: 'Input',
        rules: [{ type: 'email', message: '邮箱格式不正确' }],
      },
      {
        name: 'phonenumber',
        label: '手机号',
        component: 'Input',
        rules: [{ pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确' }],
      },
      {
        name: 'sex',
        label: '性别',
        component: 'Select',
        props: {
          options: [
            { label: '男', value: 0 },
            { label: '女', value: 1 },
            { label: '未知', value: 2 },
          ],
        },
      },
      {
        name: 'status',
        label: '状态',
        component: 'Select',
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
    [isEdit, orgTree],
  );

  useEffect(() => {
    if (!open) {
      return;
    }

    void (async () => {
      const data = await orgApi.tree();
      const transform = (nodes: typeof data): Record<string, unknown>[] =>
        nodes.map((node) => ({
          title: node.orgName,
          value: node.id,
          key: node.id,
          children: node.children ? transform(node.children) : undefined,
        }));

      setOrgTree(transform(data));
    })();
  }, [open]);

  useEffect(() => {
    if (!open) {
      form.resetFields();
      return;
    }

    if (!userId) {
      form.resetFields();
      form.setFieldsValue({ status: 0 });
      return;
    }

    void (async () => {
      setLoading(true);
      try {
        const detail = await userApi.detail(userId);
        form.setFieldsValue(detail as unknown as UserFormData);
      } finally {
        setLoading(false);
      }
    })();
  }, [form, open, userId]);

  const handleSubmit = async () => {
    const values = await form.validateFields();
    setLoading(true);
    try {
      if (userId) {
        await userApi.update(userId, values);
        message.success('用户更新成功');
      } else {
        await userApi.create(values);
        message.success('用户创建成功');
      }
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title={userId ? '编辑用户' : '新增用户'}
      width={640}
      confirmLoading={loading}
      onCancel={onCancel}
      onOk={() => void handleSubmit()}
    >
      <BasicForm
        form={form}
        schemas={schemas}
        layout="vertical"
        variant="modal"
        showActionButtons={false}
      />
    </BasicModal>
  );
}
