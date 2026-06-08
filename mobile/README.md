# gin_mam_mobile

Gin MAM 移动端（Flutter），与 Web / Go 后端共用 `/api/v1` 约定。

## antd_flutter_mobile（alpha）

本项目使用 [antd_flutter_mobile](https://pub.dev/packages/antd_flutter_mobile) 作为 Ant Design Mobile 的 Flutter 实现。

> **请注意：** 该库目前处于早期 **alpha** 阶段，组件属性可能在不另行通知的情况下变更。
> 在纳入生产前请自行评估稳定性；详见包 README 与 [文档站点](https://antd-flutter.vercel.app/)。

应用入口在 `lib/app.dart` 中通过 `AntdLayer.observer` 与 `AntdTokenBuilder` 包裹 `MaterialApp`。

## API 基址

默认 `http://10.0.2.2:8080/api/v1`（Android 模拟器访问宿主机）。真机请改 `lib/core/constants/api_constants.dart`
中的 `defaultApiBaseUrl` 为电脑局域网 IP。

## 开发

```bash
flutter pub get
flutter analyze
flutter test
flutter run
```