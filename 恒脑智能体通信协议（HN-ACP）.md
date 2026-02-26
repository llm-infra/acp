# 恒脑智能体通信协议（HN-ACP）

# 一、概述

恒脑智能体通信协议（**HN-ACP：Agent Conversation Protocol**）是一种智能体交互协议，旨在解决各类客户端与恒脑智能体的连接问题。

## 1.1 核心原则

1.  常规传输：使用标准HTTP+SSE（Server-Sent Events），ACP服务端以多个REST接口提供完整的协议服务。
    
2.  事件驱动：服务端产生协议要求的标准事件，客户端对不同事件响应处理。
    
3.  增量消息：实时流式传输中，内容必须以增量（delta）追加，不允许输出空内容。
    

## 1.2 主要目标

1.  支持同步与异步调用
    
2.  支持流式与非流式调用
    
3.  支持多智能体嵌套消息
    
4.  支持多分支并行消息
    
5.  支持多模态消息
    
6.  支持中断交互
    
7.  支持工作流+L3
    

# 二、协议

## 2.1 核心操作

### 2.1.1 发送对话-同步sse

```shell
POST /conversation/:id/completion
```

### 2.1.2 订阅对话-异步webhook

```shell
POST /conversation/:id/completion/sub
```

### 2.1.3 恢复对话

```shell
POST /conversation/:id/resume
```

### 2.1.4 取消对话

```shell
POST /conversation/:id/cancel
```

### 2.1.5 查询会话

```shell
GET /conversation/:id
```

### 2.1.6 列举会话

```shell
GET /conversations
```

## 3.1 数据结构

ACP基于流式事件（Event）推送实现智能体服务端与客户端的信息交互，所有事件都按照共同结构规范：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | 事件类型 |
| timestamp | Int64 | 毫秒时间戳 |
| ... |  | 不同事件自定义字段 |

事件类型列表：

| 事件 | 定义 |
| --- | --- |
| 运行开始 | run\_started |
| 运行完成 | run\_finished |
| 运行错误 | run\_error |
| 区块开始 | block\_start |
| 区块结束 | block\_end |
| 内容开始 | content\_start |
| 内容消息 | content\_delta |
| 内容结束 | content\_end |

*   运行事件代表智能体运行的生命周期，一次智能体运行必须以run\_started开始，以run\_finished/run\_error结束；
    
*   一次智能体运行周期内，以多个区块事件对构成（block\_start/end），表示多步骤执行流程，每个区块都有唯一区块ID以及关联父级区块信息，以此来表示区块之间的并行、缩进关系；
    
*   在每个区块内，由多个内容事件对构成（content\_start/delta/end），协议内置多种content消息类型：文本消息、思考消息、工具消息、制品消息等。
    

