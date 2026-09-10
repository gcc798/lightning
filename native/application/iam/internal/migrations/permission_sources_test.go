package migrations

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/gcc798/microservice-kit/application/iam/internal/domain"
	"github.com/gcc798/microservice-kit/application/iam/internal/domain/model"
	"github.com/gcc798/microservice-kit/internal/database"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type permissionTestLogger struct{}

func (permissionTestLogger) Get() *zap.Logger                   { return zap.NewNop() }
func (permissionTestLogger) Debug(string, ...zap.Field)         {}
func (permissionTestLogger) Info(string, ...zap.Field)          {}
func (permissionTestLogger) Warn(string, ...zap.Field)          {}
func (permissionTestLogger) Error(string, ...zap.Field)         {}
func (permissionTestLogger) Fatal(string, ...zap.Field)         {}
func (l permissionTestLogger) With(...zap.Field) logging.Logger { return l }

func TestManualAndDerivedPermissionSourcesAreIndependent(t *testing.T) {
	dsn := os.Getenv("MS_K_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("MS_K_TEST_POSTGRES_DSN is not set")
	}
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	if err := Up(sqlDB); err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Use(&database.IDGenPlugin{}); err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer tx.Rollback()

	role := model.Role{RoleKey: "permission-source-test", RoleName: "permission-source-test"}
	menu := model.Menu{MenuName: "permission-source-test", MenuType: 1}
	permissions := []model.ApiPermission{
		{Module: "test", Code: "test.one", Name: "test.one", NodeType: 2, Action: "read"},
		{Module: "test", Code: "test.two", Name: "test.two", NodeType: 2, Action: "read"},
	}
	for _, value := range []any{&role, &menu, &permissions} {
		if err := tx.Create(value).Error; err != nil {
			t.Fatal(err)
		}
	}
	links := []model.MMenuApiPermission{
		{MenuId: menu.ID, PermissionId: permissions[0].ID},
		{MenuId: menu.ID, PermissionId: permissions[1].ID},
	}
	if err := tx.Create(&links).Error; err != nil {
		t.Fatal(err)
	}

	roles := iam.NewRoleService(tx, permissionTestLogger{})
	api := iam.NewApiPermissionService(tx)
	ctx := context.Background()
	if err := roles.AssignMenusToRole(ctx, role.ID, []int64{menu.ID}); err != nil {
		t.Fatal(err)
	}
	if err := api.AssignRolePermissions(ctx, role.ID, []int64{permissions[0].ID}, 1); err != nil {
		t.Fatal(err)
	}
	assertPermissionSourceCount(t, tx, role.ID, model.PermissionSourceManual, 1)
	assertPermissionSourceCount(t, tx, role.ID, model.PermissionSourceDerived, 2)

	if err := api.AssignRolePermissions(ctx, role.ID, nil, 1); err != nil {
		t.Fatal(err)
	}
	assertPermissionSourceCount(t, tx, role.ID, model.PermissionSourceManual, 0)
	assertPermissionSourceCount(t, tx, role.ID, model.PermissionSourceDerived, 2)

	if err := api.AssignRolePermissions(ctx, role.ID, []int64{permissions[0].ID}, 1); err != nil {
		t.Fatal(err)
	}
	if err := roles.AssignMenusToRole(ctx, role.ID, nil); err != nil {
		t.Fatal(err)
	}
	assertPermissionSourceCount(t, tx, role.ID, model.PermissionSourceManual, 1)
	assertPermissionSourceCount(t, tx, role.ID, model.PermissionSourceDerived, 0)
}

func assertPermissionSourceCount(t *testing.T, db *gorm.DB, roleID int64, source int32, want int64) {
	t.Helper()
	var got int64
	if err := db.Model(&model.MRoleApiPermission{}).
		Where("role_id = ? AND source = ?", roleID, source).Count(&got).Error; err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("permission source %d count = %d, want %d", source, got, want)
	}
}
