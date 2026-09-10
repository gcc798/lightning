import { Flex, Spin, Typography } from '@/components/ui';

interface PageLoadingProps {
  tip?: string;
}

export function PageLoading({ tip = '正在加载，请稍候...' }: PageLoadingProps) {
  return (
    <Flex
      align="center"
      justify="center"
      vertical
      gap={16}
      style={{ minHeight: '50vh' }}
    >
      <Spin size="large" />
      <Typography.Text type="secondary">{tip}</Typography.Text>
    </Flex>
  );
}
