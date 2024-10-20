# Linevis后端

Linevis后端是一个基于Go + Gin + Gorm + Sqlite的项目，主要为Linevis前端提供数据支持。

## 注意事项
若为打包使用，注意对应端口需未占用，默认情况下为9999。

## 前置条件
- Go环境，推荐版本为1.22及以上。
- Postman或其他API测试工具。

## 调试之前
1. 下载项目：
```shell
git clone https://github.com/MoonWX/linevis-backend.git
cd linevis-backend
```

2. 运行项目：
```shell
go run main.go
```

[//]: # (3. 编译项目：)

[//]: # (```shell)

[//]: # (go build linevis-backend)

[//]: # (```)

## 调试方法与API文档
默认情况下，URL应为`http://localhost:9999`
### 1. 字段说明
- `name`： 商品名称
   - 类型：`String`
   - 示例：`"Test1"`
- `main_barcode`： 主条码
    - 类型：`String`
    - 示例：`"11111111"`
- `model`： 商品型号
    - 类型：`String`
    - 示例：`"test"`
- `weight`： 重量
    - 类型：`String`
    - 示例：`10kg`
- `specification`： 规格
    - 类型：`String`
    - 示例：`"test"`
- `target_address`： 目标地址
    - 类型：`String`
    - 示例：`"test"`
- `manual`： 电子指导书文件路径
    - 类型：`String`
    - 示例：`"/manual/test1.png"`
    - 注意：文件路径为相对路径，相对于项目根目录。
- `sub_materials`： 子物料
    - 类型：`JSON Array`
    - 示例：
    ```json
    [
      {
        "name": "Sub1",
        "sub_barcode": "22222222"
      },
      {
        "name": "Sub2",
        "sub_barcode": "33333333"
      },
      {
        "name": "Sub3",
        "sub_barcode": "44444444"
      }
    ]
    ```
    - 注意：子物料为数组，每个子物料包含`name`和`sub_barcode`两个字段。
    - `name`： 子物料名称
        - 类型：`String`
        - 示例：`"Sub1"`
    - `sub_barcode`： 子物料条码
        - 类型：`String`
        - 示例：`"22222222"`

### 2. 可用性检查
- 请求方式：`GET`
- 请求地址：`/ping`
- 返回示例：
```json
{
    "message": "pong"
}
```

### 3. 获取所有数据
- 请求方式：`GET`
- 请求地址：`/products`
- 返回示例：
```json
[
  {
    "ID": 1,
    "CreatedAt": "2024-10-10T13:01:40.2086408+08:00",
    "UpdatedAt": "2024-10-10T13:01:40.2086408+08:00",
    "DeletedAt": null,
    "name": "Test2",
    "main_barcode": "11111111",
    "model": "test",
    "weight": "10",
    "specification": "test",
    "target_address": "test",
    "manual": "test1.png",
    "sub_materials": [
      {
        "name": "Sub1",
        "sub_barcode": "22222222"
      },
      {
        "name": "Sub2",
        "sub_barcode": "33333333"
      },
      {
        "name": "Sub3",
        "sub_barcode": "44444444"
      }
    ]
  },
  {
    "ID": 2,
    "CreatedAt": "2024-10-10T13:19:58.1609398+08:00",
    "UpdatedAt": "2024-10-10T13:19:58.1609398+08:00",
    "DeletedAt": null,
    "name": "Test3",
    "main_barcode": "11111111",
    "model": "test",
    "weight": "10",
    "specification": "test",
    "target_address": "test",
    "manual": "test1.png",
    "sub_materials": [
      {
        "name": "Sub1",
        "sub_barcode": "22222222"
      },
      {
        "name": "Sub2",
        "sub_barcode": "33333333"
      },
      {
        "name": "Sub3",
        "sub_barcode": "44444444"
      }
    ]
  }
]
```

### 4. 获取单个数据
- 请求方式：`GET`
- 请求地址：`/products/:id`
- 地址示例：`/products/1`
- 返回示例：
```json
{
  "ID": 1,
  "CreatedAt": "2024-10-10T13:01:40.2086408+08:00",
  "UpdatedAt": "2024-10-10T13:01:40.2086408+08:00",
  "DeletedAt": null,
  "name": "Test2",
  "main_barcode": "11111111",
  "model": "test",
  "weight": "10",
  "specification": "test",
  "target_address": "test",
  "manual": "test1.png",
  "sub_materials": [
    {
      "name": "Sub1",
      "sub_barcode": "22222222"
    },
    {
      "name": "Sub2",
      "sub_barcode": "33333333"
    },
    {
      "name": "Sub3",
      "sub_barcode": "44444444"
    }
  ]
}
```

### 5. 添加数据
- 请求方式：`POST`
- 请求地址：`/products`
- 请求示例：
```json
{
  "name": "Test4",
  "main_barcode": "11111111",
  "model": "test",
  "weight": "10",
  "specification": "test",
  "target_address": "test",
  "manual": "test1.png",
  "sub_materials": [
    {
      "name": "Sub1",
      "sub_barcode": "22222222"
    },
    {
      "name": "Sub2",
      "sub_barcode": "33333333"
    },
    {
      "name": "Sub3",
      "sub_barcode": "44444444"
    }
  ]
}
```
- 返回示例：
```json
{
  "ID": 3,
  "CreatedAt": "2024-10-10T13:19:58.1609398+08:00",
  "UpdatedAt": "2024-10-10T13:19:58.1609398+08:00",
  "DeletedAt": null,
  "name": "Test4",
  "main_barcode": "11111111",
  "model": "test",
  "weight": "10",
  "specification": "test",
  "target_address": "test",
  "manual": "test1.png",
  "sub_materials": [
    {
      "name": "Sub1",
      "sub_barcode": "22222222"
    },
    {
      "name": "Sub2",
      "sub_barcode": "33333333"
    },
    {
      "name": "Sub3",
      "sub_barcode": "44444444"
    }
  ]
}
```

### 6. 更新数据（即修改数据）
- 请求方式：`PUT`
- 请求地址：`/products/:id`
- 地址示例：`/products/1`
- 请求示例：
```json
{
  "name": "Test1",
  "main_barcode": "88888888",
  "model": "test",
  "weight": "10",
  "specification": "test",
  "target_address": "test",
  "manual": "test1.png",
  "sub_materials": [
    {
      "name": "Sub1",
      "sub_barcode": "22222222"
    },
    {
      "name": "Sub2",
      "sub_barcode": "33333333"
    },
    {
      "name": "Sub3",
      "sub_barcode": "44444444"
    }
  ]
}
```
- 返回示例：
```json
{
    "ID": 3,
    "CreatedAt": "2024-10-10T13:25:28.2784443+08:00",
    "UpdatedAt": "2024-10-10T14:03:57.6864517+08:00",
    "DeletedAt": null,
    "name": "Test2",
    "main_barcode": "88888888",
    "model": "test",
    "weight": "10",
    "specification": "test",
    "target_address": "test",
    "manual": "test1.png",
    "sub_materials": [
        {
            "name": "Sub1",
            "sub_barcode": "22222222"
        },
        {
            "name": "Sub2",
            "sub_barcode": "33333333"
        },
        {
            "name": "Sub3",
            "sub_barcode": "44444444"
        }
    ]
}
```

### 7. 删除数据
- 请求方式：`DELETE`
- 请求地址：`/products/:id`
- 地址示例：`/products/1`
- 返回示例：
```json
{
  "message": "product deleted"
}
```

## 调试完毕

构建项目：

Linux:
```shell
export GIN_MODE=release
go build linevis-backend
```

Windows CMD:
```shell
set GIN_MODE=release
go build linevis-backend
```

Windows PowerShell:
```shell
$env:GIN_MODE="release"
go build linevis-backend
```

随后，与前端项目一同打包即可。
