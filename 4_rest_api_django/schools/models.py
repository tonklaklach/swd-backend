from django.db import models


class Gender(models.TextChoices):
    MALE = "male", "Male"
    FEMALE = "female", "Female"
    OTHER = "other", "Other"


class School(models.Model):
    name = models.TextField()
    abbreviation = models.TextField()
    address = models.TextField()


class Classroom(models.Model):
    school = models.ForeignKey(School, on_delete=models.PROTECT, related_name="classrooms")
    year = models.PositiveIntegerField()
    section = models.PositiveIntegerField()

    class Meta:
        constraints = [
            models.UniqueConstraint(
                fields=["school", "year", "section"], name="unique_classroom_per_school"
            ),
            models.CheckConstraint(condition=models.Q(year__gt=0), name="positive_year"),
            models.CheckConstraint(condition=models.Q(section__gt=0), name="positive_section"),
        ]


class Teacher(models.Model):
    firstname = models.TextField()
    lastname = models.TextField()
    gender = models.CharField(max_length=6, choices=Gender)
    classrooms = models.ManyToManyField(
        Classroom, through="TeachingAssignment", related_name="teachers"
    )

    class Meta:
        constraints = [
            models.CheckConstraint(
                condition=models.Q(gender__in=Gender.values), name="valid_teacher_gender"
            )
        ]


class TeachingAssignment(models.Model):
    teacher = models.ForeignKey(Teacher, on_delete=models.CASCADE)
    classroom = models.ForeignKey(Classroom, on_delete=models.PROTECT)

    class Meta:
        constraints = [
            models.UniqueConstraint(
                fields=["teacher", "classroom"], name="unique_teaching_assignment"
            )
        ]


class Student(models.Model):
    firstname = models.TextField()
    lastname = models.TextField()
    gender = models.CharField(max_length=6, choices=Gender)
    classroom = models.ForeignKey(Classroom, on_delete=models.PROTECT, related_name="students")

    class Meta:
        constraints = [
            models.CheckConstraint(
                condition=models.Q(gender__in=Gender.values), name="valid_student_gender"
            )
        ]
