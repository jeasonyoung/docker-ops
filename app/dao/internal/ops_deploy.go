// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// OpsDeployDao is the manager for logic model data accessing and custom defined data operations functions management.
type OpsDeployDao struct {
	Table   string           // Table is the underlying table name of the DAO.
	Group   string           // Group is the database configuration group name of current DAO.
	Columns OpsDeployColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// OpsDeployColumns defines and stores column names for table tbl_ops_deploy.
type OpsDeployColumns struct {
	Id           string // 部署ID
	Name         string // 部署名称
	Remark       string // 部署描述
	GroupId      string // 服务器分组ID
	RepositoryId string // 所属镜像仓库ID
	Command      string // 附加执行命令
	Status       string // 状态(-1:删除,0:停用,1:启用)
	CreateTime   string // 创建时间
	UpdateTime   string // 更新时间
}

//  opsDeployColumns holds the columns for table tbl_ops_deploy.
var opsDeployColumns = OpsDeployColumns{
	Id:           "id",
	Name:         "name",
	Remark:       "remark",
	GroupId:      "group_id",
	RepositoryId: "repository_id",
	Command:      "command",
	Status:       "status",
	CreateTime:   "create_time",
	UpdateTime:   "update_time",
}

// NewOpsDeployDao creates and returns a new DAO object for table data access.
func NewOpsDeployDao() *OpsDeployDao {
	return &OpsDeployDao{
		Group:   "default",
		Table:   "tbl_ops_deploy",
		Columns: opsDeployColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OpsDeployDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OpsDeployDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OpsDeployDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
