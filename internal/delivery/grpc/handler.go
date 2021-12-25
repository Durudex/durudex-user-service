/*
	Copyright © 2021 Durudex

	This file is part of Durudex: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as
	published by the Free Software Foundation, either version 3 of the
	License, or (at your option) any later version.

	Durudex is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with Durudex. If not, see <https://www.gnu.org/licenses/>.
*/

package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"github.com/Durudex/durudex-user-service/internal/delivery/grpc/protobuf"
	"github.com/Durudex/durudex-user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Handler struct {
	service *service.Service
}

// Creating a new grpc handler.
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Registration services handlers.
func (h *Handler) RegisterHandlers(srv *grpc.Server) {
	pb.RegisterUserServiceServer(srv, NewUserHandler(h.service))
}

// Loading TLS credentials.
func LoadTLSCredentials(caCertPath, certPath, keyPath string) (credentials.TransportCredentials, error) {
	// Load certificate on the CA who signed client's certificate.
	pemCA, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCA) {
		return nil, errors.New("error to add server CA's certificate")
	}

	// Load server's certificate and private key.
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	// Create the credentials and returning it.
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
