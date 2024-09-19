package report

var ReportStart = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>The sdakit Project | Summary</title>
    <style>
        html, body {
            height: 100%;
            background-color: white;
        }
        html {
            display: table;
            margin: auto;
        }
        body {
            display: table-cell;
            vertical-align: text-top;
        }
        #head {
            text-align: center; 
            margin: 20px; 
        }
        #logo-headline {
            width: 40%;
            height: auto; 
        }
        #headline {
            font-family: 'Courier New', Courier, monospace;
        }
        #details {
            margin-top: 25px;
            text-align: center;
            font-family: 'Courier New', Courier, monospace;
        }
        #overview-table {
            background-color: white;
            border-collapse: collapse;
            width: 100%;
        }
        #table-container {
            margin-top: 60px;
            margin-bottom: 60px;
        }
        #overview-table td, #overview-table th {
            border: 2px solid #000000;
            padding: 8px;
        }
        #overview-table th {
            padding-top: 12px;
            padding-bottom: 12px;
            text-align: left;
            background-color: rgb(235, 17, 17);
            color: white;
            font-family: 'Courier New', Courier, monospace;
        }
        #category {
            font-family: 'Courier New', Courier, monospace;
            margin-left: 2%;
        }
        #category-headline {
            font-family: 'Courier New', Courier, monospace;
            text-decoration: underline;
            margin-top: 50px;
        }
        ol li {
            font-family: 'Courier New', Courier, monospace;
            margin-left: 7%;
        }
    </style>
</head>
<body>
    <div id="head">
        <img id="logo-headline" src="https://github.com/PlagueByteSec/sdakit-project/blob/main/assets/sdakit-logo.png?raw=true" alt="logo">
        <h1 id="headline">Report</h1>
    </div>
    <hr style="border-color: rgb(214, 0, 0);">    
	<table id="overview-table">
        <tr>
            <th>Domain</th>
            <th>Date</th>
            <th>Method</th>
        </tr>`

var ReportEnd = `<hr style="border-color: rgb(214, 0, 0);">
    <div id="details">
		<h5>This report was generated automatically.</h5>
    </div>
</body>
</html>`
