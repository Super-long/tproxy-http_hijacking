这个样例将所有经过eth0网卡发送到本机80端口的TCP请求通过透明代理的方式转发到本机监听192.168.123.1:1234的进程上。

1. sh tproxy.md
2. gcc -o tproxy_captive_portal tproxy_captive_portal.c
3. ./tproxy_captive_portal 192.168.123.1