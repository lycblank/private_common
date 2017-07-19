/*
 Package xxx

 Copyright(c) 2000-2016 Quqian Technology Company Limit.
 All rights reserved.
 创建人：  longyongchun
 创建日期：2016-10-25 16:28:15
 修改记录：
 ----------------------------------------------------------------
 修改人        |  修改日期    				|    备注
 ----------------------------------------------------------------
 longyongchun  | 2016-10-25 16:28:15		|    创建文件
 ----------------------------------------------------------------
*/

package utils

const APP_ID = "10000"

// ======================= 用户相关 =============================
const (
	REDIS_KEY_USER_ATTR                   = "user.in:user:user_attr:" // 单个用户的可变属性
	REDIS_FIELD_USER_ATTR_NICKNAME        = "nick_name"               // 昵称
	REDIS_FIELD_USER_ATTR_SEX             = "sex"                     // 性别
	REDIS_FIELD_USER_ATTR_QQ              = "qq"                      // QQ号
	REDIS_FIELD_USER_ATTR_LOCATION        = "location"                // 所在地
	REDIS_FIELD_USER_ATTR_BIRTHTIME       = "birth_time"              // 出生时间戳
	REDIS_FIELD_USER_ATTR_PORTRAIT_URL    = "portrait_url"            // 头像url
	REDIS_FIELD_USER_ATTR_SIGNATURE       = "signature"               // 个性签名
	REDIS_FIELD_USER_ATTR_CONT_SIGN_DAYS  = "cont_sign_days"          // 连续签到天数
	REDIS_FIELD_USER_ATTR_BIND_MOBILE     = "bind_mobile"             // 绑定的手机号
	REDIS_FIELD_USER_ATTR_ANCHOR_VERIFIED = "anchor_verified"         // 主播校验
	REDIS_FIELD_USER_ATTR_LIVE_NOTICE     = "live_notice"             // 开始直播提醒
	REDIS_FIELD_USER_ATTR_EXP             = "exp"                     // 经验
	REDIS_FIELD_USER_ATTR_BG_IMAGE_URL    = "bg_image_url"            // 背景图片
	REDIS_FIELD_USER_ATTR_CERT_FLAG       = "cert_flag"               // 认证标识
	REDIS_FIELD_USER_ATTR_CERT_INFO       = "cert_info"               // 认证信息
	REDIS_FIELD_USER_ATTR_FANS_COUNT      = "fans_count"              // 粉丝数量
	REDIS_FIELD_USER_ATTR_LABEL           = "label"                   // 标签属性[以英文,分隔]
	REDIS_FIELD_USER_ATTR_FRONT_TAB_ID    = "front_tab_id"            // 首次进入用户首页时的tab指定

	REDIS_KEY_UID_APPID_MAP_SID  = "user.in:uid_appid_map_sid:"  // uid_app_id 映射 session id
	REDIS_KEY_UID_SOURCE_MAP_SID = "user.in:uid_source_map_sid:" // uid_source 映射 session id

	REDIS_KEY_USER_BLACKLIST             = "user.in:blacklist"                   // 黑名单
	REDIS_KEY_UGC_FAVORITE_BODIES_IN_UID = "user.in:ugc:favorite_bodies_in_uid:" // 用户收藏
	REDIS_KEY_ORDER_RECORD_IN_UID        = "user.in:order_records_in_uid:"       // 用户订单记录

	REDIS_KEY_FOLLOWED_UIDS_IN_UID        = "user.in:followed_uids_in_uid:"        // 用户关注的用户列表
	REDIS_KEY_FOLLOWED_CIDS_IN_UID        = "user.in:followed_cids_in_uid:"        // 用户关注的剧组列表
	REDIS_KEY_FOLLOWED_ANCHOR_UIDS_IN_UID = "user.in:followed_anchor_uids_in_uid:" // 主播关注的用户列表
	REDIS_KEY_FANS_UIDS_IN_UID            = "user.in:fans_uids_in_uid:"            // 用户的粉丝列表

	REDIS_KEY_USER_TASK_TYPE_DATA_IN_TTID    = "user.in:data_in_ttid:" // 任务类型属性
	REDIS_KEY_USER_TASK_TYPE_DATA_FILED_NAME = "name"                  // 任务类型名称

	REDIS_KEY_USER_TASK_DATA_IN_TID                = "user.in:data_in_tid:" // 任务属性
	REDIS_KEY_USER_TASK_DATA_FILED_NAME            = "name"                 // 任务名称
	REDIS_KEY_USER_TASK_DATA_FILED_TYPE            = "type"                 // 任务类型
	REDIS_KEY_USER_TASK_DATA_FILED_ICON_URL        = "icon_url"             // 任务图标
	REDIS_KEY_USER_TASK_DATA_FILED_BONUS_EXP       = "bonus_exp"            // 任务经验值
	REDIS_KEY_USER_TASK_DATA_FILED_BONUS_COIN_PAIN = "bonus_coin_pain"      // 任务获取的普通货币
	REDIS_KEY_USER_TASK_DATA_FILED_BONUS_COUPON    = "bonus_coupon"         // 任务获取的卷数
	REDIS_KEY_USER_TASK_DATA_FILED_BEGIN_TIME      = "begin_time"           // 任务开始时间
	REDIS_KEY_USER_TASK_DATA_FILED_END_TIME        = "end_time"             // 任务结束时间

	REDIS_KEY_USER_TASK_TYPE              = "user.in:task_type"           // 任务类型集合
	REDIS_KEY_USER_TASK_TIDS_IN_TTID      = "user.in:tids_in_ttid:"       // 任务类型下的任务集合
	REDIS_KEY_USER_TASK_FINSH_TIDS_IN_UID = "user.in:finish_tids_in_uid:" // 用户类型下完成的任务id

	REDIS_KEY_USER_CONS_SIGNIN_NUM_IN_UID = "user.in:cons_signin_num_in_uid:" //用户连续签到的值

	REDIS_KEY_USER_RECOMMEND_ANCHOR_UIDS = "user.in:recommend:anchor_uids" // 推荐的主播列表
	REDIS_KEY_USER_MIN_EXP_IN_LEVEL      = "user.in:min_exp_in_level"      // 等级与最小经验的对应关系
	REDIS_KEY_USER_CIDS_IN_UID           = "user.in:cids_in_uid:"          // 等级与最小经验的对应关系

	REDIS_KEY_USER_VIDS_IN_UID  = "user.in:vids_in_uid:"
	REDIS_KEY_USER_DIDS_IN_UID  = "user.in:dids_in_uid:"
	REDIS_KEY_USER_SVIDS_IN_UID = "user.in:svids_in_uid:"

	REDIS_KEY_USER_RED_DOT_IN_UID   = "user.in:red_dot_in_uid:"
	REDIS_KEY_FIELD_SCAN_VIDEO_TIME = "scan_video_time"

	REDIS_KEY_USER_ORDER_RECORDS_IN_UID = "user.in:order_records_in_uid:"
)

