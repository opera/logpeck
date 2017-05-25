function listPeckTasks() {
	var addr = document.getElementById('pecker-address').textContent;
	console.log('listPeckTasks: ' + addr);
	$.ajax({
		url: '/list-pecktasks?addr=' + addr,
		type: 'POST',
		success: function (response) {
		var serverHtml = '<table id="pecktask-list"  class="table"> <thead><tr>'
			+ '<th>Name</th>'
			+ '<th>Log Path</th>'
			+ '<th>Status</th>'
			+ '<th>Type</th>'
			+ '</tr></thead><tbody>';
		$.each(response, function (name, val) {
		 	serverHtml += '<tr>'
			 	+ '<td>' + name + '</td>'
			 	+ '</tr>';
	 	});
	 	serverHtml += '</tbody></table>';
	 	$('#pecktask-list').html(serverHtml);
	 	},
	 	error: function (error) {
		 	console.log(error);
	 	}
 	});
}


