package dto

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 用户相关dto
 * @Date: 2024-10-04 19:41
 */
type UserRegisterRequstDTO struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	PassWord string `json:"password"`
}
type UserRegisterResponceDTO struct {
	UserId   string `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名，唯一标识用户
}
type UserLoginRequstDTO struct {
	Email    string `json:"email"`
	PassWord string `json:"password"`
}
type UserLoginResponceDTO struct {
	UserId   string `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名，唯一标识用户
}
type UserGetInfoRequstDTO struct {
	UserId string `json:"user_id"` // 用户ID
}
type UserGetInfoResponceDTO struct {
	UserId   string `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名，唯一标识用户
	Email    string `json:"email"`    // 用户邮箱，唯一
	Credit   int64  `json:"credit"`   // 用户积分余额
}
