from django.urls import include, path
from rest_framework.routers import SimpleRouter

from schools.views import ClassroomViewSet, HealthView, SchoolViewSet, StudentViewSet, TeacherViewSet

router = SimpleRouter(trailing_slash=False)
router.register("schools", SchoolViewSet)
router.register("classrooms", ClassroomViewSet)
router.register("teachers", TeacherViewSet)
router.register("students", StudentViewSet)

urlpatterns = [
    path("healthz", HealthView.as_view()),
    path("api/v1/", include(router.urls)),
]

handler404 = "schools.errors.not_found"
