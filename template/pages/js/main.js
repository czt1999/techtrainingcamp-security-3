function getName(){
	console.log(1)
	$.ajax({
		type: "POST",
		dataType: "json",
		url: '/api/user/name',
		contentType: "application/json",
		data:JSON.stringify({
			"session_id": sessionStorage.getItem("session_id"),
		}),
		success: function (result) {
			console.log(result)
			if (result.code == 0) {
				document.getElementById("name").innerHTML = result.data.username
			}else {
				alert("登录超时")
				window.location.href='/';
			}
		}
	});
}
function logout(){

	$.ajax({
		type: "POST",
		dataType: "json",
		url: 'api/logout',
		contentType: "application/json",
		data:JSON.stringify({
			"session_id": sessionStorage.getItem("session_id"),
			"action_type" : 1
		}),
		success: function (result) {
			console.log(sessionStorage.getItem("session_id"))
			console.log(result)
			if (result.code == 0) {
				
				alert("登出成功");
				sessionStorage.removeItem("session_id")
				window.location.href='/';
			}else {
				alert("登出失败" + result.msg)
				window.location.href='/';
			}
		}
	});
	
}
function deleteL(){
	$.ajax({
		type: "POST",
		dataType: "json",
		url: 'api/logout',
		contentType: "application/json",
		data:JSON.stringify({
			"session_id": sessionStorage.getItem("session_id"),
			"action_type" : 2
		}),
		success: function (result) {
			console.log(sessionStorage.getItem("session_id"))
			console.log(result)
			if (result.code == 0) {
				
				alert("注销成功");
				sessionStorage.removeItem("session_id")
				window.location.href='/';
			}else {
				alert("注销失败" + result.msg)
				window.location.href='/';
			}
		}
	});
	
}