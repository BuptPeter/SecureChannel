# SecureChannel
自适应的SDN网元数据保护系统 V3.0-Build-20180728

用于SDN场景下，Open vSwitch与控制器之间的加密数据通信和Open vSwitch本地数据防篡改。
        2018-01-26  v1.0
          1.稳定通信数据元加解密功能
          2.API接口稳定可用
          3.稳定大并发性能
          4.性能测试稳定在10W/s

        2018-03-06  v1.1:
          1.删除测试模式
          2.添加是否启用TLS选项

        2018-03-21  v1.2:
          1.添加工作位置选项
          2.删除协议及多分发
          3.Fix Bugs

        2018-05-21 v2.0:
          1.修改转发方式为数据流加解密
          2.速度大幅提升

        2018-05-21 v2.1:
          1.加入Kerberos身份认证和密钥分发

        2018-07-21 v3.0:
          1.加入ovs-hash服务，记录流表下发
          2.加入ovs-check服务，定时校检HMAC值
