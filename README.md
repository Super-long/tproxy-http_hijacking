这个样例将所有经过eth0网卡发送到本机80端口的TCP请求通过透明代理的方式转发到本机监听192.168.123.1:1234的进程上。

### 服务器
1. sh tproxy.md
2. gcc -o tproxy_captive_portal tproxy_captive_portal.c
3. ./tproxy_captive_portal 192.168.123.1

### 客户端
1. 在客户端机器的浏览器上输入: http://42.192.201.7/whatever ，42.192.201.7为部署透明代理的服务器