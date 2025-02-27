import requests
from datetime import datetime, timedelta

api_key = "AIzaSyCb2vieDF4CbnYUvmhf0ZRgk-JTCV5PqN4"
print(f"Your API key is: {api_key}")

if not api_key:
    print("API Key not found. Please set the API_KEY in the code.")
    exit(1)

def get_tomorrow():
    today = datetime.today()
    tomorrow = today + timedelta(days=1)
    return tomorrow.strftime("%d-%m-%Y")

def call_gemini_api(input_text):
    url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + api_key

    prompt = f"""
    แปลงวันที่ในข้อความให้เป็น JSON ที่มีโครงสร้าง:
    - ถ้าข้อมูลครบถ้วนให้ return JSON รูปแบบ {{ "year": "YYYY", "month": "MM", "day": "DD" }}
    - ถ้าข้อมูลไม่ครบ เช่น ไม่รู้ปี ให้ใช้ "-" แทน
    - ตัวอย่าง: "7 เม.ย. 2565" -> {{ "year": "2565", "month": "04", "day": "07" }}
    - ตัวอย่าง: "พรุ่งนี้" -> วันพรุ่งนี้: {get_tomorrow()}
    
    ข้อความที่ต้องแปลง: "{input_text}"
    """

    headers = {
        'Content-Type': 'application/json'
    }

    data = {
        "contents": [{"parts": [{"text": prompt}]}]
    }

    try:
        response = requests.post(url, json=data, headers=headers)
        response.raise_for_status()
        
        result = response.json()
        
        if 'contents' in result:
            return result['contents'][0]['parts'][0]['text']
        else:
            print("ไม่พบข้อมูลที่ต้องการจาก API")
            return None
    except requests.exceptions.RequestException as e:
        print(f"Error calling Gemini API: {e}")
        return None

input_text = "7 เม.ย. 2565"
result = call_gemini_api(input_text)

if result:
    print("Response from Gemini API:", result)
else:
    print("ไม่สามารถดึงข้อมูลจาก API ได้")