// ======================= 视频相关 =====================================
const (
	REDIS_KEY_VIDEO_PIDS_IN_VID     = "video.in:pids_in_vid:" // 存储视频中的评论id集合 type:sortset member:pid score:time
	REDIS_KEY_DATA_IN_VID           = "video.in:data_in_vid:" // 视频的内容，hash
	REDIS_FIELD_VIDEO_TITLE         = "title"                 // 视频标题
	REDIS_FIELD_VIDEO_VIDEO_URL     = "video_url"             // 视频url
	REDIS_FIELD_VIDEO_IMAGE_URL     = "image_url"             // 封面图像
	REDIS_FIELD_VIDEO_SIZE          = "size"                  // 视频大小
	REDIS_FIELD_VIDEO_SRC_TYPE      = "src_type"              // 视频来源 1：直播回放 2：上传视频
	REDIS_FIELD_VIDEO_MD5_KEY       = "md5_key"               // 视频校验的md5值
	REDIS_FIELD_VIDEO_LOOK_COUNT    = "look_count"            // 观看次数
	REDIS_FIELD_VIDEO_LIKE_COUNT    = "like_count"            // 点赞数量
	REDIS_FIELD_VIDEO_COMMENT_COUNT = "comment_count"         // 评论数量
	REDIS_FIELD_VIDEO_SHARE_COUNT   = "share_count"           // 分享数量
	REDIS_FIELD_VIDEO_TIME          = "time"                  // 视频生成时间
	REDIS_FIELD_VIDEO_UID           = "uid"                   // 谁的视频
	REDIS_FIELD_VIDEO_CID           = "cid"                   // 谁的视频
	REDIS_FIELD_VIDEO_DURATION      = "duration"              // 视频时长

	REDIS_KEY_VIDEO_LIKE_UIDS_IN_VID = "video.in:ugc:like_uids_in_vid:"
)

