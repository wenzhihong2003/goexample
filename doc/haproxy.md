haproxy
=======

参见: https://www.jianshu.com/p/c9f6d55288c0

#### HAProxy的核心功能

- 负载均衡：L4和L7两种模式，支持RR/静态RR/LC/IP Hash/URI Hash/URL_PARAM Hash/HTTP_HEADER Hash等丰富的负载均衡算法
- 健康检查：支持TCP和HTTP两种健康检查模式
- 会话保持：对于未实现会话共享的应用集群，可通过Insert Cookie/Rewrite Cookie/Prefix Cookie，以及上述的多种Hash方式实现会话保持
- SSL：HAProxy可以解析HTTPS协议，并能够将请求解密为HTTP后向后端传输
- HTTP请求重写与重定向
- 监控与统计：HAProxy提供了基于Web的统计信息页面，展现健康状态和流量数据。基于此功能，使用者可以开发监控程序来监控HAProxy的状态

#### HAProxy的关键特性

**性能**

- 采用单线程、事件驱动、非阻塞模型，减少上下文切换的消耗，能在1ms内处理数百个请求。并且每个会话只占用数KB的内存。
- 大量精细的性能优化，如O(1)复杂度的事件检查器、延迟更新技术、Single-buffereing、Zero-copy forwarding等等，这些技术使得HAProxy在中等负载下只占用极低的CPU资源。
- HAProxy大量利用操作系统本身的功能特性，使得其在处理请求时能发挥极高的性能，通常情况下，HAProxy自身只占用15%的处理时间，剩余的85%都是在系统内核层完成的。
- HAProxy作者在8年前（2009）年使用1.4版本进行了一次测试，单个HAProxy进程的处理能力突破了10万请求/秒，并轻松占满了10Gbps的网络带宽

**稳定性**  

作为建议以单进程模式运行的程序，HAProxy对稳定性的要求是十分严苛的。按照作者的说法，HAProxy在13年间从未出现过一个会导致其崩溃的BUG，HAProxy一旦成功启动，除非操作系统或硬件故障，否则就不会崩溃（我觉得可能多少还是有夸大的成分）。

在上文中提到过，HAProxy的大部分工作都是在操作系统内核完成的，所以HAProxy的稳定性主要依赖于操作系统，作者建议使用2.6或3.x的Linux内核，对sysctls参数进行精细的优化，并且确保主机有足够的内存。这样HAProxy就能够持续满负载稳定运行数年之久。

