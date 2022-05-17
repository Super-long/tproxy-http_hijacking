这个样例将所有经过eth0网卡发送到本机80端口的TCP请求通过透明代理的方式转发到本机监听192.168.123.1:1234的进程上。

### 服务器
1. sh tproxy.md
2. gcc -o tproxy_captive_portal tproxy_captive_portal.c
3. ./tproxy_captive_portal 192.168.123.1

### 客户端
1. 在客户端机器的浏览器上输入: http://42.192.201.7/whatever ，42.192.201.7为部署透明代理的服务器


### 其他
如下配置
1. iptables -t mangle -N DIVERT **在nat表上新建名为DIVERT自定义链**
2. iptables -t mangle -A DIVERT -j MARK --set-mark 1 **进入DIVERT的数据包设置标记(skb->cb？看看源码吧)**
3. iptables -t mangle -A DIVERT -j ACCEPT **默认情况下，内核会丢弃数据包，现在要确保不会**
4. iptables -t mangle -A PREROUTING -p tcp -m socket -j DIVERT **已建立的socket的TCP数据包执行DIVERT**
5. ip rule add fwmark 1 lookup 100 **所有带有1标记的数据包都不再使用默认路由表，而是使用100**
6. ip route add local 0.0.0.0/0 dev lo table 100 **添加一个路由规则，使得所有数据包（0.0.0.0）最终都被认为是本地的包**
7. iptables -t mangle -A PREROUTING -p tcp --dport 80 -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1234 --on-ip 192.168.123.1 **所有发送到80端口的TCP请求会被标记0x1并被转发到192.168.123.1:1234**

### man
TPROXY
This target is only valid in the mangle table, in the PREROUTING chain and user-defined chains which are only called from this chain. It redirects the packet to a local socket without changing the packet header in any way. It can also change the mark value which can then be used in advanced routing rules. It takes three options:
--on-port port
This specifies a destination port to use. It is a required option, 0 means the new destination port is the same as the original. This is only valid if the rule also specifies -p tcp or -p udp.
--on-ip address
This specifies a destination address to use. By default the address is the IP address of the incoming interface. This is only valid if the rule also specifies -p tcp or -p udp.
--tproxy-mark value[/mask]
Marks packets with the given value/mask. The fwmark value set here can be used by advanced routing. (Required for transparent proxying to work: otherwise these packets will get forwarded, which is probably not what you want.)


