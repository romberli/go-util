# go-util

供 OpenCode、Codex 等通用 agent 使用的项目级工作说明。与用户明确指令冲突时，以用户指令为准。

## 沟通规范

- 默认使用中文与用户沟通，除非用户明确要求英文
- Commit message 必须使用全英文
- 代码中的日志输出（`log`/`logger` 调用）必须使用全英文

## Superpowers 简写约定

- 规范写法使用 `superpowers/{skill}`，例如 `superpowers/brainstorming`
- 简写规则默认使用 `sp` + skill 名各单词首字母缩写；如果出现缩写冲突，则补充冲突单词的后续字母以保证唯一
- `spb` 等价于 `superpowers/brainstorming`
- `spdpa` 等价于 `superpowers/dispatching-parallel-agents`
- `spep` 等价于 `superpowers/executing-plans`
- `spfadb` 等价于 `superpowers/finishing-a-development-branch`
- `spreccr` 等价于 `superpowers/receiving-code-review`
- `spreqcr` 等价于 `superpowers/requesting-code-review`
- `spsdd` 等价于 `superpowers/subagent-driven-development`
- `spwp` 等价于 `superpowers/writing-plans`
- `spws` 等价于 `superpowers/writing-skills`
- `spsd` 等价于 `superpowers/systematic-debugging`
- `sptdd` 等价于 `superpowers/test-driven-development`
- `spugw` 等价于 `superpowers/using-git-worktrees`
- `spus` 等价于 `superpowers/using-superpowers`
- `spvbc` 等价于 `superpowers/verification-before-completion`
- 当用户以这些简称开头描述任务时，agent 应将其视为对应的 `superpowers` skill 调用意图，并优先按对应工作模式处理
- 当用户输入 `superpowers/{skill}` 完整写法或上述简称时，agent 应优先尝试真实加载对应的 `superpowers` skill；若当前运行环境可直接调用该 skill，则必须先调用 skill 再继续处理任务；若当前运行环境不可直接调用该 skill，则应明确说明，并按对应 skill 的工作模式执行

## 必须遵守的编码规范

### 错误处理

- 如果返回的错误类型不是 `github.com/pingcap/errors`（通常是调用依赖包函数返回的错误），则必须用 `errors.Trace(err)` 包裹后返回
- 禁止在业务代码中硬编码错误字符串，必须在 `pkg/message/{module}/` 中定义消息常量
- Handler 层错误响应必须用 `resp.ResponseNOK(c, msgXxx.ErrXxx, err, ...args)`
- Handler 层成功响应必须用 `resp.ResponseOK(c, jsonStr, msgXxx.InfoXxx, ...args)`

### git
- 如果需要创建worktree, 则在项目根目录下创建.worktree目录

## 常用命令

## 依赖包规范

- 错误处理统一使用 `github.com/pingcap/errors`，禁止使用标准库 `errors` 或其他错误包
- 日志输出统一使用 `github.com/romberli/log`，禁止使用标准库 `log` 或其他日志包

## 禁止事项

- 如果返回的错误类型不是 `github.com/pingcap/errors`（通常是调用依赖包函数返回的错误），则必须用 `errors.Trace(err)` 包裹后返回
- 禁止硬编码错误字符串、配置值、端口号
- 禁止直接提交代码, 仅生成commit message后提醒用户来提交即可
- 禁止直接提交需求文档、设计文档、实现计划等, 仅生成文档即可, 先进行后续步骤, 最后由用户来提交
