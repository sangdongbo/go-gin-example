# OpenSpec 工作流使用指南

## 关于本文档

本文档是 OpenSpec 工作流的详细中文使用指南，专为本项目团队成员编写。如果您是新加入的开发者或需要了解如何使用 OpenSpec 进行结构化的变更管理，这份文档将帮助您快速掌握完整的工作流程。

**与其他文档的关系**:
- **README.md**: 项目整体介绍（English）
- **README1.md**: 双语学习指南
- **README2.md**: 深度中文学习指南
- **openspec.md**: OpenSpec 工作流使用规范（本文档）

本文档专注于 OpenSpec 工作流本身，不涉及项目的业务逻辑和技术实现细节。

---

## 目录

1. [OpenSpec 简介](#openspec-简介)
2. [核心概念](#核心概念)
3. [完整工作流程](#完整工作流程)
4. [命令详解](#命令详解)
5. [工件详解](#工件详解)
6. [目录结构](#目录结构)
7. [实际使用示例](#实际使用示例)
8. [工作流场景](#工作流场景)
9. [最佳实践](#最佳实践)
10. [常见问题解答](#常见问题解答)

---

## OpenSpec 简介

**OpenSpec** 是一个结构化的变更管理工作流系统，帮助团队以规范、可追溯的方式管理项目变更。它通过定义清晰的工作流程和标准化的文档格式，确保每个变更从提案、设计、实施到归档都有完整的记录。

### 为什么使用 OpenSpec？

- **结构化管理**: 每个变更都遵循统一的流程，包含提案、设计文档、需求规范和任务清单
- **可追溯性**: 所有变更决策和实施细节都有文档记录，便于回顾和审计
- **团队协作**: 标准化的格式使团队成员能够快速理解变更的目的、设计和进度
- **质量保证**: 通过需求规范和验收场景，确保变更满足预期要求
- **知识沉淀**: 归档的变更成为项目的知识库，帮助新成员了解项目演进历史

### 本项目的使用情况

本项目已经成功使用 OpenSpec 管理了多个变更，包括：
- **add-readme1-file**: 添加双语学习指南
- **add-readme2-file**: 添加深度中文学习指南

这些变更都已完成并归档在 `openspec/changes/archive/` 目录下，可作为参考示例。

---

## 核心概念

### 工作流类型

OpenSpec 支持多种工作流模式，本项目使用 **spec-driven（规范驱动）** 模式。

#### spec-driven 工作流

规范驱动工作流强调"先定义需求，再实施开发"的理念。它包含以下阶段：

1. **Propose（提案）**: 说明变更的动机和目标
2. **Design（设计）**: 制定技术方案和实施策略
3. **Specify（规范）**: 定义详细的需求和验收标准
4. **Task（任务）**: 分解为可执行的具体任务
5. **Apply（应用）**: 实施任务并标记进度
6. **Archive（归档）**: 完成后归档变更并同步规范

### 工件类型

spec-driven 工作流包含四种核心工件（Artifacts）：

#### 1. Proposal（提案）

**文件**: `proposal.md`

**作用**: 说明"为什么（Why）"需要这个变更

**内容**:
- **Why（为什么）**: 问题描述或机会说明，1-2 句话说明动机
- **What Changes（变更内容）**: 具体的变更列表，说明新增、修改或删除的功能
- **Capabilities（能力）**: 列出将要创建或修改的规范能力
  - **New Capabilities（新能力）**: 新引入的功能能力，每个对应一个新的规范文件
  - **Modified Capabilities（修改的能力）**: 现有能力的需求变更
- **Impact（影响）**: 受影响的代码、API、依赖或系统

#### 2. Design（设计）

**文件**: `design.md`

**作用**: 说明"如何（How）"实现这个变更

**内容**:
- **Context（背景）**: 当前状态、限制条件、相关方
- **Goals / Non-Goals（目标与非目标）**: 明确范围
- **Decisions（决策）**: 关键技术选择及理由，包含备选方案对比
- **Risks / Trade-offs（风险与权衡）**: 已知限制和可能的问题
- **Migration Plan（迁移计划）**: 部署步骤和回滚策略（如适用）
- **Open Questions（待定问题）**: 需要解决的未决事项

#### 3. Specs（规范）

**文件**: `specs/<capability>/spec.md`

**作用**: 定义"做什么（What）"，即详细的需求和验收标准

**内容**:
- **Requirements（需求）**: 系统必须满足的功能需求，每个需求包含：
  - 需求名称和描述（使用 SHALL/MUST 等规范性语言）
  - 一个或多个验收场景（Scenarios）
- **Scenarios（场景）**: 验收标准，采用 WHEN/THEN 格式
  - **WHEN（当...时）**: 触发条件
  - **THEN（那么...）**: 期望结果

**格式要求**:
- Requirements 使用 `### Requirement:` 标题（三个 #）
- Scenarios 使用 `#### Scenario:` 标题（四个 #，非常重要！）
- 每个需求至少包含一个场景

**Delta 操作**（针对修改现有规范）:
- `## ADDED Requirements`: 新增需求
- `## MODIFIED Requirements`: 修改需求（必须包含完整的更新内容）
- `## REMOVED Requirements`: 删除需求（必须包含原因和迁移说明）
- `## RENAMED Requirements`: 重命名需求（使用 FROM:/TO: 格式）

#### 4. Tasks（任务）

**文件**: `tasks.md`

**作用**: 将实施工作分解为可追踪的具体任务

**内容**:
- 按功能模块分组的任务列表
- 每个任务使用复选框格式：`- [ ] X.Y 任务描述`
- 任务编号格式：组号.任务号（如 1.1, 1.2, 2.1）
- 任务应足够小，能在一个工作会话内完成
- 任务按依赖顺序排列

**示例**:
```markdown
## 1. 环境准备

- [ ] 1.1 安装项目依赖
- [ ] 1.2 配置开发环境

## 2. 核心功能实现

- [ ] 2.1 实现数据模型
- [ ] 2.2 实现业务逻辑
```

### 工件之间的依赖关系

工件之间存在明确的依赖关系，必须按顺序创建：

```
proposal（提案）
    ↓
    ├─→ design（设计）
    │       ↓
    │       └─→ tasks（任务）
    │
    └─→ specs（规范）
            ↓
            └─→ tasks（任务）
```

**依赖说明**:
1. **proposal** 是起点，其他所有工件都依赖它
2. **design** 和 **specs** 都依赖 proposal，可以并行创建
3. **tasks** 依赖 design 和 specs，必须最后创建
4. 只有 tasks 完成后，才能进入 apply（实施）阶段

### 工作流的生命周期

一个完整的 OpenSpec 变更经历以下生命周期：

1. **创建（Create）**: 使用 `openspec new change <name>` 创建变更目录
2. **提案（Propose）**: 创建 proposal、design、specs、tasks 工件
3. **准备（Ready）**: 所有必需工件创建完成，状态变为 ready
4. **实施（Apply）**: 按照 tasks 逐一实施变更，标记完成进度
5. **完成（Complete）**: 所有任务完成
6. **归档（Archive）**: 将变更移动到归档目录，同步规范到主目录

---

## 完整工作流程

### 工作流程图

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenSpec 工作流程                          │
└─────────────────────────────────────────────────────────────┘

   开始
    │
    ▼
┌────────────────┐
│ 1. 创建变更      │  openspec new change <name>
│   (Create)     │  或使用 /opsx:propose
└────────────────┘
    │
    ▼
┌────────────────┐
│ 2. 编写提案      │  创建 proposal.md
│   (Propose)    │  - 说明 Why（为什么）
└────────────────┘  - 列出 What（变更内容）
    │               - 定义 Capabilities（能力）
    │
    ├──────────────┐
    │              │
    ▼              ▼
┌────────────┐  ┌────────────┐
│ 3a. 设计    │  │ 3b. 规范   │
│  (Design)  │  │  (Specs)   │
│ design.md  │  │ spec.md    │
└────────────┘  └────────────┘
    │              │
    │              │
    └──────┬───────┘
           │
           ▼
    ┌────────────┐
    │ 4. 任务清单 │  创建 tasks.md
    │   (Tasks)  │  - 分解为具体任务
    └────────────┘  - 建立依赖顺序
           │
           ▼
    ┌────────────┐
    │ 5. 实施     │  /opsx:apply
    │   (Apply)  │  - 逐一完成任务
    └────────────┘  - 标记 [x] 进度
           │
           ▼
    ┌────────────┐
    │ 6. 归档     │  /opsx:archive
    │  (Archive) │  - 移动到 archive/
    └────────────┘  - 同步规范到 specs/
           │
           ▼
         完成
```

### 阶段详解

#### 阶段 1: 创建变更（Create）

**目的**: 初始化变更目录和配置文件

**操作**:
```bash
openspec new change "变更名称"
```

**产出**:
- 在 `openspec/changes/<name>/` 创建目录
- 生成 `.openspec.yaml` 配置文件

**说明**:
- 变更名称使用 kebab-case 格式（小写字母，单词间用连字符）
- 例如: `add-user-auth`, `fix-login-bug`, `refactor-api`

#### 阶段 2: 编写提案（Propose）

**目的**: 说明变更的必要性和目标

**操作**:
- 手动创建 `proposal.md`
- 或使用 `/opsx:propose` 命令自动生成所有工件

**产出**:
- `proposal.md` 文件
- 明确的变更动机和范围

**关键点**:
- 在 Capabilities 部分明确列出需要创建的规范
- 保持简洁（1-2 页），专注于"为什么"而非"如何"

#### 阶段 3a: 编写设计（Design）

**目的**: 制定技术实施方案

**操作**:
- 创建 `design.md`
- 参考 proposal 了解需求背景

**产出**:
- 技术决策文档
- 实施方案和风险评估

**关键点**:
- 解释关键技术选择的理由
- 列出备选方案和选择依据
- 识别潜在风险和缓解措施

#### 阶段 3b: 编写规范（Specs）

**目的**: 定义详细的需求和验收标准

**操作**:
- 为每个 capability 创建 `specs/<capability>/spec.md`
- 使用 Requirement 和 Scenario 格式

**产出**:
- 可测试的需求规范
- 明确的验收标准

**关键点**:
- Scenarios 必须使用四个 # （`#### Scenario:`）
- 每个 Requirement 至少一个 Scenario
- 使用 SHALL/MUST 等规范性语言

#### 阶段 4: 编写任务清单（Tasks）

**目的**: 将工作分解为可执行的任务

**操作**:
- 创建 `tasks.md`
- 参考 design 和 specs 确定实施步骤

**产出**:
- 分组的任务清单
- 明确的实施顺序

**关键点**:
- 使用 `- [ ]` 复选框格式
- 任务编号：组号.任务号
- 任务足够小，可在一次完成

#### 阶段 5: 实施（Apply）

**目的**: 按照任务清单实施变更

**操作**:
```bash
/opsx:apply
```

**活动**:
- 逐一完成任务
- 修改 tasks.md，将 `- [ ]` 改为 `- [x]`
- 提交代码变更

**监控**:
```bash
openspec status --change "<name>"
```

**关键点**:
- 保持每个任务的变更最小化
- 完成一个任务立即标记
- 遇到问题及时暂停和沟通

#### 阶段 6: 归档（Archive）

**目的**: 完成变更并整理文档

**操作**:
```bash
/opsx:archive
```

**活动**:
1. 检查所有任务是否完成
2. 将变更目录移动到 `openspec/changes/archive/YYYY-MM-DD-<name>/`
3. 同步 delta specs 到 `openspec/specs/`

**产出**:
- 归档的完整变更记录
- 更新的主规范目录

**关键点**:
- 归档前确保所有任务完成
- 规范同步使永久的能力定义保持最新
- 归档目录以日期为前缀，便于检索

---

## 命令详解

### 核心命令

OpenSpec 提供了命令行工具（CLI）和技能命令（Skill Commands）两种方式。本项目主要使用技能命令，它们是对 CLI 的高级封装。

### CLI 命令

#### `openspec new change <name>`

**用途**: 创建新的变更目录

**语法**:
```bash
openspec new change "变更名称"
```

**参数**:
- `name`: 变更名称，使用 kebab-case 格式

**示例**:
```bash
openspec new change "add-user-profile"
```

**输出**:
```
✔ Created change 'add-user-profile' at openspec/changes/add-user-profile/ (schema: spec-driven)
```

**说明**:
- 自动创建变更目录
- 生成 `.openspec.yaml` 配置文件
- 使用项目默认的工作流模式（spec-driven）

---

#### `openspec status --change <name>`

**用途**: 查看变更的当前状态和进度

**语法**:
```bash
openspec status --change "<name>"
openspec status --change "<name>" --json  # JSON 格式输出
```

**示例**:
```bash
openspec status --change "add-openspec-doc"
```

**输出示例**:
```
Change: add-openspec-doc
Schema: spec-driven
Progress: 4/4 artifacts complete

[x] proposal
[x] design
[x] specs
[x] tasks

All artifacts complete!
```

**JSON 输出示例**:
```json
{
  "changeName": "add-openspec-doc",
  "schemaName": "spec-driven",
  "isComplete": true,
  "applyRequires": ["tasks"],
  "artifacts": [
    {"id": "proposal", "outputPath": "proposal.md", "status": "done"},
    {"id": "design", "outputPath": "design.md", "status": "done"},
    {"id": "specs", "outputPath": "specs/**/*.md", "status": "done"},
    {"id": "tasks", "outputPath": "tasks.md", "status": "done"}
  ]
}
```

**说明**:
- 显示工件完成情况
- `--json` 选项提供机器可读的格式
- 状态包括：`ready`（就绪）、`blocked`（阻塞）、`done`（完成）

---

#### `openspec instructions <artifact> --change <name>`

**用途**: 获取创建特定工件的指令和模板

**语法**:
```bash
openspec instructions <artifact-id> --change "<name>"
openspec instructions <artifact-id> --change "<name>" --json
```

**参数**:
- `artifact-id`: 工件 ID，如 `proposal`, `design`, `specs`, `tasks`, `apply`, `archive`

**示例**:
```bash
openspec instructions proposal --change "add-openspec-doc"
openspec instructions apply --change "add-openspec-doc" --json
```

**输出包含**:
- 工件的说明和用途
- 创建指令和要求
- 文档模板
- 依赖的其他工件

**说明**:
- 用于手动创建工件时获取指导
- JSON 格式适合自动化工具使用

---

#### `openspec list`

**用途**: 列出所有活动中的变更

**语法**:
```bash
openspec list
openspec list --json
```

**输出示例**:
```
Active changes:
- add-openspec-doc (spec-driven) - 4/4 artifacts complete
```

**说明**:
- 只显示 `openspec/changes/` 目录下的活动变更
- 不包括已归档的变更

---

### 技能命令

技能命令是对 CLI 的高级封装，提供更智能的交互体验。它们通常由 AI 助手（如 GitHub Copilot）执行。

#### `/opsx:propose`

**用途**: 一步创建变更并生成所有工件

**语法**:
```
/opsx:propose "变更描述"
```

**示例**:
```
/opsx:propose "为项目添加用户认证功能"
/opsx:propose "修复登录页面的样式问题"
```

**执行流程**:
1. 根据描述生成 kebab-case 变更名称
2. 创建变更目录
3. 自动生成 proposal.md
4. 自动生成 design.md
5. 自动生成 specs（根据 proposal 中的 capabilities）
6. 自动生成 tasks.md

**产出**:
- 完整的变更工件集
- 变更立即进入 ready 状态，可以开始实施

**优点**:
- 快速启动变更
- 减少手动编写工件的工作量
- AI 助手会根据上下文生成合理的文档内容

---

#### `/opsx:apply`

**用途**: 实施变更，执行任务清单中的任务

**语法**:
```
/opsx:apply
/opsx:apply <change-name>  # 指定变更名称
```

**示例**:
```
/opsx:apply
/opsx:apply add-openspec-doc
```

**执行流程**:
1. 检查变更状态，确认所有工件完成
2. 读取上下文文件（proposal、design、specs、tasks）
3. 显示当前进度
4. 逐一实施待完成的任务
5. 完成任务后标记为 `[x]`
6. 继续下一个任务，直到全部完成或遇到阻塞

**监控进度**:
- 实时显示当前正在执行的任务
- 显示已完成和剩余任务数量

**暂停条件**:
- 任务不清晰，需要澄清
- 实施中发现设计问题
- 遇到错误或阻塞
- 用户中断

---

#### `/opsx:archive`

**用途**: 归档已完成的变更

**语法**:
```
/opsx:archive
/opsx:archive <change-name>  # 指定变更名称
```

**示例**:
```
/opsx:archive
/opsx:archive add-readme1-file
```

**执行流程**:
1. 检查工件完成状态
2. 检查任务完成状态
3. 评估 delta specs 同步需求
4. 询问是否同步规范（如有 delta specs）
5. 创建归档目录（`YYYY-MM-DD-<name>`）
6. 移动变更目录到 `openspec/changes/archive/`
7. 同步 specs 到 `openspec/specs/`（如选择）

**确认项**:
- 未完成工件：警告并确认是否继续
- 未完成任务：警告并确认是否继续
- 规范同步：询问是否同步 delta specs

**产出**:
- 归档的变更记录
- 更新的主规范目录（如同步）

---

#### `/opsx:explore`

**用途**: 进入探索模式，用于思考和讨论变更

**语法**:
```
/opsx:explore
```

**说明**:
- 探索模式适合在创建变更前进行思考和讨论
- 帮助澄清需求、探索想法、调查问题
- 不直接创建或修改文件
- 可以在探索后再使用 `/opsx:propose` 创建正式变更

---

### 其他实用命令

#### 查看目录结构

```bash
# PowerShell
Get-ChildItem -Path openspec -Recurse -Directory

# 查看特定变更的文件
Get-ChildItem -Path openspec/changes/add-openspec-doc
```

#### 查看任务进度

```bash
# 查看任务文件
cat openspec/changes/add-openspec-doc/tasks.md

# 统计完成的任务
(Select-String -Path openspec/changes/add-openspec-doc/tasks.md -Pattern "- \[x\]").Count
```

---

## 工件详解

### Proposal（提案）详解

#### 作用

Proposal 是变更的起点，回答"为什么（Why）"需要这个变更。它为整个变更提供动机和方向。

#### 章节要求

##### 1. Why（为什么）

**要求**:
- 1-2 句话简明扼要说明问题或机会
- 回答：解决什么问题？为什么现在做？

**示例**:
```markdown
## Why

项目目前缺少一份详细的 OpenSpec 使用流程规范文档。虽然项目已经在使用 OpenSpec 工作流进行结构化的变更管理，但新加入的开发者和团队成员可能不熟悉完整的工作流程，包括如何提出变更、设计方案、编写规范、实施任务以及归档变更等步骤。
```

##### 2. What Changes（变更内容）

**要求**:
- 具体的变更项列表
- 明确新增、修改或删除的内容
- 标记破坏性变更（使用 **BREAKING**）

**示例**:
```markdown
## What Changes

- 在项目根目录创建 `openspec.md` 文件
- 文档将使用纯中文编写，便于中文用户理解
- 文档将详细说明 OpenSpec 的完整工作流程
- 包含每个命令的详细说明和使用示例
- 提供实际的工作流场景和最佳实践
```

##### 3. Capabilities（能力）

**要求**:
- 明确列出新增或修改的规范能力
- 使用 kebab-case 命名
- 每个能力对应一个 spec 文件

**示例**:
```markdown
## Capabilities

### New Capabilities

- `openspec-usage-documentation`: 提供 OpenSpec 工作流的完整中文使用指南，涵盖从创建变更到归档的全过程

### Modified Capabilities

<!-- 如果不修改现有能力，留空或注释 -->
```

##### 4. Impact（影响）

**要求**:
- 列出受影响的范围
- 包括代码、API、依赖、系统

**示例**:
```markdown
## Impact

- 在项目根目录添加新文件 `openspec.md`
- 为团队提供统一的 OpenSpec 使用规范参考
- 不影响现有代码或配置
- 有助于新成员快速上手 OpenSpec 工作流
```

#### 编写技巧

1. **保持简洁**: 1-2 页足够，详细内容放在 design 中
2. **关注动机**: 强调"为什么"而非"如何"
3. **明确能力**: Capabilities 部分决定后续 specs 的创建
4. **避免实现细节**: 技术决策属于 design，不属于 proposal

---

### Design（设计）详解

#### 作用

Design 说明"如何（How）"实现变更，提供技术方案和决策依据。

#### 章节要求

##### 1. Context（背景）

**要求**:
- 当前状态描述
- 限制条件和约束
- 相关方和目标用户

**示例**:
```markdown
## Context

项目已经成功使用 OpenSpec 工作流管理了多个变更（如 add-readme1-file 和 add-readme2-file），这些变更都遵循了 spec-driven 模式的完整生命周期。然而，项目缺少一份系统性的文档来说明 OpenSpec 的使用方法。

当前状态：
- 项目已初始化 OpenSpec（存在 openspec/config.yaml）
- 已有两个归档的成功案例可作为参考
- 团队已熟悉 propose → apply → archive 的基本流程
- 使用 spec-driven 工作流模式

目标用户：
- 新加入项目的开发者
- 需要了解变更管理流程的团队成员
```

##### 2. Goals / Non-Goals（目标与非目标）

**要求**:
- 明确变更要达到的目标
- 明确不包含在范围内的内容

**示例**:
```markdown
## Goals / Non-Goals

**Goals:**
- 创建一份完整的 OpenSpec 中文使用指南
- 详细说明四个核心命令：propose、apply、archive、sync
- 提供实际操作示例和完整的工作流演示

**Non-Goals:**
- 不替代 OpenSpec 官方文档（如果存在）
- 不涉及 OpenSpec CLI 的安装和配置
- 不讨论其他工作流模式（仅关注 spec-driven）
```

##### 3. Decisions（决策）

**要求**:
- 列出关键技术决策
- 说明选择的理由
- 列出备选方案和对比

**格式**:
```markdown
### 决策 N: 决策名称

**决定**: 简述决定内容

**理由**:
- 理由 1
- 理由 2

**替代方案**:
- 方案 A：为什么不选
- 方案 B：为什么不选
```

##### 4. Risks / Trade-offs（风险与权衡）

**要求**:
- 识别潜在风险
- 提供缓解措施
- 说明权衡取舍

**格式**:
```markdown
### 风险 N: 风险描述

**描述**: 详细说明风险

**缓解措施**:
- 措施 1
- 措施 2

### 权衡 N: 权衡名称

**选择**: 选择了什么

**理由**: 为什么这样选择
```

##### 5. Migration Plan（迁移计划）（可选）

**要求**:
- 部署步骤
- 回滚策略
- 数据迁移方案（如适用）

##### 6. Open Questions（待定问题）（可选）

**要求**:
- 列出需要解决的问题
- 标注倾向的答案（如有）

**示例**:
```markdown
## Open Questions

1. 是否需要包含 OpenSpec CLI 的安装说明？
   - 倾向：否，假设用户已经安装
   
2. 是否需要详细说明每个工件文件的 Markdown 格式规范？
   - 倾向：是，这对于手动编辑工件很有帮助
```

#### 编写技巧

1. **解释决策理由**: 每个关键决策都要说明"为什么"
2. **对比备选方案**: 让读者理解为什么选择当前方案
3. **识别风险**: 提前识别问题比事后处理更有价值
4. **保持客观**: 使用事实和数据支撑决策

---

### Specs（规范）详解

#### 作用

Specs 定义"做什么（What）"，通过需求（Requirements）和场景（Scenarios）提供明确的验收标准。

#### 格式规范

##### Requirements（需求）

**格式**:
```markdown
### Requirement: 需求名称
需求描述，使用规范性语言（SHALL/MUST）。

#### Scenario: 场景名称
- **WHEN** 触发条件
- **THEN** 期望结果
```

**要求**:
- 使用三个 # 标记 Requirement（`### Requirement:`）
- 使用规范性语言：SHALL（应当）、MUST（必须）、MAY（可以）
- 每个需求至少包含一个场景

##### Scenarios（场景）

**格式**:
- 使用四个 # 标记 Scenario（`#### Scenario:`） - **非常重要！**
- 采用 WHEN/THEN 格式
- WHEN 描述触发条件
- THEN 描述期望结果

**示例**:
```markdown
### Requirement: openspec.md 必须存在于项目根目录
项目必须在根目录包含一个 openspec.md 文件，作为 OpenSpec 工作流的使用指南文档。

#### Scenario: 文件位置正确
- **WHEN** 用户导航到项目根目录
- **THEN** 必须能找到 openspec.md 文件

#### Scenario: 文件内容完整
- **WHEN** 用户打开 openspec.md
- **THEN** 文件必须包含完整的工作流说明
```

#### Delta 操作

当修改现有规范时，使用 Delta 操作标记变更类型。

##### ADDED Requirements（新增需求）

用于新增功能或能力。

**格式**:
```markdown
## ADDED Requirements

### Requirement: 新需求名称
需求描述

#### Scenario: 场景1
- **WHEN** ...
- **THEN** ...
```

##### MODIFIED Requirements（修改需求）

用于修改现有需求的行为。

**重要**:
- 必须包含完整的更新内容
- 不能只列出变更部分

**格式**:
```markdown
## MODIFIED Requirements

### Requirement: 现有需求名称（完全匹配）
更新后的完整需求描述

#### Scenario: 更新的场景
- **WHEN** ...
- **THEN** ...
```

**工作流程**:
1. 从 `openspec/specs/<capability>/spec.md` 找到现有需求
2. 复制整个需求块（从 `### Requirement:` 到所有场景）
3. 粘贴到 delta spec 中的 `## MODIFIED Requirements` 下
4. 修改需求描述和场景

##### REMOVED Requirements（删除需求）

用于删除不再需要的功能。

**格式**:
```markdown
## REMOVED Requirements

### Requirement: 要删除的需求名称
**Reason**: 删除的原因
**Migration**: 迁移指导（如何替代）
```

##### RENAMED Requirements（重命名需求）

用于仅改变需求名称，不改变行为。

**格式**:
```markdown
## RENAMED Requirements

### Requirement: FROM: 旧名称
**TO**: 新名称
```

#### 编写技巧

1. **可测试性**: 每个场景应该能转化为测试用例
2. **明确性**: 使用 SHALL/MUST 等明确的语言
3. **完整性**: 覆盖正常和异常场景
4. **独立性**: 每个需求应该独立，不依赖其他需求的顺序
5. **四个 #**: Scenario 必须使用四个 #，这是解析器的要求！

---

### Tasks（任务）详解

#### 作用

Tasks 将实施工作分解为可追踪的具体任务，每个任务代表一个可完成的工作单元。

#### 格式规范

##### 任务组

使用二级标题组织相关任务：

```markdown
## 1. 任务组名称

- [ ] 1.1 任务描述
- [ ] 1.2 任务描述

## 2. 另一个任务组

- [ ] 2.1 任务描述
- [ ] 2.2 任务描述
```

##### 任务项

**格式**:
```markdown
- [ ] X.Y 任务描述
```

**要求**:
- 使用复选框：`- [ ]` 表示未完成，`- [x]` 表示已完成
- 任务编号：组号.任务号（如 1.1, 1.2, 2.1）
- 任务描述清晰、具体、可验证

##### 完整示例

```markdown
## 1. 环境配置

- [ ] 1.1 安装项目依赖包
- [ ] 1.2 配置数据库连接
- [ ] 1.3 初始化数据库表结构

## 2. 功能实现

- [ ] 2.1 创建用户模型（User model）
- [ ] 2.2 实现用户注册接口
- [ ] 2.3 实现用户登录接口
- [ ] 2.4 添加密码加密逻辑

## 3. 测试验证

- [ ] 3.1 编写单元测试
- [ ] 3.2 手动测试注册流程
- [ ] 3.3 手动测试登录流程
```

#### 任务分解原则

1. **粒度适中**: 每个任务应该能在一个工作会话（30-60分钟）内完成
2. **可验证**: 任务完成后能明确判断是否达标
3. **依赖顺序**: 按照依赖关系排序，先做基础工作
4. **功能分组**: 相关任务放在同一组
5. **清晰描述**: 避免模糊的描述，如"完成功能"

#### 任务状态管理

##### 标记完成

将 `- [ ]` 改为 `- [x]`：

```markdown
- [x] 1.1 安装项目依赖包
- [x] 1.2 配置数据库连接
- [ ] 1.3 初始化数据库表结构
```

##### 查看进度

使用 OpenSpec 命令：
```bash
openspec status --change "<name>"
openspec instructions apply --change "<name>"
```

#### 编写技巧

1. **从粗到细**: 先确定任务组，再细化每组的任务
2. **参考 specs**: 每个需求可能对应多个任务
3. **参考 design**: 设计决策会影响任务的分解
4. **包含质量检查**: 添加测试、审查、验证类任务
5. **留有余地**: 预留 bug 修复和调整的任务

---

## 目录结构

### OpenSpec 目录概览

```
openspec/
├── config.yaml                    # OpenSpec 配置文件
├── changes/                       # 活动变更目录
│   ├── add-openspec-doc/          # 示例：当前活动变更
│   │   ├── .openspec.yaml         # 变更配置
│   │   ├── proposal.md            # 提案
│   │   ├── design.md              # 设计
│   │   ├── tasks.md               # 任务清单
│   │   └── specs/                 # 规范（delta specs）
│   │       └── <capability>/
│   │           └── spec.md
│   └── archive/                   # 归档目录
│       ├── 2026-02-25-add-readme1-file/
│       │   ├── .openspec.yaml
│       │   ├── proposal.md
│       │   ├── design.md
│       │   ├── tasks.md
│       │   └── specs/
│       │       └── readme1-documentation/
│       │           └── spec.md
│       └── 2026-02-25-add-readme2-file/
│           └── ...
└── specs/                         # 主规范目录（持久化规范）
    ├── readme1-documentation/
    │   └── spec.md
    └── readme2-chinese-documentation/
        └── spec.md
```

### 目录说明

#### `openspec/config.yaml`

**作用**: OpenSpec 的项目配置文件

**内容**:
- 默认工作流模式
- 项目级别的约束和规则

#### `openspec/changes/`

**作用**: 存放所有活动中的变更

**说明**:
- 每个变更一个子目录
- 目录名即变更名（kebab-case）
- 正在开发中的变更都在这里

#### `openspec/changes/<name>/.openspec.yaml`

**作用**: 单个变更的配置文件

**内容**:
- 变更使用的工作流模式
- 变更特定的配置

#### `openspec/changes/<name>/specs/`

**作用**: 存放 delta specs（变更对规范的增量）

**说明**:
- 新能力：创建新的 spec 文件
- 修改能力：创建 delta spec，使用 MODIFIED/ADDED/REMOVED 标记

#### `openspec/changes/archive/`

**作用**: 存放已完成并归档的变更

**说明**:
- 归档目录名格式：`YYYY-MM-DD-<变更名>`
- 保留完整的变更记录
- 作为项目历史和知识库

#### `openspec/specs/`

**作用**: 主规范目录，存放当前生效的所有能力规范

**说明**:
- 每个能力一个子目录
- 只包含当前有效的规范，不包含历史版本
- 归档时从 delta specs 同步更新

### 文件命名规则

#### 变更命名

- 使用 kebab-case：小写字母，单词间用连字符
- 示例：`add-user-auth`, `fix-login-bug`, `refactor-api-layer`

#### Capability 命名

- 使用 kebab-case
- 清晰描述能力的功能域
- 示例：`user-authentication`, `data-export`, `api-rate-limiting`

#### 归档目录命名

- 格式：`YYYY-MM-DD-<变更名>`
- 示例：`2026-02-25-add-readme1-file`

---

## 实际使用示例

### 完整工作流演示

以 `add-readme1-file` 变更为例，展示完整的工作流程。

#### 1. 创建变更

**命令**:
```
/opsx:propose "给项目添加一个 README1.md 文件，介绍这是一个 Go 的学习系统"
```

**执行过程**:
- AI 助手推导变更名称：`add-readme1-file`
- 创建变更目录：`openspec/changes/add-readme1-file/`
- 生成 `proposal.md`、`design.md`、`specs/`、`tasks.md`

**查看状态**:
```bash
openspec status --change "add-readme1-file"
```

**输出**:
```
Change: add-readme1-file
Schema: spec-driven
Progress: 4/4 artifacts complete

[x] proposal
[x] design
[x] specs
[x] tasks

All artifacts complete!
```

#### 2. 查看生成的工件

##### proposal.md 内容片段

```markdown
## Why

项目需要一个学习指南文档，帮助初学者了解这是一个 Go 语言学习系统。

## What Changes

- 在项目根目录创建 README1.md 文件
- 文档包含项目介绍和学习目标
- 提供技术栈说明和学习路径

## Capabilities

### New Capabilities

- `readme1-documentation`: 提供 README1.md 学习指南能力
```

##### design.md 内容片段

```markdown
## Decisions

### 决策 1: 文档命名和位置

**决定**: 命名为 README1.md，放置在项目根目录

**理由**:
- 与现有 README.md 并存
- 提供额外的学习视角

**替代方案**:
- 替换 README.md：会丢失原有内容
- 放在 docs/ 目录：降低可见性
```

##### specs/readme1-documentation/spec.md 内容片段

```markdown
## ADDED Requirements

### Requirement: README1.md 必须存在于项目根目录
项目必须在根目录包含一个 README1.md 文件。

#### Scenario: 文件存在于正确位置
- **WHEN** 用户导航到项目根目录
- **THEN** 他们必须能找到 README1.md 文件
```

##### tasks.md 内容片段

```markdown
## 1. 文件创建

- [ ] 1.1 在项目根目录创建 README1.md 文件
- [ ] 1.2 添加文档标题和简介

## 2. 内容编写

- [ ] 2.1 编写项目介绍部分
- [ ] 2.2 编写学习目标部分
- [ ] 2.3 编写技术栈说明
```

#### 3. 实施变更

**命令**:
```
/opsx:apply
```

**执行过程**:
- AI 助手读取所有上下文文件
- 显示进度：0/20 tasks complete
- 开始实施任务 1.1：创建 README1.md 文件
- 标记任务 1.1 为完成：`- [x] 1.1 ...`
- 继续任务 1.2、2.1...
- 所有任务完成：20/20 tasks complete

**查看进度**（中途）:
```bash
openspec status --change "add-readme1-file"
```

**输出**:
```
Change: add-readme1-file
Schema: spec-driven
Progress: 4/4 artifacts complete
Tasks: 15/20 complete
```

#### 4. 归档变更

**命令**:
```
/opsx:archive
```

**执行过程**:
1. 检查工件状态：4/4 完成 ✓
2. 检查任务状态：20/20 完成 ✓
3. 发现 delta spec：`specs/readme1-documentation/`
4. 询问：是否同步规范到主目录？
5. 用户选择：现在同步
6. 创建主规范：`openspec/specs/readme1-documentation/spec.md`
7. 移动变更目录到：`openspec/changes/archive/2026-02-25-add-readme1-file/`

**最终结果**:
```
## 归档完成

**变更名称**: add-readme1-file
**工作流模式**: spec-driven
**归档位置**: openspec/changes/archive/2026-02-25-add-readme1-file/
**规范同步**: ✓ 已同步到 openspec/specs/readme1-documentation/

所有工件已完成 (4/4)
所有任务已完成 (20/20)
```

#### 5. 验证归档结果

**检查归档目录**:
```bash
Get-ChildItem -Path openspec/changes/archive/2026-02-25-add-readme1-file/
```

**输出**:
```
proposal.md
design.md
tasks.md
specs/
  readme1-documentation/
    spec.md
```

**检查主规范目录**:
```bash
Get-ChildItem -Path openspec/specs/
```

**输出**:
```
readme1-documentation/
  spec.md
```

---

## 工作流场景

### 场景 1: 添加新功能

**描述**: 为项目添加全新的功能或特性

**变更命名示例**:
- `add-user-authentication`
- `add-email-notification`
- `add-api-rate-limiting`

**工作流**:
1. 使用 `/opsx:propose "添加用户认证功能"`
2. 审查生成的 proposal 和 design
3. 检查 specs 中的需求和场景
4. 使用 `/opsx:apply` 实施
5. 完成后使用 `/opsx:archive` 归档

**注意事项**:
- 确保在 proposal 中明确列出新的 capability
- specs 使用 `## ADDED Requirements`
- 考虑与现有功能的集成

### 场景 2: 修复 Bug

**描述**: 修复代码中的错误或问题

**变更命名示例**:
- `fix-login-validation`
- `fix-memory-leak`
- `fix-api-response-format`

**工作流**:
1. 使用 `/opsx:propose "修复登录验证的空指针问题"`
2. 在 design 中说明 bug 的根本原因
3. specs 可能不需要新增需求，重点是任务清单
4. 实施并测试修复
5. 归档变更

**注意事项**:
- proposal 中清楚说明 bug 的影响
- design 中分析根本原因
- 如果修复涉及行为变更，需要更新 specs

### 场景 3: 重构代码

**描述**: 改善代码结构，不改变外部行为

**变更命名示例**:
- `refactor-auth-service`
- `refactor-database-layer`
- `refactor-api-routing`

**工作流**:
1. 使用 `/opsx:propose "重构认证服务，提高代码可维护性"`
2. 在 design 中详细说明重构方案和理由
3. specs 通常不需要新增或修改（行为不变）
4. 任务清单详细列出重构步骤
5. 实施并验证行为一致性
6. 归档变更

**注意事项**:
- 明确 Non-Goals：不改变外部行为
- 如果外部行为改变，需要更新 specs
- 包含充分的测试任务

### 场景 4: 添加文档

**描述**: 为项目添加或更新文档

**变更命名示例**:
- `add-api-documentation`
- `add-setup-guide`
- `update-architecture-doc`

**工作流**:
1. 使用 `/opsx:propose "为项目添加 API 文档"`
2. design 说明文档的结构和内容组织
3. specs 定义文档必须包含的内容
4. 任务清单分解文档的各个部分
5. 编写并审查文档
6. 归档变更

**注意事项**:
- specs 用于定义内容完整性要求
- 任务可以按文档章节分组
- 考虑文档的可维护性

### 场景 5: 修改现有功能

**描述**: 改变现有功能的行为

**变更命名示例**:
- `modify-user-permissions`
- `update-api-response-format`
- `change-caching-strategy`

**工作流**:
1. 使用 `/opsx:propose "修改用户权限模型"`
2. 在 proposal 的 Capabilities 中列出被修改的 capability
3. specs 使用 `## MODIFIED Requirements`
4. 在 design 中说明修改的理由和影响
5. 实施变更
6. 归档时 delta spec 会更新主规范

**注意事项**:
- **重要**: MODIFIED Requirements 必须包含完整的更新内容
- 评估对现有用户的影响
- 考虑是否需要迁移计划

---

## 最佳实践

### 提案编写最佳实践

1. **开门见山**: Why 部分直接说明核心问题，不要铺垫过多
2. **具体明确**: What Changes 使用具体的动作词，避免模糊表述
3. **准确列举**: Capabilities 仔细确认，每个都会生成一个 spec 文件
4. **影响评估**: Impact 不仅包括积极影响，也要提及可能的负面影响

**好的示例**:
```markdown
## Why

项目缺少用户认证机制，导致所有 API 端点对外开放，存在安全风险。

## What Changes

- 添加基于 JWT 的用户认证中间件
- 为 `/api/v1/*` 路由添加认证要求
- 创建用户登录和注册端点
- **BREAKING**: 所有现有 API 端点将需要认证
```

**不好的示例**:
```markdown
## Why

我们需要改善安全性。

## What Changes

- 添加一些安全功能
- 更新 API
```

### 设计文档编写最佳实践

1. **决策驱动**: 重点说明"为什么这样设计"，而非"设计是什么"
2. **对比方案**: 每个关键决策都列出备选方案，说明取舍
3. **风险关注**: 提前识别潜在问题，提出缓解措施
4. **保持更新**: 如果实施过程中发现设计问题，及时更新设计文档

**好的决策描述**:
```markdown
### 决策 1: 使用 JWT 而非 Session

**决定**: 采用 JWT（JSON Web Token）进行用户认证

**理由**:
- 无状态：不需要服务端存储会话，便于水平扩展
- 跨域友好：适合前后端分离架构
- 标准化：广泛使用的行业标准

**替代方案**:
- Session-based 认证：
  - 优点：服务端完全控制，可随时撤销
  - 缺点：需要共享存储（如 Redis），增加依赖
  - 不选择原因：项目目前不需要即时撤销能力，优先考虑简单性
```

### 规范编写最佳实践

1. **需求独立**: 每个需求应该独立，不依赖其他需求的顺序
2. **场景完整**: 覆盖正常流程、边界情况和错误场景
3. **可测试性**: 每个场景应该能直接转化为测试用例
4. **使用规范性语言**: SHALL（应当）、MUST（必须）、MAY（可以）
5. **四个 #**: Scenario 必须使用 `####` 四个 #，否则解析失败

**好的需求和场景**:
```markdown
### Requirement: 系统必须验证 JWT 令牌的有效性
API 端点在处理请求前必须验证 JWT 令牌的签名和过期时间。

#### Scenario: 有效令牌被接受
- **WHEN** 用户发送带有有效 JWT 令牌的请求
- **THEN** 系统必须允许请求继续处理

#### Scenario: 过期令牌被拒绝
- **WHEN** 用户发送带有过期 JWT 令牌的请求
- **THEN** 系统必须返回 401 Unauthorized 错误

#### Scenario: 无效签名被拒绝
- **WHEN** 用户发送带有无效签名的 JWT 令牌
- **THEN** 系统必须返回 401 Unauthorized 错误
```

### 任务清单编写最佳实践

1. **合理粒度**: 每个任务 30-60 分钟可完成
2. **依赖顺序**: 先做基础任务，再做依赖任务
3. **功能分组**: 相关任务放在同一组
4. **包含验证**: 每个功能组后面加测试验证任务
5. **清晰描述**: 避免模糊词汇，使用具体的动词

**好的任务清单**:
```markdown
## 1. 数据模型

- [ ] 1.1 创建 User 模型（包含 id, username, password_hash）
- [ ] 1.2 创建数据库迁移脚本
- [ ] 1.3 运行迁移，创建 users 表

## 2. JWT 工具函数

- [ ] 2.1 实现 GenerateToken 函数（生成 JWT）
- [ ] 2.2 实现 ValidateToken 函数（验证 JWT）
- [ ] 2.3 添加 JWT 工具函数的单元测试

## 3. 认证中间件

- [ ] 3.1 创建 AuthMiddleware 中间件
- [ ] 3.2 实现令牌提取逻辑（从 Authorization header）
- [ ] 3.3 实现令牌验证逻辑
- [ ] 3.4 集成中间件到路由

## 4. 测试验证

- [ ] 4.1 测试有效令牌的请求
- [ ] 4.2 测试无效令牌的请求
- [ ] 4.3 测试缺少令牌的请求
```

### 变更命名最佳实践

**格式**: 使用 kebab-case（小写字母，单词间用连字符）

**模式**:
- `add-<feature>`: 添加新功能
- `fix-<issue>`: 修复问题
- `refactor-<component>`: 重构组件
- `update-<item>`: 更新内容
- `remove-<deprecated>`: 移除废弃功能

**好的示例**:
- `add-user-authentication`
- `fix-null-pointer-in-login`
- `refactor-database-layer`
- `update-api-documentation`
- `remove-legacy-payment-gateway`

**不好的示例**:
- `authentication`（缺少动作）
- `fixBug`（不是 kebab-case）
- `ADD_USER_FEATURE`（不是 kebab-case）
- `user_auth`（使用下划线，不是连字符）

### 何时创建新规范 vs 修改现有规范

**创建新规范**（使用 ADDED Requirements）:
- 添加全新的功能或能力
- 与现有功能相互独立
- 可以单独描述和验证

**修改现有规范**（使用 MODIFIED Requirements）:
- 改变现有功能的行为
- 修改现有 API 的接口或返回格式
- 调整现有需求的约束条件

**判断标准**:
- 如果新功能可以在不修改现有行为的前提下添加 → 新规范
- 如果必须改变现有功能才能实现 → 修改现有规范

### 常见错误和避免方法

#### 错误 1: Scenario 使用三个 # 而非四个

**错误**:
```markdown
### Scenario: 测试场景
```

**正确**:
```markdown
#### Scenario: 测试场景
```

**避免方法**: 记住"Requirement 三个 #，Scenario 四个 #"

#### 错误 2: MODIFIED Requirements 只包含变更部分

**错误**:
```markdown
## MODIFIED Requirements

### Requirement: 用户登录
添加验证码验证
```

**正确**:
```markdown
## MODIFIED Requirements

### Requirement: 用户登录
系统必须验证用户的用户名、密码和验证码，三者都正确才允许登录。

#### Scenario: 成功登录
- **WHEN** 用户提供正确的用户名、密码和验证码
- **THEN** 系统必须返回 JWT 令牌

#### Scenario: 验证码错误
- **WHEN** 用户提供正确的用户名和密码，但验证码错误
- **THEN** 系统必须拒绝登录
```

**避免方法**: 从原规范文件复制完整需求，再修改

#### 错误 3: 任务没有使用复选框格式

**错误**:
```markdown
## 1. 任务组

1. 创建文件
2. 编写内容
```

**正确**:
```markdown
## 1. 任务组

- [ ] 1.1 创建文件
- [ ] 1.2 编写内容
```

**避免方法**: 始终使用 `- [ ]` 格式，OpenSpec 依赖此格式追踪进度

#### 错误 4: 变更名称不符合 kebab-case

**错误**:
- `AddUserAuth`
- `add_user_auth`
- `ADD-USER-AUTH`

**正确**:
- `add-user-auth`

**避免方法**: 全小写，单词间用连字符

#### 错误 5: Proposal 的 Capabilities 列表不准确

**错误**: 列出的 capability 与实际创建的 spec 不匹配

**避免方法**:
- 在创建 specs 前仔细规划 Capabilities
- 每个 New Capability 对应一个新的 spec 文件
- 每个 Modified Capability 对应一个 delta spec

---

## 常见问题解答

### Q1: 如何开始一个新变更？

**A**: 有两种方式：

**方式 1: 使用技能命令（推荐）**
```
/opsx:propose "变更描述"
```
AI 助手会自动创建变更并生成所有工件。

**方式 2: 手动创建**
```bash
openspec new change "变更名称"
# 然后手动创建 proposal.md、design.md 等工件
```

**推荐使用方式 1**，因为更快速且自动生成的工件质量高。

### Q2: 工件之间的依赖关系是什么？

**A**: 依赖关系如下：

```
proposal → design → tasks
proposal → specs → tasks
```

- **proposal** 必须首先创建
- **design** 和 **specs** 依赖 proposal，可以并行创建
- **tasks** 依赖 design 和 specs，必须最后创建
- 只有所有工件完成后才能执行 apply

### Q3: 如何查看变更的当前进度？

**A**: 使用以下命令：

```bash
# 查看工件完成状态
openspec status --change "<变更名称>"

# 查看任务完成情况
openspec instructions apply --change "<变更名称>"
```

或者直接查看 tasks.md 文件，统计 `[x]` 的数量。

### Q4: 何时应该归档变更？

**A**: 满足以下条件时归档：

1. 所有工件（proposal、design、specs、tasks）都已创建
2. 所有任务都已完成（tasks.md 中所有复选框都标记为 `[x]`）
3. 功能已测试验证
4. 代码已提交到版本控制系统

使用命令：
```
/opsx:archive
```

### Q5: 规范同步（Spec Sync）是什么？

**A**: 规范同步是指将变更中的 delta specs 合并到主规范目录。

**两个规范目录**:
- `openspec/changes/<name>/specs/`: Delta specs（变更的增量）
- `openspec/specs/`: 主规范（当前生效的所有规范）

**同步时机**: 在归档变更时

**同步操作**:
- 对于新 capability：复制 delta spec 到主规范目录
- 对于修改的 capability：合并 ADDED/MODIFIED/REMOVED 到主规范

**为什么需要同步**:
- 主规范目录代表项目当前的完整能力定义
- Delta specs 只存在于变更中，作为历史记录
- 同步确保主规范始终保持最新

### Q6: 如何修改已有的规范？

**A**: 使用 Delta 操作：

1. 在 proposal 的 **Modified Capabilities** 中列出要修改的 capability
2. 创建 delta spec 文件：`specs/<capability>/spec.md`
3. 使用 `## MODIFIED Requirements` 标记
4. 从主规范复制完整的需求内容，然后修改
5. 归档时选择同步规范

**重要**: MODIFIED Requirements 必须包含完整的更新内容，不能只列出变更部分。

### Q7: 任务标记的规则是什么？

**A**: 任务标记规则：

**格式要求**:
- 未完成：`- [ ] X.Y 任务描述`
- 已完成：`- [x] X.Y 任务描述`

**注意事项**:
- 必须使用 `- [ ]` 和 `- [x]`，不能用其他格式
- 方括号内的空格和 x 很重要
- OpenSpec 依赖此格式解析任务进度

**更新方式**:
- 手动修改 tasks.md 文件
- 或使用 `/opsx:apply` 命令，AI 助手会自动标记

### Q8: 变更失败或遇到阻塞如何处理？

**A**: 处理步骤：

**1. 评估问题**:
- 任务不清晰：更新 tasks.md，细化任务描述
- 设计有问题：更新 design.md，调整技术方案
- 需求有误：更新 specs，修正需求定义

**2. 更新工件**:
直接编辑相关的 md 文件，工件可以在实施过程中修改。

**3. 重新开始或继续**:
```
/opsx:apply  # 继续实施剩余任务
```

**4. 如果问题无法解决**:
- 可以暂停变更，不归档
- 可以创建新变更来修复问题
- 活动变更可以长期存在，不需要立即完成

### Q9: 可以跳过某些工件吗？

**A**: 不推荐，但 design.md 在某些情况下可选。

**必需工件**:
- proposal.md（必需）
- specs（必需）
- tasks.md（必需）

**可选工件**:
- design.md（对于简单变更可能不需要）

**判断标准**（design.md 是否需要）:
- 跨多个模块的变更：需要
- 新的架构模式：需要
- 重要的技术决策：需要
- 简单的内容添加（如文档）：可能不需要

**建议**: 初学者建议始终创建所有工件，帮助养成完整的思考习惯。

### Q10: 如何查看历史变更？

**A**: 查看归档目录：

```bash
# 列出所有归档的变更
Get-ChildItem -Path openspec/changes/archive/

# 查看特定归档变更的内容
Get-ChildItem -Path openspec/changes/archive/2026-02-25-add-readme1-file/
```

每个归档变更包含完整的工件，可以作为参考示例。

---

## 总结

OpenSpec 工作流为项目变更管理提供了结构化、可追溯的方法。通过遵循 spec-driven 模式，我们能够确保每个变更都有清晰的动机（proposal）、合理的设计（design）、明确的需求（specs）和可执行的任务（tasks）。

### 关键要点

1. **工作流程**: propose → design & specs → tasks → apply → archive
2. **四种工件**: proposal（为什么）、design（如何做）、specs（做什么）、tasks（具体步骤）
3. **技能命令**: `/opsx:propose`、`/opsx:apply`、`/opsx:archive` 简化操作
4. **规范格式**: Requirement 用三个 #，Scenario 用四个 #
5. **任务格式**: 必须使用 `- [ ]` 复选框格式
6. **归档同步**: 归档时同步 delta specs 到主规范目录

### 开始使用

如果你是第一次使用 OpenSpec，建议：

1. **阅读本文档**：理解核心概念和工作流
2. **查看归档案例**：参考 `openspec/changes/archive/` 中的实际案例
3. **尝试小变更**：从简单的文档添加开始，熟悉完整流程
4. **使用技能命令**：`/opsx:propose` 快速生成工件
5. **逐步深入**：掌握基本流程后，学习 delta specs 和规范修改

### 获得帮助

- 查看归档的变更案例
- 参考本文档的相关章节
- 使用 `/opsx:explore` 进入探索模式，与 AI 助手讨论问题
- 查看 OpenSpec CLI 的 `--help` 选项

---

**文档版本**: 1.0  
**最后更新**: 2026年2月25日  
**适用于**: OpenSpec CLI (spec-driven 工作流)
