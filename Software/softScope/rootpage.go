package main

// Serves the scope's main page.

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, page)
}


type Update struct{
	ID string
	Attr string
	Value interface{}
}

var nrx = 0

func rxHandler(w http.ResponseWriter, r *http.Request) {
	nrx++
	var updates []Update
	updates = append(updates, Update{"NRX", "innerHTML", nrx})
	check(json.NewEncoder(w).Encode(updates))
}

const TX_MAGIC = 1234567

func txHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[(len("/tx/")):]
	split := strings.SplitN(url, "/", 2)
	cmd := split[0]
	val, err := strconv.Atoi(split[1])
	if err != nil {
		log.Println(err)
	}
	fmt.Println("tx", cmd, val)

	// todo: lock tty
	serial.writeInt(TX_MAGIC)
	serial.writeInt(binCommand(cmd))
	serial.writeInt(uint32(val))
}

var binCmds = map[string]uint32{
	"nop": 0,
	"samples":  1,
	"timebase": 2,
	"triglev":  3,
	"softgain": 4 }

func binCommand(cmd string) uint32 {
	if bin, ok := binCmds[strings.ToLower(cmd)]; ok {
		return bin
	} else {
		log.Println("unknown command:", cmd)
		return 0
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const page = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<title>SoftScope</title>

	<link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.4.2/pure-min.css">
	<style media="all" type="text/css">

		body  { margin-left: 5%; margin-right:5%; font-family: sans-serif; }
		table { border-collapse: collapse; }
		hr    { border-style: none; border-top: 1px solid #CCCCCC; }
		a     { color: #375EAB; text-decoration: none; }
	</style>

	<script>

var tick = 500;

// wraps document.getElementById, shows error if not found
function elementById(id){
	var elem = document.getElementById(id);
	if (elem == null){
		alert("undefined: " + id);
		return null;
	}
	return elem;
}

function setAttr(id, attr, value){
	var elem = elementById(id);
	if (elem[attr] == null){
		alert("settAttr: undefined: " + elem + "[" + attr + "]");
		return;
	}
	elem[attr] = value;
}

var pending = false;

// onreadystatechange function for update http request.
// updates the DOM with new values received from server.
function onReqReady(req){
	if (req.readyState == 4) { // DONE
		if (req.status == 200) {	
			var resp = JSON.parse(req.responseText);	
			setAttr("Samples", "value", resp["Samples"]);
			setAttr("TimeBase", "value", resp["TimeBase"]);
			setAttr("TrigLev", "value", resp["TrigLev"]);
			setAttr("SoftGain", "value", resp["SoftGain"]);
			setAttr("screen", "src", resp["Screen"]); 
		} 
	pending = false;
	}
}

function refresh(){
	if(pending){
		return;
	}
	pending = true;
	var req = new XMLHttpRequest();
	req.open("GET", document.URL + "/event", true); 
	req.timeout = 2*tick;
	req.onreadystatechange = function(){ onReqReady(req) };
	req.setRequestHeader("Content-type","application/x-www-form-urlencoded");
	req.send("");
}

setInterval(refresh, tick);

function val(id){
	return elementById(id).value
}

function upload(id){
	var req = new XMLHttpRequest();
	req.open("GET", document.URL + "tx/" + id + "/" + val(id), false);
	req.send("");
}

	</script>

</head>

<body>
	
<h1><i>Soft</i>Scope</h1>

<div style="padding-top:2em;">
	<table>
		<tr> <td><b> NRX      </b></td> <td> <span id="NRX"> </span>   </td></tr>
	</table>
</div>


<div>
	<img id="screen" src="/screen.svg" />
</div>

<div style="padding-top:2em;">
	<table>
		<tr> <td><b> Samples  </b></td> <td> <input id=Samples  type=number min=1 value=512           onchange="upload('Samples') ;"></td></tr>
		<tr> <td><b> TrigLev  </b></td> <td> <input id=TrigLev  type=number min=0 value=2000 max=4096 onchange="upload('TrigLev') ;"></td></tr>
		<tr> <td><b> TimeBase </b></td> <td> <input id=TimeBase type=number min=1 value=100           onchange="upload('TimeBase');"></td></tr>
		<tr> <td><b> SoftGain </b></td> <td>-<input id=SoftGain type=number min=0 value=2             onchange="upload('SoftGain');"></td></tr>
	</table>
</div>

</body>
</html>
`
