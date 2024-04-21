let button = document.querySelector(".centers");
if (button) {
    button.onclick = function (e) {
        let inputs = document.querySelectorAll(".tables > input");
        let data = {};
        for (let i = 0; i < inputs.length; i++) {
            data[inputs[i].name] = inputs[i].value;
        }
        let xhr = new XMLHttpRequest();
        xhr.open("POST", "/spravochnik/edg/" + data["ID"]);
        xhr.onload = function (e) {
            let response = JSON.parse(e.currentTarget.response);
            if (response.Error == null) {
                console.log("Пользователь успешно зарегистрирован");
                window.location.url = "http://127.0.0.1:8080/spravochnik/update/" + data["ID"];
                //$("#centersH").html( "Данные успешно сохранены");

            } else {
                console.log(response.Error);
            }
        };
        xhr.send(JSON.stringify(data));
    }
}

let buttonDel = document.querySelector(".delete");
if (buttonDel) {
    buttonDel.onclick = function (e) {
        let inputs = document.querySelectorAll(".tables > input");
        let data = {};
        for (let i = 0; i < inputs.length; i++) {
            data[inputs[i].name] = inputs[i].value;
        }
        let xhrr = new XMLHttpRequest();
        xhrr.open("POST", "/spravochnik/delete", true)
        xhrr.onload = function (e) {
            if (xhrr.status !== 200) {
                return;
            }
            let response = JSON.parse(e.currentTarget.response);
            if ("Error" in response) {
                if (response.Error == null) {
                    //console.log("Изменения успешно внесены");
                    //window.location.assign("http://127.0.0.1:8080/");
                } else {
                    console.log(response.Error);
                }
            } else {
                console.log("Некорректные данные");
            }
        };
        xhrr.send(JSON.stringify(data));
        console.log(JSON.stringify(data));
    }
}