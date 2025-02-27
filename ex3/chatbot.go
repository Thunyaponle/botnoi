package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const LINE_ACCESS_TOKEN = "3yfZsB8sp7t0lOrAZb7Hsr1wQwDJcV5TT4ue98brdPbLcBnmi/be4U1mlFnfrb++oScuxOhzAm4JTzVLKLDVA7lPan7RekKm9s0R4WOQSj2eo1WAL1DArKYhyoGdhbzgkWDr6gtn8Akl4F3iwKR1wAdB04t89/1O/w1cDnyilFU="
const LINE_API_URL = "https://api.line.me/v2/bot/message/reply"

type LineEvent struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Message    struct {
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type LineMessage struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

func replyMessage(replyToken string, messages []interface{}) {
	body, _ := json.Marshal(LineMessage{ReplyToken: replyToken, Messages: messages})
	req, _ := http.NewRequest("POST", LINE_API_URL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+LINE_ACCESS_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending message:", err)
	} else {
		defer resp.Body.Close()
		log.Println("LINE API response:", resp.Status)
	}
}

func webhookHandler(c *gin.Context) {
	var lineEvent LineEvent
	if err := c.BindJSON(&lineEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		log.Println("Error binding JSON:", err)
		return
	}

	log.Println("Received webhook event:", lineEvent)

	for _, event := range lineEvent.Events {
		replyToken := event.ReplyToken
		userMessage := event.Message.Text

		log.Println("User message:", userMessage)

		var messages []interface{}

		// ถ้าไม่มีข้อความที่ผู้ใช้พิมพ์ ระบบจะแสดงคำแนะนำเบื้องต้น
		if userMessage == "options" {
			messages = []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "ยินดีต้อนรับ! คุณสามารถใช้คำสั่งเหล่านี้:\n- button: แสดงปุ่มตัวเลือก\n- quick reply: แสดงคำตอบด่วน\n- carousel: แสดงตัวเลือกแบบหมุนเวียน",
				},
				map[string]interface{}{
					"type":    "template",
					"altText": "Choose an option",
					"template": map[string]interface{}{
						"type": "buttons",
						"text": "เลือกตัวเลือก",
						"actions": []map[string]string{
							{"type": "message", "label": "button", "text": "button"},
							{"type": "message", "label": "quick reply", "text": "quick reply"},
							{"type": "message", "label": "carousel", "text": "carousel"},
						},
					},
				},
			}
		} else if userMessage == "Options" {
			messages = []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "ยินดีต้อนรับ! คุณสามารถใช้คำสั่งเหล่านี้:\n- button: แสดงปุ่มตัวเลือก\n- quick reply: แสดงคำตอบด่วน\n- carousel: แสดงตัวเลือกแบบหมุนเวียน",
				},
				map[string]interface{}{
					"type":    "template",
					"altText": "Choose an option",
					"template": map[string]interface{}{
						"type": "buttons",
						"text": "เลือกตัวเลือก",
						"actions": []map[string]string{
							{"type": "message", "label": "button", "text": "button"},
							{"type": "message", "label": "quick reply", "text": "quick reply"},
							{"type": "message", "label": "carousel", "text": "carousel"},
						},
					},
				},
			}
		} else {
			// ถ้าผู้ใช้พิมพ์คำที่รู้จัก
			switch userMessage {
			case "button", "Button":
				messages = []interface{}{
					map[string]interface{}{
						"type":    "template",
						"altText": "This is a button template",
						"template": map[string]interface{}{
							"type": "buttons",
							"text": "เลือกตัวเลือก",
							"actions": []map[string]string{
								{"type": "message", "label": "Option 1", "text": "เลือก Option 1"},
								{"type": "message", "label": "Option 2", "text": "เลือก Option 2"},
							},
						},
					},
				}
			case "quick reply", "Quick reply":
				messages = []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": "เลือกคำตอบ",
						"quickReply": map[string]interface{}{
							"items": []map[string]interface{}{
								{"type": "action", "action": map[string]string{"type": "message", "label": "A", "text": "A"}},
								{"type": "action", "action": map[string]string{"type": "message", "label": "B", "text": "B"}},
							},
						},
					},
				}
			case "carousel", "Carousel":
				messages = []interface{}{
					map[string]interface{}{
						"type":    "template",
						"altText": "Carousel Example",
						"template": map[string]interface{}{
							"type": "carousel",
							"columns": []map[string]interface{}{
								{
									"text": "Item 1",
									"actions": []map[string]string{
										{"type": "message", "label": "เลือก 1", "text": "เลือก 1"},
									},
								},
								{
									"text": "Item 2",
									"actions": []map[string]string{
										{"type": "message", "label": "เลือก 2", "text": "เลือก 2"},
									},
								},
							},
						},
					},
				}
			default:
				// ถ้าไม่ใช่คำที่รู้จักจะส่งข้อความกลับ
				messages = []interface{}{
					map[string]string{"type": "text", "text": "ข้อความของคุณ: " + userMessage},
					map[string]string{"type": "text", "text": "คุณสามารถใช้คำสั่งเหล่านี้:\n- options: เพื่อแสดงตัวเลือกทั้งหมด\n- button: แสดงปุ่มตัวเลือก\n- quick reply: แสดงคำตอบด่วน\n- carousel: แสดงตัวเลือกแบบหมุนเวียน: "},
				}
			}
		}

		log.Println("Sending reply message:", messages)

		replyMessage(replyToken, messages)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	r := gin.Default()
	r.POST("/webhook", webhookHandler)
	r.Run(":5000")

}

//ตอน run ให้พี่แก้ LINE_ACCESS_TOKEN เป็นค่าของไลน์พี่นะคะ
// run by ngrok (ngrok.exe)
// พิมพ์คำสั่งนี้เพื่อเริ่ม : ngrok http 5000
