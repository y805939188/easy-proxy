## easy-proxy 轻松代理本地服务

基于 iptables 实现, 一条命令实现本地服务代理工具, 支持 ip 转发, http 转发, https 转发~

---
### 相关连接
[浅入浅出 iptables](https://zhuanlan.zhihu.com/p/507786224)</br>
[基于 iptables 实现一个 https 代理工具](https://zhuanlan.zhihu.co)</br>

---

### 使用方法
```bash
# 先编译源代码然后 build 命令会自动给命令做硬链接
make build
```
#### 设置代理
ip 代理:
```bash
# 将所有发往 http://1.2.3.4:8080 的请求都代理到本地的 3000 端口上
easy-proxy set -s http://1.2.3.4:8080 -t 127.0.0.1:3000
```

http 代理:
```bash
# 将所有发往 http://www.baidu.com 的请求都代理到本地的 3000 端口上
easy-proxy set -s http://www.baidu.com -t 127.0.0.1:3000
```

https 代理:
```bash
# 将所有发往 https://www.baidu.com 的请求都代理到本地的 3000 端口上
easy-proxy set -s https://www.baidu.com -t 127.0.0.1:3000
```

#### 查看已设置的规则
```bash
# 列出所有规则
easy-proxy list
```

#### 删除规则
```bash
# 删除某条规则, id 可通过 list 获取
easy-proxy del <id>
```

#### 清除所有规则
```bash
# 删除所有规则
easy-proxy fresh
```

---
> P.S. 基于 iptables 开发, 暂只支持 Ubuntu
