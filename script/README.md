# CoCopilot

## 有些东西越有效就越简单，比如这个小工具，有什么用我不能说。

###  本项目只是帮你授权，并没有数据/代码泄露风险。



1. 使用`JetBrains`全家桶IDE，安装`Copilot`插件。
2. `windows系统`执行`cocopilot.bat`
3. `macOS/linux系统`执行`cocopilot.sh`
4. 看到`done. please restart your ide.`表示成功。
5. 重启你的IDE就好。此方式对`Vim/NeoVim`亦有效。
6. 这是个小玩具，可能测试不充分，别找我。

### 对于`VSCode`，步骤基本相同，执行对应`vscode.sh`/`vscode.bat`（如果是使用vscode远程连接Ubuntu服务器且副驾驶拓展安装在了远程服务器上，需要执行 `vscode-remote.sh`），**无需执行**`cocopilot.sh`/`cocopilot.bat`。
### `VSCode`中插件更新后需要重新执行脚本，`JetBrains`则不需要。


## IDE Copilot 配置
**Github Copliot 和 CoCopilot 都适用**

### Jetbrains IDE

将`vscode-proxy.bat`脚本中的 `https://xxxx/copilot_internal/v2/token` 替换为的 `https//<BaseURL>/copilot_internal/v2/token`

Windows 示例：
```shell
@echo off
setlocal

set folder=%userprofile%\AppData\Local\github-copilot
set jsonfile=%folder%\hosts.json

if not exist "%folder%" (
    mkdir "%folder%"
)

echo {"github.com":{"user":"cocopilot","oauth_token":"ccu_Vxxxxx","dev_override":{"copilot_token_url":"https://<BaseURL>/cocopilot/copilot_internal/v2/token"}}} > "%jsonfile%"
echo done. please restart your ide.
pause
```

### Vscode
插件下载：https://cocopilot.org/static/assets/files/cocopilot_scripts.zip
1. vscode 中安装 `Github Copilot` 插件, 必须完全重启 vscode！！
2. 安装 `CoCopilot` 插件
3. 使用 `CoCopilot` 插件登录授权
4. 将`vscode-proxy.bat`脚本中的 `set host=https://proxy.fungs.cn` 替换为的 `set host=<BaseURL>`
5. 使用 `vscode-proxy.bat` 脚本授权
6. 重启 vscode

## 贡献者们

> 如果你刚好有路子、头也铁，可以[进行投喂](https://zhile.io/contribute-copilot-token)造福大家哦。

> 感谢所有让这个项目变得更好的贡献者们！
