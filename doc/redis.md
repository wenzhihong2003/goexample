redis与mongodb运维
===============

- [Consul Template](https://github.com/hashicorp/consul-template)
This project provides a convenient way to populate values from Consul into the file system using the consul-template daemon.  
The daemon consul-template queries a Consul or Vault cluster and updates any number of specified templates on the file system. As an added bonus, it can optionally run arbitrary commands when the update process completes. Please see the examples folder for some scenarios where this functionality might prove useful.



### haproxy

#### ubuntu16.04安装haproxy1.8

[How to install HAProxy 1.8 on Ubuntu 16](https://blog.sleeplessbeastie.eu/2017/12/30/how-to-install-haproxy-1-8-on-ubuntu-16/)

```
# Install software-properties-common package to use add-apt-repository command.
$ sudo apt-get install software-properties-common
# Configure HAProxy 1.8 PPA by Vincent Bernat.
$ sudo add-apt-repository ppa:vbernat/haproxy-1.8
# Update package index.
$ sudo apt-get update
# Verify available HAProxy versions
$ apt-cache policy haproxy
# Upgrade HAProxy package.
$ sudo apt-get install haproxy
```

#### Haproxy反向代理WebSocket的方法
参考: https://segmentfault.com/a/1190000012779073

WebSocket请求和一般的Http请求不一样，它会长时间保持一个connection，Haproxy反向代理WebSocket请求需要用到timeout tunnel参数，否则这个链接可能就会提前关闭。例如：

```
defaults
  timeout tunnel 1h
  timeout client-fin 30s
```

参考资料：

[Websockets load-balancing with HAProxy](https://www.haproxy.com/blog/websockets-load-balancing-with-haproxy/)
[Haproxy Doc - timeout tunnel](https://cbonte.github.io/haproxy-dconv/1.6/configuration.html#timeout%20tunnel)
[InfoQ - WebSocket Proxy Servers](https://www.infoq.com/articles/Web-Sockets-Proxy-Servers)


#### mqtt haproxy

test config haproxy for gRPC loadbalancing

```
global
  tune.ssl.default-dh-param 1024
 
defaults
  timeout connect 10000ms
  timeout client 60000ms
  timeout server 60000ms
 
frontend fe_http
  mode http
  bind *:8000
  # Redirect to https
  redirect scheme https code 301
 
frontend fe_https
  mode tcp
  bind *:8443 npn spdy/2 alpn h2,http/1.1
  default_backend be_grpc

# gRPC servers running on port 8083-8084
backend be_grpc
  mode tcp
  balance roundrobin
  server srv01 127.0.0.1:8083
  server srv02 127.0.0.1:8084
```


haproxy.cfg 配制信息

可参考: https://zhuanlan.zhihu.com/p/31285285

```
global
  ulimit-n 99999
  maxconn 99999
  maxpipes 99999
  tune.maxaccept 500
  log 127.0.0.1 local0
  log 127.0.0.1 local1 notice
  chroot /var/lib/haproxy
  user haproxy
  group haproxy

defaults
  log global
  mode http
  option dontlognull
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms
  errorfile 400 /etc/haproxy/errors/400.http
  errorfile 403 /etc/haproxy/errors/403.http
  errorfile 408 /etc/haproxy/errors/408.http
  errorfile 500 /etc/haproxy/errors/500.http
  errorfile 502 /etc/haproxy/errors/502.http
  errorfile 503 /etc/haproxy/errors/503.http
  errorfile 504 /etc/haproxy/errors/504.http

listen stats :80
  stats enable
  stats uri / # must be present to see the logs
  stats auth admin:admin

listen mqtt
  bind *:1883
  bind *:8883 ssl crt /certs/lelylan-mqtt.pem
  mode tcp
  #使用http 另外开个端口进行健康检查, http://cbonte.github.io/haproxy-dconv/1.8/configuration.html#4-option%20httpchk
  option httpchk GET /health
  #记录tcp转发的信息 http://cbonte.github.io/haproxy-dconv/1.8/configuration.html#4.2-option%20tcplog
  option tcplog   
  #Use this to avoid the connection loss when client subscribed for a topic and its idle for sometime
  option clitcpka # For TCP keep-alive
  timeout client 3h #By default TCP keep-alive interval is 2hours in OS kernal, 'cat /proc/sys/net/ipv4/tcp_keepalive_time'
  timeout server 3h #By default TCP keep-alive interval is 2hours in OS kernal
  option tcplog
  balance leastconn
  server mosca_1 178.62.122.204:1883 check port 8081 inter 1s rise 2 fall 2  # 这里增加端口进行校验健康检查的    
  server mosca_2 178.62.104.172:1883 check port 8091 inter 1s rise 2 fall 2  
```


### redis

- [阿里云Redis开发规范](https://yq.aliyun.com/articles/531067)
- [Redis 内存分析方法](https://www.alibabacloud.com/help/zh/faq-detail/50037.htm) 可使用 [redis-rdb-tools](https://github.com/sripathikrishnan/redis-rdb-tools) 和 sqlite 结合来进行分析. 使用在客户端执行 bgsave 生成 rdb 文件. 或者使用[go-redis-memory-analysis](https://github.com/hhxsv5/go-redis-memory-analysis)
- [如何提取Redis中的大KEY](https://www.cnblogs.com/svan/p/7050396.html) 还有这个 [Redis4.0新特性之-大KEY删除](https://www.cnblogs.com/svan/p/7054129.html)
- [MongoShake——基于MongoDB的跨数据中心的数据复制平台](https://yq.aliyun.com/articles/603329)
- [Redis Stream——作为消息队列的典型应用场](https://yq.aliyun.com/articles/603193)
- [Why Redis 4.0?](https://yq.aliyun.com/articles/600648)
- [MongoDB · 引擎特性 · journal 与 oplog，究竟谁先写入？](https://yq.aliyun.com/articles/594727)
- [Redis用于频率限制上踩过的坑](https://github.com/aCoder2013/blog/issues/26)
- [Install and Configure HAProxy Load Balancer on Ubuntu 16.04](https://devops.profitbricks.com/tutorials/install-and-configure-haproxy-load-balancer-on-ubuntu-1604/)
- [How to Install HAProxy Load Balancer on Ubuntu](https://www.upcloud.com/support/haproxy-load-balancer-ubuntu/#layer-4-load-balancing)
- [haprox intro](http://cbonte.github.io/haproxy-dconv/1.8/intro.html)  [configuration](http://cbonte.github.io/haproxy-dconv/1.8/configuration.html) [Management](http://cbonte.github.io/haproxy-dconv/1.8/management.html)
- [HAProxy 教程](https://www.aliyun.com/jiaocheng/topic_23074_1.html)
- [我们如何使用HAProxy实现单机200万SSL连接](http://www.yunweipai.com/archives/14890.html)
- [haproxy-mqtt](https://github.com/lelylan/haproxy-mqtt)
- [HAProxy配置相关(https://github.com/chenzhiwei/linux/tree/master/haproxy)
- [HAProxy从零开始到掌握](https://www.jianshu.com/p/c9f6d55288c0)
- [Dynamic Scaling for Microservices with the HAProxy Runtime API](https://www.haproxy.com/blog/dynamic-scaling-for-microservices-with-runtime-api/)
- [Let’s Encrypt (ACMEv2) for HAProxy](https://www.haproxy.com/blog/lets-encrypt-acme2-for-haproxy/)
- [HAProxy and Consul with DNS for Service Discovery](https://help.aliyun.com/document_detail/27544.html?spm=a2c4g.11186623.6.542.N5ezhE)
- [Multithreading in HAProxy](https://www.haproxy.com/blog/multithreading-in-haproxy/)
- [HAProxy Ingress Controller for Kubernetes](https://www.haproxy.com/blog/haproxy_ingress_controller_for_kubernetes/)
- [Dynamic Configuration with the HAProxy Runtime API](https://www.haproxy.com/blog/dynamic-configuration-haproxy-runtime-api/)
- [TLS 1.3 and 0-RTT in HAProxy](https://www.haproxy.com/blog/tls-1-3-0-rtt-haproxy/)
- [What’s New in HAProxy 1.8](https://www.haproxy.com/blog/whats-new-haproxy-1-8/)
- [DNS for Service Discovery in HAProxy](https://www.haproxy.com/blog/dns-service-discovery-haproxy/)
- [Using HAProxy with the Proxy Protocol to Better Secure Your Database](https://www.haproxy.com/blog/using-haproxy-with-the-proxy-protocol-to-better-secure-your-database/)
- [HAProxy Technologies offers free hardware load balancers to students / interns](https://www.haproxy.com/blog/haproxy-technologies-offers-free-hardware-load-balancers-to-students-interns/)
- [HAProxy and ELK (Elasticsearch, Logstash and Kibana) stack](https://www.haproxy.com/blog/haproxy-and-elk-elasticsearch-logstash-and-kibana-stack/)
- [What is a slow POST Attack and how to turn HAProxy into your first line of Defense?](https://www.haproxy.com/blog/what-is-a-slow-post-attack-and-how-turn-haproxy-into-your-first-line-of-defense/)
- [HAProxy and container IP changes in Docker](https://www.haproxy.com/blog/haproxy-and-container-ip-changes-in-docker/)
- [What’s new in HAProxy 1.6](https://www.haproxy.com/blog/whats-new-in-haproxy-1-6/)
- [Serving ECC and RSA certificates on same IP with HAproxy](https://www.haproxy.com/blog/serving-ecc-and-rsa-certificates-on-same-ip-with-haproxy/)
- [HAProxy and HTTP Strict Transport Security (HSTS) header in HTTP redirects](https://www.haproxy.com/blog/haproxy-and-http-strict-transport-security-hsts-header-in-http-redirects/)
- [Packetshield: quand votre load-balancer vous protège contre les DDOS!](https://www.haproxy.com/blog/packetshield-quand-votre-load-balancer-vous-protege-contre-les-ddos/)
- [HAProxy’s load-balancing algorithm for static content delivery with Varnish](https://www.haproxy.com/blog/haproxys-load-balancing-algorithm-for-static-content-delivery-with-varnish/)
- [Microsoft Remote Desktop Services (RDS) Load-Balancing](https://www.haproxy.com/blog/microsoft-remote-desktop-services-rds-load-balancing/)
- [A HTTP monitor which matches multiple conditions in HAProxy](https://www.haproxy.com/blog/a-http-monitor-which-matches-multiple-conditions-in-haproxy/)
- [Web application name to backend mapping in HAProxy](https://www.haproxy.com/blog/web-application-name-to-backend-mapping-in-haproxy/)
- [binary health check with HAProxy 1.5: php-fpm/fastcgi probe example](https://www.haproxy.com/blog/binary-health-check-with-haproxy-1-5-php-fpmfastcgi-probe-example/)
- [Asymmetric routing, multiple default gateways on Linux with HAProxy](https://www.haproxy.com/blog/asymmetric-routing-multiple-default-gateways-on-linux-with-haproxy/)
- [How to protect application cookies while offloading SSL](https://www.haproxy.com/blog/how-to-protect-application-cookies-while-offloading-ssl/)
- [Emulating Active/passing application clustering with HAProxy](https://www.haproxy.com/blog/emulating-activepassing-application-clustering-with-haproxy/)
- [HAProxy advanced Redis health check](https://www.haproxy.com/blog/haproxy-advanced-redis-health-check/)
- [failover and worst case management with HAProxy](https://www.haproxy.com/blog/failover-and-worst-case-management-with-haproxy/)
- [Configuring HAProxy and Nginx for SPDY](https://www.haproxy.com/blog/configuring-haproxy-and-nginx-for-spdy/)
- [Howto transparent proxying and binding with HAProxy and ALOHA Load-Balancer](https://www.haproxy.com/blog/howto-transparent-proxying-and-binding-with-haproxy-and-aloha-load-balancer/)
- [SSL Client certificate information in HTTP headers and logs](https://www.haproxy.com/blog/ssl-client-certificate-information-in-http-headers-and-logs/)
- []()
- []()
- [HAProxy 反向代理的使用](http://liaoph.com/haproxy-tutorial/)
- [使用 Haproxy + Keepalived 构建基于 Docker 的高可用负载均衡服务（一）](https://blog.coding.net/blog/Haproxy&Keepalived)
- [使用 HAProxy + Keepalived 构建基于 Docker 的高可用负载均衡服务（二）](https://blog.coding.net/blog/Haproxy&Keepalived&docker)
- [haproxy cookbook](https://supermarket.chef.io/cookbooks/haproxy)
- [An Introduction to HAProxy and Load Balancing Concepts](https://www.digitalocean.com/community/tutorials/an-introduction-to-haproxy-and-load-balancing-concepts)
- [An Introduction to Networking Terminology, Interfaces, and Protocols](https://www.digitalocean.com/community/tutorials/an-introduction-to-networking-terminology-interfaces-and-protocols)
- [Floating IPs 浮动ip, 用于高可用(HA)](https://www.digitalocean.com/docs/networking/floating-ips/)
- [How To Use HAProxy As A Layer 7 Load Balancer For WordPress and Nginx On Ubuntu 14.04](https://www.digitalocean.com/community/tutorials/how-to-use-haproxy-as-a-layer-7-load-balancer-for-wordpress-and-nginx-on-ubuntu-14-04)
- [How To Implement SSL Termination With HAProxy on Ubuntu 14.04](https://www.digitalocean.com/community/tutorials/how-to-implement-ssl-termination-with-haproxy-on-ubuntu-14-04)
- [Load Testing HAProxy (Part 1)](https://medium.freecodecamp.org/load-testing-haproxy-part-1-f7d64500b75d)
- [Load Testing HAProxy (Part 2)](https://medium.freecodecamp.org/load-testing-haproxy-part-2-4c8677780df6)
- [How we fine-tuned HAProxy to achieve 2,000,000 concurrent SSL connections](https://medium.freecodecamp.org/how-we-fine-tuned-haproxy-to-achieve-2-000-000-concurrent-ssl-connections-d017e61a4d27)
- [我们是如何优化HAProxy以让其支持2,000,000个并发SSL连接的？](http://www.infoq.com/cn/articles/fine-tuned-haproxy-to-achieve-concurrent-ssl-connections)
- [Use HAProxy to load balance 300k concurrent tcp socket connections: Port Exhaustion, Keep-alive and others](https://www.linangran.com/?p=547)
- [谈一谈使用 HAProxy 构建 API 网关服务的思路](https://zhuanlan.zhihu.com/p/27155886)
- [Go library for interacting with HAProxy via command socket](https://github.com/bcicen/go-haproxy)
- [haproxycmd HAProxy Unix Socket commands CLI](https://github.com/winebarrel/haproxycmd)
- [基于HAProxy的高性能缓存服务器nuster](https://zhuanlan.zhihu.com/p/32366458)
- [Pecemaker+Corosync+Haproxy实现Openstack高可用](https://zhuanlan.zhihu.com/p/31285285)
- [有趣也有用的现代类型系统](https://zhuanlan.zhihu.com/p/33882384)
- [利用 Etcd + Confd + HAProxy 實現服務動態拓展](https://blog.toright.com/posts/5687/%E5%88%A9%E7%94%A8-etcd-confd-haproxy-%E5%AF%A6%E7%8F%BE%E6%9C%8D%E5%8B%99%E5%8B%95%E6%85%8B%E6%8B%93%E5%B1%95.html)
- [Haproxy+etcd+confd+Docker搭建节点自动发现的高可用负载均衡框架](https://www.jianshu.com/p/bc85a54f98ff)
- [微服务注册发现集群搭建——Registrator + Consul + Consul-template + nginx](https://www.jianshu.com/p/a4c04a3eeb57)
- [confd Manage local application configuration files using templates and data from etcd or consul](http://www.confd.io/) 
- [Confd+etcd实现高可用自动发现](http://www.361way.com/confd-etcd/5470.html)
- [etcd集群搭建](http://www.361way.com/etcd-cluster/5468.html)
- [HAProxy（一）之常用配置介绍，ACL详解](https://www.jianshu.com/p/60dcac5c032c)
- [Haproxy（二）之负载均衡配置详解](https://www.jianshu.com/p/a7e3199a0a09)
- [Haproxy配置文件详解](https://www.jianshu.com/p/b671610b5cea)
- [关于TCP 半连接队列和全连接队列](https://www.jianshu.com/p/569bdd440b09)
- [异步网络模型](https://www.jianshu.com/p/9b51eb95fadf)
- [HAProxy 做rabbitmq的高可用](https://www.jianshu.com/p/7a3142335518)
- [HAProxy用法详解](http://www.ttlsa.com/linux/haproxy-study-tutorial/)
- [HAProxy常用配置介绍，ACL详解](https://segmentfault.com/a/1190000007532860)
- [HAProxy配置文档（精简版）](http://blog.xiayf.cn/gitbook/tech-note/operation/haproxy/conf-manual.html)
- [haproxy 简介](https://www.kancloud.cn/hanxt/foreign-docker/186537)
- [haproxy的基本使用与常见实践](http://www.lijiaocn.com/%E6%8A%80%E5%B7%A7/2017/06/26/haproxy-usage.html)
- [confd：本地配置文件的管理工具confd](http://www.lijiaocn.com/%E6%8A%80%E5%B7%A7/2017/09/21/linux-tool-confd.html)
- [网络流量控制技术](http://www.lijiaocn.com/%E6%8A%80%E5%B7%A7/2015/03/06/%E6%B5%81%E9%87%8F%E6%8E%A7%E5%88%B6.html)
- [高可用实现方法汇总](http://www.lijiaocn.com/%E6%8A%80%E5%B7%A7/2015/07/02/%E9%AB%98%E5%8F%AF%E7%94%A8%E5%AE%9E%E7%8E%B0.html)
- [Linux系统的优化方法](http://www.lijiaocn.com/%E6%8A%80%E5%B7%A7/2014/10/16/Linux%E7%B3%BB%E7%BB%9F%E8%BF%90%E7%BB%B4.html)
- [Linux的nftables的使用](http://www.lijiaocn.com/%E6%96%B9%E6%B3%95/2018/06/15/linux-nftables-usage.html)
- [区块链在互联网借贷领域的应用探索](http://www.lijiaocn.com/%E6%96%B9%E6%B3%95/2018/05/11/blockchain-app-loan.html)
- []()
- []()
- [基于haproxy实现灰度发布](https://www.jianshu.com/p/0a54cfb2e6f4)
- [HAProxy 是什么，不是什么](https://www.jianshu.com/p/0820b45b4f75)
- [redisHA 用proxy实现redis的高可用与读写分离](https://github.com/kelgon/redisHA)
- [haproxy小结（一）基础概念篇](http://www.361way.com/haproxy-basic/5267.html)
- [haproxy小结（二）配置文件篇](http://www.361way.com/haproxy-config/5275.html)
- [haproxy小结（三）配置示例](http://www.361way.com/haproxy-config-example/5277.html)
- [haproxy小结（四）ebtree](http://www.361way.com/haproxy-ebtree/5280.html)
- [keepalived健康检查方式](http://www.361way.com/keepalived-health-check/5218.html)
- [LVS高可用（六）LVS+keepalived主从](http://www.361way.com/lvs-keepalived-dr-master-backup/5221.html)
- [LVS高可用（七）LVS+keepalived双主](http://www.361way.com/lvs-keepalived-two-master/5225.html)
- [keepalived的同步组和sorry地址](http://www.361way.com/sorry-server-sync-group/5228.html)
- [keepalived配置架构详解](http://www.361way.com/keepalived-framework/5208.html)
- [LVS 负载均衡（一）理论篇](http://www.361way.com/lvs-summary/5177.html)
- [LVS 负载均衡（二）NAT模式](http://www.361way.com/lvs-nat/5187.html)
- [LVS 负载均衡（三）DR模式](http://www.361way.com/lvs-dr/5192.html)
- [LVS-DR工作原理及答疑](http://www.361way.com/lvs-dr-theory/5195.html)
- [kubernetes-the-hard-way Bootstrap Kubernetes the hard way on Google Cloud Platform. No scripts](https://github.com/kelseyhightower/kubernetes-the-hard-way)
- [envconfig Golang library for managing configuration data from environment variables](https://github.com/kelseyhightower/envconfig)
- [event-gateway-on-kubernetes How to guide on running Serverless.com's Event Gateway on Kubernetes ](https://github.com/kelseyhightower/event-gateway-on-kubernetes)
- [Confd+Consul 配置文件动态生成](https://my.oschina.net/guol/blog/480788)
- [使用confd和Nginx做边缘服务](https://servicecomb.incubator.apache.org/cn/users/edging-service/nginx/)
- [使用etcd+confd管理nginx配置](https://www.cnblogs.com/Anker/p/6112022.html)
- [django + etcd + confd 配置管理平台](https://zhuanlan.zhihu.com/p/38353902)
- [在centos7中搭建haproxy-confd-etcd-tomcat](https://segmentfault.com/a/1190000007213942)
- [使用 Etcd 和 Haproxy 做 Docker 服务发现](https://segmentfault.com/a/1190000000730186)
- [通过Etcd+Confd自动管理Haproxy(多站点)](http://blog.51cto.com/fengwan/1899964)
- [docker服务发现——confd](https://yq.aliyun.com/articles/8708)
- [如何在Ubuntu 16.04上用Logrotate管理日志文件](https://www.howtoing.com/how-to-manage-logfiles-with-logrotate-on-ubuntu-16-04)
- [Linux 进程的管理与监控](http://liaoph.com/inux-process-management/)
- [加密、解密以及 openssl](http://liaoph.com/encrytion-and-openssl/)
- [计算机原理 —— 主板与内存映射](http://liaoph.com/motherboard-and-memory-map/)
- [计算机原理 —— 计算机是如何启动的](http://liaoph.com/how-computers-boot-up/)
- [文件共享服务 FTP，NFS 和 Samba](http://liaoph.com/ftp-nfs-samba/)
- [负载均衡集群 LVS 详解](http://liaoph.com/lvs/)
- []()


### mongodb

- [MongoDB-Elasticsearch 实时数据导入](https://zhuanlan.zhihu.com/p/26906652)
- [patroni A template for PostgreSQL High Availability with ZooKeeper, etcd, or Consul](https://github.com/zalando/patroni)
- [Java IO 及 Netty原理详解](https://my.oschina.net/tantexian/blog/775875)
- [系统原理分析架构-三--代理服务器简介及分类](https://my.oschina.net/tantexian/blog/626178)
- [系统原理分析架构-开篇 （对于架构师与开发语言及被青春饭的一些想法）](https://my.oschina.net/tantexian/blog/626177)
- [系统原理分析架构-四-squid(简介及正向代理)](https://my.oschina.net/tantexian/blog/626176)
- [keepalived原理（主从配置+haproxy）及配置文件详解](https://my.oschina.net/tantexian/blog/626175)
- [负载均衡之Haproxy配置详解（及httpd配置）](https://my.oschina.net/tantexian/blog/626174)
- [系统原理分析架构-二-CDN内容分发网络](https://my.oschina.net/tantexian/blog/626173)
- [系统原理分析架构-一-DNS负载均衡](https://my.oschina.net/tantexian/blog/626171)
- [系统原理分析架构-六-负载均衡（定义及介绍及LVS/Nginx/Haproxy比较）](https://my.oschina.net/tantexian/blog/626170)
- [squid,nginx,lighttpd反向代理的区别(同步VS异步模式) ](https://my.oschina.net/tantexian/blog/626169)
- [系统原理分析架构-五-squid(反代理即web缓存服务器)](https://my.oschina.net/tantexian/blog/626168)
- [99%海量数据处理](https://my.oschina.net/tantexian/blog/1607415)
- [高性能Server---Reactor模型](https://my.oschina.net/tantexian/blog/1563960)
- [设计Go API的管道使用原则](http://blog.jobbole.com/73700/)
- [Go并发模式：管道和显式取消](http://blog.jobbole.com/65552/)
- [用Go语言写HTTP中间件](http://blog.jobbole.com/53265/)
- [复用Go的内存buffer](http://blog.jobbole.com/48969/)
- [这篇文章很赞：我读过的最好的一篇分布式技术文章关于log](https://my.oschina.net/tantexian/blog/884577)
- [分布式跟踪](https://github.com/Yirendai/cicada/blob/master/cicada-docs/cicada_design.md)
- [JAVA集合框架中的常用集合及其特点、适用场景、实现原理简介](https://www.jianshu.com/p/b54f1df33f84)
- [Java异常控制机制和异常处理原则](https://www.jianshu.com/p/15872cba211d)
- [LeetCode算法题解：LFU Cache](https://www.jianshu.com/p/437f53341f67)
- [Redis基础、高级特性与性能调优](https://www.jianshu.com/p/2f14bc570563)
- [不依赖客户端库，实现Redis主从架构的高可用/读写分离/负载均衡](https://www.jianshu.com/p/99be0a88517d)
- [Redis + Sentinel服务搭建](https://www.jianshu.com/p/9553dee9d1fc)
- [数据库连接池到底应该设多大？这篇文章可能会颠覆你的认知](https://www.jianshu.com/p/a8f653fc0c54)
- []()

### 消息队列
- [消息队列实现概要——深度解读分区Topic的实现](https://github.com/aCoder2013/blog/issues/27)
- [基于JVM之上的并发编程模式剖析](https://github.com/aCoder2013/blog/issues/25)
- [CopyOnWriteArrayList内部工作原理剖析](https://github.com/aCoder2013/blog/issues/24)
- [Java并发工具类之LongAdder原理总结](https://github.com/aCoder2013/blog/issues/22)
- [分布式消息队列实现概要](https://github.com/aCoder2013/blog/issues/21)
- [日志: 分布式系统的核心](https://github.com/aCoder2013/blog/issues/20)
- [Java原生类型包装类初解析 ](https://github.com/aCoder2013/blog/issues/14)
- [Atomic包之FieldUpdater深度解析](https://github.com/aCoder2013/blog/issues/10)
- [JVM指令集中tableswitch和lookupswitch指令的区别](https://github.com/aCoder2013/blog/issues/7)
- [深度解析Java线程池的异常处理机制](https://github.com/aCoder2013/blog/issues/3)
- [Java线程那点事儿](https://github.com/aCoder2013/blog/issues/4)
- [MongoDB导出场景查询优化](https://github.com/aCoder2013/blog/issues/1)
- [百万级分布式开源物联网MQTT消息服务器](http://emqtt.com/)
- [Fast Topic Matching 前缀匹配](https://bravenewgeek.com/fast-topic-matching/)
- [Dissecting Message Queues](https://bravenewgeek.com/dissecting-message-queues/)
- [Introducing Liftbridge: Lightweight, Fault-Tolerant Message Streams](https://bravenewgeek.com/introducing-liftbridge-lightweight-fault-tolerant-message-streams/)
- [Building a Distributed Log from Scratch, Part 5: Sketching a New System](https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-5-sketching-a-new-system/)
- [Building a Distributed Log from Scratch, Part 4: Trade-Offs and Lessons Learned](https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-4-trade-offs-and-lessons-learned/)
- [Building a Distributed Log from Scratch, Part 3: Scaling Message Delivery](https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-3-scaling-message-delivery/)
- [Building a Distributed Log from Scratch, Part 2: Data Replication](https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-2-data-replication/)
- [Building a Distributed Log from Scratch, Part 1: Storage Mechanics](https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-1-storage-mechanics/)
- [FIFO, Exactly-Once, and Other Costs](https://bravenewgeek.com/fifo-exactly-once-and-other-costs/)
- [You Cannot Have Exactly-Once Delivery Redux](https://bravenewgeek.com/you-cannot-have-exactly-once-delivery-redux/)
- [Smart Endpoints, Dumb Pipes](https://bravenewgeek.com/smart-endpoints-dumb-pipes/)
- [Fast Topic Matching](https://bravenewgeek.com/fast-topic-matching/)
- [Take It to the Limit: Considerations for Building Reliable Systems](https://bravenewgeek.com/take-it-to-the-limit-considerations-for-building-reliable-systems/)
- [Benchmarking Message Queue Latency](https://bravenewgeek.com/benchmarking-message-queue-latency/)
- [What You Want Is What You Don’t: Understanding Trade-Offs in Distributed Messaging](https://bravenewgeek.com/what-you-want-is-what-you-dont-understanding-trade-offs-in-distributed-messaging/)
- [You Cannot Have Exactly-Once Delivery](https://bravenewgeek.com/you-cannot-have-exactly-once-delivery/)
- [Benchmark Responsibly](https://bravenewgeek.com/benchmark-responsibly/)
- [Iris Decentralized Cloud Messaging](https://bravenewgeek.com/iris-decentralized-cloud-messaging/)
- [Dissecting Message Queues](https://bravenewgeek.com/dissecting-message-queues/)
- [A Look at Nanomsg and Scalability Protocols (Why ZeroMQ Shouldn’t Be Your First Choice)](https://bravenewgeek.com/a-look-at-nanomsg-and-scalability-protocols/)
- [Distributed Messaging with ZeroMQ](https://bravenewgeek.com/distributed-messaging-with-zeromq/)
- [A Look at Nanomsg and Scalability Protocols (Why ZeroMQ Shouldn’t Be Your First Choice)](https://bravenewgeek.com/a-look-at-nanomsg-and-scalability-protocols/)
- [Everything You Know About Latency Is Wrong](https://bravenewgeek.com/everything-you-know-about-latency-is-wrong/)
- [Breaking and Entering: Lose the Lock While Embracing Concurrency](https://bravenewgeek.com/breaking-and-entering-lose-the-lock-while-embracing-concurrency/)
- [So You Wanna Go Fast?](https://bravenewgeek.com/so-you-wanna-go-fast/)
- [Breaking and Entering: Lose the Lock While Embracing Concurrency](https://bravenewgeek.com/breaking-and-entering-lose-the-lock-while-embracing-concurrency/)
- [Stream Processing and Probabilistic Methods: Data at Scale](https://bravenewgeek.com/stream-processing-and-probabilistic-methods/)
- [MySQL日志功能详解](https://jkzhao.github.io/2018/04/16/MySQL%E6%97%A5%E5%BF%97%E5%8A%9F%E8%83%BD%E8%AF%A6%E8%A7%A3/)
- [Nginx实现TCP负载均衡](https://jkzhao.github.io/2018/01/28/Nginx%E5%AE%9E%E7%8E%B0TCP%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1/)
- [Nginx配置详解及优化](https://jkzhao.github.io/2018/01/23/Nginx%E9%85%8D%E7%BD%AE%E8%AF%A6%E8%A7%A3%E5%8F%8A%E4%BC%98%E5%8C%96/)
- [安装部署Openresty](https://jkzhao.github.io/2018/01/21/%E5%AE%89%E8%A3%85%E9%83%A8%E7%BD%B2Openresty/)
- [I/O模型](https://jkzhao.github.io/2018/01/18/I-O%E6%A8%A1%E5%9E%8B/)
- [在Kubernetes上使用Traefik](https://jkzhao.github.io/2017/09/25/%E5%9C%A8Kubernetes%E4%B8%8A%E4%BD%BF%E7%94%A8Traefik/)
- [Kubernetes Ingress实战](https://jkzhao.github.io/2017/09/24/Kubernetes-Ingress%E5%AE%9E%E6%88%98/)
- [Kubernetes监控:部署Heapster、InfluxDB和Grafana](https://jkzhao.github.io/2017/09/21/Kubernetes%E7%9B%91%E6%8E%A7-%E9%83%A8%E7%BD%B2Heapster%E3%80%81InfluxDB%E5%92%8CGrafana/)
- [企业级Docker Registry —— Harbor搭建和使用](https://jkzhao.github.io/2017/09/08/%E4%BC%81%E4%B8%9A%E7%BA%A7Docker-Registry-%E2%80%94%E2%80%94-Harbor%E6%90%AD%E5%BB%BA%E5%92%8C%E4%BD%BF%E7%94%A8/)
- [overlay实现容器跨主机通信](https://jkzhao.github.io/2017/09/05/overlay%E5%AE%9E%E7%8E%B0%E5%AE%B9%E5%99%A8%E8%B7%A8%E4%B8%BB%E6%9C%BA%E9%80%9A%E4%BF%A1/)
- [Registry私有仓库搭建及认证](https://jkzhao.github.io/2017/09/01/Registry%E7%A7%81%E6%9C%89%E4%BB%93%E5%BA%93%E6%90%AD%E5%BB%BA%E5%8F%8A%E8%AE%A4%E8%AF%81/)
- []()
- []()
- []()

- [聊聊前端开发的测试](https://blog.coding.net/blog/frontend-testing)
- [深入浅出 Git](https://blog.coding.net/blog/git-from-the-inside-out)
- [关于 Git 你需要知道的一些事情](https://blog.coding.net/blog/about-git)
- [使用原理视角看 Git](https://blog.coding.net/blog/principle-of-Git)
- [Commit message 和 Change log 编写指南](https://blog.coding.net/blog/commit_message_change_log)
- [大话 Git 工作流](https://blog.coding.net/blog/git-workflow)
- [场景式解读 Git 工作流](https://blog.coding.net/blog/git-workflow-2)
- []()
- []()