// ===========================feed相关============================
const (
	REDIS_KEY_FEEDS_DISCOVER_DYNAMIC_IDS = "feeds.in:discover_dynamic_ids" //存储所有的直播与视频动态
	REDIS_KEY_LIVE_LIVE_LIST             = "live.in:live_list"             //存储所有的直播动态
	REDIS_KEY_VIDEO_VIDEO_LIST           = "video.in:video_list"           //存储所有的直播动态

	REDIS_KEY_FEEDS_FOLLOW_DYNAMIC_IDS_IN_UID  = "feeds.in:follow_dynamic_ids_in_uid:"
	REDIS_KEY_FEEDS_PUBLISH_DYNAMIC_IDS_IN_UID = "feeds.in:publish_dynamic_ids_in_uid:"
)

// ==========================社区相关==============================
const (
	// 主题帖
	REDIS_KEY_COMMUNITY_DATA_IN_TID        = "community.in:ugc:data_in_tid:"
	REDIS_FIELD_PUBLISH_TIME               = "publish_time"   // 发布时间
	REDIS_FIELD_UPDATE_TIME                = "update_time"    // 更新时间
	REDIS_FIELD_COMMENT_COUNT              = "comment_count"  // 评论数
	REDIS_FIELD_CID                        = "cid"            // 剧组id
	REDIS_FIELD_VIEW_COUNT                 = "view_count"     // 浏览数
	REDIS_FIELD_FAVORITE_COUNT             = "favorite_count" // 收藏数
	REDIS_KEY_INFO_DATA_FIELD_TITLE        = "title"          //保存资讯的属性信息
	REDIS_KEY_INFO_DATA_FIELD_SOURCE       = "source"         //保存资讯的属性信息
	REDIS_KEY_INFO_DATA_FIELD_PUBLISH_TIME = "publish_time"   //保存资讯的属性信息
	REDIS_KEY_INFO_DATA_FIELD_TYPE         = "type"           //保存资讯的属性信息
	REDIS_KEY_INFO_DATA_FIELD_CONTENT      = "content"        //保存资讯的属性信息

	// 跟贴
	REDIS_KEY_COMMUNITY_DATA_IN_PID = "community.in:ugc:data_in_pid:"
	REDIS_FIELD_POST_TIME           = "post_time" // 跟贴时间
	REDIS_FIELD_TID                 = "tid"       // 主题帖id
	REDIS_FIELD_FLOOR_NO            = "floor_no"  // 楼层号

	// 评论
	REDIS_KEY_COMMUNITY_DATA_IN_PCID = "community.in:ugc:data_in_pcid:"
	REDIS_FIELD_COMMENT_TIME         = "comment_time" // 评论时间
	REDIS_FIELD_PID                  = "pid"          // 跟贴id
	REDIS_FIELD_REPLY_CID            = "reply_cid"    // 被评论的评论id

	// 共有的field
	REDIS_FIELD_UID     = "uid"     // 内容生成者uid
	REDIS_FIELD_CONTENT = "content" //内容

	REDIS_KEY_COMMUNITY_LIKE_UIDS_IN_TID      = "community.in:ugc:like_uids_in_tid:"
	REDIS_KEY_COMMUNITY_LIKE_UIDS_IN_PID      = "community.in:ugc:like_uids_in_pid:"
	REDIS_KEY_COMMUNITY_TOPIC_COMMENT_NUMS    = "community.in:stats:ugc:topic_comment_nums"
	REDIS_KEY_COMMUNITY_PICS_IN_TID           = "community.in:ugc:pics_in_tid:"
	REDIS_KEY_COMMUNITY_PICS_IN_PID           = "community.in:ugc:pics_in_pid:"
	REDIS_KEY_COMMUNITY_CIDS_IN_PID           = "community.in:ugc:cids_in_pid:"
	REDIS_KEY_COMMUNITY_POST_LIKE_NUMS_IN_TID = "community.in:ugc:post_like_nums_in_tid:"

	REDIS_KEY_COMMUNITY_TIDS_IN_NEWEST = "community.in:ugc:tids_in_newest"

	REDIS_KEY_COMMUNITY_AWARD_UIDS_IN_TID = "community.in:ugc:award_uids_in_tid:"
)

