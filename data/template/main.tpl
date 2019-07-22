<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
	<title>PoE Stash</title>
  </head>
  <body>
    <b>Wealth:</b> {{ .Wealth }} <br />

    <div class="stash">
      {{ range $index, $element := .Stash}}
      <div class="tab"> {{print $index }} </div>
        <div class="stashtab">
          {{ range .Items }}
          <div class="item">
            {{ .Name }}<br />
            {{ .Type }}<br />
            <img src="{{ .Icon }}" />
          </div>
          {{end}}
        </div>
      </div>
      {{end}}
    </div>


  </body>
</html>
