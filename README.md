# FileNest

FileNest 是一个现代化的文件管理系统，提供直观的用户界面和强大的文件管理功能。

## 功能特点

- 📁 文件浏览：支持列表视图和网格视图
- 🔍 文件搜索：快速查找文件和文件夹
- 📊 文件统计：显示文件夹大小、文件数量等信息
- ⭐ 收藏夹：支持将常用文件和文件夹添加到收藏夹
- 📥 文件上传：支持拖拽上传和进度显示
- 📤 文件下载：支持单文件下载
- 📂 文件管理：创建文件夹、删除文件等基本操作

## 技术栈

- 前端框架：Vue 3
- 状态管理：Pinia
- UI 组件：Naive UI
- HTTP 客户端：Axios
- 构建工具：Vite
- 开发语言：TypeScript

## 项目结构

```
web/
├── src/
│   ├── api/            # API 接口定义
│   ├── assets/         # 静态资源
│   ├── components/     # 通用组件
│   │   ├── file/      # 文件相关组件
│   │   └── home/      # 主页相关组件
│   ├── stores/        # Pinia 状态管理
│   ├── utils/         # 工具函数
│   └── views/         # 页面视图
├── public/            # 公共资源
└── package.json       # 项目配置
```

## 主要组件

- `HomeView`: 主页视图，包含文件列表和操作界面
- `FileItem`: 文件项组件，支持列表和网格两种显示模式
- `HomeHeader`: 顶部操作栏，包含搜索、视图切换等功能
- `FileStats`: 文件统计信息显示

## 开发指南

### 环境要求

- Node.js >= 16
- pnpm >= 7

### 安装依赖

```bash
pnpm install
```

### 开发服务器

```bash
pnpm dev
```

### 构建生产版本

```bash
pnpm build
```

## API 接口

### 文件操作

- `GET /api/file/list` - 获取文件列表
- `POST /api/file/create-folder` - 创建文件夹
- `DELETE /api/file/delete` - 删除文件
- `GET /api/file/download` - 下载文件
- `POST /api/file/upload` - 上传文件
- `GET /api/file/stats` - 获取文件统计信息
- `GET /api/file/search` - 搜索文件

### 收藏夹操作

- `GET /api/file/favorites` - 获取收藏列表
- `POST /api/file/favorite` - 添加收藏
- `DELETE /api/file/favorite` - 移除收藏

## 状态管理

使用 Pinia 进行状态管理，主要包含：

- `fileStore`: 处理文件列表、当前路径等状态
- 支持文件操作的异步动作
- 提供计算属性用于文件排序和过滤

## 网络请求

使用封���的 `request` 工具进行 API 调用：

- 统一的错误处理
- 请求/响应拦截
- 支持 TypeScript 类型
- 专门的文件下载处理

## 贡献指南

1. Fork 本仓库
2. 创建特性分支：`git checkout -b feature/xxx`
3. 提交更改：`git commit -m 'Add xxx'`
4. 推送分支：`git push origin feature/xxx`
5. 提交 Pull Request

## 许可证

[MIT License](LICENSE) 