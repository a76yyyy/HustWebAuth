# ruijie_weblogin
锐捷WEB认证跨平台脚本

## 使用方法
[下载](https://github.com/a76yyyy/ruijie_weblogin/releases)
```bash
ruijie_web_login is a program used to implement Ruijie web authentication.

Usage:
  ruijie_weblogin_{os}_{arch} [flags]
  ruijie_weblogin_{os}_{arch} [command]

Available Commands:
  get         Get the login url from the redirect url.
  help        Help about any command
  login       Ruijie web login only once

Flags:
  -a, --account string           Account
  -f, --config string            Config file (default is localDir/ruijie.yaml)
  -c, --cycle                    Enable cycle mode
      --cycleDuration duration   Cycle duration (default 5m0s)
      --cycleRetry int           Cycle retry times, -1 means retry forever (default 3)
  -d, --daemon                   Enable daemon mode, not support windows
      --daemonPidFile string     Daemon pid file
  -h, --help                     help for main.exe
      --logAppend                Log file append mode.
                                 NOTE: if logRandom is true, it will be ignored (default true)
      --logConnected             Enable logging of "The network is connected" (default true)
      --logDir string            Log Directory (default Temp/ruijie)
  -l, --logFile string           Log file name (default means output to os.stdout)
      --logRandom                Log file name with random string.
                                 NOTE: If logFile includes a "*", the random string replaces the last "*".
                                  (default true)
  -p, --password string          Password
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
  -s, --service string           Service, options: [internet, local] (default "internet")
      --syslog                   Enable syslog, not support windows

Use "ruijie_weblogin_{os}_{arch} [command] --help" for more information about a command.
```

`internet`: 互联网

`local`: 内网

`register`: 开启无感认证
