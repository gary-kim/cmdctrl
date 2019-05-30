package server

// ControlUITemplate is the web code for cmdctrl's webui
var ControlUITemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>cmdctrl webui</title>
    <style>
        input {
            height: 2em;
            background: white;
        }
    </style>
</head>
<body>
    <div>
        <p>Clients</p>
        <div id="clients"></div>
    </div>
    <div id="controls">
        <form action="" method="post">
            <select name="action" id="input-action">
                <option value="math">Math</option>
                <option value="cmd">cmd</option>
            </select>
            <br>
                <label>Priority: <input type="number" name="priority" id="input-priority"></label><br>
                <label>Client: <input type="text" name="client" id="input-client"></label><br>
                <label>cmd: <input type="text" name="input" id="input-input"></label><br>
            <input type="submit" value="submit">
        </form>
    </div>
    <script>
        async function loadClients() {
            let list = await (await fetch('/clients')).text();
            document.querySelector("#clients").innerText = list;
        }
        loadClients();
        setTimeout(loadClients, 1000);
    </script>
</body>
</html>
`
