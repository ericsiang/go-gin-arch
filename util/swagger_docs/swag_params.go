package swagger_docs

// DocUsersLogin 用户登录参数
type DocUsersLogin struct {
	Account  string `example:"" json:"account"`
	Password string `example:"" json:"password"`
}

// DocUsersCreate 用户创建参数
type DocUsersCreate struct {
	Account  string `example:"" json:"account"`
	Password string `example:"" json:"password"`
}

// DocAdminsLogin 管理员登录参数
type DocAdminsLogin struct {
	Account  string `example:"" json:"account"`
	Password string `example:"" json:"password"`
}

// DocAdminsCreate 管理员创建参数
type DocAdminsCreate struct {
	Account  string `example:"" json:"account"`
	Password string `example:"" json:"password"`
}
