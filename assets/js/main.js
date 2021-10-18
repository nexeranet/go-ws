var url = 'ws://localhost:8080/ws';
var c = new WebSocket(url);

var send = function(data){
  c.send(data)
}

c.onmessage = function(msg){
  console.log(msg)
}

c.onopen = function(){
  setInterval( 
    function(){ send("ping") }
    , 1000 )
}
