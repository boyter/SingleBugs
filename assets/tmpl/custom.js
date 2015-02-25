$(document).ready(function () {
	var location = window.location.href;
	var projectTemplate = _.template('<li>' +
									 '<a class="project" href="<%- id %>">' +
									 '<%- name %>' +
									 '</a>' + 
									 '</li>');
	var issueTemplate   = _.template('<li>' +
									 '<a class="issue" href="<%- id %>">' +
									 '#<%- id %> <%- name %>' +
									 '</a>' + 
									 '</li>');
	var noteTemplate    = _.template('<div class="well" href="<%- id %>">' +
									 '<%= content %>' +
									 '</div>');

	var errorTemplate   = _.template('<div class="alert alert-error">' +
									 '<button type="button" class="close" data-dismiss="alert">&times;</button>' +
									 '<%- message %>' +
									 '</div>');
									 
	getAllProjects();
	getAllIssues();
	
	$("#search").focus();
	var currentIssue = -1;
	var currentProject = -1;
	
	$(document).on('click', '.issue', function (e) {
		e.preventDefault();
		var id = $(this).attr('href');
		
		$('#issues a').each( function(e) {
			$(this).removeClass('highlight-select');
		});
		
		$(this).addClass('highlight-select');
		
		issueClicked(id);
	});
	
	$(document).on('click', '.project', function (e) {
		e.preventDefault();
		var id = $(this).attr('href');
		
		
		$('#projectspeople a').each( function(e) {
			$(this).removeClass('highlight-select');
		});
		
		$(this).addClass('highlight-select');
		
		$('#issues').html('');
		$('#notes').html('');
		$('#noteForm').hide();
		projectClicked(id);
	});
	
	$('#includeclosed').change(function (e) {
		projectSearch();
		issueSearch();
	});
	
	$('#search').keyup(debounce(function(e) {
		$('#notes').html('');
		$('#issues').html('');
		$('#noteForm').hide();
		projectSearch();
		issueSearch();
	}, 200));
	
	$('#addcomment').click(function(e) {
		$.ajax({
			type: "POST",
			dataType: "json",
			url: location + 'savenote/',
			data: {
				issueid: currentIssue,
				body: $('#noteData').val()
			},
			success: function(e) {
				if(e.Success == true) {
					issueClicked(currentIssue);
					$('#noteData').val('');
				}
			}
		});
		
	});
	
	$('#createProject').click(function (e) {
		var projectName = $('#modalprojectname').val();
		
		$.ajax({
			type: "POST",
			dataType: "json",
			url: location + 'saveproject/',
			data: {
				projectname: projectName
			},
			success: function(e) {
				if(e.Success == true) {
					projectSearch();
					issueSearch();
					$('#modalprojectname').val('');
					$('#projectModal').modal('hide');
				}
				else {
					var tmp = errorTemplate({message:e.Message});
					$('#projectModalError').html(tmp);
				}
			}
		});
	});
	
	$('#closeIssue').click(function(e) {
		$('#confirmCloseIssueModal').modal();
	});
	
	$('#closeIssueConfirm').click(function(e) {
		$('#confirmCloseIssueModal').modal('hide');
		$.ajax({
			type: "POST",
			dataType: "json",
			url: location + 'closeissue/',
			data: {
				issueid: currentIssue
			},
			success: function(e) {
				if(e.Success == true) {
					currentIssue = -1;
					$('#notes').html('');
					$('#noteForm').hide();
					projectSearch();
					issueSearch();
					
				}
			}
		});
	});
	
	
	// In case they fixed the error for createProject method
	// without dismissing the message lets hide it now
	$('#projectModal').on('show', function () {
		$('#projectModalError').html('');
		
		$.getJSON(location + "allprojects/", function(data) {
			var items = [];

			$.each(data, function(key, val) {
				var id = val.Id;
				var name = val.Name;
				var selected = '';
				
				if(id == currentProject) {
					selected = 'selected';
				}
				items.push('<option value="'+id+'" '+selected+'>'+name+'</option>');
				
				selected = '';
			});

			$('#closeprojectid').html('' + items.join(''));
		});
		
	});
	
	
	$('#closeProject').click(function (e) {
		
		$.ajax({
			type: "POST",
			dataType: "json",
			url: location + 'closeproject/',
			data: {
				projectid: $('#closeprojectid').val()
			},
			success: function(e) {
				if(e.Success == true) {
					$('#projectModal').modal('hide');
					$('#myTab a:first').tab('show'); // Select first tab
					currentIssue = -1;
					currentProject = -1;
					$('#notes').html('');
					$('#noteForm').hide();
					projectSearch();
					issueSearch();
				}
			}
		});
	});
	
	
	$('#createIssue').click(function (e) {
		var projectid = $('#issueprojectid').val();
		var issuename = $('#issuename').val();
		var issuecontent = $('#issuecontent').val();
		
		$.ajax({
			type: "POST",
			dataType: "json",
			url: location + 'saveissue/',
			data: {
				projectid: projectid,
				issuename: issuename,
				issuecontent: issuecontent
			},
			success: function(e) {
				if(e.Success == true) {
					$("#search").val('#'+e.Id);
					projectSearch();
					issueSearch();
					$('#issuename').val('');
					$('#issuecontent').val('');
					$('#issueModal').modal('hide');
				}
				else {
					var tmp = errorTemplate({message:e.Message});
					$('#issueModalError').html(tmp);
				}
			}
		});
	});
	
	// In case they fixed the error for createIssue method
	// without dismissing the message lets hide it now
	$('#issueModal').on('show', function () {
		$('#issueModalError').html('');
	});
	
	

	// Get a list of the projects we have and add them to the drop
	// down. If we have a preference select that as well
	$('#issueModal').on('show', function () {
		var c = $('#includeclosed').is(':checked');
		
		$.getJSON(location + "allprojects/?c="+c, function(data) {
			var items = [];

			$.each(data, function(key, val) {
				var id = val.Id;
				var name = val.Name;
				var selected = '';
				
				if(id == currentProject) {
					selected = 'selected';
				}
				items.push('<option value="'+id+'" '+selected+'>'+name+'</option>');
				
				selected = '';
			});

			$('#issueprojectid').html('' + items.join(''));
		});
	
	});
	
	
	// Buffers requests so we only fire once on the delay
	// http://remysharp.com/2010/07/21/throttling-function-calls/
	function debounce(fn, delay) {
		var timer = null;
		return function () {
			var context = this, args = arguments;
			clearTimeout(timer);
				timer = setTimeout(function () {
					fn.apply(context, args);
				}, delay);
		};
	}
	
	function nl2br (str, is_xhtml) {
		var breakTag = (is_xhtml || typeof is_xhtml === 'undefined') ? '<br />' : '<br>';
		return (str + '').replace(/([^>\r\n]?)(\r\n|\n\r|\r|\n)/g, '$1' + breakTag + '$2');
	}
	
	
	function projectClicked(id) {
		currentProject = id;
		var c = $('#includeclosed').is(':checked');
		$.getJSON(location + "issuesbyproject/?q="+id+"&c="+c, function(data) {
			var items = [];
			
			var id;
			
			$.each(data, function(key, val) {
				id = val.Id;
				var name = val.Name;

				items.push(issueTemplate({id: id, name: name}));
			});

			// If we only have 1 then lets load the notes
			if(items.length == 1) {
				issueClicked(id);
			}

			if(data.length != 0) {
				$('#issues').html('' + items.join(''));
			}
			else {
				$('#issues').html('No issues.');
			}
		});
	}
	
	function issueClicked(id) {
		currentIssue = id;
		$('#issueid').val(id);
		$.getJSON(location + "notesbyissue/?q="+id, function(data) {
			var items = [];

			$.each(data, function(key, val) {
				var id = val.Id;
				var content = nl2br(val.Content, true);
				items.push(noteTemplate({id: id, content : content}));
			});

			$('#noteForm').show();

			if(data.length != 0) {
				$('#notes').html('' + items.join(''));
			}
			else {
				$('#notes').html('No notes.');
			}
		});
	}
	
	function getAllProjects() {
	
		var c = $('#includeclosed').is(':checked');
	
		$.getJSON(location + "allprojects/?c="+c, function(data) {
			var items = [];

			$.each(data, function(key, val) {
				var id = val.Id;
				var name = val.Name;
				items.push(projectTemplate({id: id, name : name}));
			});

			if(data.length != 0) {
				$('#projectspeople').html('' + items.join(''));
			}
			else {
				$('#projectspeople').html('No projects.');
			}
		});
	}
	
	function getAllIssues() {
		var c = $('#includeclosed').is(':checked');
		
		$.getJSON(location + "allissues/?c="+c, function(data) {
			var items = [];

			var id;
			var name;
			
			$.each(data, function(key, val) {
				id = val.Id;
				name = val.Name;

				items.push(issueTemplate({id: id, name: name}));
			});

			// If we only have 1 then lets load the notes
			if(items.length == 1) {
				issueClicked(id);
			}

			if(data.length != 0) {
				$('#issues').html('' + items.join(''));
			}
			else {
				$('#issues').html('No issues.');
			}
		});
	}
	
	function projectSearch() {
		var c = $('#includeclosed').is(':checked');
		var s = encodeURIComponent($("#search").val());
		if(s == '') {
			getAllProjects();
			return;
		}
		$.getJSON(location + "projectssearch/?q="+s+"&c="+c, function(data) {
			var items = [];

			var id;
			var name;
			
			$.each(data, function(key, val) {
				id = val.Id;
				name = val.Name;
				items.push(projectTemplate({id: id, name : name}));
			});

			if(items.length == 1) {
				currentProject = id;
			}
			
			if(data.length != 0) {
				$('#projectspeople').html('' + items.join(''));
			}
			else {
				$('#projectspeople').html('No projects.');
			}
		});
	}

	function issueSearch() {
		var c = $('#includeclosed').is(':checked');
		var s = encodeURIComponent($("#search").val());
		if(s == '') {
			getAllIssues();
			return;
		}
		$.getJSON(location + "issuessearch/?q="+s+"&c="+c, function(data) {
			var items = [];
			
			var id;
			var name;

			$.each(data, function(key, val) {
				id = val.Id;
				name = val.Name;

				items.push(issueTemplate({id: id, name: name}));
			});

			// If we only have 1 then lets load the notes
			if(items.length == 1) {
				currentIssue = id;
				issueClicked(id);
			}

			if(data.length != 0) {
				$('#issues').html('' + items.join(''));
			}
			else {
				$('#issues').html('No issues.');
			}
		});
	}
	
});