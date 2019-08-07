window.onload = function() {
    refresh_scores()
    connect()

    ctx = document.getElementById('graphc').getContext('2d');
    window.graph = new Chart(ctx, {
	// The type of chart we want to create
	type: 'bar',

	// The data for our dataset
	data: {
	    datasets: [{
		backgroundColor: '#616161',
		borderColor: '#9e9e9e',
	    }]
	},

	// Configuration options go here
	options: {
	    legend: {
		display: false,
	    },
	    responsive: true,
	    maintainAspectRatio: false,
	    scales: {
		yAxes: [{
		    ticks: {
			beginAtZero: true,
			suggestedMax: 50,
		    }
		}],
		xAxes: [{
		    gridLines: {
			display: false
		    }
		}],
	    },
	}
    });

    ctx = document.getElementById('cur_graph_c').getContext('2d');
    window.cur_graph = new Chart(ctx, {
	// The type of chart we want to create
	type: 'bar',

	// The data for our dataset
	data: {
	    labels: [""],

	    datasets: [{
		backgroundColor: '#616161',
		borderColor: '#9e9e9e',
	    }]
	},

	// Configuration options go here
	options: {
	    legend: {
		display: false,
	    },
	    responsive: true,
	    maintainAspectRatio: false,
	    scales: {
		yAxes: [{
		    ticks: {
			beginAtZero: true,
			suggestedMax: 10,
			suggestedMin: -10
		    }
		}],
		xAxes: [{
		    gridLines: {
			display: false
		    }
		}],
	    },
	}
    });

    ctx = document.getElementById('next_graph_c').getContext('2d');
    window.next_graph = new Chart(ctx, {
	// The type of chart we want to create
	type: 'bar',

	// The data for our dataset
	data: {
	    datasets: [{
		backgroundColor: '#616161',
		borderColor: '#9e9e9e',
	    }]
	},

	// Configuration options go here
	options: {
	    legend: {
		display: false,
	    },
	    responsive: true,
	    maintainAspectRatio: false,
	    scales: {
		yAxes: [{
		    ticks: {
			beginAtZero: true,
			suggestedMax: 10,
			suggestedMin: -10
		    }
		}],
		xAxes: [{
		    gridLines: {
			display: false
		    }
		}],
	    },
	}
    });

    ctx = document.getElementById('max_graph_c').getContext('2d');
    window.max_graph = new Chart(ctx, {
	// The type of chart we want to create
	type: 'bar',

	// The data for our dataset
	data: {
	    labels: [""],
	    datasets: [{
		backgroundColor: '#616161',
		borderColor: '#9e9e9e',
	    }]
	},

	// Configuration options go here
	options: {
	    legend: {
		display: false,
	    },
	    responsive: true,
	    maintainAspectRatio: false,
	    scales: {
		yAxes: [{
		    ticks: {
			beginAtZero: true,
			suggestedMax: 10,
			suggestedMin: -10
		    }
		}],
		xAxes: [{
		    gridLines: {
			display: false
		    }
		}],
	    },
	}
    });
}

function refresh_scores() {
    const scores = ["classic", "unlimited"]
    scores.forEach(item => fetch("tfh/scores/" + item)
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
	    setTimeout(function() { live_finalized(json) }, 5000)
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

    s = document.getElementById("live_static")

    parent = document.getElementById("graph");
    ctx = document.getElementById('graphc').getContext('2d');

    ctx.width = parent.clientWidth;
    ctx.height = parent.clientHeight;

    window.graph.data.labels = []
    window.graph.data.datasets[0].data = []

    cur_rank = "None"
    cur_score = "None"
    window.cur_score = 0
    if (json.cur_rank >= 0) {
	cur_rank = json.cur_rank + 1
	cur_score = json.ranks[json.cur_rank]
	cur_score = Math.round(cur_score * 1000) / 1000
        window.cur_score = cur_score
    }
    document.getElementById("cur_rank").innerText = cur_rank
    document.getElementById("cur_score").innerText = cur_score

    // set up graph
    parent = document.getElementById("cur_graph");
    ctx = document.getElementById('cur_graph_c').getContext('2d');

    ctx.width = parent.clientWidth;
    ctx.height = parent.clientHeight;

    window.graph.data.labels = []
    window.graph.data.datasets[0].data = []

    next_rank = "None"
    next_score = "None"
    window.next_score = 0
    if (json.cur_rank >= 0) {
	next_rank = Math.max(json.cur_rank, 1)
	next_score = json.ranks[Math.max(json.cur_rank-1, 0)]
	next_score = Math.round(next_score * 1000) / 1000
        window.next_score = next_score
    }
    document.getElementById("next_rank").innerText = next_rank
    document.getElementById("next_score").innerText = next_score

    max_score = "None"
    window.max_score = 0
    if (json.ranks.length > 0) {
	max_score = json.ranks[0]
	max_score = Math.round(max_score * 1000) / 1000
        window.max_score = max_score
    }
    document.getElementById("max_rank").innerText = "1"
    document.getElementById("max_score").innerText = max_score

    window.max_score = max_score
}

function live_updated(json) {
    window.graph.data.labels.push("")
    window.graph.data.datasets[0].data.push(json.score)
    window.graph.update()

    red = "#810000"
    green = "#008100"

    cur_diff = json.score - window.cur_score
    window.cur_graph.data.datasets[0].data = [cur_diff]
    window.cur_graph.data.datasets[0].backgroundColor = cur_diff > 0 ? green : red
    window.cur_graph.update()

    next_diff = json.score - window.next_score
    window.next_graph.data.datasets[0].data = [next_diff]
    window.next_graph.data.datasets[0].backgroundColor = next_diff > 0 ? green : red
    window.next_graph.update()

    max_diff = json.score - window.max_score
    window.max_graph.data.datasets[0].data = [max_diff]
    window.max_graph.data.datasets[0].backgroundColor = max_diff > 0 ? green : red
    window.max_graph.update()

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
