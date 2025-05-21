#!/bin/bash
set -e

# Directory paths
PROTO_SOURCE="../backend/api/proto/reddit_service.proto"
PROTO_DEST="./src/generated"
PROTO_INCLUDE="../backend/api/proto"

# Create the output directory if it doesn't exist
mkdir -p $PROTO_DEST

# Generate JavaScript code
protoc \
  --proto_path=$PROTO_INCLUDE \
  --js_out=import_style=commonjs:$PROTO_DEST \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:$PROTO_DEST \
  $PROTO_SOURCE

# Add a message to confirm generation
echo "✅ Successfully generated gRPC-web TypeScript/JavaScript files"
echo "Output location: $PROTO_DEST"

