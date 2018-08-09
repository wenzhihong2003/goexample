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

haproxy.cfg 配制信息

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
  #Use this to avoid the connection loss when client subscribed for a topic and its idle for sometime
  option clitcpka # For TCP keep-alive
  timeout client 3h #By default TCP keep-alive interval is 2hours in OS kernal, 'cat /proc/sys/net/ipv4/tcp_keepalive_time'
  timeout server 3h #By default TCP keep-alive interval is 2hours in OS kernal
  option tcplog
  balance leastconn
  server mosca_1 178.62.122.204:1883 check
  server mosca_2 178.62.104.172:1883 check
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
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()


### mongodb



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
- []()
- []()
- []()
- []()

### linux
- [How to log dropped connections from iptables firewall using netfilter userspace logging daemon](https://blog.sleeplessbeastie.eu/2018/08/01/how-to-log-dropped-connections-from-iptables-firewall-using-netfilter-userspace-logging-daemon/)
- [How to calculate how fast data is copied to the specified directory](https://blog.sleeplessbeastie.eu/2018/07/30/how-to-calculate-how-fast-data-is-copied-to-the-specified-directory/)
- [How to execute logrotate every hour](https://blog.sleeplessbeastie.eu/2018/07/11/how-to-execute-logrotate-every-hour/)
- [How to locate directories that contain one or multiple files with particular content](https://blog.sleeplessbeastie.eu/2018/07/09/how-to-locate-directories-that-contain-one-or-multiple-files-with-particular-content/)
- [How to log dropped connections from iptables firewall using rsyslog](https://blog.sleeplessbeastie.eu/2018/07/06/how-to-log-dropped-connections-from-iptables-firewall-using-rsyslog/)
- [How to find symbolic link by a target name](https://blog.sleeplessbeastie.eu/2018/08/06/how-to-find-symbolic-link-by-a-target-name/)
- [How to display package dependencies](https://blog.sleeplessbeastie.eu/2018/07/02/how-to-display-package-dependencies/)
- [How to create iptables firewall using custom chains](https://blog.sleeplessbeastie.eu/2018/06/21/how-to-create-iptables-firewall-using-custom-chains/)
- [How to create iptables firewall](https://blog.sleeplessbeastie.eu/2018/06/13/how-to-create-iptables-firewall/)
- [How to use HTTP host header to choose HAProxy backend](https://blog.sleeplessbeastie.eu/2018/06/06/how-to-use-evironment-variable-and-lua-to-choose-haproxy-backend/)
- [How to detect and log changes in the list of mounted filesystems](https://blog.sleeplessbeastie.eu/2018/06/04/how-to-detect-and-log-changes-in-the-list-of-mounted-filesystems/)
- [How to determine process execution time](https://blog.sleeplessbeastie.eu/2018/05/30/how-to-determine-process-execution-time/)
- [How to execute additional commands during system startup using cron](https://blog.sleeplessbeastie.eu/2018/05/28/how-to-execute-additional-commands-during-system-startup-using-cron/)
- [How to redirect every request to defined domain to particular location](https://blog.sleeplessbeastie.eu/2018/05/24/how-to-redirect-every-request-to-defined-domain-to-particular-location/)
- [How to store LeaseWeb data traffic in OpenTSDB time series database](https://blog.sleeplessbeastie.eu/2018/05/17/how-to-store-leaseweb-data-traffic-in-opentsdb-time-series-database/)
- [How to serve files from memory with a fallback using nginx](https://blog.sleeplessbeastie.eu/2018/05/14/how-to-serve-files-from-memory-with-a-fallback-using-nginx/)
- [How to create simplest possible iptables firewall](https://blog.sleeplessbeastie.eu/2018/05/09/how-to-create-simplest-possible-iptables-firewall/)
- [How to use variable to choose HAProxy backend](https://blog.sleeplessbeastie.eu/2018/05/02/how-to-use-variable-to-choose-haproxy-backend/)
- [How to setup private Docker registry](https://blog.sleeplessbeastie.eu/2018/04/16/how-to-setup-private-docker-registry/)
- [How to block defined IP addresses on HAProxy](https://blog.sleeplessbeastie.eu/2018/03/26/how-to-block-particular-ip-addresses-on-haproxy/)
- [How to define basic authentication on HAProxy](https://blog.sleeplessbeastie.eu/2018/03/08/how-to-define-basic-authentication-on-haproxy/)
- [How to define allowed HTTP methods on HAProxy](https://blog.sleeplessbeastie.eu/2018/03/01/how-to-define-allowed-http-methods-on-haproxy/)
- [How to display dependencies for deb package](https://blog.sleeplessbeastie.eu/2018/03/05/how-to-display-dependencies-for-deb-package/)
- [How to locally check SSL certificate](https://blog.sleeplessbeastie.eu/2018/02/19/how-to-check-ssl-certificate/)
- [How to parse free output to display available memory](https://blog.sleeplessbeastie.eu/2018/01/31/how-to-parse-free-output-to-display-available-memory/)
- [How to install and configure MariaDB unixODBC driver](https://blog.sleeplessbeastie.eu/2018/01/08/how-to-install-and-configure-mariadb-unixodbc-driver/)
- [How to configure HTTP/2 in http mode on HAProxy and fix bad request problem](https://blog.sleeplessbeastie.eu/2018/01/01/how-to-configure-http2-in-http-mode-on-haproxy-and-fix-bad-request-problem/)
- [How to generate private key](https://blog.sleeplessbeastie.eu/2017/12/28/how-to-generate-private-key/)
- [How to access Microsoft SQL Server instance when you are locked out](https://blog.sleeplessbeastie.eu/2017/12/25/how-to-access-microsoft-sql-server-instance-when-you-are-locked-out/)
- [How to verify IPv4 address](https://blog.sleeplessbeastie.eu/2017/12/11/how-to-verify-ipv4-address/)
- [How to determine encryption algorithm used to store password](https://blog.sleeplessbeastie.eu/2017/12/06/how-to-determine-encryption-algorithm-used-to-store-password/)
- [How to display IPv4 network information](https://blog.sleeplessbeastie.eu/2017/12/04/how-to-display-ipv4-network-information/)
- [How to remove unused dependency packages](https://blog.sleeplessbeastie.eu/2017/10/30/how-to-remove-unused-dependency-packages/)
- [How to automatically control APT cache](https://blog.sleeplessbeastie.eu/2017/10/16/how-to-automatically-control-apt-cache/)
- [How to reverse text file](https://blog.sleeplessbeastie.eu/2017/10/12/how-to-reverse-text-file/)
- [How to clear the APT cache](https://blog.sleeplessbeastie.eu/2017/10/09/how-to-clean-the-apt-cache/)
- [How to disable the APT cache](https://blog.sleeplessbeastie.eu/2017/10/02/how-to-disable-the-apt-cache/)
- [How to keep track of network latency](https://blog.sleeplessbeastie.eu/2017/09/11/how-to-keep-track-of-network-latency/)
- [How to display real destination URL](https://blog.sleeplessbeastie.eu/2017/08/28/how-to-display-real-destination-url/)
- [How to display HTTP response code](https://blog.sleeplessbeastie.eu/2017/08/07/how-to-display-http-response-code/)
- [How to restart process depending on the log file modification time](https://blog.sleeplessbeastie.eu/2017/07/24/how-to-restart-process-depending-on-the-log-file-modification-time/)
- [How to check connection to the RabbitMQ message broker](https://blog.sleeplessbeastie.eu/2017/07/10/how-to-check-connection-to-the-rabbitmq-message-broker/)
- [How to assign an IPv6 address to an interface](https://blog.sleeplessbeastie.eu/2017/06/19/how-to-assign-an-ipv6-address-to-an-interface/)
- [How to display message provided from standard input or as an parameter](https://blog.sleeplessbeastie.eu/2017/06/05/how-to-display-message-provided-from-standard-input-or-as-an-parameter/)
- [How to display network connections using lsof and GNU awk](https://blog.sleeplessbeastie.eu/2017/05/29/how-to-display-network-connections-using-lsof-and-gnu-awk/)
- [How to copy standard output and catch error code in the meantime](https://blog.sleeplessbeastie.eu/2017/05/01/how-to-copy-standard-output-and-catch-error-code-in-the-meantime/)
- [How to install missing ifconfig utility](https://blog.sleeplessbeastie.eu/2017/04/27/how-to-install-missing-ifconfig-utility/)
- [How to use a parameter or standard input inside the shell script](https://blog.sleeplessbeastie.eu/2017/04/20/how-to-use-a-parameter-or-standard-input-inside-the-shell-script/)
- [How to stop referral spam using Nginx](https://blog.sleeplessbeastie.eu/2017/04/10/how-to-stop-referral-spam-using-nginx/)
- [How to display days till certificate expiration](https://blog.sleeplessbeastie.eu/2017/04/03/how-to-display-days-till-certificate-expiration/)
- [How to list configured APT repositories](https://blog.sleeplessbeastie.eu/2017/03/27/how-to-list-configured-apt-repositories/)
- [How to display certificate issuer and dates](https://blog.sleeplessbeastie.eu/2017/03/20/how-to-display-certificate-issuer-and-dates/)
- [How to download files recursively](https://blog.sleeplessbeastie.eu/2017/02/06/how-to-download-files-recursively/)
- [How to print IP address assigned to an interface](https://blog.sleeplessbeastie.eu/2017/01/23/how-to-print-ip-address-assigned-to-an-interface/)
- [How to display process environment](https://blog.sleeplessbeastie.eu/2017/01/16/how-to-display-process-environment/)
- [How to remove invalid entires from known hosts file](https://blog.sleeplessbeastie.eu/2017/01/09/how-to-remove-invalid-entires-from-known-hosts-file/)
- [How to get the number of connections broken down by a host](https://blog.sleeplessbeastie.eu/2017/01/02/how-to-get-the-number-of-connections-broken-down-by-a-host/)
- [How to display processes using swap space](https://blog.sleeplessbeastie.eu/2016/12/26/how-to-display-processes-using-swap-space/)
- [How to sync passwords between different devices](https://blog.sleeplessbeastie.eu/2016/12/05/how-to-sync-passwords-between-different-devices/)
- [How to perform real-time performance monitoring](https://blog.sleeplessbeastie.eu/2016/11/21/how-to-perform-real-time-performance-monitoring/)
- [How to determine when process was started](https://blog.sleeplessbeastie.eu/2016/10/31/how-to-determine-when-process-was-started/)
- [How to upgrade selected packages](https://blog.sleeplessbeastie.eu/2016/10/10/how-to-upgrade-selected-packages/)
- [How to list available updates using apt](https://blog.sleeplessbeastie.eu/2016/10/03/how-to-list-available-updates-using-apt/)
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
- []()
