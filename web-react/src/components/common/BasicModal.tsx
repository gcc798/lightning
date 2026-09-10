import type { ReactNode } from 'react';
import { Modal } from '@/components/ui';

interface BasicModalProps {
  open: boolean;
  title: ReactNode;
  width?: number;
  footer?: ReactNode | null;
  confirmLoading?: boolean;
  onOk?: () => void;
  onCancel: () => void;
  children: ReactNode;
}

export function BasicModal({
  open,
  title,
  width = 720,
  footer,
  confirmLoading,
  onOk,
  onCancel,
  children,
}: BasicModalProps) {
  return (
    <Modal
      centered
      className="basic-modal"
      destroyOnClose
      open={open}
      title={title}
      width={width}
      footer={footer}
      confirmLoading={confirmLoading}
      onOk={onOk}
      onCancel={onCancel}
    >
      <div className="basic-modal-body">{children}</div>
    </Modal>
  );
}
