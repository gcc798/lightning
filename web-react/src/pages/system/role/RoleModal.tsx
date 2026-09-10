import { useEffect, useMemo, useState } from 'react';
import { App, Form } from '@/components/ui';
import type { FormSchema } from '@/types/form';
import type { SnowflakeId } from '@/types/api';
import type { RoleRecord } from '@/types/system';
import { roleApi } from '@/api/role';
import { BasicForm } from '@/components/common/BasicForm';
import { BasicModal } from '@/components/common/BasicModal';

interface RoleModalProps {
  open: boolean;
  roleId?: SnowflakeId;
  onCancel: () => void;
  onSuccess: () => void;
}

export function RoleModal({
  open,
  roleId,
  onCancel,
  onSuccess,
}: RoleModalProps) {
  const [form] = Form.useForm<Partial<RoleRecord>>();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);

  const schemas = useMemo<FormSchema[]>(
    () => [
      {
        name: 'roleName',
        label: '角色名称',
        component: 'Input',
        rules: [{ required: true, message: '请输入角色名称' }],
      },
      {
        name: 'roleKey',
        label: '角色标识',
        component: 'Input',
        rules: [{ required: true, message: '请输入角色标识' }],
        props: { disabled: Boolean(roleId) },
      },
      {
        name: 'sort',
        label: '显示顺序',
        component: 'InputNumber',
        initialValue: 0,
        props: { min: 0 },
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
    [roleId],
  );

  useEffect(() => {
    if (!open) {
      form.resetFields();
      return;
    }

    if (!roleId) {
      form.resetFields();
      form.setFieldsValue({ sort: 0, status: 0 });
      return;
    }

    void (async () => {
      setLoading(true);
      try {
        const detail = await roleApi.detail(roleId);
        form.setFieldsValue(detail);
      } finally {
        setLoading(false);
      }
    })();
  }, [form, open, roleId]);

  const handleSubmit = async () => {
    const values = await form.validateFields();
    setLoading(true);
    try {
      if (roleId) {
        await roleApi.update(roleId, values);
        message.success('角色更新成功');
      } else {
        await roleApi.create(values);
        message.success('角色创建成功');
      }
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title={roleId ? '编辑角色' : '新增角色'}
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
