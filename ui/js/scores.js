window.onload = function() {
    refresh_scores()
    connect()
}

function refresh_scores() {
    const scores = ["classic", "unlimited"]
    scores.forEach(item => fetch("scores/" + item)
		   .then(function(resp) {
		       return resp.json()
		   })
		   .then(function(json) {
		       update_scores(item, json)
		   })
		  )
}

function update_scores(item, json) {
    c = document.getElementById(item)
    // remove old scores
    old_scores = c.getElementsByClassName("row");
    for (i = old_scores.length-1; i >= 0; i--) {
	// don't delete headers
	if (old_scores[i].className.indexOf("hdr") === -1) {
	    c.removeChild(old_scores[i])
	}
    }

    tmpl = document.getElementById("score")
    json.scores.forEach((s, idx) => {
	r = tmpl.content.cloneNode(true)
	r.querySelector(".rank").innerText = idx+1
	r.querySelector(".name").innerText = s.name
	r.querySelector(".score").innerText = s.score
	r.querySelector(".attempts").innerText = s.attempts
	c.appendChild(r)
    })
}

function connect() {
    ws = new WebSocket("ws://" + location.host + "/live")
    ws.onopen = function(e) {
	console.log("Websocket connected")
    }

    ws.onmessage = function(e) {
	json = JSON.parse(e.data)
	if (json.id === "started") {
	    live_started(json)
	}
	else if (json.id === "updated") {
	    live_updated(json)
	} else if (json.id === "finalized") {
	    live_finalized(json)
	} else {
	    console.log("Unknown websocket message: ", json)
	}
    }

    ws.onclose = function(e) {
	console.log("Lost connection, reconnecting in 1 second: ", e.reason)
	setTimeout(function() {
	    connect()
	}, 1000);
    }

    ws.onerror = function(e) {
	console.error("Websocket error: ", e)
	ws.close()
    }
}

function live_started(json) {
    console.log("Started: "+JSON.stringify(json))
    live = document.getElementById("live")
    scores = document.getElementById("scores")

    scores.style.display = "none"
    live.style.display = "block"

    nm = document.getElementById("live_name")
    nm.innerText = json.name

    cls = document.getElementById("live_class")
    cls.innerText = json.class

}

function live_updated(json) {
    console.log("Updated: "+JSON.stringify(json))
}

function live_finalized(json) {
    console.log("Finalized: "+JSON.stringify(json))
    live = document.getElementById("live")
    scores = document.getElementById("scores")

    scores.style.display = "block"
    live.style.display = "none"

    refresh_scores()
}
