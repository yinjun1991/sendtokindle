# 项目介绍

sendtokindle 是在本地启动 http server，你可以在网页上传电子书，再在 kindle 打开网页下载电子书的网站

## 技术栈

1. go + gin
2. 纯 html/css/js，不要外部依赖

## 交互

1. 提供 bin 文件，一键（支持 port 参数）在本地启动 http server 并打开 admin html
2. admin html 可以预览本地文件并上传到 http server（在 ~/.sendtokindle 目录下）
3. 在 kindle 打开 http server 主页，页面显示 http server 包含的电子书，可以点击下载单个电子书，或者批量下载；注意网页一定要简单，符合 kindle 浏览器的技术要求、页面适配；但是页面还是要简约美观，notion 风格
4. 2个 html 页面和 go 代码都打包到一个 bin 文件中；用法越简单越好
