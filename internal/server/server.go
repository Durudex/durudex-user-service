/*
 * Copyright © 2021-2022 Durudex

 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.

 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package server

import (
	"net"

	"github.com/durudex/durudex-user-service/internal/config"
	"github.com/durudex/durudex-user-service/internal/delivery/grpc"

	"github.com/rs/zerolog/log"
)

// The main structure of the server.
type Server struct {
	listener *net.Listener
	grpc     *GRPCServer
	handler  *grpc.Handler
}

// Create a new server.
func NewServer(cfg *config.Config, handler *grpc.Handler) (*Server, error) {
	// Server address.
	address := cfg.Server.Host + ":" + cfg.Server.Port

	// Creating a new TCP connections.
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	// Creating a new gRPC server.
	grpcServer, err := NewGRPCServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: &lis,
		grpc:     grpcServer,
		handler:  handler,
	}, nil
}

// Run server.
func (s *Server) Run() {
	log.Debug().Msg("Running server...")

	// Register gRPC handlers.
	s.handler.RegisterHandlers(s.grpc.Server)

	// Running gRPC server.
	if err := s.grpc.Server.Serve(*s.listener); err != nil {
		log.Fatal().Msgf("error running server: %s", err.Error())
	}
}

// Stop server.
func (srv *Server) Stop() {
	log.Info().Msg("Stoping grpc server...")

	srv.grpc.Server.Stop()
}
