window.onload = refreshScores()

function refreshScores() {
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
