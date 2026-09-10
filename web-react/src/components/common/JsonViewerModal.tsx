import { useMemo, useState } from 'react';
import { App, Button, Space } from '@/components/ui';
import { BulbOutlined, CopyOutlined } from '@/utils/icons';
import { BasicModal } from './BasicModal';
import { MonacoEditor } from './MonacoEditor';

interface JsonViewerModalProps {
  open: boolean;
  title: string;
  data: unknown;
  onCancel: () => void;
}

function normalizeJsonContent(data: unknown) {
  if (data === undefined || data === null) {
    return '';
  }

  try {
    if (typeof data === 'string') {
      return JSON.stringify(JSON.parse(data), null, 2);
    }
    return JSON.stringify(data, null, 2);
  } catch {
    return String(data);
  }
}

export function JsonViewerModal({
  open,
  title,
  data,
  onCancel,
}: JsonViewerModalProps) {
  const { message } = App.useApp();
  const [theme, setTheme] = useState<'vs-dark' | 'vs'>('vs-dark');

  const content = useMemo(() => normalizeJsonContent(data), [data]);

  const handleCopy = async () => {
    if (!content) {
      message.warning('没有可复制的内容');
      return;
    }

    try {
      await navigator.clipboard.writeText(content);
      message.success('复制成功');
    } catch {
      message.error('复制失败');
    }
  };

  return (
    <BasicModal
      open={open}
      title={title}
      width={860}
      footer={null}
      onCancel={onCancel}
    >
      <div style={{ marginBottom: 12 }}>
        <Space>
          <Button icon={<CopyOutlined />} onClick={() => void handleCopy()}>
            复制 JSON
          </Button>
          <Button
            icon={<BulbOutlined />}
            onClick={() => setTheme((value) => (value === 'vs-dark' ? 'vs' : 'vs-dark'))}
          >
            切换主题
          </Button>
        </Space>
      </div>
      <MonacoEditor
        height={360}
        language="json"
        readonly
        theme={theme}
        value={content}
      />
    </BasicModal>
  );
}
