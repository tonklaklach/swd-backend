# REST API ระบบโรงเรียน

Go และ Echo v5 แยก handler/service/repository ด้วย interfaces สำหรับระบบโรงเรียน ห้องเรียน ครู และนักเรียน ใช้ Go module จาก repo root

## รันโปรเจกต์

รันคำสั่งในโฟลเดอร์ 4_rest_api_go ต้องมี Go 1.26 และ C compiler เพราะ SQLite driver ใช้ CGO บน macOS ใช้ Xcode Command Line Tools ได้

```sh
go mod download
go run .
```

API เริ่มที่ http://localhost:8080 เก็บข้อมูลใน school.db ซึ่งสร้าง schema ครั้งแรกให้อัตโนมัติจาก database/schema.sql ผ่าน GORM เพื่อกำหนด foreign key และ delete constraints ให้ชัดเจน เปลี่ยน HTTP_ADDR และ DATABASE_DSN ผ่าน environment ได้

```sh
HTTP_ADDR=:8081 DATABASE_DSN=./school.db go run .
```

## เอกสารและโครงสร้าง

- [Database diagram และ constraints](docs/database-diagram.md)
- [API spec แบบ OpenAPI 3.0](docs/openapi.json) นำเข้า Swagger Editor หรือ Postman ได้
- [Request flow, teacher transaction และ ER diagram](docs/flow.md)
- main.go ประกอบ database → repository → service → handler และ register Echo routes ทั้งหมด พร้อม graceful shutdown
- database เปิด GORM ด้วย SQLite driver และ initialize schema ผ่าน GORM Exec
- domain เก็บ entities, typed filters และ domain errors
- model เก็บ request structs และ validation tags ส่วน service แปลง writable fields เป็น domain models
- handler รับ HTTP request, bind request DTO, ตรวจ tags ด้วย go-playground/validator, เรียก service ด้วย request context และส่ง response
- service เก็บ service interfaces และ implementation รับ context.Context และ typed inputs ทำ business logic กับประกอบผลลัพธ์ผ่าน repository interfaces
- repository ใช้ GORM ทำ CRUD, query scopes, eager loading, transaction และแปลง SQLite errors เป็น domain errors
- response/struct เก็บ success/error envelope และ detail response models
- response ดูแลการส่ง success และ error response ส่วน errorcode เก็บค่ารหัสข้อผิดพลาด
- utils เก็บ Bind, Validate, ParseID และ ParseFilter ซึ่งใช้ร่วมกันใน handler ของทั้งสี่ resource

## ตัวอย่างใช้งาน

```sh
curl -i localhost:8080/healthz

curl -i -X POST localhost:8080/api/v1/schools   -H 'Content-Type: application/json'   -d '{"name":"Swift School","abbreviation":"SWD","address":"Bangkok"}'

curl -i -X POST localhost:8080/api/v1/classrooms   -H 'Content-Type: application/json'   -d '{"school_id":1,"year":1,"section":1}'

curl -i -X POST localhost:8080/api/v1/teachers   -H 'Content-Type: application/json'   -d '{"firstname":"Ann","lastname":"Smith","gender":"female","classroom_ids":[1]}'

curl -i -X POST localhost:8080/api/v1/students   -H 'Content-Type: application/json'   -d '{"firstname":"Bob","lastname":"Lee","gender":"male","classroom_id":1}'

curl -i 'localhost:8080/api/v1/teachers?school=1&classroom=1&firstname=ann'
curl -i localhost:8080/api/v1/schools/1
```

## Postman

Import [Postman collection](docs/swd-backend.postman_collection.json) แล้วตั้ง collection variable base_url เป็น http://localhost:8080 จากนั้นใช้ Collection Runner รันทุก folder ตามลำดับ

Collection จะสร้าง school → classroom → teacher/student และเก็บ ID ให้อัตโนมัติ มีตัวอย่าง filters, PUT, detail, validation errors และ cleanup ที่ลบข้อมูลตัวอย่างตอนท้าย หากรันทีละ request ให้รัน Create sample data ก่อน หรือกำหนด ID variables เอง

มี tests ตรวจ HTTP status, success/error envelope และความสัมพันธ์ข้อมูล หากหยุดกลางทางให้ใช้ Cleanup sample data กับ variables ชุดเดิมเพื่อลบข้อมูลที่สร้างแล้ว

## API conventions

ทุก resource รองรับ GET list, POST create, GET detail, PUT replace และ DELETE ที่ /api/v1/schools, /classrooms, /teachers, /students รายละเอียดฟิลด์และ filters อยู่ใน spec

- Success response ใช้ {"data": ...} ส่วน DELETE ที่สำเร็จตอบ 204 โดยไม่มี body
- List อยู่ใน data เป็น array เรียงตาม id รองรับ limit 1–100 (ค่าเริ่มต้น 50) และ offset ตั้งแต่ 0 ส่วนข้อมูลใน detail ส่งรายการที่เกี่ยวข้องทั้งหมด
- ชื่อค้นหาแบบ substring โดย SQLite lower รองรับการแปลงตัวพิมพ์ ASCII ส่วน ID และ gender เทียบตรงตัว
- PUT ส่ง writable fields ครบ ครูที่ไม่ส่ง classroom_ids จะมีรายการห้องว่าง
- ครูอาจยังไม่มีห้องและสอนข้ามโรงเรียนได้ โรงเรียนของครูคำนวณผ่านห้องที่สอน
- Student ต้องมีห้องอยู่จริง Classroom ต้องมีโรงเรียนอยู่จริง และชั้นปี/ทับเป็นจำนวนเต็มบวก
- ห้องห้ามซ้ำ school_id + year + section ส่วน gender รับ male, female, other
- Teacher update เปลี่ยนชื่อและความสัมพันธ์ใน transaction เดียว การอ้างห้องที่ไม่มีอยู่จะ rollback ทั้งชุด
- ลบ school/classroom ที่มีข้อมูลอ้างอิงตอบ 409 ส่วนลบ teacher ลบเฉพาะครูกับรายการเชื่อมห้อง
- Error ใช้ {"code":"ERROR_CODE","error":"message"} โดย malformed input ตอบ 400, missing record 404, conflict 409, content type 415 และ invalid fields/reference 422
- POST/PUT รับ JSON ไม่เกิน 1 MiB และปฏิเสธ unknown fields ส่วน id ใน body ยอมรับแต่ใช้ ID จากฐานข้อมูลหรือ URL
- Authentication ไม่ได้ระบุในโจทย์ ตัวอย่างนี้เปิด API สำหรับการตรวจงาน

## ตรวจสอบ

```sh
go test -race ./...
go vet ./...
go build -o ./bin/server .
```

เทสต์ใช้ SQLite จริงใน temporary directory ครอบคลุม CRUD, filters, pagination, จำนวนครูแบบ distinct, constraint และ rollback ฐานข้อมูลไฟล์หลักไม่ถูกเปลี่ยนโดยเทสต์
