# Appica UI 重构计划

## 技术栈与决策

- React 19、TypeScript 5.9、Vite 7、React Router 6、Zustand 5、Axios。
- UI 使用 `@appica/ui-react` 1.1、`@appica/icons-react` 1.0 和 Tailwind CSS 4。
- 保留现有 `BasicTable`、`BasicForm` 和业务页面结构，通过 `src/components/ui.tsx` 统一适配 Appica 组件语义。
- 工程实际只使用了 AntD 的表单、表格、弹窗、树、选择器、布局、消息和图标等基础能力，没有依赖 ProTable、复杂图表或高级企业组件。因此完整替换合理，迁移后不保留 AntD 兼容依赖。

## 目标

- 使用 Appica UI 重构 `web-react`，最终移除 `antd` 和 `@ant-design/icons`。
- 保持现有 HTTP 契约、动态路由、权限、CRUD、主题和日志能力不变。
- 采用适合后台操作的 Appica 视觉语言：中性色、细边框、低阴影、清晰层级和稳定密度。
- 通过真实 `native` 后端完成登录、列表、编辑和权限流程闭环。

## 实施阶段

- [x] 阶段 0：将原 AntD 版本保存到远程分支 `ui/AntD`。
- [x] 阶段 1：升级 React 19，引入 Appica UI、Tailwind CSS 4 和 Appica Icons。
- [x] 阶段 2：替换应用 Provider、主题、布局、导航、提示和通用视觉样式。
- [x] 阶段 3：重构 `BasicModal`、`BasicForm`、`BasicTable`、操作菜单和树选择等公共组件。
- [x] 阶段 4：迁移登录、仪表盘、错误页及所有业务页面。
- [x] 阶段 5：删除 AntD 依赖和残留引用，完成类型检查与生产构建。
- [x] 阶段 6：启动 `native` 后端与前端，完成核心接口闭环和桌面/移动端视觉检查。
- [x] 阶段 7：按 Appica Toolbar 组合规范统一筛选区与表格工具栏，修正字段标签、动作分组和局部 loading。
- [x] 阶段 8：按 Appica 官方 ThemeToggle 模式，将顶部主题入口改为浅色/深色一键切换。

## 验收条件

- `rg "antd|@ant-design" web-react/src web-react/package.json` 无结果。
- `pnpm build` 通过。
- 登录、退出、Token 刷新、动态菜单和按钮权限正常。
- 用户、角色、菜单、组织、字典、配置和 API 权限的列表及编辑流程正常。
- 登录日志、操作日志、权限分配和批量操作正常。
- 浅色、深色主题以及桌面、移动端布局可用。

## 进度说明

重构过程中只在完成对应阶段并通过该阶段检查后勾选，避免把“已编码”误记为“已完成”。

## 验证记录

- `pnpm build` 通过，包含 TypeScript project build 和 Vite production build。
- AntD 残留扫描通过：`src`、`package.json`、`pnpm-lock.yaml` 中无 `antd`、`@ant-design` 或 `.ant-`。
- 使用 `native` Gateway `http://127.0.0.1:9009` 完成真实登录、Token、动态菜单和权限路由闭环。
- 真实接口验证了用户新增/编辑、用户列表分页、菜单树展开、菜单编辑及多选 API 权限。
- 逐页验证用户、角色、菜单、组织、字典、配置、API 权限、登录日志和操作日志列表接口，均无加载错误。
- Playwright 验证浅色/深色仪表盘、登录页、菜单页和用户列表；1440px 与 390px 页面均无横向溢出，宽表仅在表格容器内滚动。
- 浏览器控制台无运行时错误；保留 React Router 6 的 future flag 提示。
- 当前主 chunk 约 654 kB，Vite 有大于 500 kB 的提示。功能与加载均正常，后续仅在首屏性能成为实际问题时拆分 Monaco 等大依赖。
- 搜索区改为可见标签、固定紧凑字段宽度，查询/重置紧跟筛选条件；表格刷新和固定列改为带 tooltip 的图标工具按钮。
- `BasicTable` 的 loading 仅覆盖表格区域，搜索区和工具栏保持可用；`pnpm build` 与 `git diff --check` 通过。
- 主题切换移除自绘弹层与菜单 portal，使用 Appica 官方推荐的图标按钮直接切换浅色/深色；store 仍可解析已有的跟随系统配置。
