# 故障排查指南

## 常见问题

### 1. 500 错误 - 无法连接到课程服务

**错误信息：**
```
failed to call course service edu-course: ...
```

**可能原因：**
- 课程服务未运行
- 服务名称不匹配
- 服务发现（mDNS）未找到服务

**解决方案：**

1. **检查课程服务是否运行**
   ```bash
   # 确认课程服务正在运行
   # 检查服务日志
   ```

2. **检查服务名称配置**
   查看 `config.yaml` 中的 `course_service_name` 是否与课程服务的实际名称匹配：
   ```yaml
   micro:
     course_service_name: "edu-course"
   ```

3. **检查服务发现**
   - 如果使用 mDNS，确保在同一网络环境中
   - 如果使用其他注册中心（如 Consul、etcd），确保配置正确

4. **查看详细错误日志**
   检查应用日志中的详细错误信息，通常会包含具体的连接失败原因

### 2. 服务发现失败

**症状：**
- 所有 API 调用都返回 500 错误
- 日志显示 "no such service" 或 "service not found"

**解决方案：**
1. 确认课程服务已注册到服务发现
2. 检查 `config.yaml` 中的注册中心配置
3. 如果使用 mDNS，确保防火墙允许 mDNS 流量

### 3. 端口冲突

**症状：**
- 启动时提示端口已被占用

**解决方案：**
修改 `config.yaml` 中的端口：
```yaml
http:
  address: ":8083"
```

## 调试技巧

### 启用详细日志

设置环境变量以查看更详细的日志：
```bash
export GIN_MODE=debug
```

### 测试服务连接

使用健康检查接口测试服务是否正常运行：
```bash
curl http://localhost:8083/health
```

### 检查配置

确认配置是否正确加载：
```bash
# 查看当前配置
cat config.yaml
```
