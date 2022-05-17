iptables -t mangle -A PREROUTING -i eth0 -p tcp --dport 80 -m tcp -j TPROXY --on-ip 192.168.123.1 --on-port 1234 --tproxy-mark 1/1
sysctl -w net.ipv4.ip_forward=1
ip rule add fwmark 1/1 table 1
ip route add local 0.0.0.0/0 dev lo table 1