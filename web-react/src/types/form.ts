import type { ColProps } from '@/components/ui';

export interface FormRule {
  required?: boolean;
  message?: string;
  min?: number;
  max?: number;
  type?: 'email';
  pattern?: RegExp;
  validator?: (_rule: FormRule, value: unknown) => Promise<void> | void;
}

export type SchemaComponent =
  | 'Input'
  | 'Password'
  | 'InputNumber'
  | 'Select'
  | 'TreeSelect'
  | 'RadioGroup'
  | 'TextArea'
  | 'Switch'
  | 'MonacoEditor'
  | 'IconPicker';

export interface FormSchema {
  name: string;
  label: string;
  component: SchemaComponent;
  rules?: FormRule[];
  props?: Record<string, unknown>;
  colProps?: ColProps;
  modalColProps?: ColProps;
  initialValue?: unknown;
  helpMessage?: string;
  hidden?: boolean | ((values: Record<string, unknown>) => boolean);
}
