package templates

func GetCommonHtmlTpl() string {
	return `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>%s</title>
		<link rel="stylesheet" type="text/css" href="/dist/styles/common.css?<|NONCE|>">
	</head>
	<body>
		%s
	</body>
	</html>
	`
}

func GetMenuIndexPageTpl() string {
	return `
    <!DOCTYPE html>
    <html>
    	<head>
    		<title>
    			Index of " %s "</title>
    		<meta name="viewport" content="width=device-width, initial-scale=1">
    		<link rel="stylesheet" href="/menu.css">
    	</head>

    	<body>
	`
}

func GetMenuTpl() string {
	return `
	<div class="menu">
	  <input type="checkbox" id="menu-toggle" class="menu-toggle"/>
	  <label for="menu-toggle" class="menu-button">Menu</label>
	  <nav class="menu-nav">
	    <ul class="menu-list">
	    </ul>
	  </nav>
	</div>
	`
}

func getMenuItemTpl() string {
	return `
		<li class="menu-item has-children">
			<a href="%s/index.html" class="menu-link">
				%s
			</a>
			<ul class="menu-dropdown">
	`
}