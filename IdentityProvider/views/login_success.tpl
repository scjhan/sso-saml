<!DOCTYPE HTML>
<html>
<head>
    <script src="static/js/jquery-3.3.1.min.js"></script>
    <script type="text/javascript">
        function CreateSession() {
            var domains = {{.Domains}};
            for (i = 0; i < domains.length; i++) {
                $.ajax({                
                    url: domains[i],
                    type: "GET",
                    dataType: "jsonp"
                });
            }
        }
        CreateSession();
    </script>
    <meta http-equiv="refresh" content="0;url=http://www.baidu.com/">
</head>
<body>
</body>
</html>