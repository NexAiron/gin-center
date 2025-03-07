# Gin-Center API 文档

## 用户接口

### 用户登录
- 路径: `/api/v1/user/login`
- 方法: POST
- 权限: 公开
- 描述: 处理用户登录请求，验证用户名和密码，返回JWT令牌
- 请求参数:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- 响应:
  - 200: 登录成功，返回JWT令牌
  - 400: 请求参数错误
  - 401: 登录失败

### 用户注册
- 路径: `/api/v1/user/register`
- 方法: POST
- 权限: 公开
- 描述: 处理用户注册请求，创建新用户账户
- 请求参数:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- 响应:
  - 200: 注册成功
  - 400: 请求参数错误
  - 500: 注册失败

## 管理员接口

### 管理员登录
- 路径: `/admin/login`
- 方法: POST
- 权限: 公开
- 描述: 管理员身份认证
- 请求参数:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```

### 管理员注册
- 路径: `/admin/register`
- 方法: POST
- 权限: 管理员
- 描述: 创建新管理员账号
- 请求参数:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```

### 获取管理员信息
- 路径: `/admin/info`
- 方法: GET
- 权限: 管理员
- 描述: 获取当前登录管理员的详细信息

### 更新管理员信息
- 路径: `/admin`
- 方法: PUT
- 权限: 管理员
- 描述: 更新管理员基本信息

## 系统管理接口

### 获取系统信息
- 路径: `/system/info`
- 方法: GET
- 权限: 管理员
- 描述: 获取系统基本信息

### 获取系统配置
- 路径: `/system/config`
- 方法: GET
- 权限: 管理员
- 描述: 获取系统配置详情

### 更新系统配置
- 路径: `/system/config`
- 方法: PUT
- 权限: 超级管理员
- 描述: 修改系统配置
- 请求参数: SystemConfig对象

### 获取系统指标
- 路径: `/system/metrics`
- 方法: GET
- 权限: 管理员
- 描述: 获取系统运行指标数据

### 获取系统健康状态
- 路径: `/system/health`
- 方法: GET
- 权限: 管理员
- 描述: 获取系统健康状态信息