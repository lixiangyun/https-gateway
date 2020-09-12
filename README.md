# https-gateway

## 准备工作
- https加固的云主机
- {youdomain.com}域名解析云主机IP地址成功
- 关闭80、443端口占用程序
- 配置加固的服务为本地地址，例如：192.168.0.1:8080


## 运行命令

```
docker run --net=host --name https -d -v /home/letsencrypt:/etc/letsencrypt -v /home/https:/opt/home linimbus/https-gateway 
```

## 操作

```
http://{youip}:18000/
```
