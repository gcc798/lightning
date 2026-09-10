import { MoonOutlined, SunOutlined } from '@/utils/icons';
import { Button } from '@/components/ui';
import { useResolvedThemeMode, useThemeStore } from '@/store/theme';

interface ThemeSwitcherProps {
  entryClassName?: string;
  buttonClassName?: string;
}

export function ThemeSwitcher({
  entryClassName = 'header-theme-entry',
  buttonClassName = 'header-theme-switch',
}: ThemeSwitcherProps) {
  const mode = useThemeStore((state) => state.mode);
  const setTheme = useThemeStore((state) => state.setTheme);
  const resolvedMode = useResolvedThemeMode(mode);
  const nextMode = resolvedMode === 'dark' ? 'light' : 'dark';
  const label = nextMode === 'dark' ? '切换到深色模式' : '切换到浅色模式';

  return (
    <div className={entryClassName}>
      <Button
        aria-label={label}
        className={buttonClassName}
        icon={resolvedMode === 'dark' ? <MoonOutlined /> : <SunOutlined />}
        title={label}
        type="text"
        onClick={() => setTheme(nextMode)}
      />
    </div>
  );
}
