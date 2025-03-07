package use_AdminInterface

type AdminServiceInterface interface {
	Register(username, password string) error
	Login(username string, password string) (string, map[string]interface{}, error)
	GetAdminInfo(username string) (*map[string]interface{}, error)
	UpdateAdmin(username string, updates map[string]interface{}) error
	PaginateAdmins(page, pageSize int) ([]map[string]interface{}, int64, error)
}
