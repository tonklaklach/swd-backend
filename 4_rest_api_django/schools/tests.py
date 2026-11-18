import json
from unittest.mock import PropertyMock, patch

from rest_framework.test import APITestCase

from .models import Classroom, School, Student, Teacher
from .serializers import TeacherSerializer


class SchoolAPITests(APITestCase):
    def setUp(self):
        self.school = School.objects.create(name="Alpha", abbreviation="A", address="Bangkok")
        other = School.objects.create(name="Beta", abbreviation="B", address="Bangkok")
        self.room = Classroom.objects.create(school=self.school, year=1, section=1)
        self.second_room = Classroom.objects.create(school=self.school, year=1, section=2)
        other_room = Classroom.objects.create(school=other, year=1, section=1)
        self.teacher = Teacher.objects.create(firstname="Ann", lastname="Lee", gender="female")
        self.teacher.classrooms.add(self.room, self.second_room)
        other_teacher = Teacher.objects.create(firstname="Ann", lastname="Smith", gender="male")
        other_teacher.classrooms.add(other_room)
        self.student = Student.objects.create(
            firstname="Bob", lastname="Lee", gender="male", classroom=self.room
        )
        Student.objects.create(
            firstname="Bob", lastname="Lee", gender="female", classroom=other_room
        )

    def request_data(self, method, path, expected_status, payload=None):
        response = self.client.generic(
            method, path, json.dumps(payload) if payload is not None else "",
            content_type="application/json",
        )
        self.assertEqual(response.status_code, expected_status, response.content)
        return response

    def test_lists_filter_and_paginate_by_expected_ids(self):
        cases = {
            "schools?name=ALP": [self.school.pk],
            "classrooms?school=1&limit=1&offset=1": [self.second_room.pk],
            "teachers?school=1&classroom=1&firstname=ANN&lastname=lee&gender=female": [
                self.teacher.pk
            ],
            "students?school=1&classroom=1&firstname=bob&lastname=lee&gender=male": [
                self.student.pk
            ],
            "schools?name=%25": [],
        }
        for path, expected in cases.items():
            with self.subTest(path=path):
                response = self.request_data("GET", "/api/v1/" + path, 200)
                self.assertEqual([item["id"] for item in response.json()["data"]], expected)

    def test_details_include_counts_and_relationships(self):
        school = self.request_data("GET", f"/api/v1/schools/{self.school.pk}", 200).json()["data"]
        self.assertEqual(
            school["counts"], {"classroom_count": 2, "teacher_count": 1, "student_count": 1}
        )
        room = self.request_data("GET", f"/api/v1/classrooms/{self.room.pk}", 200).json()["data"]
        self.assertEqual(room["teachers"][0]["classroom_ids"], [self.room.pk, self.second_room.pk])
        self.assertEqual(room["students"][0]["id"], self.student.pk)

    def test_errors_and_missing_update(self):
        cases = [
            ("GET", "/api/v1/students?school=1.0", None, 400, "BAD_REQUEST"),
            ("GET", "/api/v1/students?limit=0", None, 400, "BAD_REQUEST"),
            ("GET", "/api/v1/students?school=1&school=2", None, 400, "BAD_REQUEST"),
            ("GET", "/api/v1/schools/invalid", None, 400, "BAD_REQUEST"),
            ("POST", "/api/v1/schools", {"name": " "}, 422, "VALIDATION_FAILED"),
            ("POST", "/api/v1/students", {
                "firstname": "Bob", "lastname": "Lee", "gender": "male", "classroom_id": 999
            }, 422, "INVALID_REFERENCE"),
            ("POST", "/api/v1/classrooms", {
                "school_id": self.school.pk, "year": 1, "section": 1
            }, 409, "CONFLICT"),
            ("DELETE", f"/api/v1/classrooms/{self.room.pk}", None, 409, "CONFLICT"),
            ("PUT", "/api/v1/schools/999", {
                "name": "New school", "abbreviation": "N", "address": "Bangkok"
            }, 404, "NOT_FOUND"),
        ]
        for method, path, payload, status, code in cases:
            with self.subTest(method=method, path=path):
                result = self.request_data(method, path, status, payload)
                self.assertEqual(result.json()["code"], code)
        self.assertEqual(School.objects.count(), 2)

    def test_teacher_update_rolls_back_if_serialization_fails(self):
        payload = {
            "firstname": "Changed",
            "lastname": "Lee",
            "gender": "female",
            "classroom_ids": [self.second_room.pk],
        }
        with patch.object(TeacherSerializer, "data", new_callable=PropertyMock) as response:
            response.side_effect = RuntimeError("response preparation failed")
            with self.assertLogs("schools.errors", level="ERROR"):
                result = self.request_data("PUT", f"/api/v1/teachers/{self.teacher.pk}", 500, payload)
        self.assertEqual(result.json()["code"], "INTERNAL_ERROR")
        self.teacher.refresh_from_db()
        self.assertEqual(self.teacher.firstname, "Ann")
        self.assertCountEqual(
            self.teacher.classrooms.values_list("pk", flat=True),
            [self.room.pk, self.second_room.pk],
        )

    def test_put_allows_clearing_teacher_assignments(self):
        result = self.request_data("PUT", f"/api/v1/teachers/{self.teacher.pk}", 200, {
            "firstname": "Ann", "lastname": "Lee", "gender": "female"
        })
        self.assertEqual(result.json()["data"]["classroom_ids"], [])

    def test_rejects_wrong_json_field_types(self):
        school = {"name": "Gamma", "abbreviation": "G", "address": "Bangkok"}
        room = {"school_id": self.school.pk, "year": 2, "section": 1}
        teacher = {
            "firstname": "May", "lastname": "Lee", "gender": "female",
            "classroom_ids": [self.room.pk],
        }
        student = {
            "firstname": "May", "lastname": "Lee", "gender": "female",
            "classroom_id": self.room.pk,
        }
        cases = [
            ("POST", "schools", dict(school, name=123)),
            ("POST", "schools", dict(school, id="123")),
            ("POST", "schools", dict(school, extra=True)),
            ("POST", "classrooms", dict(room, school_id=str(self.school.pk))),
            ("POST", "classrooms", dict(room, year=2.0)),
            ("POST", "classrooms", dict(room, section=True)),
            ("POST", "classrooms", dict(room, section=1 << 63)),
            ("POST", "teachers", dict(teacher, gender=1)),
            ("POST", "teachers", dict(teacher, classroom_ids="1")),
            ("POST", "teachers", dict(teacher, classroom_ids=[str(self.room.pk)])),
            ("POST", "students", dict(student, classroom_id=str(self.room.pk))),
            ("PUT", f"schools/{self.school.pk}", dict(school, name=123)),
        ]
        for method, path, payload in cases:
            with self.subTest(method=method, path=path, payload=payload):
                response = self.request_data(method, f"/api/v1/{path}", 400, payload)
                self.assertEqual(response.json()["code"], "BAD_REQUEST")
        self.school.refresh_from_db()
        self.assertEqual(self.school.name, "Alpha")
        self.assertEqual(School.objects.count(), 2)
        self.assertEqual(Classroom.objects.count(), 3)
        self.assertEqual(Teacher.objects.count(), 2)
        self.assertEqual(Student.objects.count(), 2)

    def test_keeps_valid_json_values_and_accepts_null_teacher_assignments(self):
        school = self.request_data("POST", "/api/v1/schools", 201, {
            "id": 999, "name": " Gamma ", "abbreviation": "G", "address": "Bangkok"
        }).json()["data"]
        self.assertEqual(school["name"], " Gamma ")
        self.assertNotEqual(school["id"], 999)

        teacher = self.request_data("PUT", f"/api/v1/teachers/{self.teacher.pk}", 200, {
            "firstname": "Ann", "lastname": "Lee", "gender": "female",
            "classroom_ids": None,
        }).json()["data"]
        self.assertEqual(teacher["classroom_ids"], [])
