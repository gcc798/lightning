import { useEffect } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { App as AppicaApp } from '@/components/ui';
import { AppRoutes } from '@/router/AppRoutes';
import { useResolvedThemeMode, useThemeStore } from '@/store/theme';

export function App() {
  const themeMode = useThemeStore((state) => state.mode);
  const resolvedThemeMode = useResolvedThemeMode(themeMode);
  useEffect(() => {
    document.documentElement.dataset.theme = resolvedThemeMode;
    document.documentElement.dataset.themePreference = themeMode;
    document.documentElement.classList.toggle('dark', resolvedThemeMode === 'dark');
    document.documentElement.classList.toggle('light', resolvedThemeMode === 'light');
  }, [resolvedThemeMode, themeMode]);

  return (
    <AppicaApp>
      <div className="app-root">
        <BrowserRouter future={{ v7_startTransition: true }}>
          <AppRoutes />
        </BrowserRouter>
      </div>
    </AppicaApp>
  );
}
