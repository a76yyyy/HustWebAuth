HustWebAuth
===========

锐捷 WEB 认证跨平台工具

Web认证
----------
当前很多学校使用了锐捷 Web 认证的方式来实现校园内有线网和无线网的网络认证和登录。

锐捷 Web 认证是一种对用户访问网络的权限进行控制的身份认证方法，这种认证方法不需要用户安装专用的客户端认证软件，使用**普通的浏览器访问**就可以进行身份认证。

未认证用户使用浏览器上网时，网络设备会强制浏览器访问特定站点，也就是Web认证服务器，通常称为Portal服务器。当用户需要访问认证服务器以外的其它网络资源时，就必须通过浏览器在Portal服务器上进行身份认证，只有认证通过后才可以使用网络资源。参考：[Web认证概述](https://image.ruijie.com.cn/Upload/Article/fd9117df-4b38-49fb-a6ac-a8b6cb43a130/RAC&RAP%20%E5%AE%9E%E6%96%BD%E4%B8%80%E6%9C%AC%E9%80%9A%EF%BC%88%E5%B0%8F%E7%9D%BF%E5%93%A5%EF%BC%89/RAC&RAP%20%E5%AE%9E%E6%96%BD%E4%B8%80%E6%9C%AC%E9%80%9A%EF%BC%88%E5%B0%8F%E7%9D%BF%E5%93%A5%EF%BC%89/8/1/Web%E8%AE%A4%E8%AF%81%E5%8E%9F%E7%90%86.html)

锐捷 WEB 认证跨平台工具的目的是方便用户**使用 Linux 或配置 OpenWrt 路由器**时，在**不打开浏览器的前提**下，通过**命令行**直接实现 Web 认证。

使用方法
-----------
1. 安装
    > 方式一: 使用 `go install` 命令安装
    >
    > ```bash
    > go install github.com/a76yyyy/HustWebAuth@latest
    > ```
    >
    > 方式二: 下载 `release` 可执行文件
    >
    > 1. 下载指定架构的[可执行文件](https://github.com/a76yyyy/HustWebAuth/releases)
    > 2. 重命名可执行文件为 `HustWebAuth` 或 `HustWebAuth.exe`
    > 3. 将文件权限修改为可执行权限, 如 `chmod +x HustWebAuth`
    > 4. **(建议)** 将可执行文件移动到 `/usr/local/bin`目录下 或 添加到`Path`环境变量中


2. 命令行运行 `HustWebAuth -h` 查看帮助
3. 命令行运行 `HustWebAuth -a account -p password` 进行认证测试

    > ### Tips:
    >
    > 1. 请确保你的账号密码正确
    > 2. 请确保你的网络连接正常, 且使用锐捷 Web 认证方式
    > 3. 可使用 `HustWebAuth -a account -p password -o` 进行认证并保存配置文件至 `$HOME` 文件夹下
    > 4. 可使用 `HustWebAuth login -r` 开启无感认证, 需提前下线你的设备

4. **(可选)** 使用 `HustWebAuth service install` 安装系统服务

    > ### Tips:
    >
    > 1. 请确保你的配置文件 `HustWebAuth.yaml` 已正确写入至 `$HOME` 文件夹下
    > 2. Windows 系统请在配置文件中 `log` 选项下设置日志文件名以方便查看日志, 如 `File: "HustWebAuth.log"`
    > 3. 建议以服务方式运行时, `log` 选项下设置 `connected` 为 `false` 以避免无效信息导致日志过大

Help 命令
==========
```bash
> HustWebAuth -h
HustWebAuth is a program used to implement Ruijie web authentication.

Usage:
  HustWebAuth [flags]
  HustWebAuth [command]

Available Commands:
  get         Get the login url from the redirect url
  help        Help about any command
  login       Hust web auth only once
  service     System service related commands

Flags:
  -a, --account string           Account for ruijie web authentication
  -f, --config string            Config file (default is $HOME/HustWebAuth.yaml)
  -c, --cycle                    Enable cycle mode
      --cycleDuration duration   Cycle duration (default 5m0s)
      --cycleRetry int           Cycle retry times, -1 means retry forever (default 3)
  -d, --daemon                   Enable daemon mode, not support windows
      --daemonPidFile string     Daemon pid file
  -e, --encrypt bool             Password is encrypted or not(default false)
  -h, --help                     help for main.exe
      --logAppend                Log file append mode.
                                 NOTE: if logRandom is true, it will be ignored (default true)
      --logConnected             Enable logging of "The network is connected" (default true)
      --logDir string            Log Directory (default Temp/HustWebAuth)
  -l, --logFile string           Log file name (default means output to os.stdout)
      --logRandom                Log file name with random string.
                                 NOTE: If logFile includes a "*", the random string replaces the last "*".
                                  (default true)
  -p, --password string          Password for ruijie web authentication
      --pingCount int            ping count (default 3)
      --pingIP string            IP address to ping (default "202.114.0.131")
      --pingPrivilege            Sets the type of ping pinger will send.
                                 false means pinger will send an "unprivileged" UDP ping.
                                 true means pinger will send a "privileged" raw ICMP ping.
                                 NOTE: setting to true requires that it be run with super-user privileges.
                                  (default true)
      --pingTimeout duration     Ping timeout (default 3s)
      --redirectURL string       Redirect URL (default "http://123.123.123.123")
  -o, --save                     Save config file
  -s, --serviceType string       Service Type, options: [internet, local] (default "internet")
      --syslog                   Enable syslog, not support windows

Use "HustWebAuth [command] --help" for more information about a command.
```

Service Help 命令
=================
```bash
> HustWebAuth service -h
Use HustWebAuth as a system service: install, start, stop, uninstall, etc.

Usage:
  HustWebAuth service [flags]
  HustWebAuth service [command]

Available Commands:
  install     Install HustWebAuth service
  restart     Restart HustWebAuth service
  start       Start HustWebAuth service
  status      Get HustWebAuth service status
  stop        Stop HustWebAuth service
  uninstall   Uninstall HustWebAuth service from system

Flags:
  -h, --help   help for service

Use "HustWebAuth service [command] --help" for more information about a command.
```