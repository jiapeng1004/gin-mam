// 测试用 HTTP 适配器与 JSON 响应构造。

import 'dart:convert';

import 'package:dio/dio.dart';

/// 测试用 HTTP 适配器：按路径返回固定 JSON 响应。
class StubHttpAdapter implements HttpClientAdapter {
  StubHttpAdapter(this.handler);

  final Future<ResponseBody> Function(RequestOptions options) handler;

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<List<int>>? requestStream,
    Future<void>? cancelFuture,
  ) {
    return handler(options);
  }

  @override
  void close({bool force = false}) {}
}

/// 构造 JSON [ResponseBody]。
ResponseBody jsonBody(Object? data, int status) {
  return ResponseBody.fromString(
    data == null ? '' : jsonEncode(data),
    status,
    headers: {
      Headers.contentTypeHeader: [Headers.jsonContentType],
    },
  );
}
