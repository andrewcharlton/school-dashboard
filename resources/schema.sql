BEGIN TRANSACTION;
CREATE TABLE `subjects` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`subject`	TEXT,
	`level_id`	INTEGER,
	`ebacc`	TEXT,
	`tm`	INTEGER,
	`ks2_prior`	TEXT,
	FOREIGN KEY(`level_id`) REFERENCES levels(id)
);
CREATE TABLE "students_temporal" (
	`upn`	TEXT NOT NULL,
	`date_id`	INTEGER NOT NULL,
	`year`	INTEGER,
	`form`	TEXT,
	`pp`	INTEGER,
	PRIMARY KEY(upn,date_id),
	FOREIGN KEY(`upn`) REFERENCES students_permanent ( upn ),
	FOREIGN KEY(`date_id`) REFERENCES dates ( id )
);
CREATE TABLE `students_permanent` (
	`upn`	TEXT NOT NULL UNIQUE,
	`surname`	TEXT NOT NULL,
	`forename`	TEXT NOT NULL,
	`eal`	INTEGER,
	`gender`	TEXT NOT NULL,
	`ethnicity`	TEXT,
	`sen_status`	TEXT,
	`sen_info`	TEXT,
	`sen_strat`	TEXT,
	`ks2_aps`	REAL,
	`ks2_band`	TEXT,
	`ks2_en`	TEXT,
	`ks2_ma`	TEXT,
	`ks2_av`	TEXT,
	PRIMARY KEY(upn)
);
CREATE TABLE `student_results` (
	`upn`	TEXT,
	`subject_id`	INTEGER,
	`resultset_id`	INTEGER,
	`grade_id`	INTEGER,
	PRIMARY KEY(upn,subject_id,resultset_id,grade_id),
	FOREIGN KEY(`upn`) REFERENCES students_permanent(upn),
	FOREIGN KEY(`subject_id`) REFERENCES subjects(id),
	FOREIGN KEY(`resultset_id`) REFERENCES resultsets(id),
	FOREIGN KEY(`grade_id`) REFERENCES grades(id)
);
CREATE TABLE `student_classes` (
	`upn`	TEXT,
	`class`	INTEGER,
	PRIMARY KEY(upn,class),
	FOREIGN KEY(`upn`) REFERENCES students(upn),
	FOREIGN KEY(`class`) REFERENCES class_teachers(id)
);
CREATE TABLE "resultsets" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`resultset`	TEXT,
	`date`	TEXT,
	`is_exam`	INTEGER,
	`list`	INTEGER
);
CREATE TABLE `levels` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`level`	TEXT,
	`is_gcse`	INTEGER
);
CREATE TABLE `grades` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`level_id`	INTEGER,
	`grade`	TEXT,
	`points`	INTEGER,
	`att8`	REAL,
	`l1_pass`	INTEGER,
	`l2_pass`	INTEGER,
	FOREIGN KEY(`level_id`) REFERENCES levels(id)
);
CREATE TABLE "dates" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`date`	TEXT UNIQUE,
	`list`	INTEGER
);
CREATE TABLE `config` (
	`key`	TEXT UNIQUE,
	`value`	TEXT,
	PRIMARY KEY(key)
);
CREATE TABLE "class_teachers" (
	`id`	INTEGER,
	`date_id`	INTEGER,
	`subject_id`	INTEGER,
	`class`	TEXT,
	`teacher`	TEXT,
	PRIMARY KEY(id),
	FOREIGN KEY(`date_id`) REFERENCES dates ( id ),
	FOREIGN KEY(`subject_id`) REFERENCES subjects(id)
);
CREATE VIEW students AS
SELECT students_permanent.upn as upn, date_id, surname, forename, year, form, pp, eal, gender, ethnicity,
sen_status, sen_info, sen_strat, ks2_aps, ks2_band, ks2_en, ks2_ma, ks2_av
FROM students_permanent
INNER JOIN students_temporal ON students_permanent.upn = students_temporal.upn
INNER JOIN dates ON students_temporal.date_id = dates.id
WHERE dates.list = 1;
CREATE VIEW results AS
SELECT upn, subject_id, subject, resultset_id as resultset, date, is_exam, grade, points
FROM student_results
INNER JOIN grades ON grade_id = grades.id
INNER JOIN subjects ON student_results.subject_id = subjects.id
INNER JOIN resultsets ON student_results.resultset_id = resultsets.id;
CREATE VIEW classes AS
SELECT upn, date_id as date, subject, class_teachers.class as class, teacher
FROM student_classes
INNER JOIN class_teachers ON student_classes.class = class_teachers.id
INNER JOIN subjects ON class_teachers.subject_id = subjects.id;
COMMIT;
