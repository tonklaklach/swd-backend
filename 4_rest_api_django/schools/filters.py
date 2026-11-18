import re

import django_filters
from django import forms
from django.db.models import Exists, OuterRef, Value
from django.db.models.functions import Lower, StrIndex
from rest_framework.exceptions import ParseError

from .models import Classroom, Gender, School, Student, Teacher, TeachingAssignment
from .serializers import MAX_ID


class PositiveIDField(forms.IntegerField):
    def to_python(self, value):
        if value not in self.empty_values and not re.fullmatch(r"[+-]?[0-9]+", value):
            raise forms.ValidationError("Enter an integer ID")
        return super().to_python(value)


class PositiveIDFilter(django_filters.Filter):
    field_class = PositiveIDField

    def __init__(self, **kwargs):
        super().__init__(min_value=1, max_value=MAX_ID, **kwargs)


class SubstringFilter(django_filters.CharFilter):
    def filter(self, queryset, value):
        if value in (None, ""):
            return queryset
        match = StrIndex(Lower(self.field_name), Lower(Value(value)))
        return queryset.alias(substring_position=match).filter(substring_position__gt=0)


class SchoolFilter(django_filters.FilterSet):
    name = SubstringFilter(strip=False)

    class Meta:
        model = School
        fields = ()


class ClassroomFilter(django_filters.FilterSet):
    school = PositiveIDFilter(field_name="school_id")

    class Meta:
        model = Classroom
        fields = ()


class PersonFilter(django_filters.FilterSet):
    firstname = SubstringFilter(strip=False)
    lastname = SubstringFilter(strip=False)
    gender = django_filters.ChoiceFilter(choices=Gender.choices)


class TeacherFilter(PersonFilter):
    school = PositiveIDFilter(method="by_school")
    classroom = PositiveIDFilter(method="by_classroom")

    def by_school(self, queryset, name, value):
        assignments = TeachingAssignment.objects.filter(
            teacher_id=OuterRef("pk"), classroom__school_id=value
        )
        return queryset.filter(Exists(assignments))

    def by_classroom(self, queryset, name, value):
        assignments = TeachingAssignment.objects.filter(
            teacher_id=OuterRef("pk"), classroom_id=value
        )
        return queryset.filter(Exists(assignments))

    class Meta:
        model = Teacher
        fields = ()


class StudentFilter(PersonFilter):
    school = PositiveIDFilter(field_name="classroom__school_id")
    classroom = PositiveIDFilter(field_name="classroom_id")

    class Meta:
        model = Student
        fields = ()


def filter_page(filterset_class, params, queryset):
    fields = filterset_class.base_filters
    for name, values in params.lists():
        if name not in fields and name not in {"limit", "offset"}:
            raise ParseError(f"Unknown filter {name}")
        if len(values) != 1:
            raise ParseError("Query parameters must appear once")
        if values[0] == "" and name in {"school", "classroom", "gender"}:
            raise ParseError(f"{name} cannot be empty")

    def number(name, default, minimum, maximum):
        value = params.get(name, str(default))
        if not re.fullmatch(r"[+-]?[0-9]{1,19}", value):
            raise ParseError(f"Invalid {name}")
        parsed = int(value)
        if not minimum <= parsed <= maximum:
            raise ParseError(f"Invalid {name}")
        return parsed

    limit = number("limit", 50, 1, 100)
    offset = number("offset", 0, 0, MAX_ID)
    filters = filterset_class(params, queryset=queryset)
    if not filters.is_valid():
        field, messages = next(iter(filters.errors.items()))
        raise ParseError(f"{field}: {messages[0]}")
    return filters.qs[offset : min(offset + limit, MAX_ID)]
