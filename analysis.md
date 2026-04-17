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

{"self_id":3557005467,
"user_id":2274735204,
"time":1776408958,
"message_id":2005105714,
"message_seq":2005105714,
"real_id":2005105714,
"real_seq":"0",
"message_type":"private",
"sender":{"user_id":2274735204,"nickname":"￴","card":""},
"raw_message":"[CQ:onlinefile,msgId=7629618380760297283,elementId=7629618380760297282,fileName=main.txt,fileSize=74099,isDir=false]","font":14,"sub_type":"friend",
"message":[{"type":"onlinefile","data":{"msgId":"7629618380760297283","elementId":"7629618380760297282","fileName":"main.txt","fileSize":"74099","isDir":false}}],
"message_format":"array",
"post_type":"message",
"target_id":2274735204}

{"self_id":3557005467,
"user_id":2274735204,
"time":1776408962,
"message_id":740017926,
"message_seq":740017926,
"real_id":740017926,
"real_seq":"10",
"message_type":"private",
"sender":{"user_id":2274735204,"nickname":"￴","card":""},
"raw_message":"👀[CQ:face,id=311,raw=&#91;object Object&#93;][CQ:face,id=312,raw=&#91;object Object&#93;]",
"font":14,
"sub_type":"friend",
"message":[{"type":"text","data":{"text":"👀"}},
{"type":"face","data":{"id":"311","raw":{"faceIndex":311,"faceText":"[打call]","faceType":2,"packId":null,"stickerId":null,"sourceType":null,"stickerType":null,"resultId":null,"surpriseId":null,"randomType":null,"imageType":null,"pokeType":null,"spokeSummary":null,"doubleHit":null,"vaspokeId":null,"vaspokeName":null,"vaspokeMinver":null,"pokeStrength":null,"msgType":null,"faceBubbleCount":null,"oldVersionStr":null,"pokeFlag":null,"chainCount":null},"resultId":null,"chainCount":null}},
{"type":"face","data":{"id":"312","raw":{"faceIndex":312,"faceText":"[变形]","faceType":2,"packId":null,"stickerId":null,"sourceType":null,"stickerType":null,"resultId":null,"surpriseId":null,"randomType":null,"imageType":null,"pokeType":null,"spokeSummary":null,"doubleHit":null,"vaspokeId":null,"vaspokeName":null,"vaspokeMinver":null,"pokeStrength":null,"msgType":null,"faceBubbleCount":null,"oldVersionStr":null,"pokeFlag":null,"chainCount":null},"resultId":null,"chainCount":null}}],
"message_format":"array",
"post_type":"message",
"target_id":2274735204}