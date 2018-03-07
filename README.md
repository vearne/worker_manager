# 可以切换到工程目录下
# 执行
```
go build main.go
# 启动服务
./main
# 服务退出, 发出SIGTERM信号，服务优雅退出
# 请求自行替换pid的值
kill -15 <pid> 
```
