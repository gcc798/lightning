import {
  useContext,
  useEffect,
  useId,
  useMemo,
  useRef,
  useState,
  createContext,
  createElement,
  Fragment,
  type CSSProperties,
  type ButtonHTMLAttributes,
  type FormEvent,
  type HTMLAttributes,
  type InputHTMLAttributes,
  type Key,
  type ReactElement,
  type ReactNode,
} from 'react';
import { Button as AppicaButton } from '@appica/ui-react/button';
import { Card as AppicaCard } from '@appica/ui-react/card';
import { Dialog, DialogBody, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@appica/ui-react/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@appica/ui-react/dropdown-menu';
import { Input as AppicaInput } from '@appica/ui-react/input';
import { Badge } from '@appica/ui-react/badge';
import { Checkbox } from '@appica/ui-react/checkbox';
import { Progress as AppicaProgress } from '@appica/ui-react/progress';
import { Loader } from '@appica/ui-react/loader';
import {
  Select as AppicaSelect,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@appica/ui-react/select';
import { Switch as AppicaSwitch } from '@appica/ui-react/switch';
import { Table as AppicaTable } from '@appica/ui-react/table';
import { ToastProvider, Toaster, useToastManager } from '@appica/ui-react/toast';

type AnyRecord = Record<string, unknown>;

export interface ButtonProps extends Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'type'> {
  type?: 'button' | 'submit' | 'reset' | 'primary' | 'text' | 'link' | 'default';
  htmlType?: 'button' | 'submit' | 'reset';
  danger?: boolean;
  loading?: boolean;
  block?: boolean;
  size?: 'small' | 'middle' | 'large' | 'sm' | 'md' | 'lg';
  icon?: ReactNode;
}

export function Button({ type = 'default', htmlType, danger, loading, block, size = 'small', icon, className, children, disabled, ...props }: ButtonProps) {
  const variant = danger ? 'destructive' : type === 'primary' ? 'primary' : type === 'text' || type === 'link' ? 'ghost' : type === 'default' ? 'outline' : 'secondary';
  const appicaSize = size === 'small' || size === 'sm' ? 'sm' : size === 'large' || size === 'lg' ? 'lg' : 'md';
  return <AppicaButton {...props} type={htmlType ?? (type === 'submit' || type === 'reset' || type === 'button' ? type : 'button')} disabled={disabled || loading} variant={variant} size={icon && !children ? `icon-${appicaSize}` as 'icon-sm' | 'icon-md' | 'icon-lg' : appicaSize} className={`${block ? 'w-full' : ''} ${className ?? ''}`}>{loading ? <Loader currentColor variant="dots" /> : icon}{children ? <span>{children}</span> : null}</AppicaButton>;
}

export interface CardProps extends HTMLAttributes<HTMLDivElement> { variant?: 'borderless' | 'outlined'; bordered?: boolean }
export function Card({ className, children, variant: _variant, bordered: _bordered, ...props }: CardProps) { return <AppicaCard frame={false} contentProps={{ className: 'ui-card-content' }} className={`page-card ${className ?? ''}`} {...props}>{children}</AppicaCard>; }

interface SpaceProps extends HTMLAttributes<HTMLDivElement> { direction?: 'horizontal' | 'vertical'; size?: number | 'small' | 'middle' | 'large'; wrap?: boolean; align?: CSSProperties['alignItems']; justify?: CSSProperties['justifyContent'] }
export function Space({ direction = 'horizontal', size = 'middle', wrap, align, justify, style, className, children, ...props }: SpaceProps) {
  const gap = typeof size === 'number' ? size : size === 'small' ? 8 : size === 'large' ? 24 : 16;
  return <div className={`ui-space ${className ?? ''}`} style={{ display: 'flex', flexDirection: direction === 'vertical' ? 'column' : 'row', gap, alignItems: align ?? (direction === 'vertical' ? 'stretch' : 'center'), justifyContent: justify, flexWrap: wrap ? 'wrap' : undefined, ...style }} {...props}>{children}</div>;
}
export function Row({ gutter, className, children, style, ...props }: HTMLAttributes<HTMLDivElement> & { gutter?: number | [number, number] }) { const [x, y] = Array.isArray(gutter) ? gutter : [gutter ?? 0, gutter ?? 0]; return <div className={`ui-row ${className ?? ''}`} style={{ display: 'grid', gridTemplateColumns: 'repeat(24, minmax(0, 1fr))', columnGap: 0, rowGap: y, '--row-gutter-x': `${x}px`, ...style } as CSSProperties} {...props}>{children}</div>; }
export function Col({ span = 24, xs, sm, md, lg, xl, xxl, className, children, style, ...props }: HTMLAttributes<HTMLDivElement> & ColProps) { return <div className={`ui-col ${className ?? ''}`} style={{ '--col-base': span, '--col-xs': xs ?? span, '--col-sm': sm ?? xs ?? span, '--col-md': md ?? sm ?? xs ?? span, '--col-lg': lg ?? md ?? sm ?? xs ?? span, '--col-xl': xl ?? lg ?? md ?? sm ?? xs ?? span, '--col-xxl': xxl ?? xl ?? lg ?? md ?? sm ?? xs ?? span, ...style } as CSSProperties} {...props}>{children}</div>; }

function Text({ type, strong, className, children, ...props }: HTMLAttributes<HTMLElement> & { type?: 'secondary' | 'danger'; strong?: boolean }) { return <span className={`${type === 'secondary' ? 'text-foreground-muted' : type === 'danger' ? 'text-error-emphasis' : ''} ${strong ? 'font-semibold' : ''} ${className ?? ''}`} {...props}>{children}</span>; }
function Title({ level = 1, className, children, ...props }: HTMLAttributes<HTMLHeadingElement> & { level?: 1 | 2 | 3 | 4 | 5 }) { return createElement(`h${level}`, { className: `ui-title ${className ?? ''}`, ...props }, children); }
function Paragraph({ className, children, ...props }: HTMLAttributes<HTMLParagraphElement>) { return <p className={className} {...props}>{children}</p>; }
export const Typography = { Text, Title, Paragraph };

type InputProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'size' | 'prefix'> & { allowClear?: boolean; prefix?: ReactNode; suffix?: ReactNode; size?: 'small' | 'middle' | 'large'; status?: string };
function InputBase({ prefix, suffix, allowClear, size, status: _status, className, onChange, ...props }: InputProps) { return <AppicaInput startSlot={prefix} endSlot={suffix} clearable={allowClear} inputSize={size === 'large' ? 'lg' : size === 'small' ? 'sm' : 'md'} className={`ui-appica-input ${className ?? ''}`} onChange={onChange} onClear={() => onChange?.({ target: { value: '' } } as never)} {...props} />; }
function Password(props: InputProps) { return <InputBase type="password" {...props} />; }
function TextArea({ className, ...props }: InputHTMLAttributes<HTMLTextAreaElement>) { return <textarea className={`ui-textarea ${className ?? ''}`} {...props} />; }
function Search(props: InputProps) { return <InputBase type="search" {...props} />; }
export const Input = Object.assign(InputBase, { Password, TextArea, Search });

export function InputNumber({ value, onChange, className, ...props }: InputHTMLAttributes<HTMLInputElement> & { value?: number | string; onChange?: (value: number | null) => void }) { return <input {...props} className={`ui-input ${className ?? ''}`} type="number" value={value ?? ''} onChange={(event) => onChange?.(event.target.value === '' ? null : Number(event.target.value))} />; }

interface Option { label: ReactNode; value: string | number; disabled?: boolean }
type SelectValueType = string | number;
interface SelectProps { value?: SelectValueType | SelectValueType[] | null; defaultValue?: SelectValueType | SelectValueType[]; options?: Option[]; onChange?: (value: any) => void; placeholder?: string; allowClear?: boolean; className?: string; style?: CSSProperties; disabled?: boolean; children?: ReactNode; showSearch?: boolean; optionFilterProp?: string; mode?: 'multiple'; size?: 'small' | 'middle' | 'large' }
export function Select({ value, defaultValue, options = [], onChange, placeholder, allowClear, className, style, disabled, mode, size = 'middle' }: SelectProps) {
  const multiple = mode === 'multiple';
  return (
    <AppicaSelect
      size={size === 'small' ? 'sm' : size === 'large' ? 'lg' : 'md'}
      items={options}
      multiple={multiple as never}
      value={(value ?? (multiple ? [] : null)) as never}
      defaultValue={defaultValue as never}
      disabled={disabled}
      onValueChange={(next) => onChange?.(next)}
    >
      <SelectTrigger
        aria-label={placeholder ?? '请选择'}
        className={`ui-select ${className ?? ''}`}
        clearable={allowClear}
        style={style}
      >
        <SelectValue placeholder={placeholder ?? '请选择'} />
      </SelectTrigger>
      <SelectContent>
        {options.map((option) => (
          <SelectItem disabled={option.disabled} key={String(option.value)} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </AppicaSelect>
  );
}

interface TreeDataNode { title: ReactNode; value?: Key; key?: Key; children?: TreeDataNode[]; disabled?: boolean; disableCheckbox?: boolean }
function flattenTree(nodes: TreeDataNode[], depth = 0): Array<Option & { depth: number }> { return nodes.flatMap((node) => [{ label: `${'  '.repeat(depth)}${String(node.title)}`, value: String(node.value ?? node.key ?? ''), disabled: node.disabled, depth }, ...flattenTree(node.children ?? [], depth + 1)]); }
export function TreeSelect({ treeData = [], ...props }: SelectProps & { treeData?: TreeDataNode[] }) { return <Select {...props} options={flattenTree(treeData)} />; }

interface RadioGroupProps { value?: string | number; options?: Option[]; onChange?: (event: { target: { value: any } }) => void; optionType?: 'button' | 'default'; className?: string }
function RadioGroup({ value, options = [], onChange, optionType = 'default', className }: RadioGroupProps) { return <div className={`ui-radio-group ${optionType === 'button' ? 'is-button' : ''} ${className ?? ''}`}>{options.map((option) => <label key={String(option.value)}><input type="radio" checked={String(value ?? '') === String(option.value)} onChange={() => onChange?.({ target: { value: option.value } })} /><span>{option.label}</span></label>)}</div>; }
export const Radio = { Group: RadioGroup };

export function Switch({ checked, defaultChecked, onChange, className, disabled }: { checked?: boolean; defaultChecked?: boolean; onChange?: (checked: boolean) => void; className?: string; disabled?: boolean }) { return <AppicaSwitch className={`ui-switch ${className ?? ''}`} checked={checked} defaultChecked={defaultChecked} disabled={disabled} onCheckedChange={(next) => onChange?.(next)} />; }

function mapBadgeVariant(color?: string): 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info' | 'outline' | 'soft' { if (color === 'success' || color === 'green') return 'success'; if (color === 'error' || color === 'red') return 'error'; if (color === 'warning' || color === 'orange' || color === 'gold') return 'warning'; if (color === 'blue' || color === 'processing') return 'info'; return color === 'default' ? 'outline' : 'soft'; }
export function Tag({ color, className, children, ...props }: HTMLAttributes<HTMLSpanElement> & { color?: string }) { return <Badge variant={mapBadgeVariant(color)} className={className} {...props}>{children}</Badge>; }

export function Avatar({ src, icon, size = 28, shape = 'circle', className }: { src?: string; icon?: ReactNode; size?: number | string; shape?: 'circle' | 'square'; className?: string }) { return <span className={`ui-avatar ${shape === 'square' ? 'is-square' : ''} ${className ?? ''}`} style={{ width: typeof size === 'number' ? size : undefined, height: typeof size === 'number' ? size : undefined }}>{src ? <img src={src} alt="" /> : icon}</span>; }

export function Tooltip({ title, children }: { title?: ReactNode; children: ReactNode }) { return <span title={typeof title === 'string' ? title : undefined}>{children}</span>; }
export function Spin({ spinning = true, children }: { spinning?: boolean; size?: 'small' | 'default' | 'large'; children?: ReactNode }) { return <div className="ui-spin">{spinning ? <div className="ui-spin-indicator"><Loader aria-label="加载中" /></div> : null}{children}</div>; }
export function Flex({ children, style, className, vertical, gap = 0, align = 'normal', justify = 'normal' }: { children?: ReactNode; style?: CSSProperties; className?: string; vertical?: boolean; gap?: number; align?: string; justify?: string }) { return <div className={className} style={{ display: 'flex', flexDirection: vertical ? 'column' : 'row', gap, alignItems: align, justifyContent: justify, ...style }}>{children}</div>; }
export function Progress({ percent, showInfo = true, strokeColor, size }: { percent: number; showInfo?: boolean; strokeColor?: string; size?: 'small' | 'default' }) { return <div className="ui-progress"><AppicaProgress value={percent} indicatorColor={strokeColor} />{showInfo ? <span>{percent}%</span> : null}</div>; }

export function Result({ status, title, subTitle, extra }: { status?: string; title?: ReactNode; subTitle?: ReactNode; extra?: ReactNode }) { return <div className="ui-result"><div className="ui-result-code">{status}</div><Title level={2}>{title}</Title><Text type="secondary">{subTitle}</Text><div>{extra}</div></div>; }

export function Layout({ children, className, ...props }: HTMLAttributes<HTMLDivElement>) { return <div className={`ui-layout ${className ?? ''}`} {...props}>{children}</div>; }
export function Header({ children, className, ...props }: HTMLAttributes<HTMLDivElement>) { return <header className={`ui-header ${className ?? ''}`} {...props}>{children}</header>; }
export function Content({ children, className, ...props }: HTMLAttributes<HTMLDivElement>) { return <main className={`ui-content ${className ?? ''}`} {...props}>{children}</main>; }
export function Sider({ children, collapsed, width = 224, className, onCollapse: _onCollapse, collapsible: _collapsible, breakpoint: _breakpoint, trigger: _trigger, theme: _theme, ...props }: HTMLAttributes<HTMLElement> & { collapsed?: boolean; width?: number; onCollapse?: (collapsed: boolean) => void; collapsible?: boolean; breakpoint?: string; trigger?: ReactNode; theme?: string }) { return <aside className={`ui-sider ${collapsed ? 'is-collapsed' : ''} ${className ?? ''}`} style={{ width: collapsed ? 72 : width }} {...props}>{children}</aside>; }

export interface MenuItem { key: string; label: ReactNode; icon?: ReactNode; children?: MenuItem[]; onClick?: () => void; onTitleClick?: () => void }
export interface MenuProps { items?: MenuItem[]; selectedKeys?: string[]; openKeys?: string[]; onOpenChange?: (keys: string[]) => void; className?: string; inlineIndent?: number; mode?: string; theme?: string }
export function Menu({ items = [], selectedKeys = [], openKeys = [], onOpenChange, className }: MenuProps) {
  const renderItems = (nodes: MenuItem[], depth = 0): ReactNode => nodes.map((item) => { const hasChildren = Boolean(item.children?.length); const open = openKeys.includes(item.key); return <div key={item.key} className={`ui-menu-item-wrap depth-${depth}`}><button className={`ui-menu-item ${selectedKeys.includes(item.key) ? 'is-selected' : ''}`} type="button" onClick={() => { if (hasChildren) onOpenChange?.(open ? openKeys.filter((key) => key !== item.key) : [...openKeys, item.key]); item.onTitleClick?.(); item.onClick?.(); }}><span className="ui-menu-icon">{item.icon}</span><span>{item.label}</span>{hasChildren ? <span className="ui-menu-chevron">{open ? '−' : '+'}</span> : null}</button>{hasChildren && open ? <div className="ui-menu-children">{renderItems(item.children ?? [], depth + 1)}</div> : null}</div>; });
  return <nav className={`ui-menu ${className ?? ''}`}>{renderItems(items)}</nav>;
}

export function Dropdown({ menu, overlay, children, open, onOpenChange, overlayClassName }: { menu?: { items?: Array<{ key: string; icon?: ReactNode; label: ReactNode; onClick?: () => void }> }; overlay?: ReactNode; children: ReactNode; open?: boolean; onOpenChange?: (open: boolean) => void; arrow?: boolean; trigger?: string[]; overlayClassName?: string }) {
  const [internalOpen, setInternalOpen] = useState(false);
  const visible = open ?? internalOpen;
  const handleOpenChange = (next: boolean) => {
    setInternalOpen(next);
    onOpenChange?.(next);
  };

  return (
    <DropdownMenu modal={false} open={visible} size="sm" onOpenChange={handleOpenChange}>
      <DropdownMenuTrigger render={children as ReactElement} />
      <DropdownMenuContent align="end" className={`ui-dropdown-panel ${overlayClassName ?? ''}`}>
        {overlay ?? menu?.items?.map((item) => (
          <DropdownMenuItem key={item.key} onClick={item.onClick}>
            {item.icon}
            {item.label}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export function Popconfirm({ title, onConfirm, children }: { title?: ReactNode; onConfirm?: () => void | Promise<void>; children: ReactNode }) { return <span onClick={(event) => { event.stopPropagation(); if (window.confirm(String(title ?? '确定执行此操作吗？'))) void onConfirm?.(); }}>{children}</span>; }

interface ModalProps { open: boolean; title?: ReactNode; width?: number; footer?: ReactNode | null; confirmLoading?: boolean; onOk?: () => void; onCancel: () => void; children?: ReactNode; centered?: boolean; destroyOnClose?: boolean; className?: string }
function ModalComponent({ open, title, width = 720, footer, confirmLoading, onOk, onCancel, children, className }: ModalProps) { return <Dialog open={open} onOpenChange={(next) => { if (!next) onCancel(); }}><DialogContent className={`ui-modal-content ${className ?? ''}`} style={{ width: `min(${width}px, calc(100vw - 32px))` }}><DialogHeader><DialogTitle>{title}</DialogTitle></DialogHeader><DialogBody>{children}</DialogBody>{footer !== null ? <DialogFooter>{footer ?? <><Button onClick={onCancel}>取消</Button><Button loading={confirmLoading} type="primary" onClick={onOk}>确定</Button></>}</DialogFooter> : null}</DialogContent></Dialog>; }
export const Modal = Object.assign(ModalComponent, { confirm: ({ title, content, onOk }: { title?: ReactNode; content?: ReactNode; onOk?: () => void | Promise<void> }) => { if (window.confirm([title, content].filter(Boolean).join('\n'))) void onOk?.(); } });

interface UiRule { required?: boolean; message?: string; min?: number; max?: number; type?: 'email'; pattern?: RegExp; validator?: (rule: UiRule, value: unknown) => Promise<void> | void }
interface FormStore { values: AnyRecord; rules: Map<string, UiRule[]>; errors: Map<string, string>; listeners: Set<() => void>; setValues: (values: AnyRecord) => void; setValue: (name: string, value: unknown) => void; reset: () => void; notify: () => void; }
const FormContext = createContext<FormStore | null>(null);
export interface FormInstance<T extends object = AnyRecord> { getFieldsValue: () => T; setFieldsValue: (values: Partial<T>) => void; setFieldValue: (name: string, value: unknown) => void; resetFields: () => void; validateFields: () => Promise<T>; __store: FormStore; }
function useForm<T extends object = AnyRecord>(): [FormInstance<T>] {
  const storeRef = useRef<FormStore | null>(null);
  if (!storeRef.current) {
    const store: FormStore = {
      values: {}, rules: new Map(), errors: new Map(), listeners: new Set(),
      notify: () => store.listeners.forEach((listener) => listener()),
      setValues: (values) => { store.values = { ...store.values, ...values }; store.notify(); },
      setValue: (name, value) => { store.values = { ...store.values, [name]: value }; store.errors.delete(name); store.notify(); },
      reset: () => { store.values = {}; store.errors.clear(); store.notify(); },
    };
    storeRef.current = store;
  }
  const store = storeRef.current;
  return [useMemo(() => ({
    getFieldsValue: () => store.values as T,
    setFieldsValue: (values: Partial<T>) => store.setValues(values as AnyRecord),
    setFieldValue: store.setValue,
    resetFields: store.reset,
    validateFields: async () => {
      store.errors.clear();
      for (const [name, rules] of store.rules) {
        const value = store.values[name];
        for (const rule of rules) {
          try {
            const empty = value === undefined || value === null || value === '';
            if (rule.required && empty) throw new Error(rule.message ?? '此项为必填项');
            if (!empty && rule.min !== undefined && String(value).length < rule.min) throw new Error(rule.message ?? `至少输入 ${rule.min} 个字符`);
            if (!empty && rule.max !== undefined && String(value).length > rule.max) throw new Error(rule.message ?? `最多输入 ${rule.max} 个字符`);
            if (!empty && rule.type === 'email' && !/^\S+@\S+\.\S+$/.test(String(value))) throw new Error(rule.message ?? '邮箱格式不正确');
            if (!empty && rule.pattern && !rule.pattern.test(String(value))) throw new Error(rule.message ?? '格式不正确');
            await rule.validator?.(rule, value);
          } catch (error) {
            store.errors.set(name, error instanceof Error && error.message !== 'Unexpected end of JSON input' ? error.message : rule.message ?? '格式不正确');
            break;
          }
        }
      }
      store.notify();
      if (store.errors.size) throw new Error('表单校验失败');
      return store.values as T;
    },
    __store: store,
  }), [store])];
}

function useWatch(_name: string | string[] | undefined, form: FormInstance) { const [, rerender] = useState(0); useEffect(() => { const listener = () => rerender((value) => value + 1); form.__store.listeners.add(listener); return () => { form.__store.listeners.delete(listener); }; }, [form]); return form.getFieldsValue(); }
interface FormProps<T extends object = AnyRecord> { form?: FormInstance<T>; initialValues?: Partial<T>; layout?: 'horizontal' | 'vertical' | 'inline'; onFinish?: (values: T) => void; className?: string; children?: ReactNode }
function FormComponent<T extends object>({ form: providedForm, initialValues, onFinish, className, children }: FormProps<T>) { const [internalForm] = useForm<T>(); const form = providedForm ?? internalForm; const initialized = useRef(false); useEffect(() => { if (!initialized.current && initialValues) { form.setFieldsValue(initialValues); initialized.current = true; } }, [form, initialValues]); const handleSubmit = (event: FormEvent<HTMLFormElement>) => { event.preventDefault(); void form.validateFields().then((values) => onFinish?.(values)).catch(() => undefined); }; return <FormContext.Provider value={form.__store}><form className={className} onSubmit={handleSubmit}>{children}</form></FormContext.Provider>; }
interface UiFormItemProps { name?: string; label?: ReactNode; rules?: UiRule[]; initialValue?: unknown; valuePropName?: string; extra?: ReactNode; children?: ReactNode; className?: string; style?: CSSProperties }
function FormItem({ name, label, rules = [], initialValue, valuePropName = 'value', extra, children, className, style }: UiFormItemProps) { const context = useContext(FormContext); const [, rerender] = useState(0); useEffect(() => { if (name && context) { context.rules.set(name, rules); if (initialValue !== undefined && context.values[name] === undefined) context.setValue(name, initialValue); } return () => { if (name) context?.rules.delete(name); }; }, [context, initialValue, name, rules]); useEffect(() => { const listener = () => rerender((value) => value + 1); context?.listeners.add(listener); return () => { context?.listeners.delete(listener); }; }, [context]); const value = name ? context?.values[name] : undefined; const child = Array.isArray(children) ? children[0] : children; const childElement = child && typeof child === 'object' && 'type' in child ? createElement((child as { type: unknown }).type as never, { ...(child as { props: AnyRecord }).props, ...(name ? { [valuePropName]: value ?? '', onChange: (next: unknown) => context?.setValue(name, next && typeof next === 'object' && 'target' in (next as AnyRecord) ? (next as { target: { value: unknown } }).target.value : next) } : {}) }) : children; const required = rules.some((rule) => rule.required); const error = name ? context?.errors.get(name) : undefined; return <div className={`ui-form-item ${error ? 'is-invalid' : ''} ${className ?? ''}`} style={style}>{label ? <label>{label}{required ? <span className="ui-required">*</span> : null}</label> : null}{childElement}{error ? <div className="ui-form-error">{error}</div> : extra ? <div className="ui-form-extra">{extra}</div> : null}</div>; }
export const Form = Object.assign(FormComponent, { useForm, useWatch, Item: FormItem });

interface TableColumn<T> { key?: string; title?: ReactNode; dataIndex?: keyof T | string; width?: number; fixed?: 'left' | 'right' | boolean; render?: (value: any, record: T, index: number) => ReactNode; children?: TableColumn<T>[] }
export type ColumnsType<T> = TableColumn<T>[];
export type TableColumnsType<T> = ColumnsType<T>;
export interface TableProps<T> { columns: ColumnsType<T>; dataSource?: T[]; rowKey?: keyof T | ((record: T) => Key); loading?: boolean; rowSelection?: { selectedRowKeys?: Key[]; onChange?: (keys: Key[], rows: T[]) => void; getCheckboxProps?: (record: T) => { disabled?: boolean } }; pagination?: false | { current?: number; pageSize?: number; total?: number; showSizeChanger?: boolean; showQuickJumper?: boolean; showTotal?: (total: number) => ReactNode; onChange?: (page: number, pageSize: number) => void }; scroll?: { x?: number | string | true; y?: number | string }; expandable?: { expandedRowKeys?: Key[]; onExpandedRowsChange?: (keys: Key[]) => void }; size?: 'small' | 'middle' | 'large'; className?: string; }
function resolveKey<T>(record: T, rowKey: TableProps<T>['rowKey'], index: number): Key { return typeof rowKey === 'function' ? rowKey(record) : rowKey ? (record as AnyRecord)[rowKey as string] as Key : index; }
function sameKey(left: Key, right: Key) { return String(left) === String(right); }
function TableComponent<T>({ columns, dataSource = [], rowKey, loading, rowSelection, pagination, expandable, className }: TableProps<T>) { const [internalPage, setInternalPage] = useState(1); const page = pagination ? pagination.current ?? internalPage : 1; const size = pagination ? pagination.pageSize ?? 10 : dataSource.length || 1; const rows = dataSource; const selection = rowSelection?.selectedRowKeys ?? []; const selectColumn: TableColumn<T> = { key: '__select', title: '', render: (_value, record) => { const key = resolveKey(record, rowKey, 0); return <Checkbox aria-label={`选择第 ${String(key)} 行`} disabled={rowSelection?.getCheckboxProps?.(record).disabled} checked={selection.some((item) => sameKey(item, key))} onCheckedChange={(next) => { const selected = next === true ? [...selection, key] : selection.filter((item) => !sameKey(item, key)); rowSelection?.onChange?.(selected, dataSource.filter((item) => selected.some((selectedKey) => sameKey(selectedKey, resolveKey(item, rowKey, 0))))); }} />; } }; const renderColumns: TableColumn<T>[] = rowSelection ? [selectColumn, ...columns] : columns; const renderRow = (record: T, index: number, depth = 0): ReactNode => { const key = resolveKey(record, rowKey, index); const expanded = expandable?.expandedRowKeys?.some((item) => sameKey(item, key)) ?? true; const children = (record as AnyRecord).children as T[] | undefined; return <Fragment key={String(key)}><tr className={depth ? 'ui-table-child-row' : ''}>{renderColumns.map((column, columnIndex) => { const value = column.dataIndex ? (record as AnyRecord)[column.dataIndex as string] : undefined; const content = column.render ? column.render(value, record, index) : value as ReactNode; const treeToggle = expandable && columnIndex === (rowSelection ? 1 : 0) ? <button className="ui-table-tree-toggle" style={{ marginLeft: depth * 16 }} type="button" disabled={!children?.length} onClick={() => expandable.onExpandedRowsChange?.(expanded ? (expandable.expandedRowKeys ?? []).filter((item) => !sameKey(item, key)) : [...(expandable.expandedRowKeys ?? []), key])}>{children?.length ? (expanded ? '−' : '+') : ''}</button> : null; return <td key={`${String(column.key ?? column.dataIndex ?? 'column')}-${columnIndex}`} style={{ width: column.width }}>{treeToggle}{content}</td>; })}</tr>{expandable && children && expanded ? children.map((child, childIndex) => renderRow(child, childIndex, depth + 1)) : null}</Fragment>; }; return <div className={`ui-table-wrap ${className ?? ''}`}>{loading ? <div className="ui-table-loading"><Loader /></div> : null}<AppicaTable className="ui-table" hoverableRows size="sm"><thead><tr>{renderColumns.map((column, index) => <th key={`${String(column.key ?? column.dataIndex ?? 'column')}-${index}`} style={{ width: column.width }}>{column.title}</th>)}</tr></thead><tbody>{rows.map((record, index) => renderRow(record, index))}</tbody></AppicaTable>{pagination ? <div className="ui-pagination"><span>{pagination.showTotal?.(pagination.total ?? dataSource.length) ?? `共 ${pagination.total ?? dataSource.length} 条`}</span>{pagination.showSizeChanger ? <select aria-label="每页条数" value={size} onChange={(event) => { const nextSize = Number(event.target.value); setInternalPage(1); pagination.onChange?.(1, nextSize); }}>{[10, 20, 50, 100].map((value) => <option key={value} value={value}>{value} 条/页</option>)}</select> : null}<button type="button" disabled={page <= 1} onClick={() => { const next = page - 1; setInternalPage(next); pagination.onChange?.(next, size); }}>上一页</button><span>{page} / {Math.max(1, Math.ceil((pagination.total ?? dataSource.length) / size))}</span><button type="button" disabled={page * size >= (pagination.total ?? dataSource.length)} onClick={() => { const next = page + 1; setInternalPage(next); pagination.onChange?.(next, size); }}>下一页</button></div> : null}</div>; }
export const Table = TableComponent;

export function Descriptions({ children, column = 1 }: { children?: ReactNode; bordered?: boolean; column?: number }) { return <dl className={`ui-descriptions columns-${column}`}>{children}</dl>; }
Descriptions.Item = function DescriptionItem({ label, children, span = 1 }: { label: ReactNode; children?: ReactNode; span?: number }) { return <div className={`ui-description-item span-${span}`}><dt>{label}</dt><dd>{children}</dd></div>; };

interface TreeProps { treeData?: any[]; checkable?: boolean; checkedKeys?: Key[] | { checked: Key[] }; onCheck?: (keys: Key[] | { checked: Key[] }) => void; selectable?: boolean; expandedKeys?: Key[]; onExpand?: (keys: Key[]) => void; height?: number; className?: string; checkStrictly?: boolean; defaultExpandAll?: boolean; fieldNames?: { title?: string; key?: string; children?: string }; titleRender?: (node: any) => ReactNode }
export function Tree({ treeData = [], checkable, checkedKeys = [], onCheck, expandedKeys, onExpand, className, defaultExpandAll, fieldNames, titleRender }: TreeProps) { const titleField = fieldNames?.title ?? 'title'; const keyField = fieldNames?.key ?? 'key'; const childrenField = fieldNames?.children ?? 'children'; const checked = Array.isArray(checkedKeys) ? checkedKeys : checkedKeys.checked; const allKeys = (nodes: AnyRecord[]): Key[] => nodes.flatMap((node) => { const key = (node[keyField] ?? node.value) as Key; return [key, ...allKeys((node[childrenField] ?? []) as AnyRecord[])]; }); const openKeys = expandedKeys ?? (defaultExpandAll ? allKeys(treeData as AnyRecord[]) : []); const render = (nodes: AnyRecord[], depth = 0): ReactNode => nodes.map((node) => { const key = (node[keyField] ?? node.value) as Key; const children = (node[childrenField] ?? []) as AnyRecord[]; const open = openKeys.includes(key); const descendants = allKeys(children); const title = titleRender ? titleRender(node) : node[titleField] as ReactNode; return <div className="ui-tree-node" key={String(key)} style={{ paddingLeft: depth * 18 }}><div className="ui-tree-line">{children.length ? <button type="button" className="ui-tree-toggle" onClick={() => onExpand?.(open ? openKeys.filter((item) => item !== key) : [...openKeys, key])}>{open ? '−' : '+'}</button> : <span className="ui-tree-toggle" />}{checkable ? <Checkbox aria-label={`选择 ${String(node[titleField] ?? key)}`} disabled={Boolean(node.disableCheckbox ?? node.disabled)} checked={checked.includes(key)} onCheckedChange={(nextChecked) => { const affected = [key, ...descendants]; const next = nextChecked === true ? [...new Set([...checked, ...affected])] : checked.filter((item) => !affected.includes(item)); onCheck?.(next); }} /> : null}<span>{title}</span></div>{children.length && open ? render(children, depth + 1) : null}</div>; }); return <div className={`ui-tree ${className ?? ''}`}>{render(treeData as AnyRecord[])}</div>; }

export const message = { success: (content: ReactNode) => window.dispatchEvent(new CustomEvent('appica-toast', { detail: { title: String(content), type: 'success' } })), error: (content: ReactNode) => window.dispatchEvent(new CustomEvent('appica-toast', { detail: { title: String(content), type: 'error' } })), warning: (content: ReactNode) => window.dispatchEvent(new CustomEvent('appica-toast', { detail: { title: String(content), type: 'warning' } })) };
export function App({ children }: { children?: ReactNode }) { return <ToastProvider><ToastBridge />{children}<Toaster position="bottom-right" /></ToastProvider>; }
function ToastBridge() { const toast = useToastManager(); useEffect(() => { const handler = (event: Event) => { const detail = (event as CustomEvent<{ title: string; type?: string }>).detail; toast.add({ title: detail.title, type: detail.type }); }; window.addEventListener('appica-toast', handler); return () => window.removeEventListener('appica-toast', handler); }, [toast]); return null; }
App.useApp = () => ({ message });

export function ConfigProvider({ children }: { children?: ReactNode }) { return <>{children}</>; }
export const theme = { darkAlgorithm: {}, defaultAlgorithm: {} };
export type ColProps = { span?: number; xs?: number; sm?: number; md?: number; lg?: number; xl?: number; xxl?: number };
export type DataNode = TreeDataNode;
