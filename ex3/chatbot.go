package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Line API Settings
const LINE_ACCESS_TOKEN = "3yfZsB8sp7t0lOrAZb7Hsr1wQwDJcV5TT4ue98brdPbLcBnmi/be4U1mlFnfrb++oScuxOhzAm4JTzVLKLDVA7lPan7RekKm9s0R4WOQSj2eo1WAL1DArKYhyoGdhbzgkWDr6gtn8Akl4F3iwKR1wAdB04t89/1O/w1cDnyilFU="
const LINE_API_URL = "https://api.line.me/v2/bot/message/reply"

// Struct สำหรับรับ Webhook Request
type LineEvent struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Message    struct {
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

// Struct สำหรับ Message
type LineMessage struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

// ฟังก์ชันส่งข้อความไปยัง LINE
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
		log.Println("Error binding JSON:", err) // log ข้อผิดพลาดที่เกิดขึ้นในการแปลง JSON
		return
	}

	log.Println("Received webhook event:", lineEvent) // log ข้อมูลที่รับมาจาก LINE webhook

	for _, event := range lineEvent.Events {
		replyToken := event.ReplyToken
		userMessage := event.Message.Text

		log.Println("User message:", userMessage) // log ข้อความจากผู้ใช้

		var messages []interface{}

		// ตรวจสอบข้อความที่ผู้ใช้ส่งเข้ามา
		switch userMessage {
		case "button":
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
		case "quick":
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
		case "carousel":
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
			messages = []interface{}{
				map[string]string{"type": "text", "text": userMessage},
			}
		}

		log.Println("Sending reply message:", messages) // log ข้อความที่จะตอบกลับ

		replyMessage(replyToken, messages) // ส่งข้อความไปยัง LINE
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	r := gin.Default()
	r.POST("/webhook", webhookHandler)
	r.Run(":5000") // Start server on port 5000
	// run by ngrok
	// พิมพ์คำสั่งนี้เพื่อเริ่ม ngrok http 5000

}
