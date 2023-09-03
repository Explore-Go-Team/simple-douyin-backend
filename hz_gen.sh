#!/bin/sh

hz update -module simple-douyin-backend -idl idl/user_http.proto
hz update -module simple-douyin-backend -idl idl/feed_http.proto
hz update -module simple-douyin-backend --unset_omitempty --idl idl/relation_http.proto
