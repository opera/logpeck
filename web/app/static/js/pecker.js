function listPeckTasks() {
	var addr = document.getElementById('pecker-address').textContent;
	console.log('listPeckTasks: ' + addr);
	$.ajax({
		url: '/list-pecktasks?addr=' + addr,
		type: 'POST',
		success: function (response) {
		var serverHtml = '<table id="pecktask-list"  class="table table-bordered"> <thead><tr>'
			+ '<th>Name</th>'
			+ '<th>Status</th>'
			+ '<th>Log Path</th>'
			+ '<th>Filter Expression</th>'
			+ '<th>Delimiter</th>'
			+ '<th>ElasticSearch</th>'
			+ '</tr></thead><tbody>';
		$.each(response, function (name, val) {
			console.log(val);
			s = 'Running';
			es = val['ESConfig']['URL'] + '/' + val['ESConfig']['Index'] + '/' + val['ESConfig']['Type'];

		 	serverHtml += '<tr>'
			 	+ '<td>' + name + '</td>'
			 	+ '<td>' + s + '</td>'
			 	+ '<td>' + val['LogPath'] + '</td>'
			 	+ '<td>' + val['FilterExpr'] + '</td>'
			 	+ '<td>' + val['Delimiters'] + '</td>'
			 	+ '<td>' + es + '</td>'
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


