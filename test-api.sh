#!/bin/bash

echo "测试 Go 后端 API..."
echo ""

echo "1. 健康检查:"
curl -s http://localhost:8080/api/health | python3 -m json.tool
echo ""

echo "2. 获取状态:"
curl -s http://localhost:8080/api/status | python3 -m json.tool
echo ""

echo "✅ API 测试完成"
