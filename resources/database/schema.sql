BEGIN TRANSACTION;
CREATE TABLE `subjects` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`subject`	TEXT,
	`code`	TEXT,
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
	FOREIGN KEY(`date_id`) REFERENCES dates(id)
);
CREATE TABLE "students_permanent" (
	`upn`	TEXT NOT NULL UNIQUE,
	`surname`	TEXT NOT NULL,
	`forename`	TEXT NOT NULL,
	`eal`	INTEGER,
	`gender`	TEXT NOT NULL,
	`ethnicity`	TEXT,
	`sen_status`	TEXT,
	`sen_need`	TEXT,
	`sen_info`	TEXT,
	`sen_strat`	TEXT,
	`sen_access`	TEXT,
	`sen_iep`	INTEGER,
	`ks2_aps`	REAL,
	`ks2_band`	TEXT,
	`ks2_en`	TEXT,
	`ks2_ma`	TEXT,
	`ks2_av`	TEXT,
	`ks2_re`	TEXT,
	`ks2_wr`	INTEGER,
	`ks2_gps`	TEXT,
	PRIMARY KEY(upn)
);
CREATE TABLE "student_results" (
	`upn`	TEXT,
	`resultset_id`	INTEGER,
	`subject_id`	INTEGER,
	`grade_id`	INTEGER,
	`effort`	INTEGER,
	PRIMARY KEY(upn,resultset_id,subject_id,grade_id),
	FOREIGN KEY(`upn`) REFERENCES students_permanent ( upn ),
	FOREIGN KEY(`resultset_id`) REFERENCES resultsets ( id ),
	FOREIGN KEY(`subject_id`) REFERENCES subjects ( id ),
	FOREIGN KEY(`grade_id`) REFERENCES grades ( id )
);
CREATE TABLE "student_classes" (
	`upn`	TEXT,
	`date_id`	INTEGER,
	`subject_id`	INTEGER,
	`class`	TEXT,
	`teacher`	TEXT,
	PRIMARY KEY(upn,date_id,subject_id),
	FOREIGN KEY(`upn`) REFERENCES students ( upn ),
	FOREIGN KEY(`date_id`) REFERENCES dates ( id ),
	FOREIGN KEY(`subject_id`) REFERENCES subjects ( id )
);
CREATE TABLE `school_years` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`year`	TEXT
);
CREATE TABLE "resultsets" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`resultset`	TEXT UNIQUE,
	`date`	TEXT,
	`is_exam`	INTEGER,
	`list`	INTEGER
);
CREATE TABLE "news" (
	`date`	TEXT,
	`comment`	TEXT
);
CREATE TABLE `nat_years` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`year`	TEXT NOT NULL UNIQUE
);
CREATE TABLE "nat_tms" (
	`subject_id`	INTEGER,
	`ks2`	TEXT,
	`grade`	TEXT,
	`probability`	REAL,
	PRIMARY KEY(subject_id,ks2,grade),
	FOREIGN KEY(`subject_id`) REFERENCES nat_tm_subjects ( id )
);
CREATE TABLE "nat_tm_subjects" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`year_id`	INTEGER,
	`subject`	TEXT,
	`level_id`	INTEGER,
	FOREIGN KEY(`year_id`) REFERENCES nat_years ( id ),
	FOREIGN KEY(`level_id`) REFERENCES levels(id)
);
CREATE TABLE "nat_progress8" (
	`year_id`	INTEGER,
	`ks2`	TEXT,
	`att8`	REAL,
	`english`	REAL,
	`maths`	REAL,
	`ebacc`	REAL,
	`other`	REAL,
	PRIMARY KEY(year_id,ks2),
	FOREIGN KEY(`year_id`) REFERENCES nat_years ( id )
);
CREATE TABLE `levels` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`level`	TEXT,
	`is_gcse`	INTEGER,
	`keystage`	INTEGER);
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
CREATE TABLE `dates` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`date`	TEXT UNIQUE,
	`list`	INTEGER,
	`year_id`	INTEGER,
	FOREIGN KEY(`year_id`) REFERENCES school_years(id)
);
CREATE TABLE `config` (
	`key`	TEXT UNIQUE,
	`value`	TEXT,
	PRIMARY KEY(key)
);
CREATE TABLE "attendance_weeks" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`year_id`	INTEGER,
	`week_start`	TEXT UNIQUE,
	FOREIGN KEY(`year_id`) REFERENCES school_years ( id )
);
CREATE TABLE "attendance_weekly" (
	`upn`	TEXT,
	`week_id`	INTEGER,
	`att_week`	REAL,
	`poss_year`	INTEGER,
	`absence_year`	INTEGER,
	`unauth_year`	INTEGER,
	PRIMARY KEY(upn,week_id),
	FOREIGN KEY(`upn`) REFERENCES students_permanent ( upn ),
	FOREIGN KEY(`week_id`) REFERENCES attendance_weeks ( id )
);
CREATE TABLE "attendance_sessions" (
	`upn`	TEXT,
	`year_id`	INTEGER,
	`mon_am`	INTEGER,
	`mon_pm`	INTEGER,
	`tue_am`	INTEGER,
	`tue_pm`	INTEGER,
	`wed_am`	INTEGER,
	`wed_pm`	INTEGER,
	`thu_am`	INTEGER,
	`thu_pm`	INTEGER,
	`fri_am`	INTEGER,
	`fri_pm`	INTEGER,
	PRIMARY KEY(upn,year_id),
	FOREIGN KEY(`upn`) REFERENCES students_permanent ( upn ),
	FOREIGN KEY(`year_id`) REFERENCES school_years ( id )
);
CREATE VIEW tms AS
SELECT year_id, level_id, subject, ks2, grade, probability
FROM nat_tms
INNER JOIN nat_tm_subjects ON subject_id = nat_tm_subjects.id;
CREATE VIEW students AS
SELECT students_permanent.upn as upn, date_id, surname, forename, year, form, pp, eal, gender, ethnicity,
sen_status, sen_need, sen_info, sen_strat, sen_access, sen_iep, ks2_aps, ks2_band, ks2_en, ks2_ma, ks2_av, ks2_re, ks2_wr, ks2_gps
FROM students_permanent
INNER JOIN students_temporal ON students_permanent.upn = students_temporal.upn
INNER JOIN dates ON students_temporal.date_id = dates.id
WHERE dates.list = 1;
CREATE VIEW results AS
SELECT upn, subject_id, subject, resultset_id as resultset, date, is_exam, list, grade, effort, points
FROM student_results
INNER JOIN grades ON grade_id = grades.id
INNER JOIN subjects ON student_results.subject_id = subjects.id
INNER JOIN resultsets ON student_results.resultset_id = resultsets.id;
CREATE VIEW classlist AS
SELECT date_id, subject_id, class, student_classes.upn
FROM student_classes
INNER JOIN students_permanent ON student_classes.upn = students_permanent.upn
ORDER BY (surname ||  " " || forename);
CREATE VIEW classes_filter AS
SELECT students_temporal.date_id as date_id, subject_id, class, students_permanent.upn as upn, surname, forename, year, pp, eal, gender, ethnicity,
sen_status, ks2_band
FROM students_permanent
INNER JOIN students_temporal ON students_permanent.upn = students_temporal.upn
INNER JOIN student_classes ON (student_classes.upn = students_permanent.upn AND student_classes.date_id = students_temporal.date_id)
INNER JOIN dates ON students_temporal.date_id = dates.id
WHERE dates.list = 1;
CREATE VIEW classes AS
SELECT upn, date_id, subject_id, subject, class, teacher
FROM student_classes
INNER JOIN subjects ON subject_id = subjects.id;
CREATE VIEW attendance AS
SELECT attendance_weekly.upn as upn, attendance_weeks.year_id, week_start, att_week, poss_year, absence_year, unauth_year, mon_am, mon_pm, tue_am, tue_pm, wed_am, wed_pm, thu_am, thu_pm, fri_am, fri_pm
FROM attendance_weekly
INNER JOIN attendance_weeks ON attendance_weekly.week_id = attendance_weeks.id
INNER JOIN attendance_sessions ON (attendance_weekly.upn = attendance_sessions.upn AND attendance_sessions.year_id = attendance_weeks.year_id)
INNER JOIN school_years ON attendance_weeks.year_id = school_years.id;
COMMIT;
