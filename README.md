# SWD Backend Assessment

คำตอบแบบทดสอบ Backend ทั้ง 4 ส่วน ข้อ 4 มี Django REST Framework ตามโจทย์ และ Go เป็นผลงานเพิ่มเติม

| ส่วน | เนื้อหา | คะแนน |
|---|---|---|
| [1. Logic Test](1_logic_test/README.md) | ฟังก์ชัน Go ทั้ง 4 ข้อ พร้อม unit tests | 2 |
| [2. Database Concept](2_db_concept/README.md) | ACID และ Indexing | 5 |
| [3. API Concept](3_api_concept/README.md) | Idempotency พร้อมตัวอย่าง Go และ JSON/Protobuf | 5 |
| 4. REST API | ระบบโรงเรียน [Django REST Framework](4_rest_api_django/README.md) และ [Go / Echo / GORM](4_rest_api_go/README.md) | 8 |

## โครงสร้าง

```text
swd-backend/
├── 1_logic_test/
├── 2_db_concept/
├── 3_api_concept/
├── 4_rest_api_go/
│   ├── main.go
│   ├── database/
│   ├── domain/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   ├── response/
│   ├── service/
│   ├── errorcode/
│   ├── utils/
│   └── docs/
├── 4_rest_api_django/
│   ├── config/
│   ├── schools/
│   └── manage.py
├── go.mod
└── README.md
```

ใช้ Go module เดียวที่ repo root เพื่อรันเทสต์ทุกส่วนได้ด้วยคำสั่งเดียว

## เริ่มใช้งาน Django

ดูคำสั่งติดตั้งและรันใน [README ของ Django](4_rest_api_django/README.md) ทั้งสองเวอร์ชันใช้ Postman collection เดียวกัน โดยตั้ง `base_url` ให้ตรงกับพอร์ตของ API ที่จะทดสอบ

## เริ่มใช้งาน Go

ต้องมี Go 1.26 และ C compiler สำหรับ SQLite driver บน macOS ใช้ Xcode Command Line Tools ได้

```sh
go mod download
go run ./4_rest_api_go
```

API เริ่มที่ http://localhost:8080 และเก็บ SQLite ใน school.db ของ working directory เปลี่ยน HTTP_ADDR และ DATABASE_DSN ผ่าน environment ได้

- [Database diagram](4_rest_api_go/docs/database-diagram.md)
- [API spec แบบ OpenAPI 3.0](4_rest_api_go/docs/openapi.json)
- [Postman collection](4_rest_api_go/docs/swd-backend.postman_collection.json)
- [Request flow, transaction และ ER diagram](4_rest_api_go/docs/flow.md)
- [วิธีรัน API และตัวอย่าง curl](4_rest_api_go/README.md)

## ตรวจสอบจาก repo root

```sh
go test -race ./...
go vet ./...
go build -o ./bin/server ./4_rest_api_go
```

ตรวจแยกส่วนได้ด้วย go test ./1_logic_test และ go test ./4_rest_api_go/... ส่วน Database/API Concept เป็นคำตอบใน README ของแต่ละโฟลเดอร์
