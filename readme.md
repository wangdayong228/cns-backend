# Cns-Backend

## Postman
[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/22322698-cd32951a-a24f-4fd5-a9fb-2e26f057532c?action=collection%2Ffork&collection-url=entityId%3D22322698-cd32951a-a24f-4fd5-a9fb-2e26f057532c%26entityType%3Dcollection%26workspaceId%3D0df0c5b3-6c0a-47ee-ab26-8ba0139261e4)

## TODO list
1. [ ] Nginx https 反向代理
2. [X] Nginx cors 中间件
3. [ ] 数据库存中文
4. [ ] 续费
5. [X] Refresh URL
6. [X] 下单时判断前端价格是否正确
7. [X] 下单时事务调整
8. [X] 根据订单到期时间，自动关闭订单(conflux-pay来做)
9. [X] 根据流程图完善下单逻辑(现在order已存在直接返回错误，因为commithash生成成本非常低)
10. [ ] 注册失败退款
11. [X] make commit 如果已存在返回错误