import { useMemo, useState } from 'react';
import type { ReactNode } from 'react';
import {
  AppstoreOutlined,
  ApartmentOutlined,
  ApiOutlined,
  CheckCircleOutlined,
  DeleteOutlined,
  DownloadOutlined,
  EditOutlined,
  EyeOutlined,
  KeyOutlined,
  LinkOutlined,
  MoreOutlined,
  PlusOutlined,
  TeamOutlined,
  UserSwitchOutlined,
} from '@/utils/icons';
import { Button, Dropdown, Popconfirm, Space } from '@/components/ui';
import { usePermissionStore } from '@/store/permission';

export interface TableActionItem {
  key: string;
  label: string;
  onClick?: () => void;
  danger?: boolean;
  permission?: string | string[];
  hidden?: boolean;
  confirmTitle?: string;
  icon?: ReactNode;
}

interface TableActionProps {
  actions: TableActionItem[];
}

const defaultActionIcons: Record<string, ReactNode> = {
  add: <PlusOutlined />,
  addChild: <PlusOutlined />,
  apiPermission: <ApiOutlined />,
  assignRole: <UserSwitchOutlined />,
  assignUsers: <TeamOutlined />,
  delete: <DeleteOutlined />,
  download: <DownloadOutlined />,
  edit: <EditOutlined />,
  permission: <ApartmentOutlined />,
  preview: <EyeOutlined />,
  resetPassword: <KeyOutlined />,
  setDefault: <CheckCircleOutlined />,
  subItems: <AppstoreOutlined />,
  test: <LinkOutlined />,
  view: <EyeOutlined />,
};

const fallbackActionIcon = <AppstoreOutlined />;

export function TableAction({ actions }: TableActionProps) {
  const hasPermission = usePermissionStore((state) => state.hasPermission);
  const [menuOpen, setMenuOpen] = useState(false);
  const visibleActions = useMemo(
    () =>
      actions.filter(
        (action) => !action.hidden && hasPermission(action.permission),
      ),
    [actions, hasPermission],
  );
  const primaryInlineAction = useMemo(
    () => visibleActions.find((action) => action.key === 'edit') ?? visibleActions[0],
    [visibleActions],
  );
  const inlineActions =
    visibleActions.length <= 2
      ? visibleActions
      : primaryInlineAction
        ? [primaryInlineAction]
        : [];
  const overflowActions = visibleActions.filter(
    (action) => !inlineActions.includes(action),
  );

  const renderActionContent = (
    action: TableActionItem,
    variant: 'inline' | 'menu',
  ) => {
    if (variant === 'menu') {
      return (
        <button
          className={`table-action-menu-item${action.danger ? ' is-danger' : ''}`}
          type="button"
          onClick={
            action.confirmTitle
              ? undefined
              : () => {
                  action.onClick?.();
                  setMenuOpen(false);
                }
          }
        >
          <span className="table-action-menu-icon">
            {action.icon ?? defaultActionIcons[action.key] ?? fallbackActionIcon}
          </span>
          <span>{action.label}</span>
        </button>
      );
    }

    return (
      <Button
        className="table-action-btn"
        danger={action.danger}
        icon={action.icon ?? defaultActionIcons[action.key] ?? fallbackActionIcon}
        size="small"
        type="text"
        onClick={action.confirmTitle ? undefined : action.onClick}
      >
        {action.label}
      </Button>
    );
  };

  const renderAction = (
    action: TableActionItem,
    variant: 'inline' | 'menu',
  ) => {
    const content = renderActionContent(action, variant);

    if (!action.confirmTitle) {
      return content;
    }

    return (
      <Popconfirm
        title={action.confirmTitle}
        onConfirm={async () => {
          await action.onClick?.();
          setMenuOpen(false);
        }}
      >
        {content}
      </Popconfirm>
    );
  };

  return (
    <Space className="table-action-group" size={4}>
      {inlineActions.map((action) => (
        <span className="table-action-inline-item" key={action.key}>
          {renderAction(action, 'inline')}
        </span>
      ))}
      {overflowActions.length ? (
        <Dropdown
          arrow={false}
          open={menuOpen}
          overlay={(
            <div className="table-action-menu">
              {overflowActions.map((action) => (
                <div className="table-action-menu-row" key={action.key}>
                  {renderAction(action, 'menu')}
                </div>
              ))}
            </div>
          )}
          overlayClassName="table-action-dropdown-overlay"
          trigger={['click']}
          onOpenChange={setMenuOpen}
        >
          <Button
            className="table-action-btn table-action-more-btn"
            icon={<MoreOutlined />}
            size="small"
            type="text"
          >
            更多
          </Button>
        </Dropdown>
      ) : null}
    </Space>
  );
}
