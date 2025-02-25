# webssh
![](https://img.shields.io/github/v/release/Jrohy/webssh.svg) 
![](https://img.shields.io/docker/pulls/jrohy/webssh.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/Jrohy/webssh)](https://goreportcard.com/report/github.com/Jrohy/webssh)
[![Downloads](https://img.shields.io/github/downloads/Jrohy/webssh/total.svg)](https://img.shields.io/github/downloads/Jrohy/webssh/total.svg) 
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)   
简易在线ssh和sftp工具, 可在线敲命令和上传下载文件

## 运行截图
![avatar](asset/1.png)
![avatar](asset/2.png)

## 命令行
```
Usage of ./webssh_linux_amd64:
  -a string
        开启账号密码登录验证, '-a user:pass'的格式传参
  -p int
        服务运行端口 (default 5032)
  -t int
        ssh连接超时时间(min) (default 120)
  -s    保存ssh密码
  -v    显示版本号
```

## 原理
```
+---------+     http     +--------+    ssh    +-----------+
| browser | <==========> | webssh | <=======> | ssh server|
+---------+   websocket  +--------+    ssh    +-----------+
```

## 配置文件
在可执行文件同目录下创建 `config.json` 文件，可以设置默认的SSH连接信息：

```json
{
    "default_ssh": {
        "host": "192.168.1.100",
        "port": 22,
        "username": "root",
        "password": "",
        "login_type": 0,
        "key_passphrase": ""
    }
}
```

配置说明：
- `host`: 默认SSH主机地址
- `port`: 默认SSH端口
- `username`: 默认用户名
- `password`: 默认密码（可留空）
- `login_type`: 登录类型（0: 密码认证, 1: 私钥认证）
- `key_passphrase`: 私钥密码（如果私钥有密码保护）

## 运行
1. 下载[releases](https://github.com/Jrohy/webssh/releases)里不同平台的包来执行即可  

2. docker运行:  
    ```
    docker run -d --net=host --log-driver json-file --log-opt max-file=1 --log-opt max-size=100m --restart always --name webssh -e TZ=Asia/Shanghai jrohy/webssh
    ```
    支持添加的环境变量:
    ```
    port: web使用端口, 默认5032
    savePass: 是否保存密码, 默认true
    authInfo: 开启账号密码登录验证, 'user:pass'的格式设置
    ```
