# REST API flow

```mermaid
flowchart LR
    C[Client] --> E[Echo routes]
    E --> H[Handler binds and validates request]
    H --> S[Service]
    S --> R[GORM repository]
    R --> D[(SQLite)]
    D --> R
    R --> S
    S --> H
    H --> J[JSON response]
    J --> C
```

main.go เปิด SQLite และประกอบ dependencies ก่อน register routes ใน Echo จากนั้น handler อ่าน path และ query parameters, bind JSON และตรวจ validation tags ตามประเภท request แล้วส่ง request context และข้อมูลให้ service ผ่าน interface

Service ทำงานตาม use case และเรียก repository เช่น ประกอบรายละเอียดห้องเรียนจากรายชื่อครูและนักเรียน ส่วน repository ใช้ GORM สำหรับ CRUD, filters และ transactions หากเกิด error ตัวจัดการกลางของ Echo จะแปลงเป็น HTTP status และ error code

## Create และ Update

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant S as Service
    participant R as Repository
    participant D as SQLite
    C->>H: POST or PUT
    H->>H: Bind and validate
    H->>S: Create or Update
    S->>R: Save resource
    R->>D: Begin transaction and write
    opt Teacher
        R->>D: Replace classroom assignments
    end
    R->>D: Read saved data and relationships
    alt Write and read succeed
        R->>D: Commit
        R-->>S: Saved data after successful commit
        S-->>H: Saved data
        H-->>C: 201 Created or 200 OK
    else Write or read fails
        R->>D: Rollback
        R-->>S: Error
        S-->>H: Error
        H-->>C: JSON error
    end
```

Repository อ่านข้อมูลที่จะคืนใน transaction เดียวกับการบันทึก ถ้าอ่านกลับล้มเหลว การเปลี่ยนแปลงทั้งหมดจะถูก rollback รวมถึงรายการห้องที่สอนของครู

[ER diagram, keys และ delete constraints](database-diagram.md)
