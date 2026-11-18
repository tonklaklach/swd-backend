CREATE TABLE IF NOT EXISTS schools (
 id INTEGER PRIMARY KEY, name TEXT NOT NULL, abbreviation TEXT NOT NULL, address TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS classrooms (
 id INTEGER PRIMARY KEY, school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE RESTRICT,
 year INTEGER NOT NULL CHECK(year>0), section INTEGER NOT NULL CHECK(section>0),
 UNIQUE(school_id,year,section)
);
CREATE TABLE IF NOT EXISTS teachers (
 id INTEGER PRIMARY KEY, firstname TEXT NOT NULL, lastname TEXT NOT NULL,
 gender TEXT NOT NULL CHECK(gender IN ('male','female','other'))
);
CREATE TABLE IF NOT EXISTS students (
 id INTEGER PRIMARY KEY, firstname TEXT NOT NULL, lastname TEXT NOT NULL,
 gender TEXT NOT NULL CHECK(gender IN ('male','female','other')),
 classroom_id INTEGER NOT NULL REFERENCES classrooms(id) ON DELETE RESTRICT
);
CREATE TABLE IF NOT EXISTS teacher_classrooms (
 teacher_id INTEGER NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
 classroom_id INTEGER NOT NULL REFERENCES classrooms(id) ON DELETE RESTRICT,
 PRIMARY KEY(teacher_id,classroom_id)
);
CREATE INDEX IF NOT EXISTS ix_students_classroom ON students(classroom_id);
CREATE INDEX IF NOT EXISTS ix_teacher_classrooms_classroom ON teacher_classrooms(classroom_id);