// ==========================栏目相关==============================
const (
	REDIS_KEY_COLUMN_DATA_IN_CID            = "column.in:data_in_cid:"
	REDIS_KEY_COLUMN_DATA_FILED_TYPE        = "type"
	REDIS_KEY_COLUMN_DATA_FILED_TITLE       = "title"
	REDIS_KEY_COLUMN_DATA_FILED_IMAGE_URL   = "image_url"
	REDIS_KEY_COLUMN_DATA_FILED_UID         = "uid"
	REDIS_KEY_COLUMN_DATA_FILED_IS_LIVE     = "is_live"
	REDIS_KEY_COLUMN_DATA_FILED_WEIGHT      = "weight"
	REDIS_KEY_COLUMN_DATA_FILED_CREATE_TIME = "create_time"

	REDIS_KEY_COLUMN_LIKE_UIDS_IN_CID = "column.in:like_uids_in_cid:"
	REDIS_KEY_FEED_IDS_IN_COLUMN      = "column.in:cids_in_column" //栏目feed流
	REDIS_KEY_PAST_IDS_IN_COLUMN      = "column.in:past_vids_in_cid:"
)

// ===========================banner相关============================
const (
	REDIS_KEY_BANNER_BIDS_IN_LOCATION = "banner.in:bids_in_location:" // 冒号后面的内容由 广告位_平台 组成
	REDIS_KEY_BANNER_DATA_IN_BID      = "banner.in:data_in_bid:"
	REDIS_FIELD_BANNER_TYPE           = "type"       // 类型
	REDIS_FIELD_BANNER_PLATFORM       = "platform"   // 平台
	REDIS_FIELD_BANNER_OBJECT_ID      = "object_id"  // 内容id
	REDIS_FIELD_BANNER_OBJECT_URL     = "object_url" // 对象链接
	REDIS_FIELD_BANNER_IMAGE_URL      = "image_url"  // 图片地址
	REDIS_FIELD_BANNER_WEIGHT         = "weight"     // 权重
	REDIS_FIELD_BANNER_ANCHOR_ID      = "anchor_id"  // 直播类栏目有用
	REDIS_FIELD_BANNER_UID            = "uid"        // 直播类栏目有用
	REDIS_FIELD_BANNER_VID            = "vid"        // 视频类栏目有用
)

// ==========================直播相关============================
const (
	REDIS_KEY_LIVE_DATA_IN_UID                = "live.in:data_in_uid:"
	REDIS_KEY_LIVE_DATA_FILED_IMAGE_URL       = "image_url"
	REDIS_KEY_LIVE_DATA_FILED_TITLE           = "title"
	REDIS_KEY_LIVE_DATA_FILED_DESCRIPTION     = "description"
	REDIS_KEY_LIVE_DATA_FILED_ONLINE_COUNT    = "online_count"
	REDIS_KEY_LIVE_DATA_FILED_TIME            = "time"
	REDIS_KEY_LIVE_DATA_FILED_STREAM_ID       = "stream_id"
	REDIS_KEY_LIVE_DATA_FILED_CID             = "cid"
	REDIS_KEY_LIVE_DATA_FILED_ORIENTATION     = "orientation"
	REDIS_KEY_LIVE_DATA_FILED_UID             = "uid"
	REDIS_KEY_LIVE_DATA_FILED_TOTAL_LIVE_TIME = "total_live_time"
	REDIS_KEY_LIVE_DATA_FILED_IS_LIVE         = "is_live"
	REDIS_KEY_LIVE_DATA_FILED_VERFIED_INFO    = "verfied_info"

	REDIS_KEY_PAAS_ANCHOR_LIST_START_TIME = "paas:anchor:list:start_time"

	REDIS_KEY_LIVE_PREVUE_IN_ID             = "live.in:prevue_in_id:"
	REDIS_KEY_LIVE_PREVUE_FIELD_TITLE       = "title"
	REDIS_KEY_LIVE_PREVUE_FIELD_UID         = "uid"
	REDIS_KEY_LIVE_PREVUE_FIELD_START_TIME  = "start_time"
	REDIS_KEY_LIVE_PREVUE_FIELD_COVER_URL   = "cover_url"
	REDIS_KEY_LIVE_PREVUE_FIELD_CONTENT     = "content"
	REDIS_KEY_LIVE_PREVUE_FIELD_LIVE_STATUS = "live_status"

	REDIS_KEY_LIVE_LIVE_PREVUE_LIST = "live.in:live_prevue_list"
)

