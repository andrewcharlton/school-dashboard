package database

// Query strings for prepared statements
var sqlStatements = map[string]string{

	"student": `SELECT upn, surname, forename, year, form,
				pp, eal, gender, ethnicity, sen_status,
				sen_need, sen_info, sen_strat, sen_access, sen_iep, 
				ks2_aps, ks2_band, ks2_en, ks2_ma, ks2_av, ks2_re, 
				ks2_wr, ks2_gps
				FROM students
				WHERE upn=? AND date_id=?`,

	"results": `SELECT subject_id, subject, grade, effort FROM results
				WHERE upn=? AND resultset=?`,

	"classes": `SELECT subject, class, teacher FROM classes
				WHERE upn=? AND date_id=?`,

	"inClass": `SELECT upn FROM classlist
				WHERE date_id=? AND subject_id=? AND class=?`,

	"subjects": `SELECT subjects.id as id
				FROM subjects
				INNER JOIN levels ON subjects.level_id = levels.id
				WHERE keystage=?`,

	"classlist": `SELECT DISTINCT class FROM classes
					WHERE date_id=? AND subject_id=?
					ORDER BY (
						CASE WHEN substr(class, 1, 1) == "1"
						THEN class
						ELSE  "0" || class
					END)`,

	"bestExams": `SELECT subject_id, subject, grade FROM results
				  WHERE upn=? AND is_exam=1
				  ORDER BY points DESC`,

	"firstExams": `SELECT subject_id, subject, grade FROM results
				  WHERE upn=? AND is_exam=1
				  ORDER BY date`,

	"search": `SELECT upn, surname, forename, year, form
			   FROM students
			   WHERE date_id=? AND
			   (((forename || " " || surname) LIKE ?)
				OR ((surname || ", " || forename) LIKE ?))
				ORDER BY (surname || " " || forename);`,

	"attendance": `SELECT poss_year, absence_year, unauth_year, mon_am,
					mon_pm, tue_am, tue_pm, wed_am, wed_pm, thu_am,
					thu_pm, fri_am, fri_pm
					FROM attendance
					WHERE upn=? AND year_id=?
					ORDER BY week_start DESC
					LIMIT 1`,
}
