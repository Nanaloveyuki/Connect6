Connect6 项目计划书

1. 项目概述

1.1 项目名称

Connect6 —— 基于六度分隔理论的 GitHub 关系拓扑可视化工具。

1.2 项目愿景

构建一个去中心化的关系网络可视化工具，通过 GitHub 开源协作数据，展示任意用户之间的连接路径。旨在解决社交场景中“陌生人突兀感”的问题，让每个人都能直观地看到“我们是如何被开源世界联系在一起的”。

1.3 核心理念

六度分隔理论（Six Degrees of Separation）：世界上任意两个人之间，平均只需通过 6 个中间人即可建立联系。本项目致力于在 GitHub 生态中验证并可视化这一过程。

2. 核心功能需求

2.1 用户关系图谱可视化

• 以力导向图（Force-Directed Graph）形式展示用户关系网络。

• 节点（Node）：代表 GitHub 用户。

• 边（Edge）：代表协作关系（如：共同贡献项目、同属一个 Organization、互相 Follow 等）。

• 支持缩放、拖拽、Hover 查看详情、点击聚焦等交互。

2.2 路径溯源（Path Tracing）

• 输入两个 GitHub 用户名（Source & Target）。

• 计算并高亮显示两者之间的最短连接路径。

• 展示路径上下文（例如：“A → 通过 Project X → B”）。

2.3 GitHub 数据集成

• 通过 GitHub API 拉取用户基础信息（Profile, Avatar）。

• 拉取用户的社交关系（Followers/Following）。

• 拉取用户的仓库与组织信息，构建协作边。

2.4 轻量级身份卡片

• 生成用户的“关系名片”，可嵌入到其他平台（如 Discord, Telegram, 论坛签名）。

• 名片展示该用户与当前浏览者的连接度（例如：“与您相距 3 度”）。

3. 技术栈选型

3.1 整体架构

采用前后端分离架构，通过 RESTful API 进行通信。

+-------------------+      HTTPS      +-------------------+
|                   |  <------------>  |                   |
|  Svelte Frontend  |                  |  Go Backend (Gin) |
|                   |                  |                   |
+-------------------+                  +-------------------+
                                                     |
                                                     | Internal API Call
                                                     v
                                              +-------------------+
                                              |    GitHub API     |
                                              +-------------------+


3.2 前端技术栈 (Client-Side)

层级 技术选型 说明

框架 Svelte / SvelteKit 极简、高性能，适合工具类应用，无虚拟 DOM 开销。

UI 库 Skeleton / Flowbite Svelte 提供基础组件，加速开发。

图表库 Sigma.js 基于 Canvas 的高性能图渲染引擎，专为关系网络设计。

样式 Tailwind CSS 原子化 CSS，快速构建界面。

状态管理 Svelte Stores Svelte 内置状态管理，无需额外引入 Redux 等库。
3.3 后端技术栈 (Server-Side)
层级 技术选型 说明

语言 Go (Golang) 高性能、高并发，适合处理大量 GitHub API 请求。

Web 框架 Gin 轻量级、高性能的 Go Web 框架，API 开发效率高。

HTTP Client Resty Go 生态中易用的 HTTP 客户端库，用于调用 GitHub API。

路由 Gin Router 负责 API 路由分发。

中间件 Gin 自带中间件 用于 CORS、Logging、Recovery。

配置管理 Viper 处理环境变量、配置文件。

并发控制 Go Routines / Channels 核心优势，用于并发抓取 GitHub 数据。

缓存策略 In-Memory Cache (初期) 使用 sync.Map 或 ristretto 缓存 GitHub API 响应，减少 Rate Limit 消耗。

4. 开发计划 (Roadmap)

Phase 1: MVP 核心链路（预计 2-3 周）
后端：搭建 Gin 基础框架，实现 GitHub API 代理接口。

后端：实现用户基础信息获取接口 (/api/users/:username)。

后端：实现简单的 BFS 算法，计算两个用户间的最短路径。

前端：搭建 SvelteKit 项目，实现基础布局。

前端：集成 Sigma.js，能够接收后端返回的节点和边数据并渲染。

联调：完成“输入两个用户名 -> 显示连接路径”的最小闭环。

Phase 2: 体验优化与交互（预计 2 周）
前端：完善力导向图的物理参数（斥力、引力），优化视觉效果。

前端：增加节点 Hover/Click 交互，展示用户 Avatar 和简介。

后端：引入缓存机制，避免重复请求相同用户数据。

后端：增加 GitHub Token 池管理，突破单一 Token 的 Rate Limit。

Phase 3: 去中心化与分享（预计 1-2 周）
后端：实现生成分享链接的功能（将路径参数编码进 URL）。

前端：实现“嵌入模式”（Embed Mode），仅展示图谱，无导航栏。

文档：编写 README 和 API 文档。

5. 风险与应对

风险点 应对策略

GitHub API Rate Limit 后端统一代理，引入 Redis/内存缓存，支持配置多个 GitHub Token 轮询。

冷启动数据为空 初期聚焦于知名开源项目维护者，利用其丰富的连接关系作为种子数据。

前端渲染性能 限制初始渲染节点数量（如仅展示 2 度以内邻居），采用 Canvas 渲染而非 SVG。

隐私问题 默认只展示用户公开的 GitHub 数据，不存储敏感信息，提供数据清除接口。

6. 结语

Connect6 不是一个传统的社交应用，而是一个关系显微镜。通过 Go 的高并发能力和 Svelte 的极致前端性能，我们将构建一个既符合黑客精神（Hacker Ethic）又具有实用价值的开源工具。


## 前端动画
Svelte
 ├── 输入框（用户名）
 ├── 按钮（计算路径）
 ├── 侧边栏（用户信息）
 └── 控制面板（布局参数）
        ↓
  事件 / JSON
        ↓
Sigma.js (Canvas)
 ├── 节点绘制
 ├── 边绘制
 ├── 力导向布局
 ├── 拖拽 / 缩放
 └── 动画