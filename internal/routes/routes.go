package routes

import (
	config2 "bot/config"
	"bot/internal/hook"
	"bot/internal/sendMsg"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

/*
{"self_id":2748963434,  //机器人自己的QQ号
"user_id":2274735204,   //消息发送者的QQ号
"time":1776055978,		//消息发送的时间戳
"message_id":2134754754,//消息ID，消息的唯一标识
"message_seq":2134754754,//消息序列号，表示消息在当前会话中的顺序
"real_id":2134754754,//消息的真实ID，通常与message_id相同
"real_seq":"55",//消息的真实序列号，通常与message_seq相同
"message_type":"private",//消息类型，表示消息的来源和性质，可能的值包括private（私聊消息）、group（群消息）等
"sender":{"user_id":2274735204,"nickname":"￴","card":""},//发送者的信息，包括user_id（发送者的QQ号）、nickname（发送者的昵称）和card（发送者在群中的备注）
"raw_message":"你好",//消息的原始内容，包含消息中的文本和CQ码等
"font":14,//消息的字体大小，通常为14
"sub_type":"friend",//消息的子类型，表示消息的具体类型，可能的值包括friend（好友消息）、group（群消息）等
"message":[{"type":"text","data":{"text":"你好"}}],//消息的结构化内容，通常是一个数组，每个元素表示消息的一部分，包含type（消息类型）和data（消息数据）等字段
"message_format":"array",//消息的格式，表示message字段的格式，可能的值包括array（数组格式）等
"post_type":"message",//事件的类型，表示事件的性质，可能的值包括message（消息事件）等
"target_id":2274735204}//消息的目标ID，表示消息的接收者的QQ号，通常与user_id相同
*/

/*
{"self_id":2748963434,
"user_id":2274735204,
"time":1776056153,
"message_id":145065637,
"message_seq":145065637,
"real_id":145065637,
"real_seq":"41669",
"message_type":"group",
"sender":{"user_id":2274735204,"nickname":"￴","card":"","role":"owner"},
"raw_message":"[CQ:at,qq=2748963434]",
"font":14,
"sub_type":"normal",
"message":[{"type":"at","data":{"qq":"2748963434"}},{"type":"text","data":{"text":" "}}],
"message_format":"array",
"post_type":"message",
"group_id":434395907,
"group_name":"月黑风高做作业夜"}
*/

//{"self_id":2748963434,"user_id":2274735204,"time":1776056890,"message_id":1437714257,"message_seq":1437714257,"real_id":1437714257,"real_seq":"41670","message_type":"group","sender":{"user_id":2274735204,"nickname":"￴","card":"","role":"owner"},"raw_message":"[CQ:image,file=2E8289BBC260D9F93ED86E69387ACB93.jpg,sub_type=0,url=https://multimedia.nt.qq.com.cn/download?appid=1407&amp;fileid=EhRyld68OCCPnlXppM_GU4DHM3cTzRjl0gog_wooz4-xwIfqkwMyBHByb2RQgL2jAVoQuUwyWBeMqO2e2Xr3kvwRmHoCEhmCAQJneg&amp;rkey=CAMSMNSNtPNodN3RIGV9uivGRSMTpN5fOrtxxt-ORnqrj6fWA2g7jIQVGeyeZfB2KvGZuQ,file_size=174437]","font":14,"sub_type":"normal","message":[{"type":"image","data":{"summary":"","file":"2E8289BBC260D9F93ED86E69387ACB93.jpg","sub_type":0,"url":"https://multimedia.nt.qq.com.cn/download?appid=1407&fileid=EhRyld68OCCPnlXppM_GU4DHM3cTzRjl0gog_wooz4-xwIfqkwMyBHByb2RQgL2jAVoQuUwyWBeMqO2e2Xr3kvwRmHoCEhmCAQJneg&rkey=CAMSMNSNtPNodN3RIGV9uivGRSMTpN5fOrtxxt-ORnqrj6fWA2g7jIQVGeyeZfB2KvGZuQ","file_size":"174437"}}],"message_format":"array","post_type":"message","group_id":434395907,"group_name":"月黑风高做作业夜"}

/*
"message_type" 有private和group两种，分别表示私聊消息和群消息。

*/

//正确解析，再加上白名单功能
//发送给llm的消息要确保正确
//实现通过调用napcat的api进行发送消息

//实现记忆系统

//{message group 2274735204 434395907 [CQ:at,qq=2748963434] 你好 2748963434}

func Run(chatModel *ark.ChatModel, ctx context.Context, messages []*schema.Message) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	_, expectedToken := config2.LoadNapcatConfig()
	judgeEnable := config2.LoadJudgeAtConfig()
	r.Use(
		hook.CheckToken(expectedToken),
		hook.ParseMsg(),
		hook.JudgeAt(judgeEnable),
	)
	r.POST("/", sendMsg.SendMsg(chatModel, ctx, messages))
	r.POST("/test", func(c *gin.Context) {
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}
		log.Println(string(all))
		c.JSON(http.StatusOK, gin.H{})
	})
	port := config2.LoadServerConfig()
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
		return
	}
}
