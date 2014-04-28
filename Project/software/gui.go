package softscope

// Serves the scope's HTML GUI.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Start the GUI server on address (e.g. ":4000")
func RunHTTPServer(addr string) {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/screen.svg", screenHandler)
	http.HandleFunc("/event/", eventHandler)
	http.HandleFunc("/refresh/", refreshHandler)
	fmt.Println("listening on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Serves the root page content (see const page)
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, page)
}

// Serves image of the screen
func screenHandler(w http.ResponseWriter, r *http.Request) {
	ExecSync(func() {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-control", "No-Cache")
		w.Write(screenBuf.Bytes())
	})
}

// Called by javascript to notify us on events (clicks etc).  E.g.:
// 	http://localhost/event/samples/128
// is called set the number of samples to 128
func eventHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[(len("/event/")):]
	split := strings.SplitN(url, "/", 2)
	cmd := split[0]
	val := atouint32(split[1])
	debug("event", cmd, val)
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
	case "clearerr":
		SendMsg(CLEAR_ERR, val)
	}
}

// Called by javascript to refresh the page's dynamic content (settings readout etc.).
// Answers with a JSON struct telling the page what to refresh.
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	ExecSync(func() {
		nrx++
		calls := make([]jsCall, 0, 3)
		calls = append(calls, jsCall{"setAttr", []interface{}{"NRX", "innerHTML", nrx}})
		calls = append(calls, jsCall{"setAttr", []interface{}{"FrameDebug", "innerHTML", fmt.Sprint(frame.Header.String())}})
		calls = append(calls, jsCall{"setAttr", []interface{}{"screen", "src", "/screen.svg"}})
		//calls = append(calls, jsCall{"setAttr", []interface{}{"FrameRate", "innerHTML", frameRate}})
		check(json.NewEncoder(w).Encode(calls))
	})
}

// element in the JSON response by refreshHandler
type jsCall struct {
	F    string        // function to call
	Args []interface{} // function arguments
}

// html served by rootHandler
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
	req.open("GET", document.URL + "/refresh", true);
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
	req.open("GET", document.URL + "event/" + id + "/" + value, false);
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

<div>
	<table>
		<tr> <td><b>samples<b></td> <td>  <input type=range id="samples"  min=0 max=4096 step=32 value=512   onchange="upload('samples') ;" oninput="upload('samples') ;"  ></td></tr>
		<tr> <td><b>timebase<b></td> <td> <input type=range id="timebase" min=0 max=42000 step=6 value=420   onchange="upload('timebase');" oninput="upload('timebase');"  ></td></tr>
		<tr> <td><b>trigger<b></td> <td>  <input type=range id="triglev"  min=0 max=5000 step=16  value=420   onchange="upload('triglev') ;" oninput="upload('triglev') ;"  ></td></tr>
	</table>
	<input type=button value="Clear Err" onclick="message('clearerr', 0)"/>
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
