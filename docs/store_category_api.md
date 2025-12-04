# 商品分类管理接口文档

## 1. 基础信息
- **Base URL**: `/mer_admin/store_category`
- **鉴权方式**: Header `Authorization: Bearer <token>`
- **数据格式**: JSON

## 2. 数据结构 (MerStoreCategory)

| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| store_category_id | int | 分类ID |
| cate_name | string | 分类名称 |
| pic | string | 图标地址 |
| sort | int | 排序 (数值越大越靠前) |
| level | int | 层级 (默认1) |
| mer_id | int | 商户ID |
| create_at | string | 创建时间 |

## 3. 接口详情

### 3.1 创建分类
**接口地址**: `POST /mer_admin/store_category`

**请求参数 (Body)**:

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| cate_name | string | 是 | 分类名称 |
| pic | string | 否 | 分类图标URL |
| sort | int | 否 | 排序值 |

**请求示例**:
```json
{
    "cate_name": "生鲜水果",
    "pic": "https://example.com/image.png",
    "sort": 100
}
```

**响应结果**:
```json
{
    "code": 200,
    "msg": "success",
    "data": {
        "store_category_id": 1,
        "cate_name": "生鲜水果",
        "pic": "https://example.com/image.png",
        "sort": 100,
        "level": 1,
        "mer_id": 10,
        "create_at": "2023-10-27T10:00:00Z"
    }
}
```

---

### 3.2 获取分类列表
**接口地址**: `GET /mer_admin/store_category`

**请求参数 (Query)**:

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |

**响应结果**:
```json
{
    "code": 200,
    "msg": "success",
    "data": {
        "list": [
            {
                "store_category_id": 1,
                "cate_name": "生鲜水果",
                "pic": "...",
                "sort": 100,
                ...
            }
        ],
        "total": 1
    }
}
```

---

### 3.3 获取分类详情
**接口地址**: `GET /mer_admin/store_category/:id`

**路径参数**:
- `id`: 分类ID

**响应结果**:
```json
{
    "code": 200,
    "msg": "success",
    "data": {
        "store_category_id": 1,
        "cate_name": "生鲜水果",
        ...
    }
}
```

---

### 3.4 更新分类
**接口地址**: `PUT /mer_admin/store_category/:id`

**路径参数**:
- `id`: 分类ID

**请求参数 (Body)**:

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| cate_name | string | 否 | 分类名称 |
| pic | string | 否 | 分类图标URL |
| sort | int | 否 | 排序值 |

**请求示例**:
```json
{
    "cate_name": "新鲜水果",
    "sort": 99
}
```

**响应结果**:
```json
{
    "code": 200,
    "msg": "success",
    "data": null
}
```

---

### 3.5 删除分类
**接口地址**: `DELETE /mer_admin/store_category/:id`

**路径参数**:
- `id`: 分类ID

**响应结果**:
```json
{
    "code": 200,
    "msg": "success",
    "data": null
}
```
