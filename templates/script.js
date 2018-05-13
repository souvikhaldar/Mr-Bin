var request = new XMLHttpRequest();

// Open a new connection, using the GET request on the URL endpoint
request.open('GET', 'https://lit-sea-89877.herokuapp.com/getPercent', true);

request.onload = function () {
    var data = JSON.parse(this.response);
	if (request.status >= 200 && request.status < 400) {
	console.log(data.Percentage);

  }else{
	  console.log('error');
  }
}

// Send request
request.send();