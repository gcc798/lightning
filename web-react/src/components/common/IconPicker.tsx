import { Select, Space, Typography } from '@/components/ui';
import { getMenuIconNode, iconOptions } from '@/utils/icons';

interface IconPickerProps {
  value?: string;
  onChange?: (value: string) => void;
  placeholder?: string;
}

export function IconPicker({
  value,
  onChange,
  placeholder = '请选择图标',
}: IconPickerProps) {
  return (
    <Select
      allowClear
      showSearch
      optionFilterProp="label"
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      options={iconOptions.map((option) => ({
        ...option,
        label: (
          <Space>
            {getMenuIconNode(option.value)}
            <Typography.Text>{option.label}</Typography.Text>
          </Space>
        ),
      }))}
    />
  );
}
