# FileNest 开发文档

## 项目概述

FileNest 是一个现代化的全栈文件管理系统，采用前后端分离架构，提供直观的文件管理界面和强大的文件操作功能。

### 技术栈

**后端**

- 语言：Go 1.22
- 框架：Gin Web Framework
- 数据库：MySQL（通过 GORM）
- 缓存：Redis
- 日志：Zap + Lumberjack

**前端**

- 框架：Vue 3 + TypeScript
- 状态管理：Pinia
- UI 组件：Naive UI
- 构建工具：Vite
- HTTP 客户端：Axios

## 项目架构

### 后端架构

```
FileNest/
├── main.go                 # 应用入口
├── cmd/                    # 应用启动相关
│   └── cmd.go
├── router/                 # 路由配置
│   └── router.go
├── internal/               # 内部包，业务逻辑核心
│   ├── controller/         # 控制器层
│   ├── service/           # 服务层接口
│   │   └── impl/          # 服务层实现
│   ├── model/             # 数据模型
│   ├── cache/             # 缓存相关
│   ├── config/            # 配置相关
│   ├── consts/            # 常量定义
│   └── utils/             # 工具类
├── common/                 # 公共组件
│   ├── glog/              # 日志组件
│   ├── database/          # 数据库组件
│   ├── middlewares/       # 中间件
│   └── response/          # 响应工具
├── upload/                # 文件上传目录
├── temp/                  # 临时文件目录
└── logs/                  # 日志目录
```

### 前端架构

```
web/
├── src/
│   ├── api/               # API 接口定义
│   │   └── file/          # 文件相关API
│   ├── components/        # 通用组件
│   │   ├── file/          # 文件相关组件
│   │   ├── home/          # 主页相关组件
│   │   └── common/        # 公共组件
│   ├── stores/            # Pinia 状态管理
│   ├── views/             # 页面视图
│   ├── router/            # 路由配置
│   ├── utils/             # 工具函数
│   ├── types/             # TypeScript 类型定义
│   ├── config/            # 配置文件
│   └── assets/            # 静态资源
├── public/                # 公共资源
└── dist/                  # 构建输出
```

## 核心功能模块

### 1. 文件管理

- 文件列表展示（网格/列表视图）
- 文件上传（支持大文件分块上传）
- 文件下载
- 文件删除
- 文件夹创建
- 文件重命名
- 文件复制/移动

### 2. 文件搜索

- 关键词搜索
- 快速文件查找

### 3. 收藏系统

- 文件/文件夹收藏
- 收藏列表管理

### 4. 统计信息

- 文件数量统计
- 存储空间统计

## API 接口设计

### 基础路径

```
BaseURL: http://localhost:9040/api
```

### 文件操作接口

#### 1. 获取文件列表

```http
GET /file/list?path={path}
```

**参数说明：**

- `path`: 目录路径（可选，默认为根目录）

**响应示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "fileName": "example.txt",
      "filePath": "/path/to/example.txt",
      "fileSize": 1024,
      "fileType": "text/plain",
      "isDir": false,
      "modTime": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### 2. 文件上传

```http
POST /file/upload
Content-Type: multipart/form-data
```

**表单参数：**

- `file`: 文件对象
- `fileName`: 文件名（可选）
- `path`: 上传路径
- `override`: 是否覆盖（true/false）

#### 3. 分块上传

```http
POST /file/upload-chunk
Content-Type: multipart/form-data
```

**表单参数：**

- `file`: 文件分块
- `fileName`: 文件名
- `path`: 上传路径
- `chunkIndex`: 分块索引
- `totalChunks`: 总分块数
- `override`: 是否覆盖

#### 4. 合并分块

```http
POST /file/merge-chunks
Content-Type: application/json
```

**请求体：**

```json
{
  "fileName": "large-file.zip",
  "path": "/upload/path",
  "totalChunks": 10,
  "override": false
}
```

#### 5. 文件下载

```http
GET /file/download?path={filePath}
```

#### 6. 删除文件

```http
DELETE /file/delete?path={filePath}&force={true|false}
```

#### 7. 创建文件夹

```http
POST /file/create-folder?path={folderPath}
```

#### 8. 文件重命名

```http
POST /file/rename?path={oldPath}&newName={newName}
```

