function showItems(id) {
    var Http = new XMLHttpRequest();
    const url = 'http://localhost:8080/categories/' + id;
    Http.open("GET", url);


    Http.onload = function (e) {
        var resp = Http.response
        document.getElementById("items").innerHTML = ""
        JSON.parse(resp).forEach(
            e => {
                document.getElementById("items").innerHTML += "    <div class=\"card\" style=\"width: 18rem;  margin: 10px;\">\n" +
                    "        <img class=\"card-img-top\" style=\"width:auto ; height: 18rem;\" src=" + e.ImageUrl + " alt=\"Card image cap\">\n" +
                    "        <div class=\"card-body\">\n" +
                    "            <h5 class=\"card-title\">Название: " + e.Name + "</h5>\n" +
                    "            <p class=\"card-text\">Цена: " + e.Price + " ₽</p>\n" +
                    "        </div>\n" +
                    "    </div>"
            }
        )
    }

    Http.send();
}

function setSubCategory(id) {
    var Http = new XMLHttpRequest();
    const url = 'http://localhost:8080/jsonCategories';
    Http.open("GET", url);

    Http.onload = function (e) {
        var resp = Http.response
        JSON.parse(resp).forEach(
            e => {
                if (e.id === id) {
                    document.getElementById("sub" + id).innerHTML += "                <a class=\"list-group-item\" onclick=\"showItems(" + id + ")\">\n" +
                        "                    "+e.Name+"\n" +
                        "                </a>"
                }
            }
        )
    }

    Http.send();
}
