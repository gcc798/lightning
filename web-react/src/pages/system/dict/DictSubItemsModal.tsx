import { useEffect, useState } from 'react';
import type { ColumnsType } from '@/components/ui';
import { App, Button, Space, Table, Tag } from '@/components/ui';
import { PlusOutlined } from '@/utils/icons';
import type { SnowflakeId } from '@/types/api';
import type { DictRecord } from '@/types/system';
import { dictApi } from '@/api/dict';
import { BasicModal } from '@/components/common/BasicModal';
import { TableAction } from '@/components/common/TableAction';
import { isZeroStatus } from '@/utils/number';
import { DictModal } from './DictModal';

interface DictSubItemsModalProps {
  open: boolean;
  parentId?: SnowflakeId;
  onCancel: () => void;
}

export function DictSubItemsModal({
  open,
  parentId,
  onCancel,
}: DictSubItemsModalProps) {
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [subItems, setSubItems] = useState<DictRecord[]>([]);
  const [itemModalOpen, setItemModalOpen] = useState(false);
  const [currentItemId, setCurrentItemId] = useState<SnowflakeId>();

  const loadSubItems = async () => {
    if (!parentId) {
      setSubItems([]);
      return;
    }

    setLoading(true);
    try {
      const parent = await dictApi.detail(parentId);
      const items = await dictApi.getByType(parent.dictType, parentId);
      setSubItems(items);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (open) {
      void loadSubItems();
    }
  }, [open, parentId]);

  const columns: ColumnsType<DictRecord> = [
    { title: '字典标签', dataIndex: 'dictLabel', width: 180 },
    { title: '字典键值', dataIndex: 'dictValue', width: 160 },
    { title: '排序', dataIndex: 'sort', width: 100 },
    {
      title: '默认',
      dataIndex: 'isDefault',
      width: 100,
      render: (value) => <Tag color={value ? 'success' : 'default'}>{value ? '是' : '否'}</Tag>,
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
      width: 180,
      render: (_, record) => (
        <TableAction
          actions={[
            {
              key: 'edit',
              label: '编辑',
              onClick: () => {
                setCurrentItemId(record.id);
                setItemModalOpen(true);
              },
            },
            {
              key: 'delete',
              label: '删除',
              danger: true,
              confirmTitle: '确定删除该子项吗？',
              onClick: async () => {
                await dictApi.delete(record.id);
                message.success('删除成功');
                await loadSubItems();
              },
            },
          ]}
        />
      ),
    },
  ];

  return (
    <>
      <BasicModal
        open={open}
        title="编辑字典子项"
        width={920}
        footer={null}
        onCancel={onCancel}
      >
        <div className="page-toolbar">
          <Space>
            <Button
              icon={<PlusOutlined />}
              type="primary"
              onClick={() => {
                setCurrentItemId(undefined);
                setItemModalOpen(true);
              }}
            >
              新增子项
            </Button>
          </Space>
        </div>
        <Table<DictRecord>
          columns={columns}
          dataSource={subItems}
          loading={loading}
          pagination={false}
          rowKey="id"
        />
      </BasicModal>

      <DictModal
        open={itemModalOpen}
        dictId={currentItemId}
        parentId={parentId}
        onCancel={() => setItemModalOpen(false)}
        onSuccess={() => {
          setItemModalOpen(false);
          void loadSubItems();
        }}
      />
    </>
  );
}
