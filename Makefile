# 1. 生成所有 Proto 代码
gen-api:
	kratos proto client api/proto/system/v1/system.proto

# 2. 生成 Ent 代码 (每当修改 schema 后运行)
gen-ent:
	go generate ./app/system/service/internal/data/ent

# 3. 运行服务
run-system:
	cd app/system/service && kratos run

# 4. 依赖注入 (Wire)
wire:
	cd app/system/service/cmd/service && wire