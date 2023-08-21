# go-domainchecker

## 域名和 Kubernetes 节点状态监控工具
此工具用于定期记录给定域名的 HTTP 状态代码和正在运行的 Kubernetes 集群的节点状态。
## 功能
定期检查指定域名列表的 HTTP 状态代码并监控域名状态；
获取 Kubernetes 集群中所有节点的状态并将其存储在日志文件中；
自动将每个域名和节点的状态记录在相应时间的文件中。
## 使用方法
在 domains.txt 文件中按行输入要监控的域名；
使用 Go 编译并运行代码；
代码将会定期检查列表中的域名并捕获其 HTTP 状态代码；
同时，代码将执行 'kubectl get nodes' 命令并将结果保存在日志文件中；
查看生成的日志文件以获取详细的域名和节点状态信息。
## 示例
sh
go build -o monitor main.go
./monitor
## 代码结构说明
executeKubectlGetNodes(file_currentTime string): 此函数负责执行 'kubectl get nodes' 命令将结果输出到日志文件；
processDomain(domain, file_currentTime string): 此函数负责检查指定域名的 HTTP 状态，并将结果保存在日志文件中；
main(): 主函数，从 domains.txt 读取要监控的域名列表，并定时调用以上两个函数进行监控。
