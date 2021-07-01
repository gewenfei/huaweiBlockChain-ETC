/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
 */

package smstub

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"git.huawei.com/poissonsearch/wienerchain/contract/sdk/logger"
	"git.huawei.com/poissonsearch/wienerchain/proto/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var grpclog = logger.GetDefaultLogger()

// SMService smart contract grpc service implementation.
type SMService struct {
}

func streamRecv(stream contract.SMCaller_CallSMServer, errChan chan error, streamRecvChan chan contract.ContractMsg) {
	for {
		recvMsg, err := stream.Recv()
		if err != nil {
			errChan <- err
			grpclog.Infof("Stream recv err:%s", err.Error())
			return
		}
		streamRecvChan <- *recvMsg
	}
}

// CallSM callback function for grpc server
func (s SMService) CallSM(stream contract.SMCaller_CallSMServer) error {
	grpclog.Infof("Create a stream: %v", stream)
	streamSendChan := make(chan contract.ContractMsg)

	errChan := make(chan error)
	streamRecvChan := make(chan contract.ContractMsg)

	go streamRecv(stream, errChan, streamRecvChan)
	for {
		select {
		case err := <-errChan:
			grpclog.Infof("Recv or Send from stream err:%s", err.Error())
			return err

		case conMsg := <-streamRecvChan:
			go ProcessMsgFromWnode(&conMsg, streamSendChan)
		case conMsg := <-streamSendChan:
			grpclog.Debugf("conMsg := <-streamSendChan, %s", conMsg.Type)
			err := stream.Send(&conMsg)
			if err != nil {
				grpclog.Errorf("Stream send err:%s", err.Error())
			}
		}
	}
}

// NewGrpcServer new grpc server.
func NewGrpcServer(config *contract.ContractConfig) {
	var server *grpc.Server

	opts := generateOpts(config)
	if opts == nil {
		server = grpc.NewServer()
	} else {
		server = grpc.NewServer(opts...)
	}

	contract.RegisterSMCallerServer(server, &SMService{})

	lis, err := net.Listen("tcp", "0.0.0.0:"+config.Port)
	if err != nil {
		grpclog.Errorf("listen failed on 9002:%s", err.Error())
	}
	grpclog.Infof("Listen success on %s", config.Port)
	err = server.Serve(lis)
	if err != nil {
		grpclog.Errorf("Start grpc server failed:%s", err.Error())
	}
}

func generateOpts(config *contract.ContractConfig) []grpc.ServerOption {
	var serverOpts []grpc.ServerOption

	grpcCfg := config.GrpcConfig
	if grpcCfg != nil {
		grpclog.Infof("MinInterval:%d, Interval:%d, Timeout:%d", grpcCfg.MinInterval, grpcCfg.Interval, grpcCfg.Timeout)
		para := keepalive.ServerParameters{}
		if grpcCfg.Timeout > 0 {
			para.Timeout = time.Duration(grpcCfg.Timeout) * time.Second
		}

		if grpcCfg.Interval > 0 {
			para.Time = time.Duration(grpcCfg.Interval) * time.Second
		}
		serverOpts = append(serverOpts, grpc.KeepaliveParams(para))

		policy := keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}
		if grpcCfg.MinInterval > 0 {
			policy.MinTime = time.Duration(grpcCfg.MinInterval) * time.Second
		}
		serverOpts = append(serverOpts, grpc.KeepaliveEnforcementPolicy(policy))
	}

	tlsConfig := config.TlsConfig
	if tlsConfig != nil {
		grpclog.Infof("Keypath:%s, certPath:%s, root cert path:%s, server name:%s, is tls enable:%t, is multi tls enable:%t",
			tlsConfig.KeyPath, tlsConfig.CertPath, tlsConfig.PeerRootcertsPath, tlsConfig.ServerName,
			tlsConfig.IsTlsEnabled, tlsConfig.IsMultiTls)

		if tlsConfig.IsTlsEnabled {
			grpclog.Infof("tlsConfig is enabled")
			option, err := tlsServerOptions(tlsConfig)
			if err == nil {
				serverOpts = append(serverOpts, option)
			}
		} else {
			grpclog.Infof("tlsConfig is not enabled")
		}
	}

	return serverOpts
}

func tlsServerOptions(tlsCfg *contract.ContractTLSConfig) (grpc.ServerOption, error) {
	cert, err := tls.LoadX509KeyPair(tlsCfg.CertPath, tlsCfg.KeyPath)
	if err != nil {
		grpclog.Infof("TlsServerOptions LoadX509KeyPair failed")
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{cert},
		SessionTicketsDisabled:   true,
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	if tlsCfg.IsMultiTls {
		certPool, err := loadRootCAs(tlsCfg.PeerRootcertsPath)
		if err != nil {
			grpclog.Errorf("The root ca certs failed\n")
			return nil, err
		}
		tlsConfig.ClientCAs = certPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	creds := credentials.NewTLS(tlsConfig)
	return grpc.Creds(creds), nil
}

func cleanPath(filePath string) string {
	return filePath
}

func loadRootCAs(certPaths []string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	ok := false
	for _, clientCAPath := range certPaths {
		clientCA, err := ioutil.ReadFile(cleanPath(clientCAPath))
		if err != nil {
			return nil, fmt.Errorf("read client root CA error: %s", err)
		}
		if certPool.AppendCertsFromPEM(clientCA) {
			ok = true
		} else {
			grpclog.Errorf("Append certs from pem failed\n")
		}
	}
	if !ok {
		return nil, fmt.Errorf("non client root certs loaded")
	}
	return certPool, nil
}
