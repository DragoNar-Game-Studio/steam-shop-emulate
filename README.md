# Steam Shop Emulator

一个基于 Go 的 Steam 商店页模拟器，用于：

- 以接近 Steam 商品页的方式展示自己的游戏素材
- 在后台编辑标题、标签、价格、卖点与长短描述
- 上传主视觉、胶囊图、Logo 与截图组
- 实时审查素材是否具备较好的商店展示基础

## 启动

```bash
go run ./cmd/server
```

默认地址：

- 前台页面: `http://localhost:8080/`
- 后台编辑: `http://localhost:8080/admin`

## 目录结构

```text
cmd/server             程序入口
internal/app           组装依赖与默认配置
internal/domain        领域模型
internal/store         JSON 持久化存储
internal/review        素材审查规则
internal/web           HTTP handler 与模板渲染
web/templates          商店页与后台模板
web/static             样式资源
data/storefront.json   商店页内容
data/uploads           上传素材目录
```

## 当前能力

- Steam 风格首页布局
- 后台文案与标签编辑
- 图片素材上传
- 实时审查接口 `/api/review`

