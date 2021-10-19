var url = 'ws://localhost:8080/ws';
var c = new WebSocket(url);

function formatDate(date){
   return new Date(date).toLocaleString()
 }
var send = function(data){
  c.send(data)
}
var sendJson = function(obj) {
  c.send(JSON.stringify(obj))
}
c.onmessage = function(msg){
  let data = JSON.parse(msg.data)
  let el = document.getElementById(data.symbol);
  if (el != null ){ 
    let color = el.children[1].style.color;
    currPrice =  Number(el.children[1].innerHTML)
    if (data.price > currPrice) {
      color = "green"
    }  else if (data.price < currPrice) {
      color = "red"
    }
    el.children[1].innerHTML = data.price;
    el.children[2].innerHTML = formatDate(data.timestamp);
    el.children[1].style.color = color
  } else {
    let newEL = document.createElement("div");
    newEL.id = data.symbol;
    newEL.dataset.id  = data.symbol;
    newEL.appendChild(document.createElement("span"));
    newEL.appendChild(document.createElement("span"));
    newEL.appendChild(document.createElement("span"));
    newEL.children[0].className = "symbol"
    newEL.children[1].className = "price"
    newEL.children[2].className = "timestamp"
    newEL.children[0].innerHTML = data.symbol;
    newEL.children[1].innerHTML = data.price;
    newEL.children[2].innerHTML = formatDate(data.timestamp);
    document.querySelector("#list").appendChild(newEL);
  }
  //document.
}

spButton = document.querySelector("#subs")
if (spButton != null) {
  spButton.addEventListener("click", function(e){
    console.log("subs")
    sendJson({"action": "subscribe", "symbols": []})
  })
}

spButton = document.querySelector("#unsubs")
if (spButton != null) {
  spButton.addEventListener("click", function(e){
    console.log("unsub")
    sendJson({"action": "unsubscribe"})
    document.querySelector("#list").innerHTML = ""
  })
}