// ==========================消息中心相关============================
const (
	REDIS_KEY_MSGCENTER_COMMENT_ITEMS_IN_UID = "msgcenter.in:comment:items_in_uid:" // 消息中心-评论, 包括跟贴、跟贴评论、视频评论
	REDIS_KEY_MSGCENTER_LIKE_ITEMS_IN_UID    = "msgcenter.in:like:items_in_uid:"    // 消息中心-点赞, 包括帖子点赞、跟贴点赞
	REDIS_KEY_MSGCENTER_SYSTEM_ITEMS_IN_UID  = "msgcenter.in:system:items_in_uid:"  // 消息中心-系统

	REDIS_KEY_MSGCENTER_DATA_IN_SYSTEM_ID = "msgcenter.in:system:data_in_system_id:" // 系统消息数据
	REDIS_KEY_MSGCENTER_FILED_TYPE        = "type"                                   // 类型
	REDIS_KEY_MSGCENTER_FILED_TIME        = "time"                                   // 生成时间
	REDIS_KEY_MSGCENTER_FILED_UID         = "uid"                                    // 发布人
	REDIS_KEY_MSGCENTER_FILED_CONTENT     = "content"                                // 消息内容
)

// ===========================config相关============================
const (
	// 闪屏
	REDIS_KEY_CONFIG_SPLASH_IDS_IN_PLATFORM = "config.in:splash_ids_in_plaform:" // 冒号后面是_平台
	REDIS_KEY_CONFIG_DATA_IN_SPLASH_ID      = "config.in:data_in_splash_id:"
	REDIS_FIELD_SPLASH_BEGIN_TIME           = "begin_time" // 开始时间
	REDIS_FIELD_SPLASH_END_TIME             = "end_time"   // 结束时间
	REDIS_FIELD_SPLASH_LINK_URL             = "link_url"   // 跳转链接
	REDIS_FIELD_SPLASH_IMAGE_URL            = "image_url"  // 封面图片

	// 审核状态集合
	REDIS_KEY_CONFIG_AUDIT_STATUS_SET = "config.in:audit_status_set"

	// 升级
	REDIS_KEY_CONFIG_DATA_IN_UPGRADE_PLATFORM = "config.in:data_in_upgrade_platform:"
	REDIS_FIELD_UPGRADE_LOWEST_VERSION        = "lowest_version"
	REDIS_FIELD_UPGRADE_NEWEST_VERSION        = "newest_version"
	REDIS_FIELD_UPGRADE_NEWEST_URL            = "newest_url"
	REDIS_FIELD_UPGRADE_NEWEST_TIPS           = "newest_tips"
)

// ==========================内部服务相关 ===========================
const (
	REDIS_KEY_FEEDS_QUEUE_WAIT_DEAL_USER_OPERATION = "feeds.queue:wait_deal_user_operation" //等待处理的用户操作
)

