# homework3 2021-11-13

## 内容：

### 第一部分
编写 Kubernetes 部署脚本，将 httpserver 部署到 kubernetes 集群（模块8作业）

以下是可以思考的维度：
1) 优雅启动
2) 优雅终止
3) 资源需求和 QoS 保证 
4) 探活
5) 日常运维需求，日志等级
6) 配置和代码分离

### 第二部分
除了将 httpServer 应用优雅地运行在 Kubernetes 之上，我们还应该考虑如何将服务发布给对内和对外的调用方。
来尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧（模块9作业）  
在第一部分的基础上提供更加完备的部署 spec， 包括（不限于）
- Service
- Ingress

可以考虑的细节
- 如何确保整个应用的高可用
- 如何通过证书保证 httpServer 的通讯安全
