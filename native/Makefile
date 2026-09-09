.PHONY: help init-config proto proto-check swagger swagger-check swagger-fmt run run-scheduler build test test-race vet docker-build verify ci clean

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make swagger      - 生成 Swagger 文档"
	@echo "  make init-config  - 从 example 配置创建本地 dev/prod 配置"
	@echo "  make proto        - 生成 gRPC 契约代码"
	@echo "  make swagger-fmt  - 格式化 Swagger 注释"
	@echo "  make run          - 运行 IAM 服务（自动生成文档）"
	@echo "  make run-scheduler - 运行独立 Scheduler 服务"
	@echo "  make build        - 编译项目"
	@echo "  make test         - 运行测试"
	@echo "  make test-race    - 运行竞态检测"
	@echo "  make vet          - 运行 Go 静态检查"
	@echo "  make swagger-check - 检查 Swagger 是否为最新"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make verify       - 执行本地完整质量检查"
	@echo "  make ci           - 执行测试、vet、race、Swagger freshness 和 Docker build"
	@echo "  make clean        - 清理生成的文件"

# 创建本地运行配置
init-config:
	@for service in gateway iam sys resource scheduler; do \
		for kind in conf zaplogger; do \
			for profile in dev prod; do \
				target="application/$$service/$$kind.$$profile.yaml"; \
				if [ ! -e "$$target" ]; then cp "application/$$service/$$kind.example.yaml" "$$target"; fi; \
			done; \
		done; \
	done
	@echo "已从 example 配置创建缺失的 dev/prod 文件"

# 生成 Swagger 文档
swagger:
	@swag init -g main.go -d application/gateway,application/iam,application/sys,application/resource,internal/httpresponse,internal/utils -o internal/openapi --parseDependency --parseInternal
	@echo "✅ Swagger 文档已生成: internal/openapi/"

proto:
	@protoc -I . --go_out=. --go_opt=module=github.com/gcc798/lightning --go-grpc_out=. --go-grpc_opt=module=github.com/gcc798/lightning api/iam/v1/iam.proto api/sys/v1/sys.proto api/resource/v1/resource.proto

proto-check:
	@tmp_dir=$$(mktemp -d); \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	protoc -I . --go_out="$$tmp_dir" --go_opt=module=github.com/gcc798/lightning --go-grpc_out="$$tmp_dir" --go-grpc_opt=module=github.com/gcc798/lightning api/iam/v1/iam.proto api/sys/v1/sys.proto api/resource/v1/resource.proto; \
	diff -ru internal/api/iam/v1/iam.pb.go "$$tmp_dir/internal/api/iam/v1/iam.pb.go"; \
	diff -ru internal/api/iam/v1/iam_grpc.pb.go "$$tmp_dir/internal/api/iam/v1/iam_grpc.pb.go"; \
	diff -ru internal/api/sys/v1/sys.pb.go "$$tmp_dir/internal/api/sys/v1/sys.pb.go"; \
	diff -ru internal/api/sys/v1/sys_grpc.pb.go "$$tmp_dir/internal/api/sys/v1/sys_grpc.pb.go"; \
	diff -ru internal/api/resource/v1/resource.pb.go "$$tmp_dir/internal/api/resource/v1/resource.pb.go"; \
	diff -ru internal/api/resource/v1/resource_grpc.pb.go "$$tmp_dir/internal/api/resource/v1/resource_grpc.pb.go"

swagger-check:
	@tmp_dir=$$(mktemp -d); \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	swag init -g main.go -d application/gateway,application/iam,application/sys,application/resource,internal/httpresponse,internal/utils -o "$$tmp_dir/openapi" --parseDependency --parseInternal >/dev/null; \
	for file in docs.go swagger.json swagger.yaml; do diff -u internal/openapi/$$file "$$tmp_dir/openapi/$$file"; done

# 格式化 Swagger 注释
swagger-fmt:
	@echo "Swagger 注释随各服务 controller 维护"

# 运行服务（先生成文档）
run: swagger
	@echo "正在启动 IAM 服务..."
	@go run ./application/iam

run-scheduler:
	@echo "正在启动 Scheduler 服务..."
	@go run ./application/scheduler

# 编译项目
build: swagger
	@echo "正在编译项目..."
	@go build -o bin/lightning-scheduler ./application/scheduler
	@go build -o bin/lightning-gateway ./application/gateway
	@go build -o bin/lightning-iam ./application/iam
	@go build -o bin/lightning-sys ./application/sys
	@go build -o bin/lightning-resource ./application/resource
	@echo "✅ 网关、IAM、SYS、Resource、Scheduler 编译完成"

# 运行测试
test:
	@echo "正在运行测试..."
	@go test -v ./...

test-race:
	@go test -race ./...

vet:
	@go vet ./...

docker-build: init-config
	@for target in gateway iam sys resource; do docker build -f Dockerfile.service --build-arg TARGET=$$target -t lightning-native-$$target:verify . || exit 1; done
	@docker build -f Dockerfile.scheduler -t lightning-native-scheduler:verify .

verify: test vet test-race proto-check swagger-check

ci: verify docker-build

# 清理生成的文件
clean:
	@echo "正在清理..."
	@rm -rf bin/
	@echo "✅ 清理完成"

# 安装 swag 工具
install-swag:
	@echo "正在安装 swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ swag 安装完成"
