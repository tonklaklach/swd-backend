from django.db import models
from rest_framework import serializers
from rest_framework.exceptions import ParseError

from .errors import InvalidReference
from .models import Classroom, Gender, School, Student, Teacher

MAX_ID = (1 << 63) - 1


class JSONTextField(serializers.CharField):
    def to_internal_value(self, data):
        if not isinstance(data, str):
            raise ParseError(f"{self.field_name} must be a JSON string")
        if not data.strip():
            self.fail("blank")
        return data


class JSONIntegerField(serializers.IntegerField):
    def to_internal_value(self, data):
        if type(data) is not int or not -MAX_ID - 1 <= data <= MAX_ID:
            raise ParseError("Expected a JSON integer in the signed 64-bit range")
        return super().to_internal_value(data)


class GenderField(serializers.ChoiceField):
    def to_internal_value(self, data):
        if not isinstance(data, str):
            raise ParseError("gender must be a JSON string")
        return super().to_internal_value(data)


class RequestSerializer(serializers.ModelSerializer):
    serializer_field_mapping = {
        **serializers.ModelSerializer.serializer_field_mapping,
        models.TextField: JSONTextField,
    }
    id = serializers.IntegerField(read_only=True)

    def to_internal_value(self, data):
        unknown = data.keys() - self.fields.keys()
        if unknown:
            raise ParseError(f"Unknown fields: {', '.join(sorted(unknown))}")
        resource_id = data.get("id")
        if resource_id is not None:
            JSONIntegerField().run_validation(resource_id)
        return super().to_internal_value(data)


class ClassroomReference(serializers.PrimaryKeyRelatedField):
    def to_internal_value(self, data):
        try:
            return super().to_internal_value(data)
        except serializers.ValidationError as error:
            if error.get_codes() == ["does_not_exist"]:
                raise InvalidReference() from error
            raise


class ClassroomIDs(serializers.ListField):
    child = JSONIntegerField(min_value=1, max_value=MAX_ID)

    def run_validation(self, data=serializers.empty):
        if data is None:
            data = []
        return super().run_validation(data)

    def to_internal_value(self, data):
        if not isinstance(data, list):
            raise ParseError("classroom_ids must be a JSON array")
        ids = super().to_internal_value(data)
        if len(ids) != len(set(ids)):
            raise serializers.ValidationError("Classroom IDs must be unique")
        classrooms = Classroom.objects.in_bulk(ids)
        if len(classrooms) != len(ids):
            raise InvalidReference()
        return [classrooms[pk] for pk in ids]

    def to_representation(self, classrooms):
        return sorted(classroom.pk for classroom in classrooms.all())


class SchoolSerializer(RequestSerializer):
    class Meta:
        model = School
        fields = ("id", "name", "abbreviation", "address")


class ClassroomSerializer(RequestSerializer):
    school_id = JSONIntegerField(min_value=1, max_value=MAX_ID)
    year = JSONIntegerField(min_value=1, max_value=MAX_ID)
    section = JSONIntegerField(min_value=1, max_value=MAX_ID)

    def validate_school_id(self, value):
        if not School.objects.filter(pk=value).exists():
            raise InvalidReference()
        return value

    class Meta:
        model = Classroom
        fields = ("id", "school_id", "year", "section")
        validators = []  # Report duplicate rooms as 409 from the database constraint.


class TeacherSerializer(RequestSerializer):
    gender = GenderField(choices=Gender.choices)
    classroom_ids = ClassroomIDs(source="classrooms", required=False, default=list)

    class Meta:
        model = Teacher
        fields = ("id", "firstname", "lastname", "gender", "classroom_ids")


class StudentSerializer(RequestSerializer):
    gender = GenderField(choices=Gender.choices)
    classroom_id = ClassroomReference(
        source="classroom",
        queryset=Classroom.objects.all(),
        pk_field=JSONIntegerField(min_value=1, max_value=MAX_ID),
    )

    class Meta:
        model = Student
        fields = ("id", "firstname", "lastname", "gender", "classroom_id")


class SchoolDetailSerializer(SchoolSerializer):
    counts = serializers.SerializerMethodField()

    def get_counts(self, school):
        return {
            "classroom_count": school.classrooms.count(),
            "teacher_count": Teacher.objects.filter(classrooms__school=school).distinct().count(),
            "student_count": Student.objects.filter(classroom__school=school).count(),
        }

    class Meta(SchoolSerializer.Meta):
        fields = (*SchoolSerializer.Meta.fields, "counts")


class ClassroomDetailSerializer(ClassroomSerializer):
    teachers = TeacherSerializer(many=True, read_only=True)
    students = StudentSerializer(many=True, read_only=True)

    class Meta(ClassroomSerializer.Meta):
        fields = (*ClassroomSerializer.Meta.fields, "teachers", "students")


class TeacherDetailSerializer(TeacherSerializer):
    classrooms = ClassroomSerializer(many=True, read_only=True)

    class Meta(TeacherSerializer.Meta):
        fields = (*TeacherSerializer.Meta.fields, "classrooms")


class StudentDetailSerializer(StudentSerializer):
    classroom = ClassroomSerializer(read_only=True)

    class Meta(StudentSerializer.Meta):
        fields = (*StudentSerializer.Meta.fields, "classroom")
