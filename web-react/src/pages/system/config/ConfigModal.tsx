import { useEffect, useMemo, useState } from 'react';
import { App, Button, Form } from '@/components/ui';
import type { SnowflakeId } from '@/types/api';
import type { FormSchema } from '@/types/form';
import type { ConfigRecord } from '@/types/system';
import { useAuthStore } from '@/store/auth';
import { configApi } from '@/api/config';
import { BasicForm } from '@/components/common/BasicForm';
import { BasicModal } from '@/components/common/BasicModal';

interface ConfigModalProps {
  open: boolean;
  configId?: SnowflakeId;
  onCancel: () => void;
  onSuccess: () => void;
}

function safeStringify(value: unknown) {
  if (value === undefined || value === null) {
    return '';
  }
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

export function ConfigModal({
  open,
  configId,
  onCancel,
  onSuccess,
}: ConfigModalProps) {
  const [form] = Form.useForm<Partial<ConfigRecord> & { data?: string }>();
  const { message } = App.useApp();
  const userId = useAuthStore((state) => state.userInfo?.userId);
  const [loading, setLoading] = useState(false);
  const [editorTheme, setEditorTheme] = useState<'vs-dark' | 'vs'>('vs-dark');

  const schemas = useMemo<FormSchema[]>(
    () => [
      {
        name: 'name',
        label: '配置名称',
        component: 'Input',
        rules: [{ required: true, message: '请输入配置名称' }],
      },
      {
        name: 'code',
        label: '配置编码',
        component: 'Input',
        rules: [{ required: true, message: '请输入配置编码' }],
      },
      {
        name: 'data',
        label: '配置数据',
        component: 'MonacoEditor',
        props: {
          height: 320,
          language: 'json',
          theme: editorTheme,
        },
        rules: [
          {
            validator: async (_: unknown, value: unknown) => {
              if (!value) {
                return;
              }
              JSON.parse(String(value));
            },
            message: '请输入合法的 JSON',
          },
        ],
      },
      {
        name: 'remark',
        label: '备注',
        component: 'TextArea',
        props: { rows: 3 },
      },
    ],
    [editorTheme],
  );

  useEffect(() => {
    if (!open) {
      form.resetFields();
      return;
    }

    if (!configId) {
      form.resetFields();
      return;
    }

    void (async () => {
      setLoading(true);
      try {
        const detail = await configApi.detail(configId);
        form.setFieldsValue({
          ...detail,
          data: safeStringify(detail.data),
        });
      } finally {
        setLoading(false);
      }
    })();
  }, [configId, form, open]);

  const handleSubmit = async () => {
    const values = await form.validateFields();
    if (!userId) {
      message.error('用户信息不存在，请重新登录');
      return;
    }

    setLoading(true);
    try {
      const payload = {
        ...values,
        data: values.data ? JSON.parse(values.data) : undefined,
        createBy: userId,
        updateBy: userId,
      };

      if (configId) {
        await configApi.update(configId, payload);
        message.success('配置更新成功');
      } else {
        await configApi.create(payload);
        message.success('配置创建成功');
      }
      onSuccess();
    } finally {
      setLoading(false);
    }
  };

  return (
    <BasicModal
      open={open}
      title={configId ? '编辑配置' : '新增配置'}
      width={760}
      confirmLoading={loading}
      onCancel={onCancel}
      onOk={() => void handleSubmit()}
    >
      <div style={{ marginBottom: 12 }}>
        <Button onClick={() => setEditorTheme((value) => (value === 'vs-dark' ? 'vs' : 'vs-dark'))}>
          切换编辑器主题
        </Button>
      </div>
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