// ==========================剧组相关 =================================
const (
	REDIS_KEY_CREW_CIDS_IN_CREW            = "crew.in:cids_in_crew"
	REDIS_KEY_CREW_ACTOR_SIDS_IN_CID       = "crew.in:actor_sids_in_cid:"
	REDIS_KEY_CREW_IMAGE_IDS_IN_CID        = "crew.in:image_ids_in_cid:"
	REDIS_KEY_CREW_VIDS_IN_CID             = "crew.in:vids_in_cid:"
	REDIS_KEY_CREW_DIRECTOR_SIDS_IN_CID    = "crew.in:director_sids_in_cid:"
	REDIS_KEY_CREW_FOLLOW_UIDS_IN_CID      = "crew.in:follow_uids_in_cid:"
	REDIS_KEY_CREW_DYNAMIC_IDS_IN_CID      = "crew.in:dynamic_ids_in_cid:"
	REDIS_KEY_CREW_DATA_IN_CID             = "crew.in:data_in_cid:"
	REDIS_KEY_CREW_DATA_FILED_COVER_URL    = "cover_url"
	REDIS_KEY_CREW_DATA_FILED_TITLE        = "title"
	REDIS_KEY_CREW_DATA_FILED_BRIEF        = "brief"
	REDIS_KEY_CREW_DATA_FILED_SUMMARY      = "summary"
	REDIS_KEY_CREW_DATA_FILED_FANS_COUNT   = "fans_count"
	REDIS_KEY_CREW_DATA_FILED_THEME_COUNT  = "theme_count"
	REDIS_KEY_CREW_DATA_IN_IMAGE_ID        = "crew.in:data_in_image_id:"
	REDIS_KEY_CREW_IMAGE_DATA_FILED_URL    = "url"
	REDIS_KEY_CREW_IMAGE_DATA_FILED_WIDTH  = "width"
	REDIS_KEY_CREW_IMAGE_DATA_FILED_HEIGHT = "height"
	REDIS_KEY_CREW_IMAGE_DATA_FILED_TITLE  = "title"

	REDIS_KEY_CREW_DATA_IN_SID = "crew.in:data_in_sid:"
	REDIS_KEY_CREW_FILED_UID   = "uid"
	REDIS_KEY_CREW_FILED_ROLE  = "role"
)

// ==========================资讯相关 ===========================
const (
	REDIS_KEY_INFO_INFORMATION_LIST   = "info.in:information_list"    //所有的资讯列表
	REDIS_KEY_USER_LOOKED_INFO_IN_UID = "user.in:looked_info_in_uid:" //保存用户看过的所有资讯列表
	REDIS_KEY_USER_INFO_IDS_IN_UID    = "user.in:info_ids_in_uid:"    //保存属于用户的所有资讯列表
	REDIS_KEY_INFO_DATA_IN_ID         = "info.in:data_in_info_id:"    //保存资讯的属性信息

	REDIS_KEY_INFO_CONVERS_IN_INFO_ID   = "info.in:convers_in_info_id:" //资讯中保存的封面
	REDIS_KEY_INFO_CONVERS_DATA_IN_ID   = "info.in:convers_data_in_id:" //封面图片内容
	REDIS_KEY_INFO_CONVERS_FIELD_WIDTH  = "width"                       //保存资讯的封面属性信息
	REDIS_KEY_INFO_CONVERS_FIELD_HEIGHT = "height"                      //保存资讯的封面属性信息
	REDIS_KEY_INFO_CONVERS_FIELD_URL    = "url"                         //保存资讯的封面属性信息
)

// =========================短视频==========================
const (
	REDIS_KEY_SHORTVIDEO_DATA_IN_SVID      = "shortvideo.in:data_in_svid:"
	REDIS_FIELD_SHORTVIDEO_UID             = "uid"
	REDIS_FIELD_SHORTVIDEO_SECONDS         = "seconds"
	REDIS_FIELD_SHORTVIDEO_IDEA            = "idea"
	REDIS_FIELD_SHORTVIDEO_COVER_URL       = "cover_url"
	REDIS_FIELD_SHORTVIDEO_COVER_WIDTH     = "cover_width"
	REDIS_FIELD_SHORTVIDEO_COVER_HEIGHT    = "cover_height"
	REDIS_FIELD_SHORTVIDEO_VIDEO_URL       = "video_url"
	REDIS_FIELD_SHORTVIDEO_LIKE_COUNT      = "like_count"
	REDIS_FIELD_SHORTVIDEO_COMMENT_COUNT   = "comment_count"
	REDIS_FIELD_SHORTVIDEO_UPDATE_TIME     = "update_time"
	REDIS_FIELD_SHORTVIDEO_VIEW_COUNT      = "view_count"
	REDIS_FIELD_SHORTVIDEO_FAKE_VIEW_COUNT = "fake_view_count"

	REDIS_KEY_SHORTVIDEO_NEWEST_SVIDS_IN_HASHTAG = "shortvideo.in:newest_svids_in_hashtag:"
	REDIS_KEY_SHORTVIDEO_HOTEST_SVIDS_IN_HASHTAG = "shortvideo.in:hotest_svids_in_hashtag:"
	REDIS_KEY_SHORTVIDEO_NEWEST_SVIDS_IN_ALL     = "shortvideo.in:newest_svids_in_all"
	REDIS_KEY_SHORTVIDEO_HOTEST_SVIDS_IN_ALL     = "shortvideo.in:hotest_svids_in_all"
	REDIS_KEY_SHORTVIDEO_HASHTAG_ID_IN_SVID      = "shortvideo.queue:hashtag_id_in_svid:" // 短视频包含的话题标签列表 list

	REDIS_KEY_SHORTVIDEO_ADS_SVIDS_IN_HASHTAG           = "shortvideo.in:ads_svids_in_hashtag:"           // 每期活动[话题标签]的的广告视频id集合，目前只有1个
	REDIS_KEY_SHORTVIDEO_ACTIVITY_RANK_SVIDS_IN_HASHTAG = "shortvideo.in:activity_rank_svids_in_hashtag:" // 每期活动[话题标签]的排名视频id集合，目前是27个

	REDIS_KEY_SHORTVIDEO_TAGS_IN_SVID = "shortvideo.in:tags_in_svid:"
)

