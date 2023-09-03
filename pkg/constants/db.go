package constants

// connection information
const (
	MySQLDefaultHost   = "47.115.202.138"
	MySQLDefaultPort   = "3306"
	MySQLDefaultUser   = "dytest"
	MySQLDefaultPwd    = "zxcvbnm"
	MySQLDefaultDBName = "dy"
	MySQLDefaultDSN    = MySQLDefaultUser + ":" + MySQLDefaultPwd + "@tcp(" + MySQLDefaultHost + ":" + MySQLDefaultPort + ")/" + MySQLDefaultDBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	MinioEndPoint        = "localhost:18001"
	MinioAccessKeyID     = "douyin"
	MinioSecretAccessKey = "douyin123"
	MiniouseSSL          = false

	RedisAddr     = "localhost:18003"
	RedisPassword = "douyin123"
)

// constants in the project
const (
	UserTableName      = "users"
	FollowsTableName   = "follows"
	VideosTableName    = "videos"
	MessageTableName   = "messages"
	FavoritesTableName = "likes"
	CommentTableName   = "comments"

	VideoFeedCount       = 30
	FavoriteActionType   = 1
	UnFavoriteActionType = 2

	MinioVideoBucketName = "videobucket"
	MinioImgBucketName   = "imagebucket"

	TestSign       = "测试账号！ offer"
	TestAva        = "avatar/test1.jpg"
	TestBackground = "background/test1.png"
)
