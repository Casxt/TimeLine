# TimeLine

TimeLine是一个生活日记类软件。希望能成为Between类应用的一个替代品，可以用于记录生活中的每一个美好时刻。并将这些回忆串为一条时间线。

每一条Line可以加入多个用户，用户可以为Line添加Slice，Slice是记录的基本单位。每个Slice可以附带一些图片等媒体。

每一个用户也可拥有多条Line，例如和父母之间的回忆有一条Line，和女友/男友之间用另一条Line。

TimeLine也是一个开源学习项目，使用Go语言完成后端搭建。拥有自己的网络组件，所有http组件均基于官方库开发，没有引入第三方库。目前前端使用Boostrap，希望有人能在后台的基础上进行二次开发。

- [TimeLine](#timeline)
    - [核心功能](#%E6%A0%B8%E5%BF%83%E5%8A%9F%E8%83%BD)
        - [数据备份](#%E6%95%B0%E6%8D%AE%E5%A4%87%E4%BB%BD)
        - [记录日记](#%E8%AE%B0%E5%BD%95%E6%97%A5%E8%AE%B0)
        - [纪念日](#%E7%BA%AA%E5%BF%B5%E6%97%A5)
        - [设定一个可爱的昵称](#%E8%AE%BE%E5%AE%9A%E4%B8%80%E4%B8%AA%E5%8F%AF%E7%88%B1%E7%9A%84%E6%98%B5%E7%A7%B0)
    - [编译](#%E7%BC%96%E8%AF%91)
        - [环境依赖](#%E7%8E%AF%E5%A2%83%E4%BE%9D%E8%B5%96)
        - [编译环境](#%E7%BC%96%E8%AF%91%E7%8E%AF%E5%A2%83)
        - [使用到的第三方库](#%E4%BD%BF%E7%94%A8%E5%88%B0%E7%9A%84%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%93)
        - [编译](#%E7%BC%96%E8%AF%91)
    - [启动](#%E5%90%AF%E5%8A%A8)
    - [API列表](#api%E5%88%97%E8%A1%A8)
        - [SignUp](#signup)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [CheckAccount](#checkaccount)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [SignIn](#signin)
            - [登陆流程：](#%E7%99%BB%E9%99%86%E6%B5%81%E7%A8%8B%EF%BC%9A)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [CreateLine](#createline)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [AddSlice](#addslice)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [GetSlices](#getslices)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [GetLines](#getlines)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [AddUser](#adduser)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [UpdateProfilePic](#updateprofilepic)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [GetUserInfo](#getuserinfo)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
        - [ChangeNickName](#changenickname)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)
    - [文件API](#%E6%96%87%E4%BB%B6api)
            - [请求格式:](#%E8%AF%B7%E6%B1%82%E6%A0%BC%E5%BC%8F)
            - [返回格式：](#%E8%BF%94%E5%9B%9E%E6%A0%BC%E5%BC%8F%EF%BC%9A)

## 核心功能

### 数据备份

目前情侣应用的数据安全是一个大问题，Between在前段时间出现了跑路传闻，让我对之前上传的数据十分担心，出于这个直接原因，我完成了这个网站。

希望通过定期备份并发给的方式解决数据安全问题，但是在主题功能完成之前暂时还不会去做这个功能

### 记录日记

通过文字，图片，甚至是视频记录每一个幸福的时刻，并分享给你的亲人。这是增进感情的一个好方式。

### 纪念日

记录纪念日，并设定提醒周期

### 设定一个可爱的昵称

可爱的昵称是必要的

## 编译

### 环境依赖

- Mysql

### 编译环境

- git
- go 1.10-latest

### 使用到的第三方库

- Mysql for Golang
- Esay Mail

### 编译
``` shell
git clone https://github.com/Casxt/TimeLine.git
cd TimeLine
go bulid
```
## 启动

1. 将 config.template.json 重命名为 config.json

2. 将SQL的内容设置为你的数据库配置

3. Cert中配置证书，如果留空则不启用ssl

4. 填写ProjectPath为网站文件目录，例如你在/root/下执行 `git clone`，则 Path为 "/root/TimeLine"

5. 在目录下执行 `./TimeLine`

## API列表

所有的API均为json格式，Method为Post，图片除外

调用地址类似为 `example.com/api/[apiName]` 如SignUp调用地址为 `example.com/api/SignUp`

### SignUp
注册接口，Phone为用户手机，Mail为用户邮箱， HashPass 为 用户密码Hash256后的结果，
#### 请求格式:
``` golang
{
    Phone    string //用户手机
    Mail     string //用户邮箱
    HashPass string //用户密码Hash256的结果
}
```
#### 返回格式：
``` golang
{
    State  string //"Success" 成功 "Failed"失败
    Msg    string //成功或失败的相关信息
    Detail string //更详细的信息，不一定存在此字段
}
```

### CheckAccount
登陆的第一个步骤，检查手机或邮箱是否存在，并返回salt供后续步骤使用
#### 请求格式:
``` golang
{
    Account string //用户账户，可以是手机或者是邮箱
}
```
#### 返回格式：
``` golang
{
    State        string //"Success" 成功 "Failed"失败
    Msg          string //成功或失败的相关信息
    Detail       string //更详细的信息，不一定存在此字段
    NickName     string //昵称
    Salt         string //后续步骤需要使用
    SignInVerify string //后续步骤需要使用
}
```

### SignIn
登陆接口，需要先使用`CheckAccount`API获得`Salt`和`SignInVerify`后才可以使用
#### 登陆流程：
1. 对用户Pass做Hash256，获得sha256Pass
2. sha256Pass + Salt 做Hash256，获得sha256SaltPass
3. sha256SaltPass做key，并生成随机IV，对SignInVerify做AES256GCM加密，获得 Encrypt 和 Tag
4. Encrypted = HEX(Encrypt) + HEX(Tag)
#### 请求格式:
``` golang
{
    Encrypted string //加密后的字符串,16进制字符串
    IV        string //加密时用到的IV,16进制字符串(12字节原始数据HEX后得到的24字节字符串)
}
```
#### 返回格式：
``` golang
{
    State     string //"Success" 成功 "Failed"失败
    Msg       string //成功或失败的相关信息
    SessionID string //登陆成功的SessionID
    Phone     string //用户手机
    NickName  string //用户昵称
}
```

### CreateLine
创建一条新Line
#### 请求格式:
``` golang
{
		Operator  string //当前登陆用户的手机
		LineName  string //Line的名字
		SessionID string //与当前登陆用户对应的SessionID
}
```
#### 返回格式：
``` golang
{
    State  string //"Success" 成功 "Failed"失败
    Msg    string //成功或失败的相关信息
    Detail string //更详细的信息，不一定存在此字段
}
```

### AddSlice
新增一个Slice
#### 请求格式:
``` golang
{
    SessionID  string   //SessionID
    Operator   string   //用户手机
    LineName   string   //要添加到的Line名字
    Content    string   //Slice的文字内容
    Gallery    []string //Slice的附带图片
    Type       string   //种类，回忆还是纪念日
    Visibility string   //查看权限，自己可见，同Line可见，公开
    Longitude  string   //精度
    Latitude   string   //纬度
    Time       string   //记录时间
}
```
#### 返回格式：
``` golang
{
    State  string //"Success" 成功 "Failed"失败
    Msg    string //成功或失败的相关信息
    Detail string //更详细的信息，不一定存在此字段
}
```

### GetSlices
获取某个Line的GetSlices
#### 请求格式:
``` golang
{
    SessionID string
    Operator  string
    LineName  string //Line名
    PageNum   int    //页数，每页20条
}
```
#### 返回格式：
``` golang
{
    State  string //"Success" 成功 "Failed"失败
    Msg    string //成功或失败的相关信息
    Detail string //更详细的信息，不一定存在此字段
    Slices []database.SliceInfo //Slice信息数组
}
//database.SliceInfo：

{
    UserName   string    //创建者昵称
    Gallery    []string  //附带图片
    Content    string    //附带文字
    Type       string    //类型
    Visibility string    //查看权限
    Location   string    //定位
    Time       time.Time //记录时间
}
```

### GetLines
获取指定用户所属的全部Line
#### 请求格式:
``` golang
{
    SessionID string
    Operator  string
}
```
#### 返回格式：
``` golang
{
    State  string   //"Success" 成功 "Failed"失败
    Msg    string   //成功或失败的相关信息
    Detail string   //更详细的信息，不一定存在此字段
    Lines  []string //全部Line名称
}
```

### AddUser
添加用户到指定Line
#### 请求格式:
``` golang
{
    Operator  string //当前登录用户
    SessionID string //当前登录用户的SessionID
    LineName  string //Line名
    NickName  string //被添加的用户昵称
    UserPhone string //被添加用户手机
}
```
#### 返回格式：
``` golang
{
    State  string //"Success" 成功 "Failed"失败
    Msg    string //成功或失败的相关信息
    Detail string //更详细的信息，不一定存在此字段
}
```
### UpdateProfilePic
设置用户头像
#### 请求格式:
``` golang
{
    SessionID string
    Operator  string
    Picture   string //图像编码
}
```
#### 返回格式：
``` golang
{
    State  string //"Success" 成功 "Failed"失败
    Msg    string //成功或失败的相关信息
    Detail string //更详细的信息，不一定存在此字段
}
```

### GetUserInfo
获取用户账户信息
#### 请求格式:
``` golang
{
    SessionID string
    Operator  string
}
```
#### 返回格式：
``` golang
{
    State      string //"Success" 成功 "Failed"失败
    Msg        string //成功或失败的相关信息
    Detail     string //更详细的信息，不一定存在此字段
    NickName   string    //昵称
    Phone      string    //手机
    Mail       string    //邮箱
    Gender     string    //性别
    ProfilePic string    //头像
    SignInTime time.Time //注册时间
}
```

### ChangeNickName
更改昵称
#### 请求格式:
``` golang
{
    SessionID string
    Operator  string
    NewName   string //新昵称
}
```
#### 返回格式：
``` golang
{
    State  string //"Success" 成功 "Failed"失败
    Msg    string //成功或失败的相关信息
    Detail string //更详细的信息，不一定存在此字段
}
```

## 文件API
文件API使用`multipart/form-data`POST格式，URL为`/image`
#### 请求格式:
POST的Cookie中应包含Operator字段和SessionID字段

所有图片的字段名称为images

单张图片不小于512字节，一次发送的图片总大小最大为20MB，一次最多传输9张
#### 返回格式：
返回格式为json字符串
``` golang
{
    State string   //"Success" 成功 "Failed"失败
    Msg   string   //成功或失败的相关信息
    Hashs []string //图片ID数组
}
```