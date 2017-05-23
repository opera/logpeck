function listServers() {
	console.log('listPeckers');
	$.ajax({
		url: '/list-peckers',
		type: 'POST',
		success: function (response) {
		var serverHtml = '<table id="pecker-list"  class="table"> <thead><tr>'
			+ '<th>Server Address</th>'
			+ '</tr></thead><tbody>';
		$.each(response, function (name, val) {
		 	serverHtml += '<tr>'
			 	+ '<td><span style="display:inline-block;width:500px">'
				+   '<a href="/pecker?addr=' + name + '">' + name + '</a>'
				+ '</span>'
				+ '<button class="btn btn-default" type="button" style="width:80px" onclick="removeServer(\''+name+'\')">Remove</button></td>'
			 	+ '</tr>';
	 	});
		serverHtml += '<tr>'
			+ '<td><span style="display:inline-block;width:500px">'
		 	+ '<input id="add-pecker" type="text" size="50" placeholder="" style="margin-right: 30px;">'
			+ '</span>'
			+ '<button class="btn btn-default" type="button" style="width:80px" onclick="addServer()">Add</button></td>'
		 	+ '</tr>';
	 	serverHtml += '</tbody></table>';
	 	$('#pecker-list').html(serverHtml);
	 	},
	 	error: function (error) {
		 	console.log(error);
	 	}
 	});
}

function addServer() {
	var addr = document.getElementById('add-pecker').value;
	if (addr.length < 10) {
		alert("addr error: " + addr);
		return;
	}
	console.log("add pecker: " + addr);
 	$.ajax({
	 	url: '/add-pecker?addr=' + addr,
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

function removeServer(addr) {
	console.log("remove pecker: " + addr);
 	$.ajax({
	 	url: '/remove-pecker?addr=' + addr,
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
