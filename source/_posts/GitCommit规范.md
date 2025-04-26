---
type: Post
title: Git Commit 规范
tags: Git
category: 开发
category_bar: true
abbrlink: 12618
date: 2024-12-11 15:57:32
---
## git commit 规范

### **1. Commit 信息的基本格式**

```git
<type>(<scope>): <subject>

<body>  # 可选
<footer>  # 可选
```

------

### **2. 详细说明**

#### **2.1. `type`（类型）**

描述本次提交的性质，通常包括以下几类：

| 类型       | 描述                                                |
| ---------- | --------------------------------------------------- |
| `feat`     | 新功能（feature）                                   |
| `fix`      | 修复 Bug                                            |
| `docs`     | 文档（documentation）修改                           |
| `style`    | 不影响代码逻辑的修改（如格式化、漏掉分号等）        |
| `refactor` | 代码重构（既不是新增功能也不是修复 Bug 的代码改动） |
| `test`     | 添加或修改测试                                      |
| `chore`    | 其他不修改源代码或测试的事务性更改（如构建脚本）    |
| `perf`     | 提升性能                                            |
| `revert`   | 回滚某次提交                                        |

------

#### **2.2. `scope`（范围）**

- 表明提交代码的影响范围，例如模块、功能或文件夹。
- 可选项，但在大型项目中建议加上。

**示例：**

- `feat(auth): add JWT authentication`
- `fix(ui): correct button alignment`

------

#### **2.3. `subject`（主题）**

- **简明扼要**地描述提交的改动。
- 建议使用 **动词原形** 开头（如 add、fix、update 等）。
- 避免句号结尾，保持风格简洁。

------

#### **2.4. `body`（正文，选填）**

- 提供更详细的修改原因和实现细节。
- 解释代码修改的上下文，描述该改动的意义。

**示例：**

```git
feat(cache): add Redis-based caching for user sessions

The new cache layer improves performance by reducing database queries.
Sessions are now stored in Redis, which supports high availability.
```

------

#### **2.5. `footer`（页脚，选填）**

页脚通常用于：

- 关闭相关问题（例如 `Closes #123`）。
- 影响的重大变更说明（如 BREAKING CHANGES）。

**示例：**

```git
BREAKING CHANGE: The 'auth' module API has been updated. Please update your code accordingly.
```

------

### **3. 实践中的示例**

1. **新功能：**

   ```git
   feat(auth): add OAuth2.0 support
   ```

2. **修复 Bug：**

   ```git
   fix(db): handle null pointer exceptions in queries
   ```

3. **文档更新：**

   ```git
   docs(readme): update deployment instructions
   ```

4. **性能优化：**

   ```git
   perf(cache): reduce response time by implementing in-memory caching
   ```

5. **回滚提交：**

   ```git
   revert: fix(ui): correct button alignment

   This reverts commit a1b2c3d4e5.
   ```