// =======================话题标签=====================
const (
	REDIS_KEY_HASHTAG_DATA_IN_ID  = "hashtag.in:data_in_id:"
	REDIS_FIELD_HASHTAG_NAME      = "name"
	REDIS_FIELD_HASHTAG_TYPE      = "type"
	REDIS_FIELD_HASHTAG_INTRO     = "intro"
	REDIS_FIELD_HASHTAG_IMAGE_URL = "image_url"
	REDIS_FIELD_HASHTAG_LINK_URL  = "link_url"
	REDIS_FIELD_HASHTAG_LINK_TEXT = "link_text"

	REDIS_KEY_HASHTAG_HASHTAG_ID_SET = "hashtag.in:hashtag_id_set" //mem: id  score: start_time

	REDIS_KEY_TAG_DATA_IN_ID = "tag.in:data_in_id:"
	REDIS_KEY_TAG_FIELD_NAME = "name"
)

//===========================评论相关============================
const (
	REDIS_KEY_COMMENT_UGC_DATA_IN_PID            = "comment.in:ugc:data_in_pid:"
	REDIS_KEY_COMMENT_UGC_POST_FILED_TIME        = "time"
	REDIS_KEY_COMMENT_UGC_POST_FILED_PARENT_ID   = "parent_id"
	REDIS_KEY_COMMENT_UGC_POST_FILED_PARENT_TYPE = "parent_type"
	REDIS_KEY_COMMENT_UGC_POST_FILED_FLOOR_NO    = "floor_no"
	REDIS_KEY_COMMENT_UGC_POST_FILED_CONTENT     = "content"
	REDIS_KEY_COMMENT_UGC_POST_FILED_UID         = "uid"

	REDIS_KEY_FILED_DATABASE_EXISTS = "database_exists" //标识数据库中是否有该数据 1:表示存在 2:表示不存在

	REDIS_KEY_COMMENT_UGC_PICS_IN_PID  = "comment.in:ugc:pics_in_pid:"
	REDIS_KEY_COMMENT_UGC_PIDS_IN_ID   = "comment.in:ugc:pids_in_id:"  //type_id
	REDIS_KEY_COMMENT_UGC_COUNT_IN_SET = "comment.in:ugc:count_in_set" //mem"type#id" score:comment_count

	REDIS_KEY_COMMENT_UGC_DATA_IN_CID           = "comment.in:ugc:data_in_cid:"
	REDIS_KEY_COMMENT_UGC_COMMENT_FILED_TIME    = "time"
	REDIS_KEY_COMMENT_UGC_COMMENT_FILED_PID     = "pid"
	REDIS_KEY_COMMENT_UGC_COMMENT_FILED_CONTENT = "content"
	REDIS_KEY_COMMENT_UGC_COMMENT_FILED_UID     = "uid"

	REDIS_KEY_COMMENT_UGC_PICS_IN_CID = "comment.in:ugc:pics_in_cid:"
	REDIS_KEY_COMMENT_UGC_CIDS_IN_PID = "comment.in:ugc:cids_in_pid:"

	REDIS_KEY_COMMENT_UGC_LIKE_UIDS_IN_PID         = "comment.in:ugc:like_uids_in_pid:"
	REDIS_KEY_COMMENT_UGC_PIDS_LIKE_NUMS_IN_PARENT = "comment.in:ugc:pids_like_nums_in_parent:" //parenttype_parentid mem:pid 点赞数量

	REDIS_KEY_LIKE_UGC_OBJECT_LIKE_NUMS_IN_ID   = "like.in:ugc:object_like_nums_in_id:"   //parenttype_parentid mem:type#id 在父级中的点赞排序
	REDIS_KEY_LIKE_UGC_OBJECT_LIKE_NUMS_IN_TYPE = "like.in:ugc:object_like_nums_in_type:" //parenttype mem:type#id 在父类型中的点赞排序
	REDIS_KEY_LIKE_UGC_UIDS_IN_ID               = "like.in:ugc:uids_in_id:"               //type_id mem:uid score:time
)

