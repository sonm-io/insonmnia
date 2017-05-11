package miner

//go:generate protoc -I ../../proto/ ../../proto/insonmnia.proto --go_out=plugins=grpc:../../daemon/miner