#### 9. 文件复制

```http
POST /file/copy?srcPath={sourcePath}&destPath={destinationPath}
```

#### 10. 文件移动

```http
POST /file/move?srcPath={sourcePath}&destPath={destinationPath}
```

### 搜索接口

#### 搜索文件

```http
GET /file/search?keyword={searchKeyword}
```

### 统计接口

#### 获取文件统计

```http
GET /file/stats?path={path}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "totalFiles": 150,
    "totalFolders": 25,
    "totalSize": 1073741824
  }
}
```

### 收藏接口

#### 添加收藏

```http
POST /file/favorite?path={filePath}
```

#### 移除收藏

```http
DELETE /file/favorite?path={filePath}
```

#### 获取收藏列表

```http
GET /file/favorites
```

## 数据模型

### FileInfo 文件信息

```go
type FileInfo struct {
    FileName string `json:"fileName"` // 文件名
    FilePath string `json:"filePath"` // 文件路径
    FileSize int64  `json:"fileSize"` // 文件大小（字节）
    FileType string `json:"fileType"` // 文件类型
    IsDir    bool   `json:"isDir"`    // 是否为目录
    ModTime  string `json:"modTime"`  // 修改时间
}
```

### FileStats 文件统计

```go
type FileStats struct {
    TotalFiles   int64 `json:"totalFiles"`   // 文件总数
    TotalFolders int64 `json:"totalFolders"` // 文件夹总数
    TotalSize    int64 `json:"totalSize"`    // 总大小（字节）
}
```

### Favorite 收藏信息

```go
type Favorite struct {
    ID         int64     `json:"id"`         // 收藏ID
    Name       string    `json:"name"`       // 文件名
    Path       string    `json:"path"`       // 文件路径
    IsDir      bool      `json:"isDir"`      // 是否是目录
    CreateTime time.Time `json:"createTime"` // 创建时间
}
```

## 代码规范

### Go 代码规范

#### 1. 包命名

- 使用小写字母，简短且有意义
- 避免使用下划线或驼峰命名

#### 2. 文件命名

- 使用下划线分隔单词
- 例如：`file_controller.go`、`upload_consts.go`

#### 3. 函数命名

- 公共函数使用大驼峰命名（PascalCase）
- 私有函数使用小驼峰命名（camelCase）
- 接口方法名要清晰表达功能

#### 4. 错误处理

```go
// 标准错误处理模式
if err != nil {
    glog.Errorf("操作失败: %s", err)
    response.Error(ctx, err.Error())
    return
}
```

#### 5. 日志记录

```go
// 使用统一的日志格式
glog.Infof("收到请求，参数: %s", param)
glog.Errorf("操作失败: %s", err)
glog.Warnf("警告信息: %s", warning)
```

#### 6. 注释规范

```go
/**
  @author: XingGao
  @date: 2024/9/22
**/

// GetFileList 获取文件列表
// 参数 path: 目录路径
// 返回文件信息列表和错误信息
func (h *FileController) GetFileList(ctx *gin.Context) {
    // 实现代码...
}
```

### TypeScript/Vue 代码规范

#### 1. 文件命名

- 组件文件使用大驼峰命名：`FileList.vue`、`HomeView.vue`
- 普通文件使用小驼峰命名：`fileStore.ts`、`request.ts`

#### 2. 变量命名

- 使用小驼峰命名（camelCase）
- 常量使用大写下划线命名（SCREAMING_SNAKE_CASE）

#### 3. 组件结构

```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
// 导入
import { ref, computed } from "vue";

// 类型定义
interface Props {
  // 属性定义
}

// 响应式数据
const loading = ref(false);

// 计算属性
const computedValue = computed(() => {
  // 计算逻辑
});

// 方法
const handleClick = () => {
  // 处理逻辑
};
</script>

<style scoped>
/* 样式 */
</style>
```

#### 4. Pinia Store 结构

```typescript
export const useFileStore = defineStore("file", () => {
  // 状态
  const state = ref(initialState);

  // 计算属性
  const computed = computed(() => {
    // 计算逻辑
  });

  // 方法
  const action = async () => {
    // 异步操作
  };

  return {
    state,
    computed,
    action,
  };
});
```

## 开发环境设置

### 后端环境

