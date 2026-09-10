import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { App, Button, Card, Form, Input } from '@/components/ui';
import { LockOutlined, ReloadOutlined, UserOutlined } from '@/utils/icons';
import { authApi } from '@/api/auth';
import { ThemeSwitcher } from '@/components/common/ThemeSwitcher';
import { DEFAULT_CLIENT_ID } from '@/constants/auth';
import { useAuthStore } from '@/store/auth';

interface LoginFormValues {
  username: string;
  password: string;
  code?: string;
}

export default function LoginPage() {
  const [form] = Form.useForm<LoginFormValues>();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { message } = App.useApp();
  const isLoggedIn = useAuthStore((state) => state.isLoggedIn);
  const login = useAuthStore((state) => state.login);
  const [loading, setLoading] = useState(false);
  const [showCaptcha, setShowCaptcha] = useState(false);
  const [captchaId, setCaptchaId] = useState('');
  const [captchaImage, setCaptchaImage] = useState('');

  const loadCaptcha = async () => {
    try {
      const data = await authApi.getImageCaptcha();
      setCaptchaId(data.id);
      setCaptchaImage(String(data.data.image ?? ''));
    } catch {
      message.error('加载验证码失败');
    }
  };

  useEffect(() => {
    if (isLoggedIn) {
      navigate('/', { replace: true });
    }
  }, [isLoggedIn, navigate]);

  useEffect(() => {
    void (async () => {
      try {
        const types = await authApi.getCaptchaTypes();
        const enabled = types.includes('image');
        setShowCaptcha(enabled);

        if (enabled) {
          await loadCaptcha();
        }
      } catch {
        setShowCaptcha(false);
      }
    })();
  }, []);

  const handleFinish = async (values: LoginFormValues) => {
    setLoading(true);
    try {
      await login({
        clientId: DEFAULT_CLIENT_ID,
        grantType: 'password',
        username: values.username,
        password: values.password,
        code: values.code,
        uuid: captchaId || undefined,
      });

      message.success('登录成功');
      // 登录后优先回到用户原本想访问的页面，和旧后台行为保持一致。
      const redirect = searchParams.get('redirect') || '/';
      navigate(redirect, { replace: true });
    } catch (error) {
      message.error(error instanceof Error ? error.message : '登录失败');
      if (showCaptcha) {
        form.setFieldValue('code', '');
        await loadCaptcha();
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-theme-dock">
        <ThemeSwitcher
          buttonClassName="login-theme-switch"
          entryClassName="login-theme-entry"
        />
      </div>

      <div className="login-shell">
        <section className="login-brand-panel">
          <div className="login-brand-lockup">
            <span className="app-logo-text login-brand-text">
              <strong>Admin</strong>
            </span>
          </div>
        </section>

        <Card className="login-card" variant="borderless">
          <Form<LoginFormValues>
            className="login-form"
            form={form}
            layout="vertical"
            onFinish={(values) => void handleFinish(values)}
          >
            <Form.Item
              label="用户名"
              name="username"
              rules={[{ required: true, message: '请输入用户名' }]}
            >
              <Input
                allowClear
                autoComplete="username"
                prefix={
                  <span className="login-input-icon">
                    <UserOutlined />
                  </span>
                }
                placeholder="请输入用户名"
                size="large"
              />
            </Form.Item>

            <Form.Item
              label="密码"
              name="password"
              rules={[{ required: true, message: '请输入密码' }]}
            >
              <Input.Password
                allowClear
                autoComplete="current-password"
                prefix={
                  <span className="login-input-icon">
                    <LockOutlined />
                  </span>
                }
                placeholder="请输入密码"
                size="large"
              />
            </Form.Item>

            {showCaptcha ? (
              <Form.Item
                label="验证码"
                name="code"
                rules={[{ required: true, message: '请输入验证码' }]}
              >
                <div className="login-captcha">
                  <Input placeholder="请输入验证码" size="large" />
                  <div
                    className="login-captcha-image"
                    onClick={() => void loadCaptcha()}
                  >
                    {captchaImage ? (
                      <img alt="验证码" src={captchaImage} style={{ width: '100%' }} />
                    ) : (
                      <ReloadOutlined />
                    )}
                  </div>
                </div>
              </Form.Item>
            ) : null}

            <Form.Item style={{ marginBottom: 0 }}>
              <Button
                block
                className="login-submit"
                htmlType="submit"
                loading={loading}
                size="large"
                type="primary"
              >
                登录
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
    </div>
  );
}
