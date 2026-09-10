import { useEffect, useMemo, useState } from 'react';
import { App, Form } from '@/components/ui';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { OrgRecord } from '@/types/system';
import { orgApi } from '@/api/org';
import { BasicForm } from '@/components/common/BasicForm';
import { BasicModal } from '@/components/common/BasicModal';

interface OrgModalProps {
  open: boolean;
  orgId?: SnowflakeId;
  parentId?: SnowflakeId;
  onCancel: () => void;
  onSuccess: () => void;
}

function toTreeOptions(nodes: OrgRecord[]): Record<string, unknown>[] {
  return nodes.map((node) => ({
    title: node.orgName,
    value: node.id,
    key: node.id,
    children: node.children ? toTreeOptions(node.children) : undefined,
  }));
}

export function OrgModal({
  open,
  orgId,
  parentId,
  onCancel,
  onSuccess,
}: OrgModalProps) {
  const [form] = Form.useForm<Partial<OrgRecord>>();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [orgOptions, setOrgOptions] = useState<Record<string, unknown>[]>([]);

  const schemas = useMemo<FormSchema[]>(
    () => [
      {
        name: 'parentId',
        label: '上级组织',
        component: 'TreeSelect',
        initialValue: parentId ?? 0,
        props: {
          treeData: [{ title: '顶级组织', value: 0, key: 0 }, ...orgOptions],
          treeDefaultExpandAll: true,
          placeholder: '请选择上级组织',
        },
      },
      {
        name: 'orgName',
        label: '组织名称',
        component: 'Input',
        rules: [{ required: true, message: '请输入组织名称' }],
      },
      {
        name: 'orgCode',
        label: '组织编码',
        component: 'Input',
        rules: [{ required: true, message: '请输入组织编码' }],
      },
      {
        name: 'leader',
        label: '负责人',
        component: 'Input',
      },
      {
        name: 'phone',
        label: '联系电话',
        component: 'Input',
      },
      {
        name: 'email',
        label: '邮箱',
        component: 'Input',
        rules: [{ type: 'email', message: '邮箱格式不正确' }],
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
    [orgOptions, parentId],
  );

  useEffect(() => {
    if (!open) {
      form.resetFields();
      return;
    }

    void (async () => {
      const data = await orgApi.tree();
      setOrgOptions(toTreeOptions(data));
    })();

    if (!orgId) {
      form.resetFields();
      form.setFieldsValue({ parentId: parentId ?? 0, sort: 0, status: 0 });
      return;
    }

    void (async () => {
      setLoading(true);
      try {
        const detail = await orgApi.detail(orgId);
        form.setFieldsValue(detail);
      } finally {
        setLoading(false);
      }
    })();
  }, [form, open, orgId, parentId]);

  const handleSubmit = async () => {
    const values = await form.validateFields();
    setLoading(true);
    try {
      if (orgId) {
        await orgApi.update(orgId, values);
        message.success('组织更新成功');
      } else {
        await orgApi.create(values);
        message.success('组织创建成功');
      }
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title={orgId ? '编辑组织' : '新增组织'}
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