#### 1. 环境要求

- Go 1.22+
- Redis（用于缓存）
- MySQL（用于数据库，可选）

#### 2. 环境变量

```bash
# 端口配置
PORT=9040

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0

# 文件存储路径
UPLOAD_DIR=./upload
TEMP_DIR=./temp
```

#### 3. 运行项目

```bash
# 安装依赖
go mod download

# 运行项目
go run main.go

# 或指定端口
go run main.go -port=9040
```

### 前端环境

#### 1. 环境要求

- Node.js 16+
- pnpm 7+

#### 2. 安装依赖

```bash
cd web
pnpm install
```

#### 3. 开发服务器

```bash
pnpm dev
```

#### 4. 构建生产版本

```bash
pnpm build
```

### Docker 部署

#### 1. 构建镜像

```bash
docker build -t filenest .
```

#### 2. Docker Compose 部署

```bash
docker-compose up -d
```

## 项目配置

### 常量配置

```go
// internal/consts/upload_consts.go
const (
    DefaultPageSize = 10        // 默认分页大小
    MaxPageSize     = 100       // 最大分页大小
    TempDir         = "./temp"  // 临时文件目录
    UploadDir       = "./upload" // 上传文件目录
)
```

### 中间件配置

- CORS 跨域处理
- Zap 日志记录
- 错误恢复处理

### 前端路由配置

```typescript
// 路由配置示例
const routes = [
  {
    path: "/",
    name: "Home",
    component: HomeView,
  },
  {
    path: "/about",
    name: "About",
    component: AboutView,
  },
];
```

## 状态管理

### File Store (useFileStore)

#### 状态

- `currentPath`: 当前目录路径
- `files`: 文件列表
- `isLoading`: 加载状态
- `viewMode`: 视图模式（网格/列表）
- `favorites`: 收藏列表

#### 主要方法

- `fetchFiles()`: 获取文件列表
- `enterDirectory(path)`: 进入目录
- `createNewFolder(name)`: 创建文件夹
- `downloadFile(path)`: 下载文件
- `deleteFile(path, force)`: 删除文件
- `searchFile(keyword)`: 搜索文件
- `addToFavorites(path)`: 添加收藏
- `removeFromFavorites(path)`: 移除收藏

## 组件说明

### 核心组件

#### 1. HomeView.vue

- 主页面组件
- 集成文件列表、操作栏、统计信息
- 处理文件上传、创建文件夹等操作

#### 2. FileList.vue

- 文件列表展示组件
- 支持网格和列表两种视图模式
- 包含文件操作菜单

#### 3. UploadConfig.vue

- 文件上传配置组件
- 支持拖拽上传
- 进度显示和错误处理

### 组件通信

- 使用 Pinia 进行状态管理
- 组件间通过 emit/props 进行通信
- 使用 provide/inject 传递上下文数据

## 错误处理

### 后端错误处理

```go
// 标准错误响应格式
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data"`
}

// 错误处理示例
if err != nil {
    glog.Errorf("操作失败: %s", err)
    response.Error(ctx, err.Error())
    return
}
```

### 前端错误处理

```typescript
// API 调用错误处理
try {
  const result = await apiCall();
  // 处理成功响应
} catch (error) {
  message.error(error instanceof Error ? error.message : "操作失败");
  console.error("API调用失败:", error);
}
```

## 缓存策略

### Redis 缓存

- 文件列表缓存
- 搜索结果缓存
- 用户会话缓存

### 缓存键命名规范

```
file_list:{path}        # 文件列表缓存
search_result:{keyword} # 搜索结果缓存
user_session:{userid}   # 用户会话缓存
```

## 日志管理

### 日志配置

- 使用 Zap 作为日志库
- 日志轮转通过 Lumberjack 实现
- 日志级别：Debug、Info、Warn、Error

### 日志格式

```
[时间] [级别] [文件:行号] 日志内容
```

### 日志使用示例

```go
glog.Infof("用户操作: %s, 参数: %+v", operation, params)
glog.Errorf("操作失败: %s, 错误: %v", operation, err)
```

## 测试指南

### 单元测试

```bash
# 运行后端测试
go test ./...

