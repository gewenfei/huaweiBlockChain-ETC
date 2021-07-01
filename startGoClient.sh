#!/bin/bash

EXEC_COMMAND="go-client"
export CGO_LDFLAGS=-L/root/openssl/lib
export C_INCLUDE_PATH=/root
export CPLUS_INCLUDE_PATH=/root

export GOPROXY=https://goproxy.cn

# 编译目标文件
go build -o go-client

# 加载环境变量
source configure

cat >> /etc/hosts <<EOF
$EIP node-$ORGNAMEHASH-0.node-$ORGNAMEHASH.default.svc.cluster.local
$EIP node-$ORGNAMEHASH-1.node-$ORGNAMEHASH.default.svc.cluster.local
EOF

PEERS="node-$ORGNAMEHASH-0,node-$ORGNAMEHASH-1"
PEER0="node-$ORGNAMEHASH-0"
PEER1="node-$ORGNAMEHASH-1"
CONSENSUS_PEERS="node-$ORGNAMEHASH-0"
CONFIG_DIR="config"

# echo "查询证书总数"
# ./$EXEC_COMMAND contract query -n $PEER1 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getNowCertificateNum -a 'settlementKey;123'

# echo "存储车主基本信息Hash"

# ./$EXEC_COMMAND contract send -n $PEER1 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f storeBasicInfoHash -a 'settlementKey;123;pubkey123;infohash232'

# echo "查询车主积分信息"
# ./$EXEC_COMMAND contract query -n $PEER1 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getBasicInfoHash -a 'settlementKey;123;pubkey123'

# echo "创建证书"
# ./$EXEC_COMMAND contract send -n $PEER1 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f generateCertificate -a 'settlementKey;123;pubkey123;1312442'

# echo "获取初始化证书"
# ./$EXEC_COMMAND contract query -n $PEER1 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getInitCertificate -a 'settlementKey;123;0'

# echo "交管局签名证书"
# ./$EXEC_COMMAND contract send -n $PEER1 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f attachSigning -a 'trafficManagementKey;123;0;交管局签名0x111'

# echo "银行财务系统签名证书"
# ./$EXEC_COMMAND contract send -n $PEER1 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f attachSigning -a 'bankKey;123;0;银行财务系统签名0x222'

# echo "省计算中心签名证书"
# ./$EXEC_COMMAND contract send -n $PEER1 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f attachSigning -a 'settlementKey;123;0;省计算中心签名0x333'

# echo "用户获取证书"
#  ./$EXEC_COMMAND contract query -n $PEER1 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getMyCertificate -a '0'

# echo "获取0～10的证书"
# ./$EXEC_COMMAND contract query -n $PEER1 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getRangeCertificates -a 'settlementKey;123;0;10'

# # 发起交易-初始化用户ETC
# echo "初始化用户ETC"
# ./$EXEC_COMMAND contract send -n $PEER0 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f initFunds  -a 'bankKey;123;pubkey123;100'

# # 查询链上交易记录,只能在某一个节点上查询
# echo "查询"
# ./$EXEC_COMMAND contract query -n $PEER0 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getFunds  -a "bankKey;123;pubkey123"

# # 发起交易-质押资金
# echo "质押资金"
# ./$EXEC_COMMAND contract send -n $PEER0 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f depositFunds  -a 'bankKey;123;pubkey123;100'

# # 查询链上交易记录,只能在某一个节点上查询
# echo "查询"
# ./$EXEC_COMMAND contract query -n $PEER0 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getFunds  -a "bankKey;123;pubkey123"

# # 发起交易-抵扣资金
# echo "抵扣资金"
# ./$EXEC_COMMAND contract send -n $PEER0 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f deductFunds  -a 'bankKey;123;pubkey123;10'

# # 查询链上交易记录,只能在某一个节点上查询
# echo "查询"
# ./$EXEC_COMMAND contract query -n $PEER0 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getFunds  -a "bankKey;123;pubkey123"

# # 发起交易-初始化用户ETC
# echo "初始化用户ETC"
# ./$EXEC_COMMAND contract send -n $PEER0 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f initFunds  -a 'bankKey;123;pubkey456;100'

# # 发起交易-转账
# echo "转账"
# ./$EXEC_COMMAND contract send -n $PEER0 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f transferFunds  -a 'bankKey;123;pubkey123;pubkey456;10'

# # 查询链上交易记录,只能在某一个节点上查询
# echo "查询"
# ./$EXEC_COMMAND contract query -n $PEER0 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getFunds  -a "bankKey;123;pubkey123"

# # 查询链上交易记录,只能在某一个节点上查询
# echo "查询"
# ./$EXEC_COMMAND contract query -n $PEER0 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getFunds  -a "bankKey;123;pubkey456"

# echo "查询交易历史"
# ./$EXEC_COMMAND contract query -n $PEER0 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getFundsHistory  -a "bankKey;123;pubkey123"

# # 发起交易-初始化用户ETC
# echo "初始化用户ETC"
# ./$EXEC_COMMAND contract send -n $PEER0 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f initFunds  -a 'bankKey;123;pubkey456;100'

# # 查询链上交易记录,只能在某一个节点上查询
# echo "查询"
# ./$EXEC_COMMAND contract query -n $PEER0 -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f getFunds  -a "bankKey;123;pubkey456"

# # 发起交易-质押资金
# echo "质押资金"
# ./$EXEC_COMMAND contract send -n $PEER0 -s $CONSENSUS_PEERS -c $CHAIN_NAME -g $CONFIG_DIR/$SDK_CONFIG_FILE  -t $CONTRACT_NAME -f depositFunds  -a 'bankKey;123;pubkey456;100'
