#  抖音简易后端
### 字节跳动青训营 GOGOG 队项目

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



#### 建议

1. 推荐使用 Goland 进行开发，使用Goland 的 git 图形化工具操作 git 

​	2.合并分支解决冲突的时候如遇不理解的问题及时提出

3. 开发一个函数后，建议在 test 包下编写测试代码进行测试



#### 注意

1. 请求格式特别是 POST 请求的格式参照原本的代码。它里面有的POST请求不放json而使用拼接URL（我也不知道为什么），这里很坑
