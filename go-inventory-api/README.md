# 进销存后端 API — Go + SQLite

## 项目结构

```
inventory-api/
├── main.go                  # 入口：初始化 DB、路由、启动服务
├── go.mod                   # 依赖声明
├── models/
│   └── models.go            # 数据模型 + DTO（Product、StockLog）
├── handlers/
│   ├── product.go           # 商品 CRUD Handler
│   └── stock.go             # 入库/出库/流水 Handler
└── middleware/
    └── logger.go            # 请求日志中间件
```

## 快速启动

```bash
# 1. 安装依赖
go mod tidy

# 2. 运行（默认监听 :8080，数据库文件 inventory.db）
go run main.go

# 3. 可选环境变量
ADDR=:9090 DB_PATH=./data.db GIN_MODE=release go run main.go
```

## API 一览

| 方法   | 路径                         | 说明             |
|--------|------------------------------|------------------|
| GET    | /health                      | 健康检查         |
| POST   | /api/v1/products             | 新增商品         |
| GET    | /api/v1/products             | 商品列表（分页） |
| GET    | /api/v1/products/:id         | 商品详情         |
| PUT    | /api/v1/products/:id         | 更新商品         |
| DELETE | /api/v1/products/:id         | 删除商品（软删）  |
| POST   | /api/v1/stock/in             | 入库             |
| POST   | /api/v1/stock/out            | 出库             |
| GET    | /api/v1/stock/logs           | 流水查询         |
| GET    | /api/v1/stock/:product_id    | 单品库存         |

---

## curl 测试示例

### 健康检查
```bash
curl http://localhost:8080/health
```

### 新增商品
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"苹果手机壳","price":29.9,"stock":100,"unit":"个"}'
```

### 商品列表（分页 + 搜索）
```bash
curl "http://localhost:8080/api/v1/products?page=1&size=10&name=手机"
```

### 更新商品
```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{"price":25.0,"unit":"件"}'
```

### 删除商品
```bash
curl -X DELETE http://localhost:8080/api/v1/products/1
```

### 入库
```bash
curl -X POST http://localhost:8080/api/v1/stock/in \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":50,"remark":"首批到货"}'
```

### 出库
```bash
curl -X POST http://localhost:8080/api/v1/stock/out \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":20,"remark":"门店补货"}'
```

### 流水查询（按商品/类型筛选）
```bash
curl "http://localhost:8080/api/v1/stock/logs?product_id=1&type=OUT&page=1"
```

### 单品库存
```bash
curl http://localhost:8080/api/v1/stock/1
```

---

## 统一响应格式

```json
{
  "code": 0,
  "message": "ok",
  "data": { ... }
}
```

| code | 含义          |
|------|---------------|
| 0    | 成功          |
| 400  | 请求参数错误  |
| 404  | 资源不存在    |
| 409  | 数据冲突      |
| 422  | 库存不足      |
| 500  | 服务器错误    |
