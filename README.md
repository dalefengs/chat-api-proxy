# chat-api-proxy
![Static Badge](https://img.shields.io/badge/%3E=1.18-blue?label=Golang)
![Docker Pulls](https://img.shields.io/docker/pulls/dalefengs/chat-api-proxy?color=gold)

## 简单介绍 

* 更强大且高效的API转发工具,让你的开发更加便捷。 
* 支持多种转发方式：
    * Copilot2GPT4
    * CoCopilot2GPT4 (拼车版)
    * Gemini2ChatAPi
* 我们的工具支持流式转发，无需等待。
* Copilot 内置了 tokens 缓存，避免重复请求获取。
* 支持 Jetbrains IDE、VSCode、Vim/NeoVim 等编辑器代理使用 Copilot Chat。
* 部署 Copilot2GPT4 在服务器上同一IP下使用，预防动态IP被封号。
* Copilot2GPT4 谨慎多IP使用，多IP调用容易官方封号！(CoCopilot 拼车版同理)
* 可以启动为BackendAPI Proxy模式，直接使用Access Token调用/backend-api/和chat2api的接口。（即将加入）


## Docker 部署
```shell
docker run -d \
  --name chat-api-proxy \
  -p 18818:8818 \
  -v LOG_LEVEL=info \
  dalefengs/chat-api-proxy
```

### 环境变量
* PROXY_API_PREFIX API路由前缀
* LOG_LEVEL 日志等级（debug, info, error）
* GEMINI_BASE_URL Gemini pro 自定义代理地址
* GEMINI_VERSION Gemini pro Api 版本(v1, v1beta等)

## API 文档

### Github Copilot
#### 获取 Copilot Token
**<BaseURL>/copilot/copilot_internal/v2/token**
```shell
curl --location 'http://127.0.0.1:18818/copilot/copilot_internal/v2/token' \
--header 'Authorization: token ghu_xxxxxxxxxxxxxxxx'
```
* 接口为GET请求。
* Authorization: token 后面为你获取到的 Github Token。

#### Completions 官方接口
**<BaseURL>/chat/completions**
```shell
curl --location 'http://127.0.0.1:18818/chat/completions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer tid=xxxxxxxx;ol=xxb457a0be36d3;exp=1705149290;sku=copilot_for_business_seat;st=dotcom;ssc=1;chat=1;sn=1;8kp=1:18b175a4e4bbf73xx3a627e7180a6469540d8316884d4ea6713edb28' \
--data '{
    "stream": true,
    "model": "gpt-4",
    "messages": [
        {
            "role": "user",
            "content": "你好"
        }
    ]
}'
```
* 接口为POST请求。
* Authorization: Bearer 后面的token为你获取到的 Copilot Token。


### CoCopilot (拼车版)


#### 获取 Copilot Token
**<BaseURL>/cocopilot/copilot_internal/v2/token**
```shell
curl --location 'http://127.0.0.1:18818/cocopilot/copilot_internal/v2/token' \
--header 'Authorization: token cuu_xxxxxxxxxxxxxxxx'
```
* 接口为GET请求。
* cocopilot：https://cocopilot.org/dash
* Authorization: token 后面为你获取到的 cocopilot Token。

**⚠拿到的 ccu-xxx ，一定要小心保存，防止泄露❗❗❗**


#### Completions 接口
**<BaseURL>/v1/chat/completions**
```shell
curl --location 'http://127.0.0.1:18818/v1/chat/completions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer cuu_xxxxxxxxxxxxxxxx' \
--data '{
    "stream": true,
    "model": "gpt-4",
    "messages": [
        {
            "role": "user",
            "content": "你好"
        }
    ]
}'
```
* 接口为POST请求。
* Authorization: Bearer 后面的token为你获取到的 CoCopilot Token。


## Gemini2ChatAPi
### completions 接口
```shell
curl --location 'http://127.0.0.1:18818/gemini/v1/chat/completions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer AIzaSyCGpzxxxxxxxxxx' \
--data '{
    "stream": true,
    "model": "gemini-pro",
    "messages": [
        {
            "role": "user",
            "content": "你好"
        }
    ]
}'
```
* 接口为POST请求。
* Authorization: Bearer 后面的token为你获取到的 Copilot Token。
* model 为 gemini-pro 或 自定义的模型名称。


## IDE Copilot 配置  
[查看文档](https://github.com/dalefengs/chat-api-proxy/tree/main/script)


## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=dalefengs/chat-api-proxy&type=Date)](https://star-history.com/#dalefengs/chat-api-proxy&Date)
