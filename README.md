#项目名称：NAS

##基本介绍：</br>
>项目目的是为了实现局域网的存储</br>
>待实现smb共享</br>
>目前http上传下载已实现</br>
>实现了监控服务器状态</br>
>本项目目前没有前端页面（期待大佬完成）</br>
>本人刚接触go语言代码质量可怜，欢迎大佬指正</br>

##感谢以下开源库:</br>
>`github.com/Unknwon/goconfig v1.0.0`</br>
>`github.com/dgrijalva/jwt-go v3.2.0+incompatible`</br>
>`github.com/disintegration/imaging v1.6.2`</br>
>`github.com/gin-contrib/gzip v0.0.5`</br>
>`github.com/gin-gonic/gin v1.7.7`</br>
>`github.com/go-redis/redis/v8 v8.11.4`</br>
>`github.com/gorilla/websocket v1.4.2`</br>
>`github.com/jinzhu/gorm v1.9.16`</br>
>`github.com/satori/go.uuid v1.2.0`</br>
>`github.com/shirou/gopsutil/v3 v3.21.12`</br>

如果您在使用本项目的过程中遇到问题欢迎提交issue</br>
如果您也是开源爱好者，并且可以看懂我这堆shi山，欢迎优化代码pull给我</br>

##开发者API接口
>登录（有大佬可以加一个验证码）</br>
> http://127.0.0.1:8080/T/login </br>
> ``` json 
> {"Username": "XXXX", "Password": "XXXXX"}
> ```
> 
> 上传确认 </br>
> http://127.0.0.1:8080/File/Uploads </br>
> ``` json
> {
>    "Name":"图片_202124.jpg",
>    "Sha512":"f394202bd0be442c5826ba56de7b77e47e7299ba695f3fa0d3c99",
>    "Size":99731,
>    "PID":"c72b05cb-ff16-481b-814d-e0b708916ed0",
> }
> ```
> 
> 文件上传 </br>
> >http://127.0.0.1:8080/File/Upload </br>
> >body key=file
> 
> 列表 </br>
> >http://127.0.0.1:8080/File/List/:pid </br>
> >例如：http://127.0.0.1:8080/File/List/0 </br>
> >pid=上级目录uuid </br>
> >0为root
> 
> 缩略图 </br>
> >http://127.0.0.1:8080/File/Image/:uuid(图片文件uuid)
> 
> 下载 </br>
> >http://127.0.0.1:8080/File/Download/:uuid(要下载文件的uuid
> 
> 重命名文件</br>
> >http://127.0.0.1:8080/File/RenameFile </br>
> ```json
> {
>    "New":"123.jpg",
>    "Name":"图片_202124.jpg",
>    "PID":"0"
> }
> ```
> 删除 </br>
> >http://127.0.0.1:8080/File/DelFile </br>
> ```json
> {
>   "Name":"微信图片_20211215223442.jpg",
>   "Pid":"0"
> }
> ```
> 获取系统信息 </br>
> >http://127.0.0.1:8080/System/sys-info
> ###更多接口请阅读源代码，调用方法大同小异或者提交issue