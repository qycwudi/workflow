### 1. "用户登录"

1. route definition

- Url: /workflow/user/login
- Method: POST
- Request: `UserLoginRequest`
- Response: `UserLoginResponse`

2. request definition



```golang
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
```


3. response definition



```golang
type UserLoginResponse struct {
	Token string `json:"token"`
}
```

### 2. "用户信息"

1. route definition

- Url: /workflow/user/info
- Method: POST
- Request: `UserInfoRequest`
- Response: `UserInfoResponse`

2. request definition



```golang
type UserInfoRequest struct {
}
```


3. response definition



```golang
type UserInfoResponse struct {
	User User `json:"user"`
}

type User struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Status int64 `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
```

### 3. "用户注册"

1. route definition

- Url: /workflow/user/register
- Method: POST
- Request: `UserRegisterRequest`
- Response: `UserRegisterResponse`

2. request definition



```golang
type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RealName string `json:"realName,optional"`
	Phone string `json:"phone,optional"`
	Email string `json:"email,optional"`
	Status int64 `json:"status,optional"`
}
```


3. response definition



```golang
type UserRegisterResponse struct {
	Token string `json:"token"`
}
```

### 4. "用户退出登录"

1. route definition

- Url: /workflow/user/logout
- Method: POST
- Request: `UserLogoutRequest`
- Response: `UserLogoutResponse`

2. request definition



```golang
type UserLogoutRequest struct {
}
```


3. response definition



```golang
type UserLogoutResponse struct {
}
```

### 5. "获取用户列表"

1. route definition

- Url: /workflow/user/list
- Method: POST
- Request: `UserListRequest`
- Response: `UserListResponse`

2. request definition



```golang
type UserListRequest struct {
	Current int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
	Username string `json:"username,optional"`
}
```


3. response definition



```golang
type UserListResponse struct {
	Total int64 `json:"total"`
	List []User `json:"list"`
}
```

### 6. "绑定角色"

1. route definition

- Url: /workflow/user/bindrole
- Method: POST
- Request: `UserBindRoleRequest`
- Response: `UserBindRoleResponse`

2. request definition



```golang
type UserBindRoleRequest struct {
	UserId int64 `json:"userId"`
	RoleId int64 `json:"roleId"`
}
```


3. response definition



```golang
type UserBindRoleResponse struct {
}
```

