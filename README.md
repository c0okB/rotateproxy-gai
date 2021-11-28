# rotateproxy-gai

参考：https://github.com/akkuman/rotateproxy  对原作者的工具进行了修改

将fofa页面上的socks5代理的爬虫数量设置成自定义



修改结果如下：

![1.JPG](https://github.com/c0okB/rotateproxy-gai/blob/main/1.jpg)
当选择type 1 时，将会爬取socks5 url 至数据库
当选择type 2 时，将会对socks5 url 进行可用检测，并启动端口监听


![1.JPG](https://github.com/c0okB/rotateproxy-gai/blob/main/2.jpg)

![1.JPG](https://github.com/c0okB/rotateproxy-gai/blob/main/3.jpg)

默认端口为 8899
默认规则为 `protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && after="2021-01-01" && country="CN"`
默认爬取数量为 3500
