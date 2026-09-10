import { useEffect, useMemo, useState } from 'react';
import { App, Form } from '@/components/ui';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { DictRecord } from '@/types/system';
import { useAuthStore } from '@/store/auth';
import { dictApi } from '@/api/dict';
import { BasicForm } from '@/components/common/BasicForm';
import { BasicModal } from '@/components/common/BasicModal';

interface DictModalProps {
  open: boolean;
  dictId?: SnowflakeId;
  parentId?: SnowflakeId;
  onCancel: () => void;
  onSuccess: () => void;
}

export function DictModal({
  open,
  dictId,
  parentId,
  onCancel,
  onSuccess,
}: DictModalProps) {
  const [form] = Form.useForm<Partial<DictRecord>>();
  const { message } = App.useApp();
  const userId = useAuthStore((state) => state.userInfo?.userId);
  const [loading, setLoading] = useState(false);
  const [parentDictType, setParentDictType] = useState('');
  const isSubItem = Boolean(parentId);

  const schemas = useMemo<FormSchema[]>(
    () => [
      ...(isSubItem
        ? []
        : [
            {
              name: 'dictType',
              label: '字典类型',
              component: 'Input',
              rules: [{ required: true, message: '请输入字典类型' }],
            } satisfies FormSchema,
          ]),
      {
        name: 'dictLabel',
        label: '字典标签',
        component: 'Input',
        rules: [{ required: true, message: '请输入字典标签' }],
      },
      {
        name: 'dictValue',
        label: '字典键值',
        component: 'Input',
        rules: [{ required: true, message: '请输入字典键值' }],
      },
      {
        name: 'sort',
        label: '字典排序',
        component: 'InputNumber',
        initialValue: 0,
        props: { min: 0 },
      },
      {
        name: 'isDefault',
        label: '是否默认',
        component: 'Switch',
        initialValue: false,
        hidden: !isSubItem,
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
    [isSubItem],
  );

  useEffect(() => {
    if (!open) {
      form.resetFields();
      return;
    }

    if (parentId) {
      void (async () => {
        const parent = await dictApi.detail(parentId);
        setParentDictType(parent.dictType);
      })();
    }

    if (!dictId) {
      form.resetFields();
      form.setFieldsValue({ sort: 0, status: 0, isDefault: false });
      return;
    }

    void (async () => {
      setLoading(true);
      try {
        const detail = await dictApi.detail(dictId);
        form.setFieldsValue(detail);
      } finally {
        setLoading(false);
      }
    })();
  }, [dictId, form, open, parentId]);

  const handleSubmit = async () => {
    const values = await form.validateFields();
    setLoading(true);
    try {
      const payload = {
        ...values,
        parentId,
        dictType: parentId ? parentDictType : values.dictType,
        createBy: userId,
        updateBy: userId,
      };

      if (dictId) {
        await dictApi.update(dictId, payload);
        message.success('字典更新成功');
      } else {
        await dictApi.create(payload);
        message.success('字典创建成功');
      }
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title={dictId ? '编辑字典' : '新增字典'}
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