![image.png](https://alidocs.oss-cn-zhangjiakou.aliyuncs.com/res/2M9qP5j0pRgpmO01/img/5106610c-6d5c-4327-9384-cd686c26024c.png)

### 3.1.1 运行事件

在会话中进行一次智能体对话（completion），仅生成一对运行事件，例如：

```shell
{
    "type": "run_started",
    "timestamp": 1764744729334,
    "run_id": "19892ca0-46e1-4aed-9faa-1b709a33c2d2"
}

{
    "type": "run_finished",
    "timestamp": 1764744736681,
    "run_id": "19892ca0-46e1-4aed-9faa-1b709a33c2d2"
}
```

#### run\_started/run\_finished

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| run\_id | String | 运行ID，UUID |

#### run\_error

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| run\_id | String | 运行ID，UUID |
| error | String | 失败信息 |

### 3.1.2 区块事件

在会话中进行一次智能体对话（completion），可以生成多对区块事件，不同的区块之间通过事件元素关联，例如：

```shell
# 区块A
{
    "type": "block_start",
    "timestamp": 1764744729381,
    "block_id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522"
}
{
    "type": "block_end",
    "timestamp": 1764744736681,
    "usage": {
        "prompt_tokens": 6711,
        "completion_tokens": 119
    },
    "block_id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522"
}

# 区块A1
{

    "type": "block_start",
    "timestamp": 1764744729399,
    "block_id": "4c0b054f-ccad-4a25-91a7-8c1616cdf5ee",
    "is_subagent": true,
    "parent_block_id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522"
}
{
    "type": "block_end",
    "timestamp": 1764744736681,
    "block_id": "4c0b054f-ccad-4a25-91a7-8c1616cdf5ee"
}

# 区块B
{
    "type": "block_start",
    "timestamp": 1764744729399,
    "block_id": "ed3f5be9-ea8e-47d9-b8e5-9dab420d8a0d",
    "is_parallel": true,
    "parent_block_id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522"
}
{
    "type": "block_end",
    "timestamp": 1764744736681,
    "block_id": "ed3f5be9-ea8e-47d9-b8e5-9dab420d8a0d"
}

# 区块C
{
    "type": "block_start",
    "timestamp": 1764744729399,
    "block_id": "89050c8c-d83d-42be-8622-09350a74713c",
    "is_parallel": true,
    "parent_block_id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522"
}
{
    "type": "block_end",
    "timestamp": 1764744736681,
    "block_id": "89050c8c-d83d-42be-8622-09350a74713c"
}

```

#### block\_start

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| block\_id | String | 区块ID，UUID |
| is\_parallel | Bool | 是否并行块 |
| is\_subagent | Bool | 是否为子智能体块 |
| metadata | map\[string\]any | 附加信息，输出块的名称、描述、所属者等信息 |
| parent\_block\_id | String | 父区块ID |

#### block\_end

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| block\_id | String | 区块ID，UUID |
| usage | Object | Token用量 |
| usage.prompt\_tokens | Int64 | 输入Token |
| usage.completion\_tokens | Int64 | 输出Token |

### 3.1.3 内容事件

在每个区块内，包含多组内容事件，这是最终面向用户的内容输出，例如：

```shell
# 思考内容
{
    "type": "content_start",
    "timestamp": 1764744733137,
    "content_id": "f9f6a58a-7279-45bd-a0fe-1a59250d3f3f",
    "block_id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522"
}
{
    "type": "content_delta",
    "timestamp": 1764744733138,
    "content": {
        "type": "thinking",
        "delta": "**Counting "
    },
    "content_id": "f9f6a58a-7279-45bd-a0fe-1a59250d3f3f"
}
{
    "type": "content_delta",
    "timestamp": 1764744733149,
    "content": {
        "type": "thinking",
        "delta": "files in "
    },
    "content_id": "f9f6a58a-7279-45bd-a0fe-1a59250d3f3f"
}
{
    "type": "content_delta",
    "timestamp": 1764744733213,
    "content": {
        "type": "thinking",
        "delta": "directory**"
    },
    "content_id": "f9f6a58a-7279-45bd-a0fe-1a59250d3f3f"
}
{
    "type": "content_end",
    "timestamp": 1764744733137,
    "content_id": "f9f6a58a-7279-45bd-a0fe-1a59250d3f3f"
}

# 命令执行
{
    "type": "content_start",
    "timestamp": 1764744736656,
    "content_id": "31de563d-fa1f-4764-b66b-48b35ef5c38e",
    "block_id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522"
}
{
    "type": "content_delta",
    "timestamp": 1764744736656,
    "content": {
        "type": "command_execution",
        "delta": "/bin/sh ",
        "call_id": "item_1"
    },
    "content_id": "31de563d-fa1f-4764-b66b-48b35ef5c38e"
}
{
    "type": "content_delta",
    "timestamp": 1764744736657,
    "content": {
        "type": "command_execution",
        "delta": "-lc 'ls -1A | wc -l'",
        "call_id": "item_1"
    },
    "content_id": "31de563d-fa1f-4764-b66b-48b35ef5c38e"
}
{
    "type": "content_delta",
    "timestamp": 1764744736659,
    "content": {
        "type": "command_result",
        "delta": "3\n",
        "exit_code": 0,
        "call_id": "item_1"
    },
    "content_id": "31de563d-fa1f-4764-b66b-48b35ef5c38e"
}
{
    "type": "content_end",
    "timestamp": 1764744736659,
    "content_id": "31de563d-fa1f-4764-b66b-48b35ef5c38e"
}
```

#### content\_start

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| content\_id | String | 内容ID，UUID |
| related\_block\_id | String | 关联的区块ID |

#### content\_end

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| content\_id | String | 内容ID，UUID |

#### content\_delta

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| content\_id | String | 内容ID，UUID |
| content | Object | 内容消息 |

内容消息包含以下类型，根据不同的类别客户端分类展示：

| 内容类别 | 枚举变量 |
| --- | --- |
| 文本消息 | text |
| 思考消息 | thinking |
| 工具调用 | tool\_call |
| 工具参数 | tool\_args |
| 工具输出 | tool\_result |
| 文件输出 | file |
| 数据输出 | data |
| 制品输出 | artifact |
| 变量消息 | variable |
| 交互消息 | interaction |
| 自定义消息 | custom |
|  |  |
| MCP调用 | mcp\_call |
| MCP参数 | mcp\_args |
| MCP结果 | mcp\_result |
| 命令执行 | command\_execution |
| 命令执行结果 | command\_execution\_result |
| 网络搜索 | web\_search |
| 网络搜索结果 | web\_search\_result |
| 代办列表 | todo\_list |

### 3.1.4 会话记录

会话是智能体的运行基础，每个会话都拥有独立b的ID，每个会话中包含多轮用户与智能体的输入输出。

```json
{
    "id": "c2132979-340e-4d0d-b680-58d4d2ed52ee",
    "title": "session A",
    "created_at": 1764744728784251,
    "updated_at": 1764744736684511,
    "messages": [
        {
            "id": "019ad7e9-9b48-76c8-8e6c-a688cdea0799",
            "role": "user",
            "blocks": [
                {
                    "id": "019ad7e9-9b48-76c8-8e6c-a689198651e2",
                    "contents": [
                        {
                            "type": "text",
                            "text": "查询目录下有多少文件"
                        }
                    ]
                }
            ],
            "created_at": 1764744728785250,
            "updated_at": 1764744736684511
        },
        {
            "id": "19892ca0-46e1-4aed-9faa-1b709a33c2d2",
            "role": "assistant",
            "blocks": [
                {
                    "id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522",
                    "contents": [
                        {
                            "type": "thinking",
                            "text": "**Counting files in directory**"
                        },
                        {
                            "type": "command_execution",
                            "call_id": "item_1",
                            "command": "/bin/sh -lc 'ls -1A | wc -l'",
                            "result": "3\n",
                            "exit_code": 0
                        },
                        {
                            "type": "thinking",
                            "text": "**Preparing to respond**"
                        },
                        {
                            "type": "text",
                            "text": "当前目录下共有 3 个文件/子项。"
                        }
                    ],
                    "usage": {
                        "prompt_tokens": 6711,
                        "completion_tokens": 119
                    }
                }
            ],
            "created_at": 1764744728785250,
            "updated_at": 1764744736684511
        }
    ]
}
```

#### conversation

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | String | 会话ID |
| title | String | 会话标题 |
| created\_at | Int64 | 会话创建时间 |
| updated\_at | Int64 | 会话更新时间 |
| messages | List Object | 消息列表 |

#### message

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | String | 消息ID，非user消息的ID就是智能体执行的run\_id |
| role | String | 角色，可选user/assistant |
| created\_at | Int64 | 消息创建时间 |
| updated\_at | Int64 | 消息更新时间 |
| blocks | List Object | 区块列表 |

#### block

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| id | String | 区块ID |
| contents | List Object | 内容列表 |
| is\_parallel | Bool | 是否并行块 |
| is\_subagent | Bool | 是否为子智能体块 |
| metadata | map\[string\]any | 附加信息，输出块的名称、描述、所属者等信息 |
| parent\_block\_id | String | 父区块ID |

### 3.1.5 内容类别

:::
所有的内容结构都包含type字段，通过type指示具体的内容类型。

相同的内容类型，在事件流式输出与会话存储结构设计上有一定差异。

流式输出需要遵循delta输出模式，会话存储中是完整的内容。
:::

#### 文本⭐

*   流式输出
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | text |
| delta | String | 增量文本内容 |

```json
{
    "type": "content_delta",
    "timestamp": 1764835191439,
    "content_id": "9d8561e9-5d07-4208-8fda-33da6fb4ca55",
    "content": {
        "type": "text",
        "delta": "Hi there!"
    }
}
{
    "type": "content_delta",
    "timestamp": 1764835191439,
    "content_id": "9d8561e9-5d07-4208-8fda-33da6fb4ca55",
    "content": {
        "type": "text",
        "delta": " What's up?"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | text |
| text | String | 完整文本内容 |

```json
{
    "id": "19892ca0-46e1-4aed-9faa-1b709a33c2d2",
    "role": "assistant",
    "blocks": [
        {
            "id": "b76e7ead-5d2e-4f1e-ac10-2a356b418522",
            "contents": [
                {
                    "type": "text",
                    "text": "Hi there! What's up?"
                }
            ],
            "usage": {
                "prompt_tokens": 6711,
                "completion_tokens": 119
            }
        }
    ],
    "created_at": 1764744728785250,
    "updated_at": 1764744736684511
}
```

#### 思考⭐

*   流式输出
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | thinking |
| delta | String | 增量文本内容 |

```json
{
    "type": "content_delta",
    "timestamp": 1765333625806,
    "content_id": "9f69add2-d2dc-4af7-9555-c1bebee76eba",
    "content": {
        "type": "thinking",
        "delta": "嗯"
    }
}
{
    "type": "content_delta",
    "timestamp": 1765333625848,
    "content_id": "9f69add2-d2dc-4af7-9555-c1bebee76eba",
    "content": {
        "type": "thinking",
        "delta": "，"
    }
}
{
    "type": "content_delta",
    "timestamp": 1765333625894,
    "content_id": "9f69add2-d2dc-4af7-9555-c1bebee76eba",
    "content": {
        "type": "thinking",
        "delta": "用户"
    }
}
{
    "type": "content_delta",
    "timestamp": 1765333625944,
    "content_id": "9f69add2-d2dc-4af7-9555-c1bebee76eba",
    "content": {
        "type": "thinking",
        "delta": "问"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | thinking |
| text | String | 完整文本内容 |

```json
{
    "id": "30f811d5-c2ee-4252-8cac-10bb72894431",
    "role": "assistant",
    "blocks": [
        {
            "id": "86cbfebc-2576-4d59-a15f-9f89d35f0b2a",
            "contents": [
                {
                    "type": "thinking",
                    "text": "嗯，用户问“你是谁”，我需要详细分析这个问题。首先，我需要确定用户的需求是什么。他们可能想了解我的基本功能、用途，或者更深入的技术细节。接下来，我应该考虑用户的背景。他们可能是普通用户，想了解AI助手的能力，或者是开发者，想了解技术实现。然后，我需要组织回答的结构，确保涵盖关键点：身份、功能、技术基础、应用范围，以及互动方式。同时，要保持回答简洁易懂，避免使用过于专业的术语，但也要准确。还要注意用户可能的后续问题，比如询问具体功能或使用场景，所以回答中可以适当引导进一步交流。最后，确保语气友好，符合阿里巴巴的品牌形象。\n"
                }
            ],
            "parent_block_id": "5d88b0cb-b430-4d92-8c30-4c2fb83e8c7b"
        }
    ],
    "created_at": 1765333609218139,
    "updated_at": 1765333645188532
}
```

#### 工具⭐

*   流式输出
    

工具调用

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | tool\_call |
| tool\_name | String | 工具名称 |

工具参数

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | tool\_args |
| delta | String | 增量参数文本 |

工具结果

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | tool\_result |
| delta | String | 增量结果文本 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "type": "content_delta",
    "timestamp": 1765335449611,
    "content_id": "67a21f49-611a-485c-ad3d-d09690ccbc74",
    "content": {
        "type": "tool_call",
        "tool_name": "APT组织uuid查询"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765335449611,
    "content_id": "67a21f49-611a-485c-ad3d-d09690ccbc74",
    "content": {
        "type": "tool_args",
        "delta": "{\"name\":"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765335449611,
    "content_id": "67a21f49-611a-485c-ad3d-d09690ccbc74",
    "content": {
        "type": "tool_args",
        "delta": "\"Lazarus\"}"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765335449699,
    "content_id": "67a21f49-611a-485c-ad3d-d09690ccbc74",
    "content": {
        "type": "tool_result",
        "delta": "{\n  \"apt_name\": \"Lazarus\",\n  \"uuid\": \"dcdf3bc7-307f-11eb-9593-ac1f6b480078\"\n}\n"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | tool\_call |
| tool\_name | String | 工具名称 |
| tool\_args | String | 工具参数 |
| tool\_result | String | 工具输出 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "id": "85905aeb-6fb6-442e-a96b-487eb43e3b2d",
    "role": "assistant",
    "blocks": [
        {
            "id": "ed318d06-c629-47e4-bfcd-1b5ec1e810ef",
            "contents": [
                {
                    "type": "tool_call",
                    "tool_name": "APT组织uuid查询",
                    "tool_args": "{\"name\":\"Lazarus\"}",
                    "tool_result": "{\n  \"apt_name\": \"Lazarus\",\n  \"uuid\": \"dcdf3bc7-307f-11eb-9593-ac1f6b480078\"\n}\n"
                }
            ]
        }
    ],
    "created_at": 1765335448266004,
    "updated_at": 1765335451824690
}
```

#### 文件⭐

*   流式输出
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | file |
| mime\_type | String | 文件媒体类型 |
| file\_id | String | 文件ID |

```json
{
    "type": "content_delta",
    "timestamp": 1765866232089,
    "content_id": "ef9d399e-105b-4e68-a92f-a2a2d9be489c",
    "content": {
        "type": "file",
        "mime_type": "text/markdown; charset=utf-8",
        "file_id": "file:test.md<52cbfb86-53dd-4053-a160-ff61a43f4b22>"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | file |
| mime\_type | String | 文件媒体类型 |
| file\_id | String | 文件ID |

```json
{
    "id": "97e62391-16cc-4998-ad37-0135cca7394f",
    "role": "assistant",
    "blocks": [
        {
            "id": "fcb1b6a1-0923-4a1f-9c18-001e0d969eac",
            "contents": [
                {
                    "type": "file",
                    "mime_type": "text/markdown; charset=utf-8",
                    "file_id": "file:test.md\u003c52cbfb86-53dd-4053-a160-ff61a43f4b22\u003e"
                }
            ]
        }
    ],
    "created_at": 1765866209188020,
    "updated_at": 1765866232090928
}
```

#### 数据⭐

*   流式输出
    
    | 字段 | 类型 | 说明 |
    | --- | --- | --- |
    | type | String | data |
    | mime\_type | String | 制品媒体类型 |
    | delta | String | 增量数据Base64编码 |
    
    ```json
    {
        "type": "content_delta",
        "timestamp": 1765866232089,
        "content_id": "ef9d399e-105b-4e68-a92f-a2a2d9be489c",
        "content": {
            "type": "data",
            "mime_type": "text/markdown; charset=utf-8",
            "delta": "5L2g5aW9"
        }
    }
    
    {
        "type": "content_delta",
        "timestamp": 1765866232089,
        "content_id": "ef9d399e-105b-4e68-a92f-a2a2d9be489c",
        "content": {
            "type": "data",
            "mime_type": "text/markdown; charset=utf-8",
            "delta": "5ZCX"
        }
    }
    ```
    
*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | data |
| mime\_type | String | 制品媒体类型 |
| data | String | 完整数据Base64编码 |

```json
{
    "id": "97e62391-16cc-4998-ad37-0135cca7394f",
    "role": "assistant",
    "blocks": [
        {
            "id": "fcb1b6a1-0923-4a1f-9c18-001e0d969eac",
            "contents": [
                {
                    "type": "data",
                    "mime_type": "text/markdown; charset=utf-8",
                    "data": "5L2g5aW95ZCX"
                }
            ]
        }
    ],
    "created_at": 1765866209188020,
    "updated_at": 1765866232090928
}
```

#### 制品⭐

*   流式输出
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | artifact |
| mime\_type | String | 制品媒体类型 |
| file\_id | String | 制品ID |

```json
{
    "type": "content_delta",
    "timestamp": 1765866232089,
    "content_id": "ef9d399e-105b-4e68-a92f-a2a2d9be489c",
    "content": {
        "type": "artifact",
        "mime_type": "text/markdown; charset=utf-8",
        "file_id": "file:test.md<52cbfb86-53dd-4053-a160-ff61a43f4b22>"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | artifact |
| mime\_type | String | 制品媒体类型 |
| file\_id | String | 制品ID |

```json
{
    "id": "97e62391-16cc-4998-ad37-0135cca7394f",
    "role": "assistant",
    "blocks": [
        {
            "id": "fcb1b6a1-0923-4a1f-9c18-001e0d969eac",
            "contents": [
                {
                    "type": "artifact",
                    "mime_type": "text/markdown; charset=utf-8",
                    "file_id": "file:test.md\u003c52cbfb86-53dd-4053-a160-ff61a43f4b22\u003e"
                }
            ]
        }
    ],
    "created_at": 1765866209188020,
    "updated_at": 1765866232090928
}
```

#### 变量⭐

*   流式输出
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | variable |
| delta | Map\[String\]any | 变量Map |

```json
{
    "type": "content_delta",
    "timestamp": 1765866232089,
    "content_id": "ef9d399e-105b-4e68-a92f-a2a2d9be489c",
    "content": {
        "type": "variable",
        "delta": {
            "a": "1"
        }
    }
}

{
    "type": "content_delta",
    "timestamp": 1765866232089,
    "content_id": "ef9d399e-105b-4e68-a92f-a2a2d9be489c",
    "content": {
        "type": "variable",
        "delta": {
            "b": "2"
        }
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | variable |
| variables | Map\[String\]any | 变量Map |

```json
{
    "id": "97e62391-16cc-4998-ad37-0135cca7394f",
    "role": "assistant",
    "blocks": [
        {
            "id": "fcb1b6a1-0923-4a1f-9c18-001e0d969eac",
            "contents": [
                {
                    "type": "variable",
                    "variable": {
                        "a": "1",
                        "b": "2"
                    }
                }
            ]
        }
    ],
    "created_at": 1765866209188020,
    "updated_at": 1765866232090928
}
```

#### 交互⭐

交互消息用于承载用户可操作交互的卡片/表单等UI信息，基于[A2UI](https://a2ui.org/)消息结构，实现生成式UI交互回传。

*   流式输出
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | interaction |
| interaction\_id | String | 交互ID，用于resume关联 |
| a2ui\_version | String | A2UI版本 |
| a2ui\_message | Object | A2UI事件消息<br>[https://a2ui.org/reference/messages](https://a2ui.org/reference/messages) |

```json
{
  "type": "interaction",
  "a2ui_version": "0.8",
  "a2ui_message": {
    "surfaceUpdate": {
      "surfaceId": "booking",
      "components": [
        {
          "id": "root",
          "component": {
            "Column": {
              "children": {
                "explicitList": [
                  "header",
                  "guests-field",
                  "submit-btn"
                ]
              }
            }
          }
        },
        {
          "id": "header",
          "component": {
            "Text": {
              "text": {
                "literalString": "Confirm Reservation"
              },
              "usageHint": "h1"
            }
          }
        },
        {
          "id": "guests-field",
          "component": {
            "TextField": {
              "label": {
                "literalString": "Guests"
              },
              "text": {
                "path": "/reservation/guests"
              }
            }
          }
        },
        {
          "id": "submit-btn",
          "component": {
            "Button": {
              "child": "submit-text",
              "action": {
                "name": "confirm",
                "context": [
                  {
                    "key": "details",
                    "value": {
                      "path": "/reservation"
                    }
                  }
                ]
              }
            }
          }
        }
      ]
    }
  }
}

{
  "type": "interaction",
  "a2ui_version": "0.8",
  "a2ui_message": {
    "dataModelUpdate": {
      "surfaceId": "booking",
      "path": "/reservation",
      "contents": [
        {
          "key": "datetime",
          "valueString": "2025-12-16T19:00:00Z"
        },
        {
          "key": "guests",
          "valueString": "2"
        }
      ]
    }
  }
}

{
  "type": "interaction",
  "a2ui_version": "0.8",
  "a2ui_message": {
    "beginRendering": {
      "surfaceId": "booking",
      "root": "root"
    }
  }
}

{
  "type": "interaction",
  "a2ui_version": "0.8",
  "a2ui_message": {
    "deleteSurface": {
      "surfaceId": "booking"
    }
  }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | interaction |
| a2ui\_version | String | A2UI版本 |
| a2ui\_messages | List Object | A2UI事件消息列表<br>[https://a2ui.org/reference/messages](https://a2ui.org/reference/messages) |

```json
{
    "id": "97e62391-16cc-4998-ad37-0135cca7394f",
    "role": "assistant",
    "blocks": [
        {
            "id": "fcb1b6a1-0923-4a1f-9c18-001e0d969eac",
            "contents": [
                {
                    "type": "interaction",
                    "a2ui_version": "0.8",
                    "a2ui_messages": [
                        {
                            "surfaceUpdate": {
                                "surfaceId": "booking",
                                "components": [
                                    {
                                        "id": "root",
                                        "component": {
                                            "Column": {
                                                "children": {
                                                    "explicitList": [
                                                        "header",
                                                        "guests-field",
                                                        "submit-btn"
                                                    ]
                                                }
                                            }
                                        }
                                    },
                                    {
                                        "id": "header",
                                        "component": {
                                            "Text": {
                                                "text": {
                                                    "literalString": "Confirm Reservation"
                                                },
                                                "usageHint": "h1"
                                            }
                                        }
                                    },
                                    {
                                        "id": "guests-field",
                                        "component": {
                                            "TextField": {
                                                "label": {
                                                    "literalString": "Guests"
                                                },
                                                "text": {
                                                    "path": "/reservation/guests"
                                                }
                                            }
                                        }
                                    },
                                    {
                                        "id": "submit-btn",
                                        "component": {
                                            "Button": {
                                                "child": "submit-text",
                                                "action": {
                                                    "name": "confirm",
                                                    "context": [
                                                        {
                                                            "key": "details",
                                                            "value": {
                                                                "path": "/reservation"
                                                            }
                                                        }
                                                    ]
                                                }
                                            }
                                        }
                                    }
                                ]
                            }
                        },
                        {
                            "dataModelUpdate": {
                                "surfaceId": "booking",
                                "path": "/reservation",
                                "contents": [
                                    {
                                        "key": "datetime",
                                        "valueString": "2025-12-16T19:00:00Z"
                                    },
                                    {
                                        "key": "guests",
                                        "valueString": "2"
                                    }
                                ]
                            }
                        },
                        {
                            "beginRendering": {
                                "surfaceId": "booking",
                                "root": "root"
                            }
                        },
                        {
                            "deleteSurface": {
                                "surfaceId": "booking"
                            }
                        }
                    ]
                }
            ]
        }
    ],
    "created_at": 1765866209188020,
    "updated_at": 1765866232090928
}
```

#### 自定义

自定义消息用于承载协议未内置类型但业务侧需要透传的文本/结构化片段。

*   流式输出
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | custom |
| raw | String | 自定义全量内容 |

```json
{
    "type": "content_delta",
    "timestamp": 1765866232089,
    "content_id": "ef9d399e-105b-4e68-a92f-a2a2d9be489c",
    "content": {
        "type": "custom",
        "raw": "{\"step\":1,\"status\":\"running\"}"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | custom |
| raw | String | 完整自定义内容 |

```json
{
    "id": "97e62391-16cc-4998-ad37-0135cca7394f",
    "role": "assistant",
    "blocks": [
        {
            "id": "fcb1b6a1-0923-4a1f-9c18-001e0d969eac",
            "contents": [
                {
                    "type": "custom",
                    "raw": "{\"step\":1,\"status\":\"running\"}"
                }
            ]
        }
    ],
    "created_at": 1765866209188020,
    "updated_at": 1765866232090928
}
```

同一 `content_id` 下，`custom` 类型仅允许出现一次；若重复发送，服务端应返回内容结构错误。

#### MCP⭐

*   流式输出
    

MCP调用

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | mcp\_call |
| server | String | MCP服务名称 |
| tool\_name | String | MCP工具名称 |

MCP参数

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | mcp\_args |
| delta | String | 增量参数文本 |

MCP结果

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | mcp\_result |
| delta | String | 增量结果文本 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "type": "content_delta",
    "timestamp": 1765446799514,
    "content_id": "fbbfdc4d-a1a8-4899-ab8a-79b23593c53d",
    "content": {
        "type": "mcp_call",
        "server": "crypto-mcp",
        "tool_name": "md5"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765446799521,
    "content_id": "fbbfdc4d-a1a8-4899-ab8a-79b23593c53d",
    "content": {
        "type": "mcp_args",
        "delta": "{\"input\":\"123457\"}"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765446799570,
    "content_id": "fbbfdc4d-a1a8-4899-ab8a-79b23593c53d",
    "content": {
        "type": "mcp_result",
        "delta": "[{\"type\":\"text\",\"text\":\"f1887d3f9e6ee7a32fe5e76f4ab80d63\"}]"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | mcp\_call |
| server | String | MCP服务名称 |
| tool\_name | String | MCP工具名称 |
| tool\_args | String | MCP工具参数 |
| tool\_result | String | MCP工具输出 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "id": "6d7c0134-79ff-450d-ac91-ec4f63ed67e0",
    "role": "assistant",
    "blocks": [
        {
            "id": "b9a7f7f7-fd4e-415f-ac88-cadad0527390",
            "contents": [
                {
                    "type": "mcp_call",
                    "server": "crypto-mcp",
                    "tool_name": "md5",
                    "tool_args": "{\"input\":\"123457\"}",
                    "tool_result": "[{\"type\":\"text\",\"text\":\"f1887d3f9e6ee7a32fe5e76f4ab80d63\"}]"
                }
            ],
            "is_parallel": true,
            "parent_block_id": "194302b1-83c6-46eb-9fc3-d6154663db75"
        }

    ],
    "created_at": 1765446791274733,
    "updated_at": 1765446836790942
}
```

#### 代码执行⭐

*   流式输出
    

代码调用

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | code\_execution |
| lang | String | 语言：python、javascript |
| delta | String | 增量代码文本 |

结果返回

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | code\_execution\_result |
| delta | String | 增量结果文本 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "type": "content_delta",
    "timestamp": 1765333617113,
    "content_id": "1e9dea80-64c1-41d4-96f0-284ab0948a53",
    "content": {
        "type": "code_execution",
        "lang": "python",
        "delta": "import hashlib\n\n# Create a md5 hash object\nhash_object = hashlib.md5()\n\n# Define the string to be hashed\nstring = '93jds2ffsa'\n\n# Update the hash object with the string encoded as bytes\nhash_object.update(string.encode('utf-8'))\n\n# Get the hexadecimal digest of the hash\nhash_hex = hash_object.hexdigest()\n\nhash_hex"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765333623141,
    "content_id": "1e9dea80-64c1-41d4-96f0-284ab0948a53",
    "content": {
        "type": "code_execution_result",
        "delta": "{\"output\":\"5261141cd2d621d12a503346f6fb1a09\",\"console\":\"\"}"
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | code\_execution |
| lang | String | 语言：python、javascript |
| code | String | 代码 |
| result | String | 结果 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "id": "30f811d5-c2ee-4252-8cac-10bb72894431",
    "role": "assistant",
    "blocks": [
        {
            "id": "5d88b0cb-b430-4d92-8c30-4c2fb83e8c7b",
            "contents": [
                {
                    "type": "code_execution",
                    "lang": "python",
                    "code": "import hashlib\n\n# Create a md5 hash object\nhash_object = hashlib.md5()\n\n# Define the string to be hashed\nstring = '93jds2ffsa'\n\n# Update the hash object with the string encoded as bytes\nhash_object.update(string.encode('utf-8'))\n\n# Get the hexadecimal digest of the hash\nhash_hex = hash_object.hexdigest()\n\nhash_hex",
                    "result": "{\"output\":\"5261141cd2d621d12a503346f6fb1a09\",\"console\":\"\"}"
                }
            ],
            "parent_block_id": "658a0e80-df89-47a5-82a3-a48e10aab77e"
        }
    ],
    "created_at": 1765333609218139,
    "updated_at": 1765333645188532
}
```

#### 命令执行⭐

*   流式输出
    

命令调用

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | command\_execution |
| delta | String | 增量命令文本 |

命令结果

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | command\_execution\_result |
| delta | String | 增量命令结果文本 |
| exit\_code | Int | 命令退出码 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "type": "content_delta",
    "timestamp": 1765348870040,
    "content_id": "8689bf9d-f60d-4d53-be42-da496d190a87",
    "content": {
        "type": "command_execution",
        "delta": "pwd"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765348870041,
    "content_id": "8689bf9d-f60d-4d53-be42-da496d190a87",
    "content": {
        "type": "command_execution_result",
        "delta": "\"/home/mel2oo\"",
        "exit_code": 0
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | command\_execution |
| command | String | 命令 |
| result | String | 命令结果 |
| exit\_code | Int | 命令退出码 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "id": "dbf92954-0b45-4352-bfc3-42be1adc7c18",
    "role": "assistant",
    "blocks": [
        {
            "id": "71204faa-a08e-4332-9691-4370dddcbf76",
            "contents": [
                {
                    "type": "command_execution",
                    "command": "pwd",
                    "result": "\"/home/mel2oo\"",
                    "exit_code": 0
                }
            ],
            "usage": {
                "prompt_tokens": 12701,
                "completion_tokens": 91
            }
        }
    ],
    "created_at": 1765348862280859,
    "updated_at": 1765348870141019
}
```

#### 网络搜索⭐

*   流式输出
    

搜索调用

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | web\_search |
| delta | String | 增量命令文本 |

搜索结果

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | web\_search\_result |
| answer | String | 结果文本 |
| results | List Object | 搜索内容列表 |
| results\[\].title | String | 内容标题 |
| results\[\].url | String | 内容链接 |
| results\[\].snippet | String | 内容片段 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "type": "content_delta",
    "timestamp": 1765939809550,
    "content_id": "8a3dcaed-1ff8-4665-a568-6a4436cc0585",
    "content": {
        "type": "web_search",
        "delta": "2024年最佳React状态管理库推荐"
    }
}

{
    "type": "content_delta",
    "timestamp": 1765939810536,
    "content_id": "8a3dcaed-1ff8-4665-a568-6a4436cc0585",
    "content": {
        "type": "web_search_result",
        "answer": "我已经搜索到了一些结果",
        "results": [
            {
                "title": "2024 年学状态管理库？ : r/reactjs",
                "url": "https://www.reddit.com/r/reactjs/comments/1db5go3/learning_state_management_libraries_in_2024/?tl=zh-hans",
                "snippet": "@tanstack/react-query 主要是一个数据获取和缓存库，从我所看到的来看，它被用于大多数前端React 仓库的生产级别。它使数据获取代码非常干净，并解决了许多 ...Read more"
            },
            {
                "title": "2024 React 状态管理库对比",
                "url": "https://juejin.cn/post/7325069743144239145",
                "snippet": "# 2024 React 状态管理库对比. React 状态管理库的意义在于提供一种机制来集中管理和维护 React 应用中的状态，并且使得这些状态能够跨组件共享。随着应用的复杂度提升，组件之间共享状态和进行状态通信变得越来越困难，这时状态管理库就显得尤为重要。包括：. * 许多状态管理库如 Redux 提供了中间件和开发工具，帮助开发者更好地进行状态的追踪、调试和异步处理。. * 一些状态管理库..."
            },
            {
                "title": "聊一聊2024 年React 生态系统",
                "url": "https://cloud.tencent.com/developer/article/2381404",
                "snippet": "如果对状态机有特别的兴趣，XState 和Zag 也是不错的选择。 如果需要一个全局存储，但不满意Zustand 或Redux，Jotai、Recoil 或Nano Stores 等本地状态管理 ...Read more"
            },
            {
                "title": "React 生態系2024 年推薦總整理",
                "url": "https://codelove.tw/@tony/post/gqB053",
                "snippet": "🔴 **YT 直播問答！**每週二晚上開講，聊聊科技、軟體新聞！ ➡️ 訂閱 YouTube 頻道 ➡️ 加入 Discord 社群. # React 生態系 2024 年推薦總整理. # React 生態系 2024 年推薦總整理. 客戶端狀態管理是現代 Web 開發的一個重要方面，可以在前端應用程式中實現高效的資料處理。 Redux Toolkit 和 Zustand 是兩種流行的用戶端狀態..."
            },
            {
                "title": "2025/2026 年React 组件库与相关库推荐",
                "url": "https://zhuanlan.zhihu.com/p/546697951",
                "snippet": "⭐ TanStack Table : 无样式的table 操作, 支持所有UI 组件库; swr : Vercel - 请求状态管理库; pmndrs/drei : ( 封装react-three-fiber ); aidenybai/react-scan: 自动 ...Read more"
            },
            {
                "title": "2024热门的几个React状态管理库",
                "url": "https://juejin.cn/post/7390646355028377627",
                "snippet": "// B组件 import from \"recoil\" import from\"@/store/user\" export default function BComponent const useRecoilState const useRecoilState const changeUserNameval: string string setName const changeUserAgeval..."
            },
            {
                "title": "React 第三方状态管理库的比较与选择原创",
                "url": "https://blog.csdn.net/weixin_40629244/article/details/148659559",
                "snippet": "最新推荐文章于 2025-10-22 18:21:36 发布. 于 2025-06-15 00:04:24 发布. CC 4.0 BY-SA版权. 版权声明：本文为博主原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接和本声明。. import from 'mobx' import from'mobx-react' class CounterStore 0 construc..."
            },
            {
                "title": "2024 React 生态工具最能打的组合！",
                "url": "https://www.51cto.com/article/792984.html",
                "snippet": "# 2024 React 生态工具最能打的组合！. * Vite：适用于**客户端渲染**的 React 应用。. * Next.js：适用于**服务端渲染**的 React 应用。. * Astro：适用于**静态生成**的 React 应用。. Next.js 是一个成熟度很高的 React 框架，也是 React 官方推荐的创建新的 React 项目的首选方式。Next.js 凭借其丰富的内..."
            },
            {
                "title": "React 狀態管理套件比較與原理實現feat. Redux, Zustand ...",
                "url": "https://medium.com/%E6%89%8B%E5%AF%AB%E7%AD%86%E8%A8%98/react-%E5%90%84%E5%80%8B%E7%8B%80%E6%85%8B%E7%AE%A1%E7%90%86%E5%A5%97%E4%BB%B6%E6%AF%94%E8%BC%83%E8%88%87%E5%8E%9F%E7%90%86%E5%AF%A6%E7%8F%BE-ba61db07332b",
                "snippet": "Redux, Zustand, Jotai, Recoil, MobX, Valtio. 而我們從下載量來看目前是 Redux 跟 Zustand 的下載量最多，再來是 Mobx，而另外兩個實現 atomic 機制的 Jotai 跟 Recoil 每週下載量大約是 50 萬左右，最少人用的 Valtio 目前差不多是每週 30 萬。. 在 2021 年的之前 react-redux 還是使用 `u..."
            },
            {
                "title": "打造卓越UI：2024 年不容错过的9 个React UI 组件库",
                "url": "https://developer.aliyun.com/article/1627881",
                "snippet": "简介： 本文介绍了2024年最受欢迎的9个React UI组件库，每一个都在设计、功能和定制化上有独特的优势，包括Material UI、Ant Design、Chakra UI等。Read more"
            }
        ]
    }
}

// 错误示例
{
    "type": "content_delta",
    "timestamp": 1765939810536,
    "content_id": "8a3dcaed-1ff8-4665-a568-6a4436cc0585",
    "content": {
        "type": "web_search_result",
        "error": {
            "type": "API_RATE_LIMIT"
            "message": "API rate limit exceeded"
        }
    }
}
```

*   会话存储
    

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| type | String | web\_search |
| query | String | 搜索输入 |
| answer | String | 结果文本 |
| results | List Object | 搜索内容列表 |
| results\[\].title | String | 内容标题 |
| results\[\].url | String | 内容链接 |
| results\[\].snippet | String | 内容片段 |
| error | Object | 错误信息 |
| error.type | String | 错误类型 |
| error.message | String | 错误消息 |

```json
{
    "id": "8f68cafc-faaf-44c9-adf0-75646eae8efc",
    "role": "assistant",
    "blocks": [
        {
            "id": "15fd4706-171b-4995-9623-a3ffd00d75f0",
            "contents": [
                {
                    "type": "web_search",
                    "query": "2024年最佳React状态管理库推荐",
                    "answer": "我已经搜索到了一些结果",
                    "results": [
                        {
                            "title": "2024 年学状态管理库？ : r/reactjs",
                            "url": "https://www.reddit.com/r/reactjs/comments/1db5go3/learning_state_management_libraries_in_2024/?tl=zh-hans",
                            "snippet": "@tanstack/react-query 主要是一个数据获取和缓存库，从我所看到的来看，它被用于大多数前端React 仓库的生产级别。它使数据获取代码非常干净，并解决了许多 ...Read more"
                        },
                        {
                            "title": "2024 React 状态管理库对比",
                            "url": "https://juejin.cn/post/7325069743144239145",
                            "snippet": "# 2024 React 状态管理库对比. React 状态管理库的意义在于提供一种机制来集中管理和维护 React 应用中的状态，并且使得这些状态能够跨组件共享。随着应用的复杂度提升，组件之间共享状态和进行状态通信变得越来越困难，这时状态管理库就显得尤为重要。包括：. * 许多状态管理库如 Redux 提供了中间件和开发工具，帮助开发者更好地进行状态的追踪、调试和异步处理。. * 一些状态管理库..."
                        },
                        {
                            "title": "聊一聊2024 年React 生态系统",
                            "url": "https://cloud.tencent.com/developer/article/2381404",
                            "snippet": "如果对状态机有特别的兴趣，XState 和Zag 也是不错的选择。 如果需要一个全局存储，但不满意Zustand 或Redux，Jotai、Recoil 或Nano Stores 等本地状态管理 ...Read more"
                        },
                        {
                            "title": "React 生態系2024 年推薦總整理",
                            "url": "https://codelove.tw/@tony/post/gqB053",
                            "snippet": "🔴 **YT 直播問答！**每週二晚上開講，聊聊科技、軟體新聞！ ➡️ 訂閱 YouTube 頻道 ➡️ 加入 Discord 社群. # React 生態系 2024 年推薦總整理. # React 生態系 2024 年推薦總整理. 客戶端狀態管理是現代 Web 開發的一個重要方面，可以在前端應用程式中實現高效的資料處理。 Redux Toolkit 和 Zustand 是兩種流行的用戶端狀態..."
                        },
                        {
                            "title": "2025/2026 年React 组件库与相关库推荐",
                            "url": "https://zhuanlan.zhihu.com/p/546697951",
                            "snippet": "⭐ TanStack Table : 无样式的table 操作, 支持所有UI 组件库; swr : Vercel - 请求状态管理库; pmndrs/drei : ( 封装react-three-fiber ); aidenybai/react-scan: 自动 ...Read more"
                        },
                        {
                            "title": "2024热门的几个React状态管理库",
                            "url": "https://juejin.cn/post/7390646355028377627",
                            "snippet": "// B组件 import from \"recoil\" import from\"@/store/user\" export default function BComponent const useRecoilState const useRecoilState const changeUserNameval: string string setName const changeUserAgeval..."
                        },
                        {
                            "title": "React 第三方状态管理库的比较与选择原创",
                            "url": "https://blog.csdn.net/weixin_40629244/article/details/148659559",
                            "snippet": "最新推荐文章于 2025-10-22 18:21:36 发布. 于 2025-06-15 00:04:24 发布. CC 4.0 BY-SA版权. 版权声明：本文为博主原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接和本声明。. import from 'mobx' import from'mobx-react' class CounterStore 0 construc..."
                        },
                        {
                            "title": "2024 React 生态工具最能打的组合！",
                            "url": "https://www.51cto.com/article/792984.html",
                            "snippet": "# 2024 React 生态工具最能打的组合！. * Vite：适用于**客户端渲染**的 React 应用。. * Next.js：适用于**服务端渲染**的 React 应用。. * Astro：适用于**静态生成**的 React 应用。. Next.js 是一个成熟度很高的 React 框架，也是 React 官方推荐的创建新的 React 项目的首选方式。Next.js 凭借其丰富的内..."
                        },
                        {
                            "title": "React 狀態管理套件比較與原理實現feat. Redux, Zustand ...",
                            "url": "https://medium.com/%E6%89%8B%E5%AF%AB%E7%AD%86%E8%A8%98/react-%E5%90%84%E5%80%8B%E7%8B%80%E6%85%8B%E7%AE%A1%E7%90%86%E5%A5%97%E4%BB%B6%E6%AF%94%E8%BC%83%E8%88%87%E5%8E%9F%E7%90%86%E5%AF%A6%E7%8F%BE-ba61db07332b",
                            "snippet": "Redux, Zustand, Jotai, Recoil, MobX, Valtio. 而我們從下載量來看目前是 Redux 跟 Zustand 的下載量最多，再來是 Mobx，而另外兩個實現 atomic 機制的 Jotai 跟 Recoil 每週下載量大約是 50 萬左右，最少人用的 Valtio 目前差不多是每週 30 萬。. 在 2021 年的之前 react-redux 還是使用 `u..."
                        },
                        {
                            "title": "打造卓越UI：2024 年不容错过的9 个React UI 组件库",
                            "url": "https://developer.aliyun.com/article/1627881",
                            "snippet": "简介： 本文介绍了2024年最受欢迎的9个React UI组件库，每一个都在设计、功能和定制化上有独特的优势，包括Material UI、Ant Design、Chakra UI等。Read more"
                        }
                    ]
                }
            ],
            "usage": {
                "prompt_tokens": 16595,
                "completion_tokens": 601
            }
        }
    ],
    "created_at": 1765939787058084,
    "updated_at": 1765939826550000
}
```

#### 代办列表

## 三、参考协议

[https://a2ui.org/](https://a2ui.org/)

[https://docs.ag-ui.com/introduction](https://docs.ag-ui.com/introduction)

[https://a2a-protocol.org/latest/specification/](https://a2a-protocol.org/latest/specification/)

[https://agentcommunicationprotocol.dev/introduction/welcome](https://agentcommunicationprotocol.dev/introduction/welcome)

[https://agentnetworkprotocol.com/](https://agentnetworkprotocol.com/)
