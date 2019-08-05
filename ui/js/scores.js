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
