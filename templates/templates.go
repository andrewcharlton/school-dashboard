package templates

var allTemplates = map[string]string{

	"attainmentgroups.tmpl": `
<h2>Attainment Group Summary</h2>
<br>

{{ $q := .Query }}
<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">

	<table class="table table-condensed table-striped table-hover sortable">
	  <thead>
		<th>Group</th>
		<th style="text-align:center;">Cohort</th>
		<th style="text-align:center;">KS2 APS</th>
		<th style="text-align:center;">English</th>
		<th style="text-align:center;">Maths</th>
		<th style="text-align:center;">Science</th>
		<th style="text-align:center;">Humanities</th>
		<th style="text-align:center;">Language</th>
		<th style="text-align:center;">Basics</th>
		<th style="text-align:center;">EBacc</th>
	  </thead>
	  <tbody>
		{{ range .Groups }}
		  <tr>
			<td>{{ .Name }}</td>
			<td style="text-align:center;">{{ .Group.Cohort }}</td>
			<td style="text-align:center;">{{ printf "%.1f" .Group.KS2APS }}</td>
			<td style="text-align:center;">{{ Percent (.Group.EBaccArea "E").PctCohort 1 }}</td>
			<td style="text-align:center;">{{ Percent (.Group.EBaccArea "M").PctCohort 1 }}</td>
			<td style="text-align:center;">{{ Percent (.Group.EBaccArea "S").PctCohort 1 }}</td>
			<td style="text-align:center;">{{ Percent (.Group.EBaccArea "H").PctCohort 1 }}</td>
			<td style="text-align:center;">{{ Percent (.Group.EBaccArea "L").PctCohort 1 }}</td>
			<td style="text-align:center;">{{ Percent .Group.Basics.Percent 1 }}</td>
			<td style="text-align:center;">{{ Percent .Group.EBacc.PctCohort 1 }}</td>
		  </tr>
		{{ end }}
	  </tbody>
	</table>
  </div>
  <div class="col-sm-1"></div>
</div>



`,

	"attendance.tmpl": `
{{ define "AttendanceHeader" }}
<th style="text-align:center;">Cohort</th>
<th style="text-align:center;">This Week</th>
<th style="text-align:center;">Possible</th>
<th style="text-align:center;">Attendance %</th>
<th style="text-align:center;">Authorised %</th>
<th style="text-align:center;">Unauthorised %</th>
<th style="text-align:center;">% PA</th>
{{ end }}

{{ define "AttendanceRow" }}
<td style="text-align:center;">{{ .Cohort }}</td>
<td style="text-align:center;">{{ Percent .Week 1}}</td>
<td style="text-align:center;">{{ .Possible }}</td>
<td style="text-align:center;">{{ Percent .PercentAttendance 1}}</td>
<td style="text-align:center;">{{ Percent .PercentAuthorised 1}}</td>
<td style="text-align:center;">{{ Percent .PercentUnauthorised 1}}</td>
<td style="text-align:center;">{{ Percent .PercentPA 1}}</td>
{{ end }}

{{ define "AttendanceSessionTable" }}
<table class="table table-condensed table-hover">
  <thead>
	<th>Session</th>
	<th style='text-align:center;vertical-align:middle'>Monday</th>
	<th style='text-align:center;vertical-align:middle'>Tuesday</th>
	<th style='text-align:center;vertical-align:middle'>Wednesday</th>
	<th style='text-align:center;vertical-align:middle'>Thursday</th>
	<th style='text-align:center;vertical-align:middle'>Friday</th>
  </thead>
  <tbody>
	<tr class="active">
	  <td>Morning</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 0}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 2}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 4}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 6}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 8}}</td>
	</tr>
	<tr class="active">
	  <td>Afternoon</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 1}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 3}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 5}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 7}}</td>
	  <td style='text-align:center;vertical-align:middle'>{{index .Sessions 9}}</td>
	</tr>
  </tbody>
</table>
{{ end }}

<h2>Attendance Explorer</h2>
<h4>Week Beginning {{ .Week }}</h4>
<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">

	{{ $q := .Query }}
	{{ with .Group }}
	<h3>Summary</h3>
	<table class="table table-condensed table-striped table-hover">
	  <thead>
		{{ template "AttendanceHeader" }}
	  </thead>
	  <tbody>
		<tr>
		  {{ template "AttendanceRow" .Attendance }}
		</tr>
	  </tbody>
	</table>
	<br>

	<h3>Absence Patterns</h3>
	{{ template "AttendanceSessionTable" .Attendance }}
	<br>

	<h3>Students</h3>
	<table class="table table-condensed table-striped table-hover">
	  <thead>
		<th>Name</th>
		<th style="text-align:center;">Year</th>
		<th style="text-align:center;">Gender</th>
		<th style="text-align:center;">PP</th>
		<th style="text-align:center;">KS2</th>
		<th style="text-align:center;">Possible</th>
		<th style="text-align:center;">Absences</th>
		<th style="text-align:center;">Unauthorised</th>
		<th style="text-align:center;">Attendance %</th>
		<th style="text-align:center;">Week %</th>
	  </thead>
	  <tbody>
	  {{ range .Students }}
		{{ with .Attendance }}
		{{ if ge .Latest 0.95 }}<tr class="success">
		{{ else if ge .Latest 0.90 }}<tr class="warning">
		{{ else if eq .Possible 0 }}<tr>
		{{ else }}<tr class="danger">
		{{ end }}
		{{ end }}
		  <td><a href="/student/{{ .UPN }}/?{{ $q }}">{{ .Name }}</a></td>
		  <td style="text-align:center;">{{ .Year }}</td>
		  <td style="text-align:center;">{{ .Gender }}</td>
		  <td style="text-align:center;">{{ template "PP" .PP }}</td>
		  <td style="text-align:center;">{{ .KS2.Av }}</td>
		  {{ with .Attendance }}
		  <td style="text-align:center;">{{ .Possible }}</td>
		  <td style="text-align:center;">{{ .Absences }}</td>
		  <td style="text-align:center;">{{ .Unauthorised }}</td>
		  <td style="text-align:center;">{{ Percent .Latest 1}}</td>
		  <td style="text-align:center;">{{ Percent .Week 1 }}</td>
		  {{ end }}
		</tr>
	  {{ end }}
	  </tbody>
	</table>

	{{ end }}
  </div>
  <div class="col-sm-1"></div>
</div>
`,

	"attendancegroups.tmpl": `
<h2>Attendance Group Summary</h2>
<h4>Week Beginning {{ .Week }}</h4>
<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">

	<h3>Summary</h3>
	<table class="table table-condensed table-striped table-hover">
	  <thead>
		<th>Group</th>
		{{ template "AttendanceHeader" }}
	  </thead>
	  <tbody>
		{{ range .YearGroups }}
		<tr>
		  <td><a href="#{{ .Name }}">{{ .Name }}</a></td>
		  {{ with index .Groups 0 }}
		  {{ template "AttendanceRow" .Group.Attendance }}
		  {{ end }}
		</tr>
		{{ end }}
	  </tbody>
	</table>
	<br>
	
	{{ $q := .Query }}
	{{ range .YearGroups }}
	<h3><a name="{{ .Name }}"></a>{{ .Name }}</h3>
	<table class="table table-condensed table-hover sortable">
	  <thead>
		<th>Group</th>
		{{ template "AttendanceHeader" }}
	  </thead>
	  <tbody>
		{{ $yq := .Query }}
		{{ range .Groups }} 
		{{ with .Group.Attendance }}
		{{ if ge .PercentAttendance 0.975 }}<tr class="success">
		{{ else if ge .PercentAttendance 0.95 }}<tr class="warning">
		{{ else if eq .Possible 0 }}<tr>
		{{ else }}<tr class="danger">
		{{ end }}
		{{ end }}
		<td><a href="/attendance/?{{ $q }}{{ $yq }}{{ .Query }}">{{ .Name }}</a></td>
		  {{ template "AttendanceRow" .Group.Attendance }}
		</tr>
		{{ end }}
	  </tbody>
	</table>
	{{ end }}
  </div>
  <div class="col-sm-1"></div>
</div>


`,

	"common.tmpl": `
{{ define "TickCross" }}
{{if .}}<span class="glyphicon glyphicon-ok" style="color: #009933;"></span>
{{else}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
{{end}}
{{ end }}

{{ define "PP" }}
  {{ if . }}<span class="glyphicon glyphicon-star" style="color:#ff9900;"></span>
  {{ end }}
{{ end }}

{{ define "StudentAttendance" }}
  {{ if gt . 0.95 }}<span class="text-success">{{ Percent . 1}}</span>
  {{ else if lt . 0.90 }}<span class="text-danger">{{ Percent . 1}}</span>
  {{ else }}<span class="text-warning">{{ Percent . 1}}</span>
  {{ end }}
{{ end }}

{{ define "StudentVA" }}
  {{ if gt . 0.67 }}<span class="text-success">{{ printf "%+.2f" . }}</span>
  {{ else if lt . -0.33 }}<span class="text-danger">{{ printf "%+.2f" . }}</span>
  {{ else }}<span class="text-warning">{{ printf "%+.2f" . }}</span>
  {{ end }}
{{ end }}

{{ define "StudentProgress8" }}
  {{ if gt . 0.2 }}<span class="text-success">{{ printf "%+.2f" . }}</span>
  {{ else if lt . -0.2 }}<span class="text-danger">{{ printf "%+.2f" . }}</span>
  {{ else }}<span class="text-warning">{{ printf "%+.2f" . }}</span>
  {{ end }}
{{ end }}

{{ define "KS3VA" }}
  {{ if gt . 0.33 }}<span class="text-success">{{ printf "%+.2f" . }}</span>
  {{ else if lt . -0.33 }}<span class="text-danger">{{ printf "%+.2f" . }}</span>
  {{ else }}<span class="text-warning">{{ printf "%+.2f" . }}</span>
  {{ end }}
{{ end }}
`,

	"ebacc.tmpl": `
{{ define "EBaccResult" }}
  {{ if .Entered }}
	{{ if gt (len .Results) 0 }}
	  <td style="text-align:center;"
		  data-container="body"
		  data-toggle="popover"
		  data-placement="right"
		  data-html="true"
		  data-trigger="hover"
		  data-content="{{ range .Results }} {{.}}<br> {{ end }}">
		{{ template "TickCross" .Achieved }}
	  </td>
	{{ else }}
	  <td style="text-align:center;">{{ template "TickCross" .Achieved }}</td>
	{{ end }}
  {{ else }}
	<td></td>
  {{ end }}
{{ end }}

{{ $g := .Group }}
{{ $q := .Query }}
{{ $areas := .Areas }}

<h2>English Baccalaureate</h2>
<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<h3>Summary</h3>
	<table class="table table-condensed table-striped table-hover">
	  <thead>
		<th>&nbsp;</th>
		<th style="text-align:center;">English</th>
		<th style="text-align:center;">Maths</th>
		<th style="text-align:center;">Science</th>
		<th style="text-align:center;">Humanities</th>
		<th style="text-align:center;">Language</th>
		<th style="text-align:center;">EBacc</th>
	  </thead>
	  <tbody>
		<tr>
		  <td>Entered</td>
		  {{ range $a := .Areas }}
			<td style="text-align:center;">{{ ($g.EBaccArea $a).Entered }}</td>
		  {{ end }}
		  <td style="text-align:center;">{{ $g.EBacc.Entered }}</td>
		</tr>
		<tr>
		  <td>Achieved</td>
		  {{ range $a := .Areas }}
			<td style="text-align:center;">{{ ($g.EBaccArea $a).Achieved }}</td>
		  {{ end }}
		  <td style="text-align:center;">{{ $g.EBacc.Achieved }}</td>
		</tr>
		<tr>
		  <td>% of Cohort Achieved</td>
		  {{ range $a := .Areas }}
			<td style="text-align:center;">{{ Percent ($g.EBaccArea $a).PctCohort 1 }}</td>
		  {{ end }}
		  <td style="text-align:center;">{{ Percent $g.EBacc.PctCohort 1 }}</td>
		</tr>
		<tr>
		  <td>% of Entries Achieved</td>
		  {{ range $a := .Areas }}
			<td style="text-align:center;">{{ Percent ($g.EBaccArea $a).PctEntries 1 }}</td>
		  {{ end }}
		  <td style="text-align:center;">{{ Percent $g.EBacc.PctEntries 1 }}</td>
		</tr>
		</body>
	</table>
	<br>

	<h3>Students</h3>
	<br>
	<p>A tick/cross indicates whether a student achieved a pass in that area. A blank indicates the
	student was not entered for that area.</p>
	<br>
	<ul class="nav nav-tabs">
	  {{ range $n, $name := .GroupHeaders }}
		<li {{if eq $n 0}}class="active"{{end}}><a href="#{{ $n }}" data-toggle="tab" aria-expanded="{{if eq $n 0}}true{{else}}false{{end}}">{{ $name }}</a></li>
	  {{ end }}
	</ul>

	{{ $headers := .GroupHeaders }}
	<div id="EBaccPanes"  class="tab-content">
	  {{ range $n, $g := .SubGroups }}
		{{ if eq $n 0 }}
		  <div class="tab-pane fade active in" id="{{ $n }}">
		{{ else }}
		  <div class="tab-pane fade" id="{{ $n }}">
		{{ end }}
		  <br>
		  <table class="table table-condensed table-striped table-hover sortable">
			<thead>
			  <th>Name</th>
			  <th style="text-align:center;">KS2</th>
			  <th style="text-align:center;">Gender</th>
			  <th style="text-align:center;">PP</th>
			  <th style="text-align:center;">English</th>
			  <th style="text-align:center;">Maths</th>
			  <th style="text-align:center;">Science</th>
			  <th style="text-align:center;">Humanities</th>
			  <th style="text-align:center;">Language</th>
			  <th style="text-align:center;">EBacc</th>
			  <th style="text-align:center;">Attendance</th>
			</thead>
			<tbody>
			  {{ range $s := $g.Students }}
				<tr>
				  <td><a href="/student/{{ $s.UPN }}/?{{ $q }}">{{ $s.Name }}</a></td>
				  <td style="text-align:center;">{{ $s.KS2.Av }}</td>
				  <td style="text-align:center;">{{ $s.Gender }}</td>
				  <td style="text-align:center;">{{ template "PP" $s.PP }}</td>
				  {{ range $a := $areas }}
					{{ template "EBaccResult" ($s.EBaccArea $a) }}
				  {{ end }}
				    {{ template "EBaccResult" $s.EBacc}}
				  <td style="text-align:center;">{{ template "StudentAttendance" $s.Attendance.Latest }}</td>
				</tr>
			  {{ end }}
			</tbody>
		  </table>
		</div>
	  {{ end }}
	</div>
  </div>
  <div class="col-sm-1"></div>
</div>

<script>
$(function () {
  $('[data-toggle="popover"]').popover()
})
</script>
`,

	"em.tmpl": `
<h2>English and Maths</h2>
<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<h3>Headlines</h3>

	<table class="table table-hover">
	  <thead>
		<th>&nbsp;</th>
		{{ range .Names }}
		  <th style="text-align:center;">{{ . }}</th>
		{{ end }}
	  </thead>
	  <tbody>
		<tr>
		  <td>Students</td>
		  {{ range .Groups }}
			<td style="text-align:center;">{{ len .Students }}</td>
		  {{ end }}
		</tr>
		<tr>
		  <td>Percentage</td>
		  {{ range .Pcts }}
			<td style="text-align:center;">{{ Percent . 1 }}</td>
		  {{ end }}
		</tr>
	  </tbody>
	</table>

	<br>

	<h3>Students</h3>

	<ul class="nav nav-tabs">
	  {{ range $n, $name := .Names }}
		<li {{if eq $n 0}}class="active"{{end}}><a href="#{{ $n }}" data-toggle="tab" aria-expanded="{{if eq $n 0}}true{{else}}false{{end}}">{{ $name }}</a></li>
	  {{ end }}
	</ul>

	{{ $names := .Names }}
	{{ $q := .Query }}
	<div id="EMGroupPanes" class="tab-content">
	  {{ range $n, $g := .Groups }}
		{{ if eq $n 0 }}
		  <div class="tab-pane fade active in" id="{{ $n }}">
		{{ else }}
		  <div class="tab-pane fade" id="{{ $n }}">
		{{ end }}
		<br>
		<table class="table table-condensed table-striped table-hover sortable">
		  <thead>
			<th>Name</th>
			<th style="text-align:center;">KS2</th>
			<th style="text-align:center;">Gender</th>
			<th style="text-align:center;">PP</th>
			<th style="text-align:center;">Language</th>
			<th style="text-align:center;">Literature</th>
			<th style="text-align:center;">Maths</th>
			<th style="text-align:center;">Progress 8</th>
			<th style="text-align:center;">Attendance</th>
		  </thead>
		  <tbody>
			{{ range $g.Students }}
			  <tr>
				<td><a href="/student/{{.UPN}}/?{{$q}}">{{.Name}}</a></td>
				<td style="text-align:center;">{{ .KS2.Av }}</td>
				<td style="text-align:center;">{{ .Gender }}</td>
				<td style="text-align:center;">{{ template "PP" .PP }}</td>
				<td style="text-align:center;">{{ .SubjectGrade "English" }}</td>
				<td style="text-align:center;">{{ .SubjectGrade "English Literature" }}</td>
				<td style="text-align:center;">{{ .SubjectGrade "Mathematics" }}</td>
				<td style="text-align:center;">{{ template "StudentProgress8" .Basket.Overall.Progress8 }}</td>
				<td style="text-align:center;">{{ template "StudentAttendance" .Attendance.Latest }}</td>
			  </tr>
			{{ end }}
		  </tbody>
		</table>
		</div>
	  {{ end }}
	</div>

  </div>
  <div class="col-sm-1"></div>
</div>



`,

	"filter.tmpl": `
<div class="row" id="filter_bar" style="display: block;">
  <div class="col-sm-10">
	{{range .Labels}}
	<a href="#" class="btn btn-{{ .Format }} btn-sm disabled" style="margin-bottom: 2px;">{{ .Text }}</a>
	{{end}}
  </div>
  <div class="col-sm-2">
	<button type="button" class="btn btn-primary btn-sm pull-right" onclick="toggle_visibility('filter_bar'); toggle_visibility('filter_form');">
	  <span class="glyphicon glyphicon-cog"> Options</span>
	</button>
  </div>
</div>

<div class="row" id="filter_form" style="display: none;">
  <div class="col-sm-2"></div>
  <div class="col-sm-8">
	<form role="form" class="form-horizontal" style="margin-bottom: 0px">
	  <div class="row">
		<div class="form-group form-group-sm">
		  <label for="f_natyear" class="control-label col-sm-2">National Data</label>
		  <div class="col-sm-10">
			<select name="natyear" id="f_resultset" class="form-control input-sm">
			  {{$rs := .NatYear}}
			  {{range .NatYears}}
			  <option value="{{.ID}}"{{if eq .ID $rs}} selected="yes"{{end}}>{{.Name}}</option>
			  {{end}}
			</select>
		  </div>
		</div>
	  </div>

	  <div class="row">
		<div class="form-group form-group-sm">
		  <label for="f_date" class="control-label col-sm-2">Effective Date</label>
		  <div class="col-sm-10">
			<select name="date" id="f_date" class="form-control input-sm">
			  {{$d := .Date}}
			  {{range .Dates}}
			  <option value="{{.ID}}"{{if eq .ID $d}} selected="yes"{{end}}>{{.Name}}</option>
			  {{end}}
			</select>
		  </div>
		</div>
	  </div>

	  <div class="row">
		<div class="form-group form-group-sm">
		  <label for="f_resultset" class="control-label col-sm-2">Resultset</label>
		  <div class="col-sm-10">
			<select name="resultset" id="f_resultset" class="form-control input-sm">
			  {{$rs := .Resultset}}
			  {{range .Resultsets}}
			  <option value="{{.ID}}"{{if eq .ID $rs}} selected="yes"{{end}}>{{.Name}}</option>
			  {{end}}
			</select>
		  </div>
		</div>
	  </div>

	  {{ if ge .Detail 1 }}
	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <label for="f_yeargroup" class="control-label col-sm-2">Yeargroup</label>
		  <div class="col-sm-10">
			<select name="year" id="f_resultset" class="form-control input-sm">
			  <option value=""{{if eq .Year ""}} selected="yes"{{end}}>All Students</option>
			  <option value="7"{{if eq .Year "7"}} selected="yes"{{end}}>7</option>
			  <option value="8"{{if eq .Year "8"}} selected="yes"{{end}}>8</option>
			  <option value="9"{{if eq .Year "9"}} selected="yes"{{end}}>9</option>
			  <option value="10"{{if eq .Year "10"}} selected="yes"{{end}}>10</option>
			  <option value="11"{{if eq .Year "11"}} selected="yes"{{end}}>11</option>
			</select>
		  </div>
		</div>
	  </div>
	  {{ end }}

	  {{ if ge .Detail 2 }}
	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <label for="f_gender" class="control-label col-sm-2">Gender</label>
		  <div class="col-sm-10">
			<select name="gender" id="f_gender" class="form-control input-sm">
			  <option value=""{{if eq .Gender ""}} selected="yes"{{end}}>All Students</option>
			  <option value="1"{{if eq .Gender "1"}} selected="yes"{{end}}>Male</option>
			  <option value="0"{{if eq .Gender "0"}} selected="yes"{{end}}>Female</option>
			</select>
		  </div>
		</div>
	  </div>

	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <label for="f_pp" class="control-label col-sm-2">Pupil Premium</label>
		  <div class="col-sm-10">
			<select name="pp" id="f_pp" class="form-control input-sm">
			  <option value=""{{if eq .PP ""}} selected="yes"{{end}}>All Students</option>
			  <option value="1"{{if eq .PP "1"}} selected="yes"{{end}}>Disadvantaged Students</option>
			  <option value="0"{{if eq .PP "0"}} selected="yes"{{end}}>Non-Disadvantaged Students</option>
			</select>
		  </div>
		</div>
	  </div>

	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <label for="f_eal" class="control-label col-sm-2">EAL</label>
		  <div class="col-sm-10">
			<select name="eal" id="f_eal" class="form-control input-sm">
			  <option value=""{{if eq .EAL ""}} selected="yes"{{end}}>All Students</option>
			  <option value="1"{{if eq .EAL "1"}} selected="yes"{{end}}>English as an Additional Language</option>
			  <option value="0"{{if eq .EAL "0"}} selected="yes"{{end}}>English as a First Language</option>
			</select>
		  </div>
		</div>
	  </div>

	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <label class="control-label col-sm-2">KS2 Band</label>
		  <div class="col-sm-10">
			<label class="checkbox-inline"><input type="checkbox" name="ks2band" value="High"{{if index .B "High"}} checked="yes"{{end}}>High</input></label>
			<label class="checkbox-inline"><input type="checkbox" name="ks2band" value="Middle"{{if index .B "Middle"}} checked="yes"{{end}}>Middle</input></label>
			<label class="checkbox-inline"><input type="checkbox" name="ks2band" value="Low"{{if index .B "Low"}} checked="yes"{{end}}>Low</input></label>
			<label class="checkbox-inline"><input type="checkbox" name="ks2band" value="None"{{if index .B "None"}} checked="yes"{{end}}>None</input></label>
		  </div>
		</div>
	  </div>

	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <label class="control-label col-sm-2">SEN</label>
		  <div class="col-sm-10">
			<label class="checkbox-inline"><input type="checkbox" name="sen" value=""{{if index .S ""}} checked="yes"{{end}}>No SEN</input></label>
			<label class="checkbox-inline"><input type="checkbox" name="sen" value="A"{{if index .S "A"}} checked="yes"{{end}}>School Action</input></label>
			<label class="checkbox-inline"><input type="checkbox" name="sen" value="P"{{if index .S "P"}} checked="yes"{{end}}>School Action Plus</input></label>
			<label class="checkbox-inline"><input type="checkbox" name="sen" value="S"{{if index .S "S"}} checked="yes"{{end}}>SEN with a Statement</input></label>
		  </div>
		</div>
	  </div>

	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <label class="control-label col-sm-2">Ethnicities</label>
		  <div class="col-sm-10">
			{{ $E := .E }}
			{{ $O := .O }}
			{{range .Ethnicities}}
			{{if not (index $O .)}}
			<label class="checkbox-inline"><input type="checkbox" name="ethnicity" value="{{.}}"{{if index $E .}} checked="yes"{{end}}>{{.}}</input></label>
			{{end}}
			{{end}}
		  </div>
		</div>
	  </div>

	  {{end}}
	  <div class="row" style="margin-top: -6px;">
		<div class="form-group form-group-sm">
		  <div class="col-sm-10"></div>
		  <div class="col-sm-2">
			<input type="Submit" class="btn btn-primary btn-sm btn-block" value="Apply Filter">
		  </div>
		</div>
	  </div>
	</form>

	<form role="form" class="form-horizontal" style="margin-top: -12px;">
	  <div class="row">
		<div class="form-group form-group-sm">
		  <div class="col-sm-10"></div>
		  <div class="col-sm-2">
			<input type="Submit" class="btn btn-default btn-sm btn-block" value="Reset Filter">
		  </div>
		</div>
	  </div>
	</form>

  </div>
  <div class="col-sm-1"></div>
  <div class="col-sm-1">
	<button type="button" class="btn btn-primary btn-sm" onclick="toggle_visibility('filter_bar'); toggle_visibility('filter_form');">
	  <span class="glyphicon glyphicon-cog"> Options</span>
	</button>
  </div>
</div>

<hr>
`,

	"footer.tmpl": `
</div>
</body>
</html>
`,

	"header.tmpl": `
<!DOCTYPE html>
<html lang="en">
  <head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
	<title>Venerable Bede - Data Dashboard</title>

	<!-- Scripts -->
	<script src="/static/js/jquery.min.js"></script>
	<script src="/static/js/bootstrap.min.js"></script>
	<script src="/static/js/custom.js"></script>
	<script src="/static/js/sorttable.js"></script>

	<!-- Bootstrap -->
	<link href="/static/css/bootstrap-spacelab.min.css" rel="stylesheet">

	<!-- CSS Overrides -->
	<link href="/static/css/custom.css" rel="stylesheet">

  </head>
  <body>

	<nav class="navbar navbar-default">
	  <div class="container">
		<div class="navbar-header">
		  <a class="navbar-brand" href="/index/?{{.Query}}">{{.School}}</a>
		</div>

		<form class="navbar-form navbar-right form-horizontal" action="/search/" role="search">
		  {{with .F}}
		  <input type="text" name="natyear" value="{{.NatYear}}" style="display: none;">
		  <input type="text" name="date" value="{{.Date}}" style="display: none;">
		  <input type="text" name="resultset" value="{{.Resultset}}" style="display: none;">
		  {{end}}
		  <div class="form-group">
			<input type="text" name="name" class="form-control input-sm" style="position: relative; top: 3px" placeholder="Search students">
		  </div>
		  <input type="Submit" style="position: absolute; top: -1000px">
		</form>

		<ul class="nav navbar-nav navbar-right">
		  <li class="dropdown">
			<a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">Departments <span class="caret"></span></a>
			<ul class="dropdown-menu" role="menu">
			  <li><a href="/subjects/?{{.Query}}">Summary List</a></li>
			  <li class="divider"></li>
			  <li><a href="/progressgrid/?{{.Query}}">Progress Grid</a></li>
			  <li><a href="/subjectgroups/?{{.Query}}">Group Comparisons</a></li>
			</ul>
		  </li>
		</ul>

		<ul class="nav navbar-nav navbar-right">
		  <li class="dropdown">
			<a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">Whole School <span class="caret"></span></a>
			<ul class="dropdown-menu" role="menu">
			  <li><a href="/progress8/?{{.Query}}">Progress 8</a></li>
			  <li><a href="/progress8groups/?{{.Query}}">Progress 8 Groups</a></li>
			  <li class="divider"></li>
			  <li><a href="/basics/?{{.Query}}">English and Maths</a></li>
			  <li><a href="/ebacc/?{{.Query}}">English Baccalaureate</a></li>
			  <li><a href="/attainmentgroups/?{{.Query}}">Attainment Groups</a></li>
			  <li class="divider"></li>
			  <li><a href="/ks3summary/?{{.Query}}">KS3 Summary</a></li>
			  <li><a href="/ks3groups/?{{.Query}}">KS3 Groups</a></li>
			  <li class="divider"></li>
			  <li><a href="/attendancegroups/?{{.Query}}">Attendance Summary</a></li>
			  <li><a href="/attendance/?{{.Query}}">Attendance Explorer</a></li>
			</ul>
		  </li>
		</ul>
	  </div>
	</nav>

	<div class="container">
`,

	"ks3groups.tmpl": `
<h2>Key Stage 3 Group Comparison</h2>
<br>

{{ $q := .Query }}
{{ $y := .Year }}
{{ $subj := .Subjects }}

<h4>Group Details</h4>
<br>

<table class="table table-condensed table-striped table-hover sortable">
  <thead>
	<th>Group</th>
	<th style="text-align:center;">Cohort</th>
	<th style="text-align:center;">KS2 APS</th>
	{{ range .Subjects }}<th style="text-align:center;">{{ .Code }}</th>{{ end }}
	  <th style="text-align:center;">Overall</th>
	  <th style="text-align:center;">Attendance</th>
  </thead>
  <tbody>
	{{ range $g := .Groups }}
	  <tr>
		<td><a href="/ks3summary/?{{ $q }}{{ $g.Query }}">{{ $g.Name }}</a></td>
		<td style="text-align:center;">{{ $g.Group.Cohort }}</td>
		<td style="text-align:center;">{{ printf "%.1f" $g.Group.KS2APS }}</td>
		{{ range $subj }}
		  <td style="text-align:center;">
			{{ $va := ($g.Group.SubjectVA .Subj).VA }}
			{{ if gt $va 0.33 }}<a class="text-success"
			{{ else if lt $va -0.33 }}<a class="text-danger"
			{{ else }}<a class="text-warning"
			{{ end }}
				href="/progressgrid/{{ .Subj }}/{{ .SubjID }}/All Year {{ $y }}/?{{ $q }}{{ $g.Query }}">
			  {{ printf "%+.2f" $va }}
			</a>
		  </td>
		{{ end }}
		<td style="text-align:center;">{{ template "KS3VA" $g.Group.AverageVA }}</td>
		<td style="text-align:center;">{{ template "StudentAttendance" $g.Group.Attendance.PercentAttendance }}</td>
	  </tr>
	{{ end }}
  </tbody>
</table>

<br>
<h4>Group Matrix</h4>
<br>

{{ with .Matrix }}
  <table class="table table-condensed table-hover" style="table-layout:fixed;">
	<thead>
	  <th></th>
	  {{ range .Headers }}
		<th style="text-align:center;">{{ . }}</th>
	  {{ end }}
	</thead>
	<tbody>
	  {{ $headers := .Headers }}
	  {{ range $i, $g := .Groups }}
		<tr>
		  <th>{{ index $headers $i }}</th>
		  {{ range $g }}
			{{ $va := .Group.AverageVA }}
			<td style="text-align:center;">
			  {{ if eq (len .Group.Students) 0 }}
				-
			  {{ else if gt $va 0.2 }}
				<a class="text-success" href="/ks3summary/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $va }}</a>
			  {{ else if lt $va -0.2 }}
				<a class="text-danger"  href="/ks3summary/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $va }}</a>
			  {{ else }}
				<a class="text-warning" href="/ks3summary/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $va }}</a>
			  {{ end }}
			</td>
		  {{ end }}
		</tr>
	  {{ end }}
	</tbody>
  </table>
{{ end }}
`,

	"ks3summary.tmpl": `
<h2>Key Stage 3 Summary</h2>
<br>

{{ $q := .Query }}
{{ $subj := .Subjects }}
{{ $ks3 := .KS3 }}
{{ $g := .Group }}

<h4>Cohort Overview</h4>
<br>
<table class="table table-striped table-condensed table-hover">
  <thead>
	<th></th>
	{{ range .Subjects }}<th style="text-align:center;">{{ .Code }}</th>{{ end }}
	<th>Average</th>
  </thead>
  <tbody>
	<tr>
	  <td>Average Points</td>
	  {{ range .Subjects }}
		<td style="text-align:center;">{{ printf "%.1f" ($g.SubjectPoints .Subj) }}</td>
	  {{ end }}
	  <td style="text-align:center;">{{ printf "%.1f" $g.APS }}</td>
	</tr>
	<tr>
	  <td>Average Grade</td>
	  {{ range .Subjects }}
		{{ $pts := ($g.SubjectPoints .Subj) }}
		<td style="text-align:center;">{{ $ks3.GradeEquivalent $pts }}</td>
	  {{ end }}
	  <td style="text-align:center;">{{ $ks3.GradeEquivalent $g.APS }}</td>
	</tr>
	<tr>
	  <td>Value Added</td>
	  {{ range .Subjects }}
		<td style="text-align:center;">{{ printf "%.1f" ($g.SubjectVA .Subj).VA }}</td>
	  {{ end }}
	  <td style="text-align:center;">{{ printf "%.1f" $g.AverageVA }}</td>
	</tr>
	<tr>
	  <td>Average Effort</td>
	  {{ range .Subjects }}
		<td style="text-align:center;">{{ printf "%.1f" ($g.SubjectEffort .Subj) }}</td>
	  {{ end }}
	  <td style="text-align:center;">{{ printf "%.1f" $g.AverageEffort }}</td>
	</tr>
  </tbody>
</table>

<br>
<h4>Students</h4>
<br>

<table class="table table-striped table-condensed table-hover sortable">
  <thead>
	<th>Name</th>
	<th style="text-align:center;">Gender</th>
	<th style="text-align:center;">PP</th>
	<th style="text-align:center;">KS2</th>
	{{ range .Subjects }}<th style="text-align:center;">{{ .Code }}</th>{{ end }}
	  <th style="text-align:center;">Current APS</th>
	  <th style="text-align:center;">Average Grade</th>
	  <th style="text-align:center;">Effort</th>
	  <th style="text-align:center;">Attendance</th>
  </thead>
  <tbody>
	{{ range .Group.Students }}
	  <tr>
		<td><a href="/student/{{ .UPN }}/?{{ $q }}">{{ .Name }}</a></td>
		<td style="text-align:center;">{{ .Gender }}</td>
		<td style="text-align:center;">{{ template "PP" .PP }}</td>
		<td style="text-align:center;">{{ .KS2.Av }}</td>
		{{ $s := . }}
		{{ range $subj }}
		  <td style="text-align:center;">{{ $s.SubjectGrade .Subj }}</td>
		{{ end }}
		<td style="text-align:center;">{{ printf "%.1f" .APS }}</td>
		<td style="text-align:center;">{{ $ks3.GradeEquivalent .APS }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .AverageEffort.Effort }}</td>
		<td style="text-align:center;">{{ template "StudentAttendance" .Attendance.Latest }}</td>
	  </tr>
	{{ end }}
  </tbody>
</table>






`,

	"progress8.tmpl": `
{{ define "P8Block" }}
  <td style="text-align:center;"
	data-container="body"
	data-toggle="popover"
	data-placement="right"
	data-html="true"
	data-trigger="hover"
	data-content="{{ range .Subjects }} {{.}}<br> {{ end }}">
	{{printf "%+.2f" .Progress8 }}
  </td>
{{ end }}

<h2>Progress 8</h2>
<br>

<div width="100%" id="chart"></div>

<br>


<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">

	{{ with .Group.Progress8 }}
	<table class="table">
	  <thead>
		<th width="25%">&nbsp;</th>
		<th width="15%" style="text-align:center;">English</th>
		<th width="15%" style="text-align:center;">Maths</th>
		<th width="15%" style="text-align:center;">EBacc</th>
		<th width="15%" style="text-align:center;">Other</th>
		<th width="15%" style="text-align:center;">Overall</th>
	  </thead>
	  <tbody>
		<tr>
		  <td>Attainment 8</td>
		  {{ range .Attainment }}
		  <td style="text-align:center;">{{printf "%.2f" . }}</td>
		  {{ end }}
		</tr>
		<tr>
		  <td>Points per Slot</td>
		  {{ range .AttainmentPerSlot }}
		  <td style="text-align:center;">{{printf "%.2f" . }}</td>
		  {{ end }}
		</tr>
		<tr>
		  <td>Entries</td>
		  {{ range .Entries }}
		  <td style="text-align:center;">{{printf "%.2f" . }}</td>
		  {{ end }}
		</tr>
		<tr>
		  <td>Progress 8</td>
		  {{ range .Progress }}
		  <td style="text-align:center;">{{printf "%+.2f" . }}</td>
		  {{ end }}
		</tr>
	  </tbody>
	</table>
	{{ end }}

	<br>

	<table class="table table-condensed table-striped table-hover sortable">
	  <thead>
		<th>Name</th>
		<th style="text-align:center;">KS2</th>
		<th style="text-align:center;">Gender</th>
		<th style="text-align:center;">PP</th>
		<th style="text-align:center;">Entries</th>
		<th style="text-align:center;">Attainment 8</th>
		<th style="text-align:center;">English</th>
		<th style="text-align:center;">Maths</th>
		<th style="text-align:center;">EBacc</th>
		<th style="text-align:center;">Other</th>
		<th style="text-align:center;">Overall</th>
		<th style="text-align:center;">Attendance</th>
	  </thead>
	  <tbody>
		{{ $q := .Query }}
		{{ range .Group.Students }}
		<tr>
		  <td><a href="/student/{{ .UPN }}/?{{ $q }}">{{ .Name }}</a></td>
		  <td style="text-align:center;">{{ .KS2.Av }}</td>
		  <td style="text-align:center;">{{ .Gender }}</td>
		  <td style="text-align:center;">{{ template "PP" .PP }}</td>

		  {{ with .Basket }}
		  {{ with .Overall }}
		  <td style="text-align:center;">{{printf "%v"  .Entries }}</td>
		  <td style="text-align:center;">{{printf "%.1f"  .Attainment }}</td>
		  {{ end }}
		  {{ template "P8Block" .English }}
		  {{ template "P8Block" .Maths }}
		  {{ template "P8Block" .EBacc }}
		  {{ template "P8Block" .Other }}
		  {{ template "P8Block" .Overall }}
		  {{ end }}

		  <td style="text-align:center;">{{printf "%.1f"  .Attendance.Latest }}</td>
		</tr>
		{{ end }}
	  </tbody>
	</table>
  </div>
  <div class="col-sm-1"></div>
</div>


<script>
$(function () {
  $('[data-toggle="popover"]').popover()
})
</script>

<script src="/static/js/plotly.min.js"></script>

<script>
var national = {
  x: [{{ range .NatLine }}{{ .X }}, {{ end }}],
  y: [{{ range .NatLine }}{{ .Y }}, {{ end }}],
  mode: 'lines',
  name: 'National',
  line: {shape: 'spline'},
  type: 'scatter',
  hoverinfo: 'none',
};

var male_pp = {
  x: [{{ range (index .PupilData 0) }}{{ .X }}, {{ end }}],
  y: [{{ range (index .PupilData 0) }}{{ .Y }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Male: Disadvantaged',
  text: [{{ range (index .PupilData 0) }}"<b>{{.Name}}</b><br>{{ range .P8.Subjects }}<br>{{.}}{{ end }} ", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};

var male_non = {
  x: [{{ range (index .PupilData 1) }}{{ .X }}, {{ end }}],
  y: [{{ range (index .PupilData 1) }}{{ .Y }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Male: Non-Disadvantaged',
  text: [{{ range (index .PupilData 1) }}"<b>{{.Name}}</b><br>{{ range .P8.Subjects }}<br>{{.}}{{ end }} ", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};

var female_pp = {
  x: [{{ range (index .PupilData 2) }}{{ .X }}, {{ end }}],
  y: [{{ range (index .PupilData 2) }}{{ .Y }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Female: Disadvantaged',
  text: [{{ range (index .PupilData 2) }}"<b>{{.Name}}</b><br>{{ range .P8.Subjects }}<br>{{.}}{{ end }} ", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};

var female_non = {
  x: [{{ range (index .PupilData 3) }}{{ .X }}, {{ end }}],
  y: [{{ range (index .PupilData 3) }}{{ .Y }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Female: Non-Disadvantaged',
  text: [{{ range (index .PupilData 3) }}"<b>{{.Name}}</b><br>{{ range .P8.Subjects }}<br>{{.}}{{ end }} ", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};

var data = [ national, male_pp, male_non, female_pp , female_non, ];

var layout = {
  hovermode: 'closest',
  scene: {
	aspectratio: {
	  x: 3,
	  y: 2,
	},
  },
  xaxis: {
	title: "KS2 Average Point Score",
	range: [ 0.75, 6.25 ]
  },
  yaxis: {
	title: "Attainment 8 Score",
	range: [0, 80]
  },
  legend: {
	x: 0.05,
	y: 1,
	xanchor: "left",
	yanchor: "top",
  },
};

Plotly.newPlot('chart', data, layout);
</script>

`,

	"progress8groups.tmpl": `
<h2>Progress 8 Group Summary</h2>
<br>

{{ $q := .Query }}
<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">

	<h4>Group Details</h4>
	<br>

	<table class="table table-condensed table-hover">
	  <thead>
		<th>Group</th>
		<th style="text-align:center;">Cohort</th>
		<th style="text-align:center;">KS2 APS</th>
		<th style="text-align:center;">Entries</th>
		<th style="text-align:center;">Attainment 8</th>
		<th style="text-align:center;">English</th>
		<th style="text-align:center;">Mathematics</th>
		<th style="text-align:center;">EBacc</th>
		<th style="text-align:center;">Other</th>
		<th style="text-align:center;">Progress 8</th>
		<th style="text-align:center;">Attendance</th>
	  </thead>
	  <tbody>
		{{ range .Groups }}
		  {{ with index .Group.Progress8.Progress 4 }}
			{{ if gt . 0.2 }}<tr class="success">
			  {{ else if lt . -0.2 }}<tr class="danger">
				{{ else }}<tr class="warning">
				  {{ end }}
				{{ end }}
				<td><a href="/progress8/?{{ $q }}{{ .Query }}">{{ .Name }}</a></td>
				{{ with .Group }}
				  <td style="text-align:center;">{{ .Cohort }}</td>
				  <td style="text-align:center;">{{ printf "%.1f" .KS2APS }}</td>
				  {{ with .Progress8 }}
					<td style="text-align:center;">{{ printf "%.1f" (index .Entries 4) }}</td>
					<td style="text-align:center;">{{ printf "%.1f" (index .Attainment 4) }}</td>
					{{ range .Progress }}
					  <td style="text-align:center;">{{ printf "%.2f" . }}</td>
					{{ end }}
				  {{ end }}
				  <td style="text-align:center;">{{ Percent .Attendance.PercentAttendance 1 }}</td>
				{{ end }}
				</tr>
			  {{ end }}
	  </tbody>
	</table>

	<br>
	<h4>Group Matrix</h4>
	<br>

	{{ with .Matrix }}
	  <table class="table table-condensed table-hover">
		<thead>
		  <th>&nbsp;</th>
		  {{ range .Headers }}
			<th style="text-align:center;">{{ . }}</th>
		  {{ end }}
		</thead>
		<tbody>
		  {{ $headers := .Headers }}
		  {{ range $i, $g := .Groups }}
			<tr>
			  <th>{{ index $headers $i }}</th>
			  {{ range $g }}
				{{ $p8 := (index .Group.Progress8.Progress 4) }}
				{{ if eq (len .Group.Students) 0 }}
				  <td style="text-align:center;"> - </td>
				{{ else if gt $p8 0.2 }}
				  <td style="text-align:center;" class="success"><a href="/progress8/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $p8 }}</a></td>
				{{ else if lt $p8 -0.2 }}
				  <td style="text-align:center;" class="danger"><a href="/progress8/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $p8 }}</a></td>
				{{ else }}
				  <td style="text-align:center;" class="warning"><a href="/progress8/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $p8 }}</a></td>
				{{ end }}
				</td>
			  {{ end }}
			</tr>
		  {{ end }}
		</tbody>
	  </table>
	{{ end }}

  </div>
  <div class="col-sm-1"></div>
</div>

<br>

`,

	"progressgrid.tmpl": `
<div class="row">
  <div class="col-sm-9"><h2>Progress Grid</h2></div>
  <div class="col-sm-3">
	<a href="/export/subject/{{.SubjID}}/?{{.Query}}" class="btn btn-primary btn-sm pull-right">
	  <span class="glyphicon glyphicon-download-alt"> Download</span>
	</a>
  </div>
</div>

<ul class="breadcrumb">
  <li><a href="/progressgrid/?{{.Query}}">Subjects</a></li>
  <li><a href="/progressgrid/{{.Subject}}/?{{.Query}}">{{.Subject}}</a></li>
  <li><a href="/progressgrid/{{.Subject}}/{{.SubjID}}/?{{.Query}}">{{.Level}}</a></li>
  <li class="active">{{.Class}}</li>
</ul>

{{ $q := .Query }}

{{ with .ProgressGrid }}
  {{ $cells := .Cells }}
  {{ $cellVA := .CellVA }}
  {{ $cohorts := .Cohorts }}
  {{ $rowVA := .RowVA }}

  <table class="table table-condensed" style="table-layout:fixed;">
	<thead>
	  <th>KS2</th>
	  {{ range .Grades }}
		<th style="text-align:center;">{{.}}</th>
	  {{ end }}
	  <th style="text-align:center;">Cohort</th>
	  <th style="text-align:center;">VA</th>
	</thead>
	<tbody>
	  {{ range $i, $ks2 := .KS2 }}
		<tr>
		  <td>{{ $ks2 }}</td>
		  {{ with index $cells $i }}
			{{ range $j, $c := . }}
			  {{ $va := index $cellVA $i $j }}
			  <td style="text-align:center; border:1px solid #888888"
		 {{ if gt $va 0.67  }}class="success"
		 {{ else if lt $va -0.33 }}class="danger"
		 {{ else }}class="warning" {{ end }}

		 {{ if eq (len $c.Students) 0 }}></td>
			  {{ else }}
				data-container="body"
				data-toggle="popover"
				data-placement="right"
				title="VA {{ printf "%+.2f" $va }}"
				data-html="true"
				data-trigger="hover click"
				data-content="{{range $c.Students}}
				<a href='/student/{{.UPN}}/?{{$q}}' target='_blank'>{{.Name}}</a><br>
			  {{end}}">
			  {{ len $c.Students }}</td>
		  {{ end }}
		{{ end }}
	  {{ end }}
	  <td style="text-align:center;">{{ index $cohorts $i }}</td>
	  <td style="text-align:center;">{{ printf "%+.2f" (index $rowVA $i) }}</td>
		</tr>
	  {{ end }}
	</tbody>
	<tfoot>
	  <th>Total</th>
	  {{ range .Counts }}
		<th style="text-align:center;">{{ . }} </th>
	  {{ end }}
	{{ end }}
	  <th style="text-align:center;">{{ .Group.Cohort }}</th>
	{{ with .Group.SubjectVA .Subject }}
	  <th style="text-align:center;">{{ printf "%.2f" .VA }}</th>

	{{ end }}
	</tfoot>
  </table>


  <script>
$(function () {
  $('[data-toggle="popover"]').popover()
})
  </script>

  <br>
  <h4>Students</h4>
  <table class="table table-condensed table-striped table-hover sortable">
	<thead>
	  <th>Name</th>
	  <th>Class</th>
	  <th style="text-align:center;">KS2</th>
	  <th style="text-align:center;">Gender</th>
	  <th style="text-align:center;">PP</th>
	  <th style="text-align:center;">Grade</th>
	  <th style="text-align:center;">Effort</th>
	  <th style="text-align:center;">Value Added</th>
	  <th style="text-align:center;">Attendance</th>
	</thead>
	<tbody>
	  {{ $subj := .Subject }}
	  {{ $prior := .KS2Prior }}
	  {{range .Group.Students}}
		<tr>
		  <td><a href="/student/{{.UPN}}/?{{$q}}">{{.Name}}</a></td>
		  <td>{{ .Class $subj }}</td>
		  <td style="text-align:center;">{{ .KS2.Score $prior }}</td>
		  <td style="text-align:center;">{{ .Gender }}</td>
		  <td style="text-align:center;">{{ template "PP" .PP }}</td>
		  <td style="text-align:center;">{{ .SubjectGrade $subj }}</td>
		  <td style="text-align:center;">{{ .SubjectEffort $subj }}</td>
		  <td style="text-align:center;">{{ template "StudentVA" (.SubjectVA $subj).Score }}</td>
		  <td style="text-align:center;">{{ template "StudentAttendance" .Attendance.Latest }}</td>
		</tr>
	  {{end}}
	</tbody>
  </table>
`,

	"select-class.tmpl": `
<h2>{{.Heading}}</h2>

<ul class="breadcrumb">
  <li><a href="{{.BasePath}}/?{{.Query}}">Subjects</a></li>
  <li><a href="{{.BasePath}}/{{.Subject}}/?{{.Query}}">{{.Subject}}</a></li>
  <li class="active">{{.Level}}</li>
</ul>

{{ $p := .Path }}
{{ $c := .Classes }}
{{ $q := .Queries }}

{{ range .Years }}
<h4>Year {{.}}</h4>
<div class="row">
  <div class="col-sm-3"></div>
  <div class="col-sm-6">
	<ul class="list-group">
	  <a class="list-group-item" href="{{ $p }}/All Year {{.}}/?{{ index $q . }}">All Year {{.}}</a>
	  {{ $y := . }}
	  {{ range index $c . }}
	  <a class="list-group-item" href="{{ $p }}/{{ . }}/?{{ index $q $y }}">{{ . }}</a>
	  {{ end }}
	</ul>
  </div>
  <div class="col-sm-3"></div>
</div>
{{ end }}

`,

	"select-level.tmpl": `
<h2>{{.Heading}}</h2>

<ul class="breadcrumb">
  <li><a href="{{.BasePath}}/?{{.Query}}">Subjects</a></li>
  <li class="active">{{.Subject}}</li>
</ul>

<div class="row">
  <div class="col-sm-3"></div>
  <div class="col-sm-6">
	{{ $p := .Path }}
	{{ $q := .Query }}
	<ul class="list-group">
	  {{ range .Levels }}
	  <a class="list-group-item" href="{{ $p }}/{{.SubjID}}/?{{ $q }}">{{.Level}}</a>
	  {{ end }}
	</ul>
  </div>
  <div class="col-sm-3"></div>
</div>
`,

	"select-subject.tmpl": `
<h2>{{.Heading}}</h2>

<ul class="breadcrumb">
  <li class="active">Subjects</li>
</ul>

<div class="row">
  <div class="col-sm-3"></div>
  <div class="col-sm-6">
  {{ $p := .Path }}
  {{ $q := .Query }}
  <ul class="list-group">
    {{ range .Subjects }}
    <a class="list-group-item" href="{{ $p }}/{{.}}/?{{ $q }}">{{.}}</a>
    {{ end }}
  </ul>
  </div>
  <div class="col-sm-3"></div>
</div>

`,

	"select-year.tmpl": `
<h2>{{.Heading}}</h2>

<ul class="breadcrumb">
  <li><a href="{{.BasePath}}/?{{.Query}}">Subjects</a></li>
  <li><a href="{{.BasePath}}/{{.Subject}}/?{{.Query}}">{{.Subject}}</a></li>
  <li class="active">{{.Level}}</li>
</ul>

{{ $p := .Path }}
{{ $q := .Queries }}

<div class="row">
  <div class="col-sm-3"></div>
  <div class="col-sm-6">
	<ul class="list-group">
	  {{ range .Years }}
	  <a class="list-group-item" href="{{ $p }}/All Year {{.}}/?{{ index $q . }}">Year {{.}}</a>
	  {{ end }}
	</ul>
  </div>
  <div class="col-sm-3"></div>
</div>

`,

	"student.tmpl": `
{{ define "StudentKS2Pane" }}
  <br>
  <table class="table table-hover">
	<tbody>
	  <tr>
		<td>Average Point Score</td>
		<td>{{.APS}}</td>
	  </tr>
	  <tr>
		<td>Band</td>
		<td>{{.Band}}</td>
	  </tr>
	  <tr>
		<td>Mathematics</td>
		<td>{{.Ma}}</td>
	  </tr>
	  <tr>
		<td>English</td>
		<td>{{.En}}</td>
	  </tr>
	  <tr>
		<td>Reading</td>
		<td>{{.Re}}</td>
	  </tr>
	  <tr>
		<td>Writing</td>
		<td>{{.Wr}}</td>
	  </tr>
	  <tr>
		<td>GPS</td>
		<td>{{.GPS}}</td>
	  </tr>
	</tbody>
  </table>
{{ end }}

{{ define "StudentSENPane" }}
  <br>
  <table class="table table-condensed table-hover">
	<tbody>
	  {{ $upn := .UPN }}
	  {{with .SEN}}
		<tr>
		  <td>Status</td>
		  <td>{{.Status}}</td>
		</tr>
		<tr>
		  <td>Needs</td>
		  <td>{{.Need}}</td>
		</tr>
		<tr>
		  <td>IEP</td>
		  <td>{{if .IEP}}<a href="/static/files/iep/{{ $upn}}.pdf">Download <span class="glyphicon glyphicon-download" style="color: #009933;"></span></a>
		  {{else}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
		  {{end}}</td>
		</tr>
	  {{end}}
	  <tbody>
  </table>
{{ end }}

{{ define "StudentResultsPane" }}
  <table class="table table-condensed table-hover sortable">
	<thead>
	  <th>Subject</th>
	  <th style="text-align:center;">Level</th>
	  <th style="text-align:center;">Grade</th>
	  <th style="text-align:center;">Effort</th>
	  <th style="text-align:center;">Class</th>
	  <th style="text-align:center;">Teacher</th>
	</thead>
	<tbody>
	  {{range .Results}}
		<tr>
		  <td>{{.Subj}}</td>
		  <td style="text-align:center;">{{.Lvl}}</td>
		  <td style="text-align:center;">{{.Grd}}</td>
		  <td style="text-align:center;">{{.Effort}}</td>
		  <td style="text-align:center;">{{.Class}}</td>
		  <td style="text-align:center;">{{.Teacher}}</td>
		</tr>
	  {{end}}
	</tbody>
  </table>
{{ end }}

{{ define "StudentHeadlinesPane" }}
  <br>
  <h4>Headline Figures</h4><br>
  <table class="table table-condensed table-hover">
	<tbody>
	  <tr>
		<td>Basics</td>
		<td>{{ template "TickCross" .Basics }}</td>
	  </tr>
	  {{ with .Basket }}
		<tr>
		  <td>Attainment 8</td>
		  <td>{{ printf "%.1f" .Overall.Attainment }}</td>
		</tr>
		<tr>
		  <td>Progress 8</td>
		  <td>{{ printf "%+.1f" .Overall.Progress8 }}</td>
		</tr>
	</tbody>
  </table>

  <h4>Progress 8</h4><br>
  <table class="table table-condensed table-hover">
	<thead>
	  <th>Basket</th>
	  <th style="text-align:center;">Entries</th>
	  <th style="text-align:center;">Attainment Per Slot</th>
	  <th style="text-align:center;">Progress 8</th>
	</thead>
	<tbody>
	  <tr>
		<td>Mathematics</td>
		<td style="text-align:center;">{{ .English.Entries }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .English.AttainmentPerSlot }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .English.Progress8 }}</td>
	  </tr>
	  <tr>
		<td>English</td>
		<td style="text-align:center;">{{ .Maths.Entries }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .Maths.AttainmentPerSlot }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .Maths.Progress8 }}</td>
	  </tr>
	  <tr>
		<td>EBacc</td>
		<td style="text-align:center;">{{ .EBacc.Entries }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .EBacc.AttainmentPerSlot }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .EBacc.Progress8 }}</td>
	  </tr>
	  <tr>
		<td>Other</td>
		<td style="text-align:center;">{{ .Other.Entries }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .Other.AttainmentPerSlot }}</td>
		<td style="text-align:center;">{{ printf "%.1f" .Other.Progress8 }}</td>
	  </tr>
	{{ end }}
	</tbody>
  </table>

  <h4>EBacc</h4><br>
  <table class="table table-condensed table-hover">
	<thead>
	  <th>Area</th>
	  <th style="text-align:center;">Entered</th>
	  <th style="text-align:center;">Achieved</th>
	</thead>
	<tbody>

	  {{ with .EBaccArea "E" }}
		<tr>
		  <td>English</td>
		  <td style="text-align:center;">{{ template "TickCross" .Entered }}</td>
		  <td style="text-align:center;">{{ template "TickCross" .Achieved }}</td>
		</tr>
	  {{ end }}

	  {{ with .EBaccArea "M" }}
		<tr>
		  <td>Mathematics</td>
		  <td style="text-align:center;">{{ template "TickCross" .Entered }}</td>
		  <td style="text-align:center;">{{ template "TickCross" .Achieved }}</td>
		</tr>
	  {{ end }}
	  {{ with .EBaccArea "S" }}
		<tr>
		  <td>Science</td>
		  <td style="text-align:center;">{{ template "TickCross" .Entered }}</td>
		  <td style="text-align:center;">{{ template "TickCross" .Achieved }}</td>
		</tr>
	  {{ end }}
	  {{ with .EBaccArea "H" }}
		<tr>
		  <td>Humanities</td>
		  <td style="text-align:center;">{{ template "TickCross" .Entered }}</td>
		  <td style="text-align:center;">{{ template "TickCross" .Achieved }}</td>
		</tr>
	  {{ end }}
	  {{ with .EBaccArea "L" }}
		<tr>
		  <td>Languages</td>
		  <td style="text-align:center;">{{ template "TickCross" .Entered }}</td>
		  <td style="text-align:center;">{{ template "TickCross" .Achieved }}</td>
		</tr>
	  {{ end }}
	</tbody>
  </table>
{{ end }}

{{ define "StudentProgress8Pane" }}
  <table class="table table-condensed table-striped table-hover">
	<thead>
	  <th>Slot</th>
	  <th>Subject</th>
	  <th>Grade</th>
	  <th>Points</th>
	</thead>
	<tbody>
	  <tr>
		<td>English</td>
		{{with (index .Slots 0)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>English</td>
		{{with (index .Slots 1)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>Mathematics</td>
		{{with (index .Slots 2)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>Mathematics</td>
		{{with (index .Slots 3)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>EBacc 1</td>
		{{with (index .Slots 4)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>EBacc 2</td>
		{{with (index .Slots 5)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>EBacc 3</td>
		{{with (index .Slots 6)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>Other 1</td>
		{{with (index .Slots 7)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>Other 2</td>
		{{with (index .Slots 8)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	  <tr>
		<td>Other 3</td>
		{{with (index .Slots 9)}}
		  <td>{{.Subject}}</td>
		  <td>{{.Grade}}</td>
		  <td>{{.Points}}</td>
		{{end}}
	  </tr>
	</tbody>
	<tfoot>
	  {{ with .Overall }}
		<tr>
		  <th>Total</th>
		  <th>{{ .Entries }} Entries</th>
		  <th></th>
		  <th>{{ printf "%.1f" .Attainment }}</th>
		</tr>
		<tr>
		  <th>Expected Score</th>
		  <th></th>
		  <th></th>
		  <th>{{ printf "%.1f" .Expected }}</th>
		</tr>
		<tr>
		  <th>Progress 8 Score</th>
		  <th></th>
		  <th></th>
		  <th>{{ printf "%+.1f" .Progress8 }}</th>
		</tr>
	  {{ end }}
	</tfoot>
  </table>
{{ end }}

{{ define "StudentAttendancePane" }}
  <h4>Attendance data from 1st September.</h4>
  <table class="table">
	<thead>
	  <th style='text-align:center;vertical-align:middle'>Possible</th>
	  <th style='text-align:center;vertical-align:middle'>Absences</th>
	  <th style='text-align:center;vertical-align:middle'>Unauthorised</th>
	  <th style='text-align:center;vertical-align:middle'>% Attendance</th>
	</thead>
	<tbody>
	  <tr>
		<td style='text-align:center;vertical-align:middle'>{{.Possible}}</td>
		<td style='text-align:center;vertical-align:middle'>{{.Absences}}</td>
		<td style='text-align:center;vertical-align:middle'>{{.Unauthorised}}</td>
		<td style='text-align:center;vertical-align:middle'>{{Percent .Latest 1}}</td>
	  </tr>
	</tbody>
  </table>

  <br>
  <h4>Absences per Session</h4>
  <table class="table">
	<thead>
	  <th>Session</th>
	  <th style='text-align:center;vertical-align:middle'>Monday</th>
	  <th style='text-align:center;vertical-align:middle'>Tuesday</th>
	  <th style='text-align:center;vertical-align:middle'>Wednesday</th>
	  <th style='text-align:center;vertical-align:middle'>Thursday</th>
	  <th style='text-align:center;vertical-align:middle'>Friday</th>
	</thead>
	<tbody>
	  <tr>
		<td>Morning</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 0}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 2}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 4}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 6}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 8}}</td>
	  </tr>
	  <tr>
		<td>Afternoon</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 1}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 3}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 5}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 7}}</td>
		<td style='text-align:center;vertical-align:middle'>{{index .Sessions 9}}</td>
	  </tr>
	</tbody>
  </table>
{{ end }}

{{/* Start of Template Proper */}}
{{ with .Student }}
  <h2>{{ .Forename }} {{ .Surname }}</h2>
  <br>

  <div class="row">
	<div class="col-sm-1"></div>
	<div class="col-sm-5">
	  <table class="table table-hover">
		<tbody>
		  <tr>
			<td>Form class</td>
			<td>{{.Year}} {{.Form}}</td>
		  </tr>
		  <tr>
			<td>Gender</td>
			<td>{{.Gender}}</td>
		  </tr>
		  <tr>
			<td>KS2 Average</td>
			<td>{{with .KS2}}{{.Av}}{{end}}</td>
		  </tr>
		  <tr>
			<td>Pupil Premium</td>
			<td>{{ template "PP" .PP }}</td>
		  </tr>
		  <tr>
			<td>SEN</td>
			<td>{{ template "TickCross" (ne .SEN.Status "N") }}</td>
		  </tr>
		  <tr>
			<td>Ethnicity</td>
			<td>{{.Ethnicity}}</td>
		  </tr>
		  <tr>
			<td>EAL</td>
			<td>{{ template "TickCross" .EAL }}</td>
		  </tr>
		</tbody>
	  </table>
	</div>
	<div class="col-sm-2"></div>
	<div class="col-sm-2">
	  <img src="/images/photos/{{.UPN}}.png" alt="No photo found" style="width: 180px; ">
	</div>
	<div class="col-sm-2"></div>
  </div>

  <br>
  <div class="row" style="min-height:700px;">
	<div class="col-sm-1"></div>
	<div class="col-sm-10">
	  <ul class="nav nav-tabs">
		<li class="active"><a href="#ks2" data-toggle="tab" aria-expanded="true">KS2</a></li>
		<li><a href="#sen" data-toggle="tab" aria-expanded="false">SEN</a></li>
		<li><a href="#attendance" data-toggle="tab" aria-expanded="false">Attendance</a></li>
		<li><a href="#results" data-toggle="tab" aria-expanded="false">Latest Assessments</a></li>
		{{ if eq .Year 10 11 }}
		  <li><a href="#headlines" data-toggle="tab" aria-expaanded="false">Headlines</a></li>
		  <li><a href="#progress8" data-toggle="tab" aria-expaanded="false">Progress 8</a></li>
		{{ end }}
	  </ul>

	  <br>

	  <div id="myTabContent" class="tab-content">
		<div class="tab-pane fade active in" id="ks2"> 
		  {{ template "StudentKS2Pane" .KS2 }}
		</div>
		<div class="tab-pane fade" id="sen">
		  {{ template "StudentSENPane" . }}
		</div>
		<div class="tab-pane fade" id="results">
		  {{ template "StudentResultsPane" . }}
		</div>

		{{ if eq .Year 10 11 }}
		  <div class="tab-pane fade" id="headlines">
			{{ template "StudentHeadlinesPane" . }}
		  </div>
		  <div class="tab-pane fade" id="progress8">
			{{ template "StudentProgress8Pane" .Basket }}
		  </div>
		{{ end }}

		<div class="tab-pane fade" id="attendance">
		  {{ template "StudentAttendancePane" .Attendance }}
		</div>
	  </div>
	</div>
	<div class="col-sm-1"></div>
  </div>
{{ end }}


`,

	"studentsearch.tmpl": `
<h2>Student Search Results</h2>

<p>Searching for: <b>{{.Name}}</b></p>

<div class="row">
  <div class="col-sm-3"></div>
  <div class="col-sm-6">
	{{ $q := .Query }}
	<ul class="list-group">
	  {{ range .Group.Students }}
	  <a class="list-group-item" href="/student/{{.UPN}}/?{{$q}}">{{.Name}}	 ({{.Year}} {{.Form}})</a>
	  {{end}}
	</ul>
  </div>
  <div class="col-sm-3"></div>
</div>

<br>
<p>If you can't find the student you are looking for, try using
just their surname or forename.<br>
If you're unsure of the spelling, a '*' can be used to replace
one or more characters.</p>
`,

	"subject-overview.tmpl": `
<h2>Subject Summaries</h2>
<br>
<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<table class="table table-striped table-hover sortable">
	  <thead>
		<th>Subject</th>
		<th style="text-align:center;">Cohort</th>
		<th style="text-align:center;">% PP</th>
		<th style="text-align:center;">KS2 APS</th>
		<th style="text-align:center;">Current APS</th>
		<th style="text-align:center;">Average Grade</th>
		<th style="text-align:center;">Value Added</th>
	  </thead>
	  <tbody>
		{{ $y := .Year }}
		{{ $q := .Query }}
		{{ range .Summaries }}
		<tr>
		  <td><a href="/progressgrid/{{ .Subject.Subj }}/{{ .Subject.SubjID }}/All Year {{ $y }}/?{{ $q}}">{{ .Subject.Subj }}</a></td>
		  <td style="text-align:center;">{{ .Group.Cohort }}</td>
		  <td style="text-align:center;">{{ Percent .Group.PP 1 }}</td>
		  <td style="text-align:center;">{{ printf "%.1f" .Group.KS2APS }}</td>
		  {{ $pts := (.Group.SubjectPoints .Subject.Subj) }}
		  <td style="text-align:center;">{{ printf "%.1f" $pts }}</td>
		  <td style="text-align:center;">{{ .Subject.Level.GradeEquivalent $pts }}</td>
		  <td style="text-align:center;">{{ printf "%+.2f" (.Group.SubjectVA .Subject.Subj).VA }}</td>
		</tr>
		{{ end }}
	</table>
  </div>
  <div class="col-sm-1"></div>
</div>
`,

	"subjectgroups.tmpl": `
<h2>{{ .Subj.Subj }} Group Summary</h2>

<ul class="breadcrumb">
  <li><a href="/subjectgroups/?{{.Query}}">Subjects</a></li>
  <li><a href="/subjectgroups/{{.Subj.Subj}}/?{{.Query}}">{{.Subj.Subj}}</a></li>
  <li><a href="/subjectgroups/{{.Subj.Subj}}/{{.Subj.SubjID}}/?{{.Query}}">{{.Subj.Lvl}}</a></li>
  <li class="active">Year {{ .Year }}</li>
</ul>

<br>

{{ $s := .Subj }}
{{ $q := .Query }}
{{ $y := .Year }}

<div class="row">
  <div class="col-md-1"></div>
  <div class="col-md-10">
	<h4>Pupil Groups</h4>
	<br>

	<table class="table table-condensed table-striped table-hover sortable">
	  <thead>
		<th>Group</th>
		<th style="text-align:center;">Cohort</th>
		<th style="text-align:center;">% PP</th>
		<th style="text-align:center;">KS2 APS</th>
		<th style="text-align:center;">Current APS</th>
		<th style="text-align:center;">Average Grade</th>
		<th style="text-align:center;">Value Added</th>
	  </thead>
	  <tbody>
		{{ range .SubGroups }}
		  {{ $va := (.Group.SubjectVA $s.Subj).VA }}
		  {{ if gt $va 0.33 }}<tr class="success">
		  {{ else if lt $va -0.33 }}<tr class="danger">
		  {{ else }}<tr class="warning">
		  {{ end }}
			<td><a href="/progressgrid/{{ $s.Subj }}/{{ $s.SubjID }}/All Year {{ $y }}/?{{ $q }}{{ .Query }}">{{ .Name }}</a></td>
			<td style="text-align:center;">{{ .Group.Cohort }}</td>
			<td style="text-align:center;">{{ Percent .Group.PP 1 }}</td>
			<td style="text-align:center;">{{ printf "%.1f" .Group.KS2APS }}</td>
			{{ $pts := (.Group.SubjectPoints $s.Subj) }}
			<td style="text-align:center;">{{ printf "%.1f" $pts }}</td>
			<td style="text-align:center;">{{ $s.Level.GradeEquivalent $pts }}</td>
			<td style="text-align:center;">{{ printf "%+.2f" $va }}</td>
		  </tr>
		{{ end }}
	  </body>
	</table>

	<br>
	<h4>Group Matrix</h4>
	<br>

	{{ with .Matrix }}
	  <table class="table table-condensed table-hover" style="table-layout:fixed;">
		<thead>
		  <th style="width:20%;"></th>
		  {{ range .Headers }}
			<th style="text-align:center;">{{ . }}</th>
		  {{ end }}
		</thead>
		<tbody>
		  {{ $headers := .Headers }}
		  {{ range $i, $g := .Groups }}
			<tr>
			  <th>{{ index $headers $i }}</th>
			  {{ range $g }}
				{{ $va := (.Group.SubjectVA $s.Subj).VA }}
				{{ if eq (len .Group.Students) 0 }}
				  <td style="text-align:center;"> - </td>
				{{ else if gt $va 0.33 }}
				  <td style="text-align:center;" class="success"><a href="/progressgrid/{{ $s.Subj }}/{{ $s.SubjID }}/All Year {{ $y }}/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $va }}</a></td>
				{{ else if lt $va -0.33 }}
				  <td style="text-align:center;" class="danger"><a href="/progressgrid/{{ $s.Subj }}/{{ $s.SubjID }}/All Year {{ $y }}/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $va }}</a></td>
				{{ else }}
				  <td style="text-align:center;" class="warning"><a href="/progressgrid/{{ $s.Subj }}/{{ $s.SubjID }}/All Year {{ $y }}/?{{ $q }}{{ .Query }}">{{ printf "%+.2f" $va }}</a></td>
				{{ end }}
				</td>
			  {{ end }}
			</tr>
		  {{ end }}
		</tbody>
	  </table>
	{{ end }}

	<br>
	<h4>Classes</h4>
	<br>

	<table class="table table-condensed table-striped table-hover sortable">
	  <thead>
		<th>Group</th>
		<th style="text-align:center;">Cohort</th>
		<th style="text-align:center;">% PP</th>
		<th style="text-align:center;">KS2 APS</th>
		<th style="text-align:center;">Current APS</th>
		<th style="text-align:center;">Average Grade</th>
		<th style="text-align:center;">Value Added</th>
	  </thead>
	  <tbody>
		{{ range .Classes }}
		  {{ $va := (.Group.SubjectVA $s.Subj).VA }}
		  {{ if gt $va 0.33 }}<tr class="success">
		  {{ else if lt $va -0.33 }}<tr class="danger">
		  {{ else }}<tr class="warning">
		  {{ end }}
			<td><a href="/progressgrid/{{ $s.Subj }}/{{ $s.SubjID }}/{{ .Name }}/?{{ $q }}">{{ .Name }}</a></td>
			<td style="text-align:center;">{{ .Group.Cohort }}</td>
			<td style="text-align:center;">{{ Percent .Group.PP 1 }}</td>
			<td style="text-align:center;">{{ printf "%.1f" .Group.KS2APS }}</td>
			{{ $pts := (.Group.SubjectPoints $s.Subj) }}
			<td style="text-align:center;">{{ printf "%.1f" $pts }}</td>
			<td style="text-align:center;">{{ $s.Level.GradeEquivalent $pts }}</td>
			<td style="text-align:center;">{{ printf "%+.2f" $va }}</td>
		  </tr>
		{{ end }}
	  </body>
	</table>

  </div>
  <div class="col-md-1"></div>
</div>
`,
}
