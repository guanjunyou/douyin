#  抖音简易后端
### 字节跳动青训营 GOGOG O队项目

## 前期准备进度
by 关竣佑 谢声儒

* 集成 GORM 框架
* 准备实体类 Model  
  所有的实体类均继承 CommonEntity (Id , CreateTime , IsDelete)
* 数据库设计  
  所有删除均采用逻辑删除 （规定删除后is_delete字段为 1 ）  
  所有的主键 id  采用bingint 对应的 go 实体类使用 int64 , 生成的时候使用雪花算法

* 分级结构 ： controller -   service -  serviceImpl - model  （MVC架构）
* 使用 jwt 生成 token , 注册/登录时将 token 存在 redis 中 



接口文档

https://apifox.com/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707521

## 开发规约

#### 禁止/必须

1. 主体逻辑代码必须放在service层中的impl层，禁止在controller层写过多大的业务的代码，controller层应尽量调用service层的方法实现业务逻辑

2. model 层的函数禁止调用其它model 层相同包下不同 model 的函数

3. 返回给前端的数据若要组装成一个 stauct  必须使用 xxxDVO来命名，参见 models.VideoDVO

4. model中 禁止进行sql字符串拼接，避免造成sql注入风险，如需使用参数拼接必须使用  ？ 传参   如  

   ```go
   err := utils.DB.Where("is_deleted != ?", 1).Find(&videolist).Error
   ```

5. 遇到的所有 error 返回都必须进程处理或返回给上级（如使用 log 输出日志）

   ```go
   		if err1 != nil {
   			log.Printf("Can not get the token!")
   		}
   ```

   

6. 所需用到的参数均放在config.go中，禁止在代码中出现魔法值。（所谓魔法值，是代码中莫名其妙出现的数字，数字意义必须通过阅读其他代码才能推断出来，这样给后期维护或者其他人员阅读代码，带来了极大不便。)如以下代码便出现了魔法值

   ```go
   // 遍历查询出的审查人对象集合
           for(AuditPersonInfoDTO adp : auditPersonInfoDTO){
               // 判断审查结果是否为空
               if(adp.getStatus()!=null){
                   // 设置审查状态，status为2代表审核通过，为3代表退回修改
                   switch (adp.getStatus()){
                       case "2" :
                           adp.setStatus("审查通过");
                           break;
                       case "3" :
                           adp.setStatus("退回修改");
                           break;
   ......
   ```

7. 每次开发前都必须pull代码！！！不然可能会造成冲突，很难解决。尽量先新建一个分支，测试功能正常后再与main分支合并

8. 禁止对已有文件进行移动（比如说移到其它包内），如需对结构有较大修改请提前说明

9. 每次 push 代码时禁止直接提交到 Master 分支 ！必须新建分支，运行测试正常后再提交分支！合并分支时遇到冲突需慎重解决，不明白的及时提出或让其他人帮忙合并

10. 所有实体类的成员必须使用**首字母大写**的驼峰命名法，Go 语言只用大写首字母才能被其它包访问。

11. 如需更改数据库请提前说明！
12. 如需提交更改后的数据库禁止删掉之前的数据库文件，以 日期-版本号.sql命名 (如：2023-7-21-v1douyin.sql)
13. 分支合并之后必须删除GitHub上的分支，每个人在GitHub上最多拥有一个分支
14. 校验token是否存在且合法使用  utils 包的 AuthAdminCheck 函数
15. 编写接口时返回的数据一定要按照接口文档要求返回的数据



#### 建议

1. 推荐使用 Goland 进行开发，使用Goland 的 git 图形化工具操作 git 

​	2.合并分支解决冲突的时候如遇不理解的问题及时提出

3. 开发一个函数后，建议在 test 包下编写测试代码进行测试
3. 如果业务操作间没有太多的关联，建议开启协程，使用 channel 通信。
3. 创建切片数组前，如果能估计大小，建议预先设置好大小，减少后期扩容开销



#### 注意

1. 请求格式特别是 POST 请求的格式参照原本的代码。它里面有的POST请求不放json而使用拼接URL（我也不知道为什么），这里很坑

# 接口基本思路

## 互动接口

### 赞操作   （王奕丹）

URL：**POST****/douyin/favorite/action/**

基本思路：主要操控like表，当action_type=1时写入对应的一条点赞关系记录，反之则删除

### 喜欢列表（王奕丹）

URL:  **GET****/douyin/favorite/list/**

基本思路：主要操控like表、user表和veido表，查出用户所喜欢的所有视频id，根据视频id进一步查询作者信息、视频信息，具体字段需求查看api文档，建议封装成DTO类

### 评论操作  （邱祥凯）

URL: **POST****/douyin/comment/action/**

基本思路：主要操控comment表，将对应的评论信息添加到数据库中即可。

### 评论列表（邱祥凯）

URL: **GET****/douyin/comment/list/**

基本思路：查询出conmment表中对应视频id的所有记录，按照创建时间进行倒序排序。如果想要通过redis进行优化，可以使用zset数据类型，该类型可以存入键值对，可以根据值进行排序，为一个有序集合

## 社交接口

### 关注操作（杨伟宁）

URL: **POST****/douyin/relation/action/**

基本思路：与点赞操作类似，只是操控的数据表变成follo表

### 关注列表（杨伟宁）

URL: **GET****/douyin/relation/follow/list/**

基本思路：根据请求中的user_id联合查询follow表和user表，返回对应的关注用户信息集合

### 粉丝列表（杨伟宁）

URL: **GET****/douyin/relation/follow/list/**

基本思路：与关注列表其实逻辑类似的,操控表也一样

### 好友列表（杨伟宁）

URL: **GET****/douyin/relation/friend/list/**

基本思路：根据现在的抖音的定义，只有两个用户相互关注才是好友，那么我们就可以先查询当前用户的关注列表，对列表中的每个用户判断其是否也关注了当前用户，将不符合条件的用户过滤。最后得到的列表就是好友列表

### 发送消息（邱祥凯）

URL: **POST****/douyin/message/action/**

基本思路：逻辑很简单，我们只需要将记录 存入3张表：message、message_push_event、message_send_event

### 聊天记录  （邱祥凯）

URL: **GET****/douyin/message/chat/**

基本思路：根据当前用户id以及to_user_id从message_send_event表中查出对应的聊天记录返回即可

# 登录鉴权

​	全局鉴权采用两层中间件(即spring中的拦截器)完成。第一层拦截器用于刷新redis中的token有效期，无论什么请求都放行到第二个中间件处理；第二层拦截器用于真正的用户鉴权，此时若用户在未登录状态访问了非法资源则会立刻拒绝该请求。现在所有请求的鉴权操作都会在拦截器中自动完成。
