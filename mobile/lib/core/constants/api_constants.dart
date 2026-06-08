/// API 基础地址与路径约定。
///
/// 与 Web（Vite 代理 `/api/v1`）及 Go 后端路由前缀保持一致。
library;

/// localStorage / SharedPreferences 中访问令牌的键名（与 Web `auth.ts` 对齐）。
const String authTokenKey = 'gin_mam_auth_token';

/// 默认 REST 前缀：Android 模拟器通过 `10.0.2.2` 访问宿主机 `localhost`。
///
/// 真机调试时请改为电脑局域网 IP，例如 `http://192.168.1.100:8080/api/v1`，
/// 并确保手机与电脑在同一网段、后端监听 `0.0.0.0:8080`。
const String defaultApiBaseUrl = 'http://10.0.2.2:8080/api/v1';

/// HTTP 请求超时（毫秒），与 Web axios 默认 30s 一致。
const Duration apiTimeout = Duration(seconds: 30);