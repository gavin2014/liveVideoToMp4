##m3u8直播流转储工具
自动下载m3u8直播流中的ts文件，并将ts文件合并，转换成mp4，注意不支持加密的数据流，工具比较简单

###使用说明
```cassandraql
转储ts文件到mp4.

Usage:
  直播流转储工具 [flags]
  直播流转储工具 [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  down        A brief description of your command
  help        Help about any command

Flags:
  -d, --duration int        持续时间小时 (default 9)
  -h, --help                help for 直播流转储工具
  -u, --m3u8-url string     m3u8的链接地址
  -m, --merge-time int      持续时间小时 (default 600)
  -o, --out string          mp4的保存路径 (default "./out")
  -s, --start-time string   开始时间 (default "08:00")
  -t, --temp string         缓存文件路径 (default "./temp")

Use "直播流转储工具 [command] --help" for more information about a command.
```
> demo   

`main.exe -s=08:00 -d=9 -m=600 -u=https://liveplaybdbj.wkbaobao.com/live/4CBD8F4622C7.m3u8 -o=./out/huodong1 -t=./temp_huodong1`

