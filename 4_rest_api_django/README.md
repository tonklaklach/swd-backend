# REST API ระบบโรงเรียนด้วย Django

เวอร์ชันสำหรับข้อ 4 ที่ใช้ Django REST Framework, ModelSerializer และ django-filter รองรับ CRUD ของโรงเรียน ห้องเรียน ครู และนักเรียน รวมถึง filters และรายละเอียดความสัมพันธ์ตามโจทย์

## รันในเครื่อง

ใช้ Python 3.12 ขึ้นไป ตัวอย่างนี้รันจาก root ของ repo

```sh
cd 4_rest_api_django
python3.12 -m venv .venv
source .venv/bin/activate
python -m pip install -r requirements.txt
python manage.py migrate
python manage.py runserver 127.0.0.1:8080
```

Django เก็บข้อมูลใน `4_rest_api_django/db.sqlite3` กำหนด `DJANGO_DB_PATH` เพื่อใช้ไฟล์อื่นได้ ถ้า Go ใช้พอร์ต 8080 อยู่ ให้รัน Django ด้วย `python manage.py runserver 127.0.0.1:8081`

## ทดสอบด้วย Postman

Import [Postman collection เดิม](../4_rest_api_go/docs/swd-backend.postman_collection.json) แล้วตั้ง `base_url` เป็น `http://127.0.0.1:8080` หรือพอร์ตที่ Django เปิดอยู่ ใช้ Collection Runner รันจาก Health check ถึง Cleanup ตามลำดับ สคริปต์จะเก็บ ID ของข้อมูลที่สร้างและลบข้อมูลตัวอย่างในขั้นสุดท้าย

Django ใช้ path และรูปแบบ response จาก [OpenAPI spec เดิม](../4_rest_api_go/docs/openapi.json) โดย path ไม่มี slash ท้าย เช่น `GET /api/v1/schools/1` รายการอยู่ใน `data` และ DELETE สำเร็จตอบ 204 พร้อม body ว่าง

## โครงสร้างและการตรวจสอบ

- `schools/models.py` และ `schools/migrations/` กำหนดความสัมพันธ์และ constraints
- `schools/serializers.py` ตรวจ request และจัดรูปแบบ list/detail response
- `schools/filters.py` ใช้ FilterSet สำหรับ query parameters และจำกัดจำนวนผลลัพธ์
- `schools/views.py` ลงมือบันทึกและเตรียม response ภายใน transaction เดียวกัน หากขั้นตอนใดล้มเหลว การเขียนจะถูก rollback

```sh
python manage.py check
python manage.py makemigrations --check --dry-run
python manage.py test
```
