# Logic Test

ฟังก์ชัน Go และ unit tests สำหรับโจทย์ทั้ง 4 ข้อ อยู่ใน logic.go และ logic_test.go

## Logic Test และสมมติฐาน

- TrailingZeroes นับตัวประกอบ 5 ใช้ O(log n)
- MaxIndex ใช้ O(n) และเลือก index แรกเมื่อค่าสูงสุดซ้ำ
- NumberToThai รับ 0–10,000,000 รองรับ ยี่สิบ และ เอ็ด
- NumberToRoman ใช้ช่วงมาตรฐาน 1–3999 เพราะโจทย์ไม่ได้กำหนดขอบเขต และ 0 ไม่มีสัญลักษณ์ในรูปแบบที่เลือก
- Input ที่ผิดเงื่อนไขคืน error ตามข้อความในโจทย์เมื่อมีระบุไว้

## รันเทสต์จาก repo root

```sh
go test ./1_logic_test
```
