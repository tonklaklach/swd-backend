from django.db import DatabaseError, IntegrityError, connection, transaction
from django.db.models.deletion import ProtectedError
from rest_framework import status, viewsets
from rest_framework.exceptions import ParseError, UnsupportedMediaType
from rest_framework.response import Response
from rest_framework.views import APIView

from .errors import Conflict, InvalidReference, error_body
from .filters import ClassroomFilter, SchoolFilter, StudentFilter, TeacherFilter, filter_page
from .models import Classroom, School, Student, Teacher
from .serializers import (
    MAX_ID,
    ClassroomDetailSerializer,
    ClassroomSerializer,
    SchoolDetailSerializer,
    SchoolSerializer,
    StudentDetailSerializer,
    StudentSerializer,
    TeacherDetailSerializer,
    TeacherSerializer,
)


class ResourceViewSet(viewsets.GenericViewSet):
    detail_serializer_class = None
    http_method_names = ["get", "post", "put", "delete", "options", "head"]

    def get_object(self):
        value = self.kwargs["pk"]
        try:
            resource_id = int(value, 10)
            if not value.removeprefix("+").isascii() or not value.removeprefix("+").isdecimal():
                raise ValueError
        except ValueError:
            raise ParseError("ID must be positive") from None
        if not 1 <= resource_id <= MAX_ID:
            raise ParseError("ID must be positive")
        return super().get_object()

    def list(self, request):
        items = filter_page(self.filterset_class, request.query_params, self.get_queryset())
        return Response({"data": self.get_serializer(items, many=True).data})

    @transaction.atomic
    def retrieve(self, request, **kwargs):
        item = self.get_object()
        return Response({"data": self.detail_serializer_class(item).data})

    def create(self, request):
        return self.save_resource(request, status.HTTP_201_CREATED)

    def update(self, request, **kwargs):
        return self.save_resource(request, status.HTTP_200_OK)

    def save_resource(self, request, response_status):
        if request.content_type.split(";", 1)[0].strip().lower() != "application/json":
            raise UnsupportedMediaType(request.content_type)
        if request.stream is None:
            raise ParseError("Invalid JSON body")
        payload = request.data
        try:
            with transaction.atomic():
                instance = self.get_object() if self.action == "update" else None
                serializer = self.get_serializer(instance, data=payload)
                serializer.is_valid(raise_exception=True)
                serializer.save()
                data = serializer.data
        except IntegrityError as error:
            if "FOREIGN KEY" in str(error):
                raise InvalidReference() from error
            raise Conflict() from error
        return Response({"data": data}, status=response_status)

    def destroy(self, request, **kwargs):
        try:
            with transaction.atomic():
                self.get_object().delete()
        except (ProtectedError, IntegrityError) as error:
            raise Conflict() from error
        return Response(status=204)



class SchoolViewSet(ResourceViewSet):
    queryset = School.objects.order_by("id")
    filterset_class = SchoolFilter
    serializer_class = SchoolSerializer
    detail_serializer_class = SchoolDetailSerializer


class ClassroomViewSet(ResourceViewSet):
    queryset = Classroom.objects.order_by("id")
    filterset_class = ClassroomFilter
    serializer_class = ClassroomSerializer
    detail_serializer_class = ClassroomDetailSerializer

    def get_queryset(self):
        queryset = super().get_queryset()
        if self.action == "retrieve":
            return queryset.prefetch_related("teachers__classrooms", "students")
        return queryset


class TeacherViewSet(ResourceViewSet):
    queryset = Teacher.objects.order_by("id").prefetch_related("classrooms")
    filterset_class = TeacherFilter
    serializer_class = TeacherSerializer
    detail_serializer_class = TeacherDetailSerializer


class StudentViewSet(ResourceViewSet):
    queryset = Student.objects.order_by("id").select_related("classroom")
    filterset_class = StudentFilter
    serializer_class = StudentSerializer
    detail_serializer_class = StudentDetailSerializer


class HealthView(APIView):
    def get(self, request):
        try:
            with connection.cursor() as cursor:
                cursor.execute("SELECT 1")
        except DatabaseError:
            return Response(error_body(503, "Database unavailable"), status=503)
        return Response({"data": {"status": "ok"}})
