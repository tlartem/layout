package v1

//nolint:unused
const homePage = `<!doctype html>
<html lang="ru">
	<head>
		<meta charset="UTF-8" />
		<meta
			name="viewport"
			content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"
		/>
		<meta http-equiv="X-UA-Compatible" content="ie=edge" />
		<link rel="preload" href="/file/font/Roboto-Regular.ttf" as="font" type="font/ttf" crossorigin />
		<link rel="preload" href="/file/font/Roboto-Medium.ttf" as="font" type="font/ttf" crossorigin />
		<link id="theme-stylesheet" rel="stylesheet" href="/file/css/theme-dark.css" />
		<link rel="stylesheet" href="{{.PathCSS}}" />
		<title>{{.Title}}</title>
	</head>
	<body>
		{{.Body}}
		<script src="/file/js/main.js"></script>
	</body>
</html>`
