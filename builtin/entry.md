# 内置词条使用指南

## 📊 词条统计

- **总词条数**: 230
- **支持语言**: 6 种（en, zh-CN, zh-TW, ja, ko, ru）
- **带参数词条**: 3 个
- **复数形式词条**: 1 个

---

## 📋 目录

1. [基础操作](#1-基础操作)
2. [状态标识](#2-状态标识)
3. [表单验证](#3-表单验证)
4. [日期时间](#4-日期时间)
5. [文件操作](#5-文件操作)
6. [用户账户](#6-用户账户)
7. [列表数据](#7-列表数据)
8. [消息通知](#8-消息通知)
9. [确认提示](#9-确认提示)
10. [内容分类](#10-内容分类)
11. [错误处理](#11-错误处理)
12. [操作反馈](#12-操作反馈)
13. [设置偏好](#13-设置偏好)
14. [其他通用](#14-其他通用)

---

## 1. 基础操作

> 37 个词条

### 增删改查

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `add` | Add | 添加 |
| `create` | Create | 创建 |
| `new` | New | 新建 |
| `view` | View | 查看 |
| `details` | Details | 详情 |
| `show` | Show | 显示 |
| `hide` | Hide | 隐藏 |
| `edit` | Edit | 编辑 |
| `update` | Update | 更新 |
| `save` | Save | 保存 |
| `apply` | Apply | 应用 |
| `delete` | Delete | 删除 |
| `remove` | Remove | 移除 |
| `clear` | Clear | 清空 |

### 通用操作

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `submit` | Submit | 提交 |
| `confirm` | OK | 确定 |
| `cancel` | Cancel | 取消 |
| `reset` | Reset | 重置 |
| `search` | Search | 搜索 |
| `refresh` | Refresh | 刷新 |
| `reload` | Reload | 重新加载 |
| `close` | Close | 关闭 |
| `back` | Back | 返回 |
| `next` | Next | 下一步 |
| `prev` | Previous | 上一步 |

### 剪贴板

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `copy` | Copy | 复制 |
| `paste` | Paste | 粘贴 |
| `cut` | Cut | 剪切 |
| `copy_link` | Copy Link | 复制链接 |

### 选择操作

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `select` | Select | 选择 |
| `deselect` | Deselect | 取消选择 |
| `select_all` | Select All | 全选 |
| `deselect_all` | Deselect All | 取消全选 |

### 展开折叠

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `expand` | Expand | 展开 |
| `collapse` | Collapse | 折叠 |
| `more` | More | 更多 |
| `less` | Less | 收起 |

---

## 2. 状态标识

> 19 个词条

### 处理状态

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `pending` | Pending | 待处理 |
| `processing` | Processing | 处理中 |
| `loading` | Loading... | 加载中... |
| `completed` | Completed | 已完成 |
| `cancelled` | Cancelled | 已取消 |
| `expired` | Expired | 已过期 |

### 启用状态

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `active` | Active | 活跃 |
| `inactive` | Inactive | 不活跃 |
| `enabled` | Enabled | 已启用 |
| `disabled` | Disabled | 已禁用 |

### 连接状态

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `online` | Online | 在线 |
| `offline` | Offline | 离线 |
| `available` | Available | 可用 |
| `unavailable` | Unavailable | 不可用 |

### 结果状态

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `success` | Success | 操作成功 |
| `failed` | Failed | 操作失败 |
| `error` | Error | 错误 |
| `warning` | Warning | 警告 |
| `info` | Info | 信息 |

---

## 3. 表单验证

> 22 个词条

### 基础字段

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `username` | Username | 用户名 |
| `password` | Password | 密码 |
| `email` | Email | 邮箱 |
| `phone` | Phone | 手机号 |
| `name` | Name | 名称 |
| `title` | Title | 标题 |
| `description` | Description | 描述 |
| `content` | Content | 内容 |

### 字段属性

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `required` | Required | 必填 |
| `optional` | Optional | 可选 |
| `invalid` | Invalid | 格式错误 |

### 验证错误

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `required_field` | This field is required | 此字段为必填项 |
| `invalid_email` | Invalid email address | 邮箱地址格式错误 |
| `invalid_phone` | Invalid phone number | 手机号码格式错误 |
| `invalid_format` | Invalid format | 格式错误 |
| `password_too_short` | Password is too short | 密码太短 |
| `password_too_weak` | Password is too weak | 密码强度太弱 |
| `passwords_not_match` | Passwords do not match | 两次输入的密码不一致 |
| `email_exists` | Email already exists | 该邮箱已被使用 |
| `username_taken` | Username is already taken | 该用户名已被占用 |

### 长度验证（带参数）

| 词条 ID | 英文 | 中文 | 参数 |
|---------|------|------|------|
| `min_length` | Minimum length is {{.Min}} | 最小长度为 {{.Min}} | Min |
| `max_length` | Maximum length is {{.Max}} | 最大长度为 {{.Max}} | Max |

```go
gi18n.Tf("min_length", "Min", 6)   // "最小长度为 6"
gi18n.Tf("max_length", "Max", 20)  // "最大长度为 20"
```

---

## 4. 日期时间

> 13 个词条

### 相对时间

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `today` | Today | 今天 |
| `yesterday` | Yesterday | 昨天 |
| `tomorrow` | Tomorrow | 明天 |
| `now` | Now | 现在 |

### 时间字段

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `date` | Date | 日期 |
| `time` | Time | 时间 |
| `datetime` | Date & Time | 日期时间 |
| `created_at` | Created At | 创建时间 |
| `updated_at` | Updated At | 更新时间 |
| `deleted_at` | Deleted At | 删除时间 |
| `start_date` | Start Date | 开始日期 |
| `end_date` | End Date | 结束日期 |
| `duration` | Duration | 时长 |

---

## 5. 文件操作

> 16 个词条

### 文件类型

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `file` | File | 文件 |
| `folder` | Folder | 文件夹 |
| `image` | Image | 图片 |
| `video` | Video | 视频 |
| `audio` | Audio | 音频 |

### 文件操作

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `upload` | Upload | 上传 |
| `download` | Download | 下载 |
| `upload_file` | Upload File | 上传文件 |
| `select_file` | Select File | 选择文件 |
| `drag_drop` | Drag & Drop | 拖拽上传 |
| `print` | Print | 打印 |
| `export` | Export | 导出 |
| `import` | Import | 导入 |

### 文件属性

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `file_size` | File Size | 文件大小 |
| `max_file_size` | Max File Size | 最大文件大小 |
| `allowed_types` | Allowed Types | 允许的文件类型 |

---

## 6. 用户账户

> 27 个词条

### 认证操作

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `login` | Login | 登录 |
| `logout` | Logout | 退出登录 |
| `register` | Register | 注册 |

### 密码管理

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `change_password` | Change Password | 修改密码 |
| `forgot_password` | Forgot Password | 忘记密码 |
| `reset_password` | Reset Password | 重置密码 |
| `remember_me` | Remember Me | 记住我 |

### 用户资料

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `profile` | Profile | 个人资料 |
| `account` | Account | 账户 |
| `welcome` | Welcome | 欢迎 |
| `greeting` | Hello, {{.Name}}! | 你好，{{.Name}}！ |

```go
gi18n.Tf("greeting", "Name", "张三")  // "你好，张三！"
```

### 角色权限

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `admin` | Admin | 管理员 |
| `user` | User | 用户 |
| `guest` | Guest | 访客 |
| `role` | Role | 角色 |
| `permission` | Permission | 权限 |

### 验证码

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `verification_code` | Verification Code | 验证码 |
| `send_code` | Send Code | 发送验证码 |
| `resend_code` | Resend Code | 重新发送 |
| `verify` | Verify | 验证 |
| `code_expired` | Code expired | 验证码已过期 |
| `invalid_code` | Invalid code | 验证码错误 |
| `captcha` | Captcha | 验证码 |

### 账户状态

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `account_not_found` | Account not found | 账号不存在 |
| `incorrect_password` | Incorrect password | 密码错误 |
| `account_locked` | Account locked | 账号已锁定 |
| `account_disabled` | Account disabled | 账号已禁用 |

---

## 7. 列表数据

> 16 个词条

### 排序筛选

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `all` | All | 全部 |
| `filter` | Filter | 筛选 |
| `sort` | Sort | 排序 |
| `sort_by` | Sort By | 排序方式 |
| `order_by` | Order By | 排序依据 |
| `ascending` | Ascending | 升序 |
| `descending` | Descending | 降序 |

### 统计计数

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `count` | Count | 数量 |
| `total` | Total | 总计 |
| `total_count` | Total Count | 总数 |
| `selected_count` | Selected Count | 已选数量 |

### 项目数量（复数，带参数）

| 词条 ID | 英文 | 中文 | 参数 |
|---------|------|------|------|
| `items` | {{.Count}} items | {{.Count}} 个项目 | Count |

```go
gi18n.Tp("items", 1)   // "1 item" / "1 个项目"
gi18n.Tp("items", 5)   // "5 items" / "5 个项目"
```

### 显示方式

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `list` | List | 列表 |
| `grid` | Grid | 网格 |
| `table` | Table | 表格 |
| `card` | Card | 卡片 |

---

## 8. 消息通知

> 7 个词条

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `message` | Message | 消息 |
| `notification` | Notification | 通知 |
| `send` | Send | 发送 |
| `receive` | Receive | 接收 |
| `inbox` | Inbox | 收件箱 |
| `outbox` | Outbox | 发件箱 |
| `draft` | Draft | 草稿 |

---

## 9. 确认提示

> 17 个词条

### 确认对话框

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `yes` | Yes | 是 |
| `no` | No | 否 |
| `ok` | OK | 好的 |
| `got_it` | Got it | 知道了 |
| `are_you_sure` | Are you sure? | 您确定吗？ |
| `confirm_delete` | Are you sure you want to delete this? | 您确定要删除此项吗？ |
| `confirm_action` | Are you sure you want to proceed? | 您确定要继续操作吗？ |
| `operation_cancelled` | Operation cancelled | 操作已取消 |

### 空状态提示

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `no_data` | No data | 暂无数据 |
| `no_results` | No results found | 未找到结果 |
| `not_found` | Not found | 未找到 |
| `page_not_found` | Page not found | 页面未找到 |

### 权限会话

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `access_denied` | Access denied | 访问被拒绝 |
| `permission_denied` | Permission denied | 权限不足 |
| `session_expired` | Session expired | 会话已过期 |
| `unauthorized` | Unauthorized | 未授权 |
| `forbidden` | Forbidden | 禁止访问 |

---

## 10. 内容分类

> 11 个词条

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `type` | Type | 类型 |
| `status` | Status | 状态 |
| `category` | Category | 分类 |
| `tag` | Tag | 标签 |
| `action` | Action | 操作 |
| `comment` | Comment | 评论 |
| `note` | Note | 备注 |
| `latest` | Latest | 最新 |
| `popular` | Popular | 热门 |
| `featured` | Featured | 精选 |
| `recommended` | Recommended | 推荐 |

---

## 11. 错误处理

> 8 个词条

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `network_error` | Network Error | 网络错误 |
| `timeout` | Timeout | 请求超时 |
| `server_error` | Server Error | 服务器错误 |
| `bad_request` | Bad Request | 请求错误 |
| `request_failed` | Request Failed | 请求失败 |
| `please_try_again` | Please try again | 请重试 |
| `something_went_wrong` | Something went wrong | 出错了 |
| `try_again_later` | Please try again later | 请稍后重试 |

---

## 12. 操作反馈

> 16 个词条

### 成功反馈

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `login_success` | Login successful | 登录成功 |
| `logout_success` | Logout successful | 退出成功 |
| `register_success` | Registration successful | 注册成功 |
| `saved_successfully` | Saved successfully | 保存成功 |
| `deleted_successfully` | Deleted successfully | 删除成功 |
| `updated_successfully` | Updated successfully | 更新成功 |
| `created_successfully` | Created successfully | 创建成功 |
| `operation_successful` | Operation successful | 操作成功 |

### 失败反馈

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `login_failed` | Login failed | 登录失败 |
| `operation_failed` | Operation failed | 操作失败 |

### 处理中

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `please_wait` | Please wait... | 请稍候... |
| `processing_request` | Processing... | 处理中... |
| `syncing` | Syncing... | 同步中... |
| `synced` | Synced | 已同步 |
| `draft_saved` | Draft saved | 草稿已保存 |
| `auto_saved` | Auto saved | 自动保存 |

---

## 13. 设置偏好

> 6 个词条

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `settings` | Settings | 设置 |
| `preferences` | Preferences | 偏好设置 |
| `language` | Language | 语言 |
| `theme` | Theme | 主题 |
| `dark_mode` | Dark Mode | 暗黑模式 |
| `light_mode` | Light Mode | 明亮模式 |

---

## 14. 其他通用

> 15 个词条

### 收藏分享

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `favorite` | Favorite | 收藏 |
| `bookmark` | Bookmark | 书签 |
| `share` | Share | 分享 |

### 归档回收

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `trash` | Trash | 回收站 |
| `archive` | Archive | 归档 |
| `restore` | Restore | 恢复 |

### 批量操作

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `batch_delete` | Batch Delete | 批量删除 |
| `batch_operation` | Batch Operation | 批量操作 |

### 导航控制

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `continue` | Continue | 继续 |
| `skip` | Skip | 跳过 |
| `finish` | Finish | 完成 |

### 其他

| 词条 ID | 英文 | 中文 |
|---------|------|------|
| `link` | Link | 链接 |
| `url` | URL | 网址 |
| `version` | Version | 版本 |
| `size` | Size | 大小 |

---

## 🌍 支持的语言

| 代码 | 语言 |
|-----|------|
| `en` | English（英语） |
| `zh-CN` | 简体中文 |
| `zh-TW` | 繁體中文 |
| `ja` | 日本語（日语） |
| `ko` | 한국어（韩语） |
| `ru` | Русский（俄语） |

---

## 📝 使用示例

### 简单翻译

```go
gi18n.SetLang("zh-CN")
gi18n.T("confirm")     // "确定"
gi18n.T("cancel")      // "取消"
```

### 带参数翻译

```go
gi18n.Tf("greeting", "Name", "张三")     // "你好，张三！"
gi18n.Tf("min_length", "Min", 6)         // "最小长度为 6"
```

### 复数翻译

```go
gi18n.Tp("items", 1)   // "1 个项目"
gi18n.Tp("items", 5)   // "5 个项目"
```

### 指定语言翻译

```go
gi18n.TL("en", "confirm")      // "OK"
gi18n.TL("ja", "confirm")      // "確認"
```

---

**词条总数**: 230  
**版本**: 1.0.0  
**更新时间**: 2026-02-14
