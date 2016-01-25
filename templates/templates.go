
// Package templates contains all of the templates needed for the dashboard.
// This is automatically generated from the .tmpl files using the "embed-template"
// script.
package templates

var Templates = map[string]string{

"attendance.tmpl" : `
<h2>Attendance</h2>
<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<table class="table table-condensed table-hover">
	  <thead>
		<th>Group</th>
		<th style="text-align:center;">Cohort</th>
		<th style="text-align:center;">Possible</th>
		<th style="text-align:center;">Absences</th>
		<th style="text-align:center;">Unauthorised</th>
		<th style="text-align:center;">Attendance %</th>
		<th style="text-align:center;">Under 85%</th>
		<th style="text-align:center;">Under 90%</th>
	  </thead>
	  <tbody>
		{{ $data := .AttGroups }}
		{{ range .Headers }}
		<tr>
		  <td>{{ . }}</td>
		  {{ with index $data . }}
		  <td style="text-align:center;">{{ .Cohort }}</td>
		  <td style="text-align:center;">{{ .Possible }}</td>
		  <td style="text-align:center;">{{ .Absences }}</td>
		  <td style="text-align:center;">{{ .Unauthorised }}</td>
		  <td style="text-align:center;">{{ printf "%.1f" .PctAttendance }}</td>
		  <td style="text-align:center;">{{ .Under85 }}</td>
		  <td style="text-align:center;">{{ .Under90 }}</td>
		  {{ end }}
		</tr>
		{{ end }}
	  </tbody>
	</table>
  </div>
  <div class="col-sm-1"></div>
</div>


`,

"classlist.tmpl" : `
<h2>Class Lists</h2>

<ul class="breadcrumb">
  <li><a href="/classlist/?{{.Query}}">Subjects</a></li>
  <li><a href="/classlist/{{.Subject}}/?{{.Query}}">{{.Subject}}</a></li>
  <li><a href="/classlist/{{.Subject}}/{{.SubjID}}/?{{.Query}}">{{.Level}}</a></li>
  <li class="active">{{.Class}}</li>
</ul>


<div class="row" style="min-height:700px;">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<ul class="nav nav-tabs">
	  <li class="active"><a href="#info" data-toggle="tab" aria-expanded="true">Pupil Info</a></li>
	  <li class=""><a href="#photos" data-toggle="tab" aria-expanded="false">Photos</a></li>
	  <li class=""><a href="#assessments" data-toggle="tab" aria-expanded="false">Assessment History</a></li>
	</ul>

	<div id="myTabContent" class="tab-content">
	  <div class="tab-pane fade active in" id="info">
		<table class="table table-hover">
		  <thead>
			<th>Name</th>
			<th>Gender</th>
			<th>KS2 Average</th>
			<th>Pupil Premium</th>
			<th>SEN</th>
			<th>Attendance</th>
		  </thead>
		  <tbody>
			{{ $q := .Query }}
			{{ range .Students }}
			<tr>
			  <td><a href="/students/{{.UPN}}/?{{$q}}">{{.Surname}}, {{.Forename}}</a></td>
			  <td>{{.Gender}}</td>
			  <td>{{.KS2.Av}}</td>
			  <td>{{if .PP}}<span class="glyphicon glyphicon-ok" style="color: #009933;"></span>
				{{else}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
				{{end}}</td>
			  <td>{{.SEN.Status}}</td>
			  <td></td>
			</tr>
			{{ end}}
		  </tbody>
		</table>
	  </div>
	  <div class="tab-pane" id="photos">
		Photos to go here!
	  </div>
	  <div class="tab-pane" id="assessments">
		Assessment History to go here
	  </div>
	</div>
  </div>
  <div class="col-sm-1"></div>
</div>
`,

"effort.tmpl" : `
<h2>Attitude to Learning</h2>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">

	<table class="table table-hover sortable">
	  <thead>
		<th>Name</th>
		<th style="text-align:center;">1's</th>
		<th style="text-align:center;">2's</th>
		<th style="text-align:center;">3's</th>
		<th style="text-align:center;">4's</th>
		<th style="text-align:center;">Average</th>
		<th style="text-align:center;">Progress 8 Score</th>
	  </thead>
	  <tbody>
		{{ $q := .Query }}
		{{ range .Efforts }}
		<tr>
		  <td><a href="/students/{{.UPN}}/?{{$q}}">{{.Name}}</a></td>
		  <td style="text-align:center;">{{index .Scores 1}}</td>
		  <td style="text-align:center;">{{index .Scores 2}}</td>
		  <td style="text-align:center;">{{index .Scores 3}}</td>
		  <td style="text-align:center;">{{index .Scores 4}}</td>
		  <td style="text-align:center;">{{printf "%1.1f" .Average}}</td>
		  <td style="text-align:center;">{{printf "%1.1f" .Prog8}}</td>
		</tr>
		{{ end }}
	  </tbody>
	</table>
  </div>
  <div class="col-sm-1"></div>
</div>
`,

"em.tmpl" : `
<h2>English and Maths</h2>
<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<h4>Headlines</h4>

	<table class="table table-hover">
	  <thead>
		<th></th>
		<th style='text-align:center;vertical-align:middle'>Number</th>
		<th style='text-align:center;vertical-align:middle'>Percent</th>
	  </thead>
	  <tbody>
		<tr>
		  <td>Cohort</td>
		  <td style='text-align:center;vertical-align:middle'>{{.Cohort}}</td>
		  <td style='text-align:center;vertical-align:middle'></td>
		</tr>
		<tr>
		  <td>English</td>
		  <td style='text-align:center;vertical-align:middle'>{{.EnPass}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{printf "%1.f" .EnPassPct}}</td>
		</tr>
		<tr>
		  <td>Mathematics</td>
		  <td style='text-align:center;vertical-align:middle'>{{.MaPass}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{printf "%1.f" .MaPassPct}}</td>
		</tr>
		<tr>
		  <td>Both</td>
		  <td style='text-align:center;vertical-align:middle'>{{.BothPass}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{printf "%1.f" .BothPassPct}}</td>
		</tr>
	  </tbody>
	</table>

	<br>

	<h4>Students</h4>

	<table class="table table-condensed table-striped table-hover sortable">
	  <thead>
		<th>Name</th>
		<th>KS2</th>
		<th style='text-align:center;vertical-align:middle'>Eng Grade</th>
		<th style='text-align:center;vertical-align:middle'>Eng Effort</th>
		<th style='text-align:center;vertical-align:middle'>Maths Grade</th>
		<th style='text-align:center;vertical-align:middle'>Maths Effort</th>
		<th style='text-align:center;vertical-align:middle'>Both</th>
		<th style='text-align:center;vertical-align:middle'>Av Effort</th>
		<th style='text-align:center;vertical-align:middle'>Progress 8</th>
		<th style='text-align:center;vertical-align:middle'>Attendance</th>
	  </thead>
	  <tbody>
		{{ $q := .Query }}
		{{ range .Students }}
		<tr>
		  <td><a href="/students/{{.UPN}}/?{{$q}}">{{.Name}}</a></td>
		  <td style='text-align:center;vertical-align:middle'>{{.KS2.Av}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{.EnGrd}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{.EnEff}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{.MaGrd}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{.MaEff}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{if .Basics}}<span style="color: #009933;">Y</span>
			{{else}}<span style="color: #cc0000;">N</span>
			{{end}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{printf "%1.1f" .AvEff}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{printf "%1.2f" .P8}}</td>
		  <td style='text-align:center;vertical-align:middle'>{{printf "%1.1f" .Att}}%</td>
		</tr>
		{{ end }}
	  </tbody>
	  <table>
  </div>
  <div class="col-sm-1"></div>
</div>

`,

"filter.tmpl" : `
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

	  {{if not .Short}}
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
			<label class="checkbox-inline"><input type="checkbox" name="ethnicity" value="Other"{{if index $E "Other"}} checked="yes"{{end}}>Other</input></label>
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

"footer.tmpl" : `
</div>
</body>
</html>
`,

"group.tmpl" : `
<h2>{{.Title}}</h2>
<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<table class="table table-striped table-hover table-condensed">
	  {{ with .Summary }}
	  <thead>
		<th></th>
		<th style="text-align:center;">Cohort</th>
		{{ range .Headers }}
		<th style="text-align:center;">{{ . }}</th>
		{{ end }}
	  </thead>
	  <tbody>
		{{ range .Groups }}
		<tr>
		  <td>{{ .Name }}</td>
		  <td style="text-align:center;">{{ .Cohort }}</td>
		  {{ range .Scores }}
		  {{ if .Error }}
		  <td style="text-align:center;">-</td>
		  {{ else }}
		  <td style="text-align:center;">{{ printf "%.1f" .Score }}</td>
		  {{ end }}
		  {{ end }}
		</tr>
		{{ end }}
	  </tbody>
	  {{ end }}
	</table>
  </div>
  <div class="col-sm-1"></div>
</div>

<br>

`,

"header.tmpl" : `
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

  </head>
  <body>

	<nav class="navbar navbar-default">
	  <div class="container">
		<div class="navbar-header">
		  <a class="navbar-brand" href="/index/?{{.Query}}">{{.School}}</a>
		</div>

		<form class="navbar-form navbar-right form-horizontal" action="/studentsearch/" role="search">
		  {{with .F}}
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
			  <li><a href="/subjects/?{{.Query}}">Summary</a></li>
			  <li class="divider"></li>
			  <li><a href="/progressgrid/?{{.Query}}">Progress Grid</a></li>
			  <li class="divider"></li>
			  <li><a href="/classlist/?{{.Query}}">Class Lists</a></li>
			</ul>
		  </li>
		</ul>

		<ul class="nav navbar-nav navbar-right">
		  <li class="dropdown">
			<a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">Whole School <span class="caret"></span></a>
			<ul class="dropdown-menu" role="menu">
			  <li class="disabled"><a href="/summary/?{{.Query}}">Summary</a></li>
			  <li class="divider"></li>
			  <li><a href="/progress8/?{{.Query}}">Progress 8</a></li>
			  <li><a href="/basics/?{{.Query}}">English and Maths</a></li>
			  <li class="disabled"><a href="/ebacc/?{{.Query}}">English Baccalaureate</a></li>
			  <li class="divider"></li>
			  <li><a href="/effort/?{{.Query}}">Effort</a></li>
			  <li><a href="/attendance/?{{.Query}}">Attendance</a></li>
			</ul>
		  </li>
		</ul>
	  </div>
	</nav>

	<div class="container">
`,

"progress8.tmpl" : `
<div class="row">
  <div class="col-sm-9"><h2>Progress 8</h2></div>
  <div class="col-sm-3">
	<a href="/export/headlines/?{{.Query}}" class="btn btn-primary btn-sm pull-right">
	  <span class="glyphicon glyphicon-download-alt"> Download</span>
	</a>
  </div>
</div>

<br>

<div width="100%" id="chart"></div>

<br>

<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">

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
		  {{ range .Slots }}
		  <td style="text-align:center;">{{printf "%.2f" .Att8 }}</td>
		  {{ end }}
		</tr>
		<tr>
		  <td>Points per Slot</td>
		  {{ range .Slots }}
		  <td style="text-align:center;">{{printf "%.2f" .Att8PerSlot }}</td>
		  {{ end }}
		</tr>
		<tr>
		  <td>Entries</td>
		  {{ range .Slots }}
		  <td style="text-align:center;">{{printf "%.2f" .Entries }}</td>
		  {{ end }}
		</tr>
		<tr>
		  <td>Progress 8</td>
		  {{ range .Slots }}
		  <td style="text-align:center;">{{printf "%+.2f" .Prog8 }}</td>
		  {{ end }}
		</tr>
	  </tbody>
	</table>

	<br>

	<table class="table table-condensed table-striped table-hover sortable">
	  <thead>
		<th>Name</th>
		<th style="text-align:center;">KS2</th>
		<th style="text-align:center;">PP</th>
		<th style="text-align:center;">English</th>
		<th style="text-align:center;">Maths</th>
		<th style="text-align:center;">EBacc</th>
		<th style="text-align:center;">Other</th>
		<th style="text-align:center;">Overall</th>
		<th style="text-align:center;">Attendance</th>
	  </thead>
	  <tbody>
		{{ $q := .Query }}
		{{ range .Students }}
		<tr>
		  <td><a href="/students/{{ .UPN }}/?{{ $q }}">{{ .Name }}</a></td>
		  <td style="text-align:center;">{{ .KS2.Av }}</td>
		  <td style="text-align:center;">
			{{if .PP}}<span class="glyphicon glyphicon-ok" style="color: #009933;"></span>
			{{else}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
			{{end}}
		  </td>

		  {{ range .Slots }}
		  <td style="text-align:center;"
		data-container="body"
  data-toggle="popover"
  data-placement="right"
  data-html="true"
  data-trigger="hover"
  data-content="{{ .Text }}"> {{printf "%+.2f" .Prog8 }}
		  </td>
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
  x: [{{ range .GraphNat.X }}{{.}}, {{end}}], 
  y: [{{ range .GraphNat.Y }}{{.}}, {{end}}], 
  mode: 'lines',
  name: 'National',
  line: {shape: 'spline'},
  type: 'scatter',
  hoverinfo: 'none',
};

{{ with index .GraphPupils 0 }}
var male_pp = {
  x: [{{ range .X }}{{ . }}, {{ end }}],
  y: [{{ range .Y }}{{ . }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Male - Disadvantaged',
  text: [{{ range .Text }}"{{.}}", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};
{{ end }}

{{ with index .GraphPupils 1 }}
var male_non = {
  x: [{{ range .X }}{{ . }}, {{ end }}],
  y: [{{ range .Y }}{{ . }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Male - Non-Disadvantaged',
  text: [{{ range .Text }}"{{.}}", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};
{{ end }}

{{ with index .GraphPupils 2 }}
var female_pp = {
  x: [{{ range .X }}{{ . }}, {{ end }}],
  y: [{{ range .Y }}{{ . }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Female - Disadvantaged',
  text: [{{ range .Text }}"{{.}}", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};
{{ end }}

{{ with index .GraphPupils 3 }}
var female_non = {
  x: [{{ range .X }}{{ . }}, {{ end }}],
  y: [{{ range .Y }}{{ . }}, {{ end }}],
  mode: 'markers',
  type: 'scatter',
  name: 'Female - Non-Disadvantaged',
  text: [{{ range .Text }}"{{.}}", {{ end }}],
  marker: { size: 8 },
  hoverinfo: 'text',
};
{{ end }}

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

"progressgrid.tmpl" : `
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

{{ with .Grid }}
{{ $g := .Grades }}
{{ $c := .Cells }}
{{ $va := .VA }}
{{ $tm := .TMExists }}

<table class="table table-condensed" style="table-layout:fixed;">
  <thead>
	<th>KS2</th>
	{{ range .Grades }}
	<th style="text-align:center;">{{.}}</th>
	{{ end }}
	<th style="text-align:center;">VA</th>
  </thead>
  <tbody>
	{{ range .KS2 }}
	{{ $ks2 := . }}
	<tr>
	  <td>{{.}}</td>
	  {{ range $g }}
	  {{ with index $c $ks2 . }}
	  <td style="text-align:center; border:1px solid #888888"
	   {{ if ge .VA 0.67 }}class="success"{{ end }}
	   {{ if lt .VA -0.33 }}class="danger"{{ end }}
	   {{ if (ge .VA -0.33) and (lt .VA 0.67) }}class="warning"{{ end }}
	   {{ if eq (len .Students) 0 }}></td>
	  {{ else }}
	  data-container="body"
	  data-toggle="popover"
	  data-placement="right"
	  title="VA {{if $tm}}{{printf "%0.2f" .VA}}{{end}}"
	  data-html="true"
	  data-trigger="hover click"
	  data-content="{{range .Students}}
	  <a href='/students/{{.UPN}}/?{{$q}}' target='_blank'>{{.Name}}</a><br>
	  {{end}}"
	  >{{len .Students}}</td>
	{{ end }}
	{{ end }}
	{{ end }}
	<td style="text-align:center;">{{printf "%+.2f" (index $va .)}}</td>
	</tr>
	{{ end }}
  </tbody>
  <tfoot>
	<th>Total</th>
	{{ $counts := .Counts }}
	{{ range $g }}<th style="text-align:center;">{{ index $counts . }} </th>{{ end }}
	<th style="text-align:center;">{{printf "%+.2f" .TotalVA }}</th>
  </tfoot>
</table>

{{ end }}

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
	<th style="text-align:center;">PP</th>
	<th style="text-align:center;">KS2</th>
	<th style="text-align:center;">Grade</th>
	<th style="text-align:center;">Effort</th>
	<th style="text-align:center;">Value Added</th>
	<th style="text-align:center;">Attendance</th>
  </thead>
  <tbody>
	{{range .Students}}
	<tr>
	  <td><a href="/students/{{.UPN}}/?{{$q}}">{{.Name}}</a></td>
	  <td>{{.Class}}</td>
	  <td style="text-align:center;">
		{{if .PP}}<span class="glyphicon glyphicon-ok" style="color: #009933;"></span>
		{{else}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
		{{end}}
	  </td>
	  <td style="text-align:center;">{{.KS2}}</td>
	  <td style="text-align:center;">{{.Grade}}</td>
	  <td style="text-align:center;">{{.Effort}}</td>
	  <td style="text-align:center;">{{if .VAExists}}{{printf "%.1f" .VA}}{{end}}</td>
	  <td style="text-align:center;">{{printf "%.1f" .Attendance}}</td>
	</tr>
	{{end}}
  </tbody>
</table>
`,

"select-class.tmpl" : `
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

"select-level.tmpl" : `
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

"select-subject.tmpl" : `
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

"student.tmpl" : `
{{ $nat := .Nat }}
{{ with .Student }}
<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-5">
	<h3>{{.Forename}} {{.Surname}}</h3>

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
		  <td>{{if .PP}}<span class="glyphicon glyphicon-ok" style="color: #009933;"></span>
			{{else}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
			{{end}}</td>
		</tr>
		<tr>
		  <td>SEN</td>
		  <td>{{with .SEN}}
			{{if eq .Status "N"}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
			{{else}}<span class="glyphicon glyphicon-ok" style="color: #009933;"></span>
			{{end}}
			{{end}}</td>
		</tr>
		<tr>
		  <td>Ethnicity</td>
		  <td>{{.Ethnicity}}</td>
		</tr>
		<tr>
		  <td>EAL</td>
		  <td>{{if .EAL}}<span class="glyphicon glyphicon-ok" style="color: #009933;"></span>
			{{else}}<span class="glyphicon glyphicon-remove" style="color: #cc0000;"></span>
			{{end}}</td>
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

<div class="row" style="min-height:700px;">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<ul class="nav nav-tabs">
	  <li class="active"><a href="#ks2" data-toggle="tab" aria-expanded="true">KS2</a></li>
	  <li><a href="#sen" data-toggle="tab" aria-expanded="false">SEN</a></li>
	  <li><a href="#attendance" data-toggle="tab" aria-expanded="false">Attendance</a></li>
	  <li><a href="#results" data-toggle="tab" aria-expanded="false">Latest Assessments</a></li>
	  {{ if eq .Year 10 11 }}
	  <li class="disabled"><a href="#headlines" data-toggle="tab" aria-expaanded="false">Headlines</a></li>
	  <li><a href="#progress8" data-toggle="tab" aria-expaanded="false">Progress 8</a></li>
	  {{ end }}
	  <!--
	 <li class="dropdown">
	 <a class="dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="false">Subjects <span class="caret"></span></a>
	 <ul class="dropdown-menu">
	 <li class=""><a href="#Mathematics" data-toggle="tab" aria-expanded="false">Mathematics</a></li>
	 <li class=""><a href="#English" data-toggle="tab" aria-expanded="false">English</a></li>
	 </ul>
	 </li>
	  -->
	</ul>

	<br>

	<div id="myTabContent" class="tab-content">
	  <div class="tab-pane fade active in" id="ks2"> 
		<table class="table table-hover">
		  <tbody>
			{{with .KS2}}
			<tr>
			  <td>Average Point Score</td>
			  <td>{{.APS}}</td>
			</tr>
			<tr>
			  <td>Band</td>
			  <td>{{.Band}}</td>
			  </td>
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
				<tr>
				  {{end}}
		  </tbody>
		</table>
	  </div>

	  <div class="tab-pane fade" id="sen">
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
			  <td>Information</td>
			  <td>{{.Info}}</td>
			</tr>
			<tr>
			  <td>Strategies</td>
			  <td>{{.Strategies}}</td>
			</tr>
			<tr>
			  <td>Access Arrangements</td>
			  <td>{{.Access}}</td>
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
	  </div>

	  <div class="tab-pane fade" id="results">
		<table class="table table-condensed table-hover sortable">
		  <thead>
			<th>Subject</th>
			<th>Level</th>
			<th>Grade</th>
			<th>Effort</th>
			<th>Class</th>
			<th>Teacher</th>
			<tbody>
			  {{range .Courses}}
			  <tr>
				<td>{{.Subj}}</td>
				<td>{{.Lvl}}</td>
				<td>{{.Grd}}</td>
				<td>{{.Effort}}</td>
				<td>{{.Class}}</td>
				<td>{{.Teacher}}</td>
			  </tr>
			  {{end}}
			</tbody>
		</table>
	  </div>

	  {{ if eq .Year 10 11 }}
	  <div class="tab-pane fade" id="headlines">
		Headlines!
	  </div>

	  <div class="tab-pane fade" id="progress8">
		{{ with .Basket }}
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
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>English</td>
			  {{with (index .Slots 1)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>Mathematics</td>
			  {{with (index .Slots 2)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>Mathematics</td>
			  {{with (index .Slots 3)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>EBacc 1</td>
			  {{with (index .Slots 4)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>EBacc 2</td>
			  {{with (index .Slots 5)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>EBacc 3</td>
			  {{with (index .Slots 6)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>Other 1</td>
			  {{with (index .Slots 7)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>Other 2</td>
			  {{with (index .Slots 8)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
			<tr>
			  <td>Other 3</td>
			  {{with (index .Slots 9)}}
			  <td>{{.Subj}}</td>
			  <td>{{.Grade}}</td>
			  <td>{{.Points}}</td>
			  {{end}}
			</tr>
		  </tbody>
		  <tfoot>
			{{ with .Progress8 $nat }}
			<tr>
			  <th>Total</th>
			  <th>{{.EntN}} Entries</th>
			  <th></th>
			  <th>{{printf "%2.1f" .Ach}}</th>
			</tr>
			<tr>
			  <th>Expected Score</th>
			  <th></th>
			  <th></th>
			  <th>{{printf "%2.1f" .Exp}}</th>
			</tr>
			<tr>
			  <th>Progress 8 Score</th>
			  <th></th>
			  <th></th>
			  <th>{{printf "%2.1f" .Pts}}</th>
			</tr>
			{{ end }}
		  </tfoot>
		</table>
		{{ end }}
	  </div>
	  {{ end }}

	  <div class="tab-pane fade" id="attendance">
		{{with .Attendance}}
		<br>
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
			  <td style='text-align:center;vertical-align:middle'>{{printf "%.1f" .Latest}}</td>
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
	  </div>
	</div>
  </div>
  <div class="col-sm-1"></div>
</div>
{{ end }}


`,

"studentsearch.tmpl" : `
<h2>Student Search Results</h2>

<p>Searching for: <b>{{.Name}}</b></p>

<div class="row">
  <div class="col-sm-3"></div>
  <div class="col-sm-6">
	{{if .Results}}
	{{ $q := .Query }}
	<ul class="list-group">
	  {{range .Students}}
	  <a class="list-group-item" href="/students/{{.UPN}}/?{{$q}}">{{.Name}}	 ({{.Year}} {{.Form}})</a>
	  {{end}}
	</ul>
	{{else}}
	<p class="text-info">No results found.</p>
	{{end}}
  </div>
  <div class="col-sm-3"></div>
</div>

<br>
<p>If you can't find the student you are looking for, try using
just their surname or forename.<br>
If you're unsure of the spelling, a '*' can be used to replace
one or more characters.</p>
`,

"subject-overview.tmpl" : `
<h2>Subject Summaries</h2>
<br>
<div class="row">
  <div class="col-sm-1"></div>
  <div class="col-sm-10">
	<table class="table table-striped table-hover">
	  <thead>
		<th>Subject</th>
		<th style="text-align:center;">Cohort</th>
		<th style="text-align:center;">PP</th>
		<th style="text-align:center;">KS2 APS</th>
		<th style="text-align:center;">KS4 APS</th>
		<th style="text-align:center;">Av. Grade</th>
		<th style="text-align:center;">Value Added</th>
	  </thead>
	  <tbody>
		{{ $y := .Year }}
		{{ $q := .Query }}
		{{ range .Subjects }}
		<tr>
		  <td><a href="/progressgrid/{{ .Subject.Subj }}/{{ .Subject.SubjID }}/All Year {{ $y }}/?{{ $q}}">{{ .Subject.Subj }}</a></td>
		  <td style="text-align:center;">{{ .Cohort }}</td>
		  <td style="text-align:center;">{{ printf "%.0f" .PP }}%</td>
		  <td style="text-align:center;">{{ printf "%.1f" .KS2 }}</td>
		  <td style="text-align:center;">{{ printf "%.1f" .Points }}</td>
		  <td style="text-align:center;">{{ .AvGrade }}</td>
		  <td style="text-align:center;">{{ printf "%+.2f" .VA }}</td>
		</tr>
		{{ end }}
	</table>
  </div>
  <div class="col-sm-1"></div>
</div>
`,

}