// =======================回放=====================
const (
	REDIS_KEY_REPLAY_NEWEST_DISTINCT = "replay.in:newest_distinct" // 最新不重复用户的回放
	REDIS_KEY_REPLAY_TOP_VIDS        = "replay.in:top_vids"        // 置顶回放列表
)

// ==========================内部服务相关 ===========================
const (
	REDIS_KEY_FEEDS_QUEUE_WAIT_DEAL_IDS = "feeds.queue:wait_deal_ids"
	REDIS_KEY_USER_QUEUE_SINGLE_PUSH    = "user.queue:single_push"
)

//=========================tab管理==============================
const (
	REDIS_KEY_TAB_CONTEXT_IN_TAB      = "tab.in:context_in_tab"
	REDIS_KEY_FILED_NEWEST_VIDEO_TIME = "newest_video_time" // 最新加入的视频时间 主要用于红点操作
)

//=========================归类标签==============================
const (
	REDIS_KEY_LABEL_DATA_IN_LID = "label.in:data_in_lid:"
	REDIS_FIELD_NAME            = "name"
	REDIS_FIELD_ADD_TIME        = "add_time"
	REDIS_FIELD_EDIT_TIME       = "edit_time"
	REDIS_FIELD_FLAG            = "flag"

	REDIS_KEY_LABEL_RECOMMEND_LIDS = "label.in:recommend_lids" // 推荐标签

	REDIS_KEY_LABEL_UIDS_IN_LID = "label.in:uids_in_lid:" // 标签下的uid
)

//=========================活动信息==============================
const (
	REDIS_KEY_ACT_BASE_DATA_IN_ID     = "activity.in:base_data_in_id:" // 活动数据
	REDIS_FIELD_ACT_BASE_TYPE         = "type"                         // 类型
	REDIS_FIELD_ACT_BASE_TITLE        = "title"                        // 活动标题
	REDIS_FIELD_ACT_BASE_UID          = "uid"                          // 发起人uid
	REDIS_FIELD_ACT_BASE_IMAGE_URL    = "image_url"                    // 封面
	REDIS_FIELD_ACT_BASE_BEGIN_TIME   = "begin_time"                   // 开始时间
	REDIS_FIELD_ACT_BASE_END_TIME     = "end_time"                     // 结束时间
	REDIS_FIELD_ACT_BASE_COVER_URL    = "cover_url"                    // 封面url
	REDIS_FIELD_ACT_BASE_COVER_WIDTH  = "cover_width"                  // 封面宽
	REDIS_FIELD_ACT_BASE_COVER_HEIGHT = "cover_height"                 // 封面高

	// CrowFunding(众筹)数据
	REDIS_KEY_CF_DATA_IN_ACT_ID  = "activity.in:cf_data_in_act_id:" // 众筹数据
	REDIS_FIELD_CF_TARGET_AMOUNT = "target_amount"                  // 目标金额 单位为分
	REDIS_FIELD_CF_LIMIT_AMOUNT  = "limit_amount"                   // 上限金额 单位为分
	REDIS_FIELD_CF_PAID_AMOUNT   = "paid_amount"                    // 已筹金额 单位为分
	REDIS_FIELD_CF_PAID_NUM      = "paid_num"                       // 支付数量

	REDIS_KEY_USER_CF_NUM_IN_ACT_ID = "activity.in:user_cf_num_in_act_id:" // 存储用户的订单数量 mem:uid score:num
	REDIS_KEY_ACT_IDS_IN_UID        = "activity.in:act_ids_in_uid:"        // 存储用户参与的活动 mem:actid score: time

	REDIS_KEY_QUEUE_REFUND_ACT_IDS = "queue.in:refund_act_ids" // 存储需要退款的活动id
)
