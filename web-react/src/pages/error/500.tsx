import { Button, Result } from '@/components/ui';
import { Link } from 'react-router-dom';

export default function ServerErrorPage() {
  return (
    <Result
      status="500"
      title="500"
      subTitle="服务端发生错误，请稍后再试。"
      extra={
        <Button type="primary">
          <Link to="/">返回首页</Link>
        </Button>
      }
    />
  );
}
