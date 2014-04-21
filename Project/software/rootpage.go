package softscope

// Serves the scope's main page.

import (
	"fmt"
	"net/http"
	"strings"
	//"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, page)
}

type jsCall struct {
	F    string        // function to call
	Args []interface{} // function arguments
}

const TX_MAGIC = 1234567

func txHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[(len("/tx/")):]
	split := strings.SplitN(url, "/", 2)
	cmd := split[0]
	val := atouint32(split[1])
	fmt.Println("tx", cmd, val)
	switch cmd {
	default:
		panic(cmd)
	case "samples":
		SendMsg(SET_SAMPLES, val)
	case "timebase":
		SendMsg(SET_TIMEBASE, val)
	case "triglev":
		SendMsg(SET_TRIGLEV, val)
	case "reqFrames":
		SendMsg(REQ_FRAMES, val)
	}
}

//var binCmds = map[string]uint32{
//	"nop":      0,
//	"samples":  1,
//	"timebase": 2,
//	"triglev":  3,
//	"softgain": 4}

//func binCommand(cmd string) uint32 {
//	if bin, ok := binCmds[strings.ToLower(cmd)]; ok {
//		return bin
//	} else {
//		log.Println("unknown command:", cmd)
//		return 0
//	}
//}

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

var tick = 100;

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
			for(var i=0; i<resp.length; i++){
				var r = resp[i];
				var func = window[r.F];
				if (func == null) {
					showErr("undefined: " + r.F);
				}else{
					func.apply(this, r.Args);
				}
			}
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
	req.open("GET", document.URL + "/rx", true);
	//req.timeout = 2*tick;
	req.onreadystatechange = function(){ onReqReady(req) };
	req.setRequestHeader("Content-type","application/x-www-form-urlencoded");
	req.send("");
}

setInterval(refresh, tick);

function val(id){
	return elementById(id).value
}

function message(id, value){
	var req = new XMLHttpRequest();
	req.open("GET", document.URL + "tx/" + id + "/" + value, false);
	req.send("");
}

function upload(id){
	message(id, val(id));
}

</script>

</head>

<body>

<h1><i>Soft</i>Scope</h1>

<span id="errorBox"> &nbsp; </span>


<div>
	<table><tr>
	<td> <img id="screen" height=265 src="/screen.svg" />  </td>
	<td> <pre style="font-size:0.7em;" ><span id="FrameDebug"> Waiting for frame </span></pre> </td>
	</tr></table>
</div>
	<table>
		<tr> <td><b>samples<b></td> <td>  <input type=range id="samples"  min=0 max=4096 step=32 value=512   onchange="upload('samples') ;" oninput="upload('samples') ;"  ></td></tr>
		<tr> <td><b>timebase<b></td> <td> <input type=range id="timebase" min=0 max=42000 step=6 value=420   onchange="upload('timebase');" oninput="upload('timebase');"  ></td></tr>
		<tr> <td><b>trigger<b></td> <td>  <input type=range id="triglev"  min=0 max=5000 step=16  value=420   onchange="upload('triglev') ;" oninput="upload('triglev') ;"  ></td></tr>
	</table>
<div>

<input type=button onclick="message('reqFrames', 1);" value="Req. Frame" /> 

</div>


<div style="padding-top:2em;">
	<table>
		<tr> <td><b> Screen refreshes  </b></td> <td> <span id="NRX">  </span>   </td></tr>
		<tr> <td><b> TTY Frames in     </b></td> <td> <span id="FrameRate">  </span> /s  </td></tr>
	</table>
</div>


</body>
</html>
`
