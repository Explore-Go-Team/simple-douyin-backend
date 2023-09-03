#!/bin/sh

# generate kitex_gen code
echo "Generating/Updating kitex_gen folder..."
kitex -module simple-douyin-backend -I idl idl/user.proto
kitex -module simple-douyin-backend -I idl idl/feed.proto
kitex -module simple-douyin-backend -I idl idl/relation_service.proto

# generate service code and no generate kitex_gen
echo "Generating/Updating services code..."
cwd=$(pwd)

cd $cwd/service/user
kitex -module simple-douyin-backend -I ../../idl -use simple-douyin-backend/kitex_gen/ -service user_service ../../idl/user.proto

cd $cwd/service/feed
kitex -module simple-douyin-backend -I ../../idl -use simple-douyin-backend/kitex_gen/ -service feed_service ../../idl/feed.proto

cd $cwd/service/relation
kitex -module relation_service -I ../../idl -use simple-douyin-backend/kitex_gen/ -service relation_service ../../idl/relation_service.proto

cd $cwd
