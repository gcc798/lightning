// Package main 提供 Lightning 首个管理员创建和用户密码重置工具。
//
// 密码只能通过 LIGHTNING_USERMGR_PASSWORD 环境变量传入，不接受命令行参数，
// 避免密码出现在 shell history 和进程列表中。
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/database"
	"github.com/gcc798/lightning/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type user struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement:false" autogen:"int64"`
	UserName string `gorm:"column:user_name"`
	NickName string `gorm:"column:nick_name"`
	UserType int32  `gorm:"column:user_type"`
	Status   int32  `gorm:"column:status"`
	Sex      int32  `gorm:"column:sex"`
	Password string `gorm:"column:password"`
}

func (user) TableName() string { return "s_user" }

type role struct {
	ID int64 `gorm:"column:id"`
}

func (role) TableName() string { return "s_role" }

type userRole struct {
	Id     int64 `gorm:"column:id;primaryKey;autoIncrement:false" autogen:"int64"`
	UserId int64 `gorm:"column:user_id"`
	RoleId int64 `gorm:"column:role_id"`
}

func (userRole) TableName() string { return "m_user_role" }

const passwordEnv = "LIGHTNING_USERMGR_PASSWORD"

type options struct {
	operation string
	username  string
	nickname  string
	role      string
	configDir string
}

func main() {
	opts, err := parseOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	password := os.Getenv(passwordEnv)
	if err := validateInput(opts, password); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	cfg, _, err := config.Load(opts.configDir, config.ServiceUserManager)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	db, err := initDB(cfg.Database.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接数据库失败: %v\n", err)
		os.Exit(1)
	}

	switch opts.operation {
	case "create":
		err = createUser(db, opts, password)
	case "reset":
		err = resetPassword(db, opts.username, password)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "操作失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("操作成功: operation=%s username=%s\n", opts.operation, opts.username)
}

func parseOptions(args []string) (options, error) {
	var opts options
	set := flag.NewFlagSet("usermgr", flag.ContinueOnError)
	set.SetOutput(os.Stderr)
	set.StringVar(&opts.operation, "operation", "", "操作类型: create 或 reset")
	set.StringVar(&opts.username, "username", "", "用户名")
	set.StringVar(&opts.nickname, "nickname", "管理员", "创建用户时的昵称")
	set.StringVar(&opts.role, "role", "super_admin", "创建用户时分配的角色标识")
	set.StringVar(&opts.configDir, "config-dir", "application/iam", "IAM 配置文件目录")
	if err := set.Parse(args); err != nil {
		return options{}, err
	}
	return opts, nil
}

func validateInput(opts options, password string) error {
	if opts.operation != "create" && opts.operation != "reset" {
		return errors.New("--operation 必须为 create 或 reset")
	}
	if strings.TrimSpace(opts.username) == "" {
		return errors.New("--username 不能为空")
	}
	if len(password) < 8 {
		return fmt.Errorf("%s 必须至少包含 8 个字符", passwordEnv)
	}
	if opts.operation == "create" {
		if strings.TrimSpace(opts.nickname) == "" {
			return errors.New("--nickname 不能为空")
		}
		if strings.TrimSpace(opts.role) == "" {
			return errors.New("--role 不能为空")
		}
	}
	return nil
}

func initDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.Use(&database.IDGenPlugin{}); err != nil {
		return nil, fmt.Errorf("注册 ID 生成器失败: %w", err)
	}
	return db, nil
}

func createUser(db *gorm.DB, opts options, password string) error {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&user{}).Where("user_name = ?", opts.username).Count(&count).Error; err != nil {
			return fmt.Errorf("检查用户名失败: %w", err)
		}
		if count != 0 {
			return fmt.Errorf("用户名 %q 已存在", opts.username)
		}

		var role role
		if err := tx.Where("role_key = ? AND status = 0", opts.role).First(&role).Error; err != nil {
			return fmt.Errorf("查找角色 %q 失败: %w", opts.role, err)
		}

		user := user{
			UserName: opts.username,
			NickName: opts.nickname,
			Password: passwordHash,
			UserType: 0,
			Status:   0,
			Sex:      2,
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
		if err := tx.Create(&userRole{UserId: user.ID, RoleId: role.ID}).Error; err != nil {
			return fmt.Errorf("分配角色失败: %w", err)
		}
		return nil
	})
}

func resetPassword(db *gorm.DB, username, password string) error {
	var existing user
	err := db.Where("user_name = ?", username).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户名 %q 不存在", username)
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}
	if err := db.Model(&user{}).Where("id = ?", existing.ID).Update("password", passwordHash).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}
	return nil
}