# 运行前端测试
cd web
pnpm test
```

### API 测试

推荐使用 Postman 或类似工具测试 API 接口

### 测试数据

- 准备测试文件和文件夹
- 模拟不同文件类型和大小
- 测试边界条件

## 性能优化

### 后端优化

- 文件列表分页加载
- Redis 缓存热点数据
- 大文件分块上传
- 异步文件操作

### 前端优化

- 虚拟滚动处理大量文件
- 图片懒加载
- 组件按需加载
- 状态管理优化

## 安全考虑

### 文件上传安全

- 文件类型检查
- 文件大小限制
- 路径遍历防护
- 恶意文件扫描

### API 安全

- 参数验证
- 路径权限检查
- 错误信息脱敏

## 监控与运维

### 健康检查

```go
// 健康检查接口
func (h *HealthController) Check(ctx *gin.Context) {
    ctx.JSON(200, gin.H{
        "status": "ok",
        "timestamp": time.Now(),
    })
}
```

### 监控指标

- API 响应时间
- 文件操作成功率
- 存储空间使用率
- 错误日志统计

## 扩展开发

### 添加新的 API 接口

1. **定义模型**（如需要）

```go
// internal/model/new_model.go
type NewModel struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}
```

2. **定义服务接口**

```go
// internal/service/new_service.go
type NewService interface {
    GetData() (*NewModel, error)
}
```

3. **实现服务**

```go
// internal/service/impl/new_service_impl.go
func (s *NewServiceImpl) GetData() (*NewModel, error) {
    // 实现业务逻辑
}
```

4. **添加控制器**

```go
// internal/controller/new_controller.go
func (h *NewController) GetData(ctx *gin.Context) {
    data, err := h.service.GetData()
    if err != nil {
        response.Error(ctx, err.Error())
        return
    }
    response.Success(ctx, data)
}
```

5. **注册路由**

```go
// router/router.go
newController := controller.NewController(service.NewService())
api.GET("/new/data", newController.GetData)
```

### 添加新的前端功能

1. **定义 API 接口**

```typescript
// src/api/new/new.ts
export const getNewData = () => {
  return request.get("/new/data");
};
```

2. **创建 Store**

```typescript
// src/stores/new.ts
export const useNewStore = defineStore("new", () => {
  const data = ref(null);

  const fetchData = async () => {
    const result = await getNewData();
    data.value = result.data;
  };

  return { data, fetchData };
});
```

3. **创建组件**

```vue
<!-- src/components/new/NewComponent.vue -->
<template>
  <div>
    <!-- 组件内容 -->
  </div>
</template>

<script setup lang="ts">
import { useNewStore } from "@/stores/new";

const newStore = useNewStore();
</script>
```

## 故障排除

### 常见问题

#### 1. 文件上传失败

- 检查上传目录权限
- 确认文件大小限制
- 查看磁盘空间

#### 2. API 调用超时

- 检查网络连接
- 确认服务器状态
- 查看日志文件

#### 3. 前端编译错误

- 检查依赖版本
- 清除 node_modules 重新安装
- 检查 TypeScript 配置

### 调试技巧

#### 后端调试

```go
// 添加调试日志
glog.Debugf("调试信息: %+v", debugData)

// 使用 IDE 断点调试
// 配置 launch.json 文件
```

#### 前端调试

```typescript
// 浏览器控制台调试
console.log("调试信息:", debugData);

// Vue DevTools 使用
// 安装 Vue DevTools 浏览器扩展
```

## 版本更新

### 更新流程

1. 备份数据库（如使用）
2. 备份上传文件
3. 更新代码
4. 重启服务
5. 验证功能

### 数据迁移

```go
// 数据迁移示例
func migrateData() error {
    // 执行数据迁移逻辑
    return nil
}
```

## 贡献指南

### 开发流程

1. Fork 项目
2. 创建功能分支
3. 提交代码更改
4. 编写测试
5. 提交 Pull Request

### 代码审查

- 遵循代码规范
- 添加必要的注释
- 确保测试通过
- 更新相关文档

---

## 结语

本文档涵盖了 FileNest 项目的完整开发指南，包括架构设计、API 接口、代码规范、开发环境等各个方面。开发者可以根据此文档快速上手项目开发，并按照规范进行功能扩展。

如有疑问或建议，请通过 Issue 或 Pull Request 与我们联系。