先查看阿里云的文章 [ECS做负载均衡需要用户做额外的配置吗？](https://help.aliyun.com/knowledge_detail/39428.html)
[Linux 实例常用内核网络参数介绍与常见问题处理](https://help.aliyun.com/knowledge_detail/41334.html)

**个人的建议：**

- 使用3.x内核的Linux操作系统运行HAProxy
- 运行HAProxy的主机上不要部署其他的应用，确保HAProxy独占资源，同时避免其他应用引发操作系统或主机的故障
- 至少为HAProxy配备一台备机，以应对主机硬件故障、断电等突发情况（搭建双活HAProxy的方法在后文中有描述）
- sysctl的建议配置（并不是万用配置，仍然需要针对具体情况进行更精细的调整，但可以作为首次使用HAProxy的初始配置使用）：
```
net.ipv4.tcp_tw_reuse = 1
net.ipv4.ip_local_port_range = 1024 65023
net.ipv4.tcp_max_syn_backlog = 10240
net.ipv4.tcp_max_tw_buckets = 400000
net.ipv4.tcp_max_orphans = 60000
net.ipv4.tcp_synack_retries = 3
net.core.somaxconn = 10000
```
以上的设置可通过下面的方式进行
For the change in a running system use
```
sudo sysctl -w net.ipv4.ip_local_port_range="1024 65000"
```
or  
```
echo "1024 65000" > /proc/sys/net/ipv4/ip_local_port_range
```

To survive a reboot, open the configuration file. 可重启的情况
```
sudo nano /etc/sysctl.conf
```
and add the lines below
```
# increase system IP port limits
net.ipv4.ip_local_port_range = 1024 65000
```
After the next reboot you can see, the correct values are set.


### 添加日志
HAProxy不会直接输出文件日志，所以我们要借助Linux的rsyslog来让HAProxy输出日志

修改haproxy.cfg
在global域和defaults域中添加：
```
global
    ...
    log 127.0.0.1 local0 info
    log 127.0.0.1 local1 warning
    ...

defaults
    ...
    log global
    ...
```

意思是将info级（及以上）的日志推送到rsyslog的local0接口，将warn级（及以上）的日志推送到rsyslog的local1接口，并且所有frontend都默认使用global中的日志配置。

注：info级的日志会打印HAProxy处理的每一条请求，会占用很大的磁盘空间，在生产环境中，建议将日志级别调整为notice

为rsyslog添加haproxy日志的配置

```
vi /etc/rsyslog.d/haproxy.conf
$ModLoad imudp
$UDPServerRun 514
$FileCreateMode 0644  #日志文件的权限
$FileOwner ha  #日志文件的owner
local0.*     /var/log/haproxy.log  #local0接口对应的日志输出文件
local1.*     /var/log/haproxy_warn.log  #local1接口对应的日志输出文件
```
修改rsyslog的启动参数
```
vi /etc/default/rsyslog
# Options for rsyslogd
# Syslogd options are deprecated since rsyslog v3.
# If you want to use them, switch to compatibility mode 2 by "-c 2"
# See rsyslogd(8) for more details
SYSLOGD_OPTIONS="-c 2 -r -m 0"
```
重启rsyslog和HAProxy
```
service rsyslog restart
service haproxy restart
```
此时就应该能在/var/log目录下看到haproxy的日志文件了

#### 用logrotate进行日志切分
通过rsyslog输出的日志是不会进行切分的，所以需要依靠Linux提供的logrotate来进行切分工作

使用root用户，创建haproxy日志切分配置文件：
```
vim /etc/logrotate.d

/var/log/haproxy.log /var/log/haproxy_warn.log {  #切分的两个文件名
    daily        #按天切分
    rotate 7     #保留7份
    create 0644 ha ha  #创建新文件的权限、用户、用户组
    compress     #压缩旧日志
    delaycompress  #延迟一天压缩
    missingok    #忽略文件不存在的错误
    dateext      #旧日志加上日志后缀
    sharedscripts  #切分后的重启脚本只运行一次
    postrotate   #切分后运行脚本重载rsyslog，让rsyslog向新的日志文件中输出日志
      /bin/kill -HUP $(/bin/cat /var/run/syslogd.pid 2>/dev/null) &>/dev/null
    endscript
}
```
并配置在crontab中运行：
```
0 0 * * * /usr/sbin/logrotate /etc/logrotate.conf
```
logrotate 请参考: [如何在Ubuntu 16.04上用Logrotate管理日志文件](https://www.howtoing.com/how-to-manage-logfiles-with-logrotate-on-ubuntu-16-04)


#### 使用HAProxy搭建L4负载均衡器
HAProxy作为L4负载均衡器工作时，不会去解析任何与HTTP协议相关的内容，只在传输层对数据包进行处理。也就是说，以L4模式运行的HAProxy，无法实现根据URL向不同后端转发、通过cookie实现会话保持等功能。

同时，在L4模式下工作的HAProxy也无法提供监控页面。

但作为L4负载均衡器的HAProxy能够提供更高的性能，适合于基于套接字的服务（如数据库、消息队列、RPC、邮件服务、Redis等），或不需要逻辑规则判断，并已实现了会话共享的HTTP服务。


### HAProxy关键配置详解
##### 总览  
HAProxy的配置文件共有5个域

- global：用于配置全局参数
- default：用于配置所有frontend和backend的默认属性
- frontend：用于配置前端服务（即HAProxy自身提供的服务）实例
- backend：用于配置后端服务（即HAProxy后面接的服务）实例组
- listen：frontend+backend的组合配置，可以理解成更简洁的配置方法

##### global域的关键配置
- daemon：指定HAProxy以后台模式运行，通常情况下都应该使用这一配置
- user \[username\] ：指定HAProxy进程所属的用户
- group \[groupname\] ：指定HAProxy进程所属的用户组
- log \[address\] \[device\] \[maxlevel\] \[minlevel\]：日志输出配置，如log 127.0.0.1 local0 info warning，即向本机rsyslog或syslog的local0输出info到warning级别的日志。其中\[minlevel\]可以省略。HAProxy的日志共有8个级别，从高到低为 `emerg/alert/crit/err/warning/notice/info/debug`
- pidfile ：指定记录HAProxy进程号的文件绝对路径。主要用于HAProxy进程的停止和重启动作。
- maxconn ：HAProxy进程同时处理的连接数，当连接数达到这一数值时，HAProxy将停止接收连接请求

##### frontend域的关键配置
- acl \[name\] \[criterion\] \[flags\] \[operator\] \[value\]：定义一条ACL，ACL是根据数据包的指定属性以指定表达式计算出的true/false值。如"acl url_ms1 path_beg -i /ms1/"定义了名为url_ms1的ACL，该ACL在请求uri以/ms1/开头（忽略大小写）时为true
- bind \[ip\]:\[port\]：frontend服务监听的端口
- default_backend \[name\]：frontend对应的默认backend
- disabled：禁用此frontend
- http-request \[operation\] \[condition\]：对所有到达此frontend的HTTP请求应用的策略，例如可以拒绝、要求认证、添加header、替换header、定义ACL等等。
- http-response \[operation\] \[condition\]：对所有从此frontend返回的HTTP响应应用的策略，大体同上
- log：同global域的log配置，仅应用于此frontend。如果要沿用global域的log配置，则此处配置为log global
- maxconn：同global域的maxconn，仅应用于此frontend
- mode：此frontend的工作模式，主要有http和tcp两种，对应L7和L4两种负载均衡模式
- option forwardfor：在请求中添加X-Forwarded-For Header，记录客户端ip
- option http-keep-alive：以KeepAlive模式提供服务
- option httpclose：与http-keep-alive对应，关闭KeepAlive模式，如果HAProxy主要提供的是接口类型的服务，可以考虑采用httpclose模式，以节省连接数资源。但如果这样做了，接口的调用端将不能使用HTTP连接池
- option httplog：开启httplog，HAProxy将会以类似Apache HTTP或Nginx的格式来记录请求日志
- option tcplog：开启tcplog，HAProxy将会在日志中记录数据包在传输层的更多属性
- stats uri \[uri\]：在此frontend上开启监控页面，通过\[uri\]访问
- stats refresh \[time\]：监控数据刷新周期
- stats auth \[user\]:\[password\]：监控页面的认证用户名密码
- timeout client \[time\]：指连接创建后，客户端持续不发送数据的超时时间
- timeout http-request \[time\]：指连接创建后，客户端没能发送完整HTTP请求的超时时间，主要用于防止DoS类攻击，即创建连接后，以非常缓慢的速度发送请求包，导致HAProxy连接被长时间占用
- use_backend \[backend\] if|unless \[acl\]：与ACL搭配使用，在满足/不满足ACL时转发至指定的backend

#### backend域的关键配置
- acl：同frontend域
- balance \[algorithm\]：在此backend下所有server间的负载均衡算法，常用的有roundrobin和source，完整的算法说明见官方文档configuration.html#4.2-balance
- cookie：在backend server间启用基于cookie的会话保持策略，最常用的是insert方式，如cookie HA_STICKY_ms1 insert indirect nocache，指HAProxy将在响应中插入名为HA_STICKY_ms1的cookie，其值为对应的server定义中指定的值，并根据请求中此cookie的值决定转发至哪个server。indirect代表如果请求中已经带有合法的HA_STICK_ms1 cookie，则HAProxy不会在响应中再次插入此cookie，nocache则代表禁止链路上的所有网关和缓存服务器缓存带有Set-Cookie头的响应。
- default-server：用于指定此backend下所有server的默认设置。具体见下面的server配置。
- disabled：禁用此backend
- http-request/http-response：同frontend域
- log：同frontend域
- mode：同frontend域
- option forwardfor：同frontend域
- option http-keep-alive：同frontend域
- option httpclose：同frontend域
- option httpchk \[METHOD\] \[URL\] \[VERSION\]：定义以http方式进行的健康检查策略。如option httpchk GET /healthCheck.html HTTP/1.1
- option httplog：同frontend域
- option tcplog：同frontend域
- server \[name\] \[ip\]:\[port\] \[params\]：定义backend中的一个后端server，\[params\]用于指定这个server的参数，常用的包括有：
> check：指定此参数时，HAProxy将会对此server执行健康检查，检查方法在option httpchk中配置。同时还可以在check后指定inter, rise, fall三个参数，分别代表健康检查的周期、连续几次成功认为server UP，连续几次失败认为server DOWN，默认值是inter 2000ms rise 2 fall 3
> cookie \[value\]：用于配合基于cookie的会话保持，如cookie ms1.srv1代表交由此server处理的请求会在响应中写入值为ms1.srv1的cookie（具体的cookie名则在backend域中的cookie设置中指定）
> maxconn：指HAProxy最多同时向此server发起的连接数，当连接数到达maxconn后，向此server发起的新连接会进入等待队列。默认为0，即无限
> maxqueue：等待队列的长度，当队列已满后，后续请求将会发至此backend下的其他server，默认为0，即无限
> weight：server的权重，0-256，权重越大，分给这个server的请求就越多。weight为0的server将不会被分配任何新的连接。所有server默认weight为1 
- timeout connect \[time\]- ：指HAProxy尝试与backend server创建连接的超时时间
- timeout check \[time\]：默认情况下，健康检查的连接+响应超时时间为server命令中指定的inter值，如果配置了timeout check，HAProxy会以inter作为健康检查请求的连接超时时间，并以timeout check的值作为健康检查请求的响应超时时间
- timeout server \[time\]：指backend server响应HAProxy请求的超时时间

#### default域
上文所属的frontend和backend域关键配置中，除acl、bind、http-request、http-response、use_backend外，其余的均可以配置在default域中。default域中配置了的项目，如果在frontend或backend域中没有配置，将会使用default域中的配置。

#### listen域
listen域是frontend域和backend域的组合，frontend域和backend域中所有的配置都可以配置在listen域下

