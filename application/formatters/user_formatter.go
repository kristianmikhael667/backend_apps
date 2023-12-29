package formatters

type (
	RequestBodyAuthAdmin struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	ResponseAdminUser struct {
		Status string `json:"status"`
		Token  string `json:"token"`
		User   User   `json:"user"`
	}

	User struct {
		Uid       string `json:"uid"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		Status    int8   `json:"status" gorm:"int2"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)
