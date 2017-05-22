function listServers() {
	console.log('listServers');
	$.ajax({
		url: '/list-servers',
		type: 'POST',
		success: function (response) {
		var serverHtml = '<table id="server-list"  class="table"> <thead><tr>'
			+ '<th>Name</th>'
			+ '</tr></thead><tbody>';
		$.each(response, function (name, val) {
		 	serverHtml += '<tr>'
			 	+ '<td><a href="/servers/' + name + '">' + name + '</a></td>'
			 	+ '</tr>';
	 	});
		serverHtml += '<tr>'
		 	+ '<td><input id="add-server" type="text" size="50" placeholder="" style="margin-right: 30px;"><button class="btn btn-default" type="button" onclick="addServer()">Add</button></td>'
		 	+ '</tr>';
	 	serverHtml += '</tbody></table>';
	 	$('#server-list').html(serverHtml);
	 	},
	 	error: function (error) {
		 	console.log(error);
	 	}
 	});
}

function addServer() {
	var addr = document.getElementById('add-server').value;
	console.log("add server: " + addr);
 	$.ajax({
	 	url: '/add-server?server_addr=' + addr,
	 	type: 'POST',
	 	success: function (response) {
			location.reload();
			console.log(response['status']);
	 	},
	 	error: function (error) {
		 	console.log(error);
	 	}
 	});
}
