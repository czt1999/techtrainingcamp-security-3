	function login(){
		if($("#method").val() === "phone"){phoneLogin()}
		if($("#method").val() === "username"){usernameLogin()}
	}
	//获取手机验证码
	function getCode(){
		if($("#phone").val() === ""){alert("请填写电话号码");return;}
		$("#getCode").unbind("click")
		var Time = 60;
		$("#getCode").val("冷却时间:"+Time)
		clearTime=setInterval(function count(){
			if(Time > 0){
				Time --;
				$("#getCode").val("冷却时间:"+Time)
			}else{
				Time = 0
				$("#getCode").val("获取")
				$("#getCode").click(c => getCode());
				clearInterval(clearTime);
			}
		},1000);
		
		$.ajax({
			type: "POST",
			dataType: "json",
			url: '/api/applycode',
			contentType: "application/json",
			data:JSON.stringify({
				"phone_number": $("#phone").val(),
			}),
			success: function (result) {
				console.log()
				if (result.code === 200) {
				   
				}else {

				}
			}
		});
	}
	function getCode2(){
		if($("#Rphone").val() === ""){alert("请填写电话号码");return;}
		$("#RgetCode").unbind("click")
		var Time = 60;
		$("#RgetCode").val("冷却时间:"+Time)
		clearTime=setInterval(function count(){
			if(Time > 0){
				Time --;
				$("#RgetCode").val("冷却时间:"+Time)
			}else{
				Time = 0
				$("#RgetCode").val("获取")
				$("#RgetCode").click(c => getCode());
				clearInterval(clearTime);
			}
		},1000);
		
		$.ajax({
			type: "POST",
			dataType: "json",
			url: '/api/applycode',
			contentType: "application/json",
			data:JSON.stringify({
				"phone_number": $("#Rphone").val(),
			}),
			success: function (result) {
				console.log()
				if (result.code == 200) {
				   
				}else {
	
				}
			}
		});
	}
	//手机登录
	function phoneLogin(){
		console.log($("#codeInput").val())
		if($("#slideCode").val() === "false"){
			return
		}
		console.log(1)
		$("#entry_btn").unbind("click")
		var Time = 5;
		$("#entry_btn").val("冷却时间:"+Time)
		clearTime=setInterval(function count(){
			if(Time > 0){
				Time --;
				$("#entry_btn").val("冷却时间:"+Time)
			}else{
				Time = 0
				$("#entry_btn").val("登录")
				$("#entry_btn").click(c => login());
				clearInterval(clearTime);
			}
		},1000);
		$.ajax({
			type: "POST",
			dataType: "json",
			url: 'api/login/phone',
			contentType: "application/json",
			data:JSON.stringify({
				"phone_number": $("#phone").val(),
				"verify_code" : $("#codeInput").val(),
				"ip": $("#ip").val()
			}),
			success: function (result) {
				console.log(result)
				if (result.code === 0) {
					console.log(result)

					alert("登陆成功");
					sessionStorage.setItem("session_id",result.data.session_id)
				    window.location.href = "/main";
				}else {
					alert("登陆失败:" + result.msg)
				}
			}
		});
		
		
	}
	//用户名登录
	function usernameLogin(){

		if($("#slideCode").val() === "false"){
			return
		}
		$("#entry_btn").unbind("click")
		var Time = 5;
		$("#entry_btn").val("冷却时间:"+Time)
		clearTime=setInterval(function count(){
			if(Time > 0){
				Time --;
				$("#entry_btn").val("冷却时间:"+Time)
			}else{
				Time = 0
				$("#entry_btn").val("登录")
				$("#entry_btn").click(c => login());
				clearInterval(clearTime);
			}
		},1000);
		$.ajax({
			type: "POST",
			dataType: "json",
			url: 'api/login/name',
			contentType: "application/json",
			data:JSON.stringify({
				"username": $("#username").val(),
				"password": $("#password").val(),
				"ip": $("#ip").val()
			}),
			success: function (result) {
				console.log(result)
				if (result.code === 0) {
					alert("登陆成功");
					console.log(result)
					sessionStorage.setItem("session_id",result.data.session_id)
				    window.location.href = "/main";
				}else {
					alert("登陆失败:" + result.msg)
				}
			}
		});
		
		
	}
	
	//用户注册
	function register(){
		if($("#Rusername").val() === ""){alert("请输入用户名"); return}
		if($("#Rpassword").val() === ""){alert("请输入密码"); return}
		if($("#Apassword").val() === ""){alert("请再次输入密码"); return}
		if($("#Rpassword").val() !== $("#Apassword").val()){alert("两次密码不相同"); return}
		if($("#Rphone").val() === ""){alert("请输入手机号码"); return}
		
		$("#register_btn").unbind("click")
		var Time = 5;
		$("#register_btn").val("冷却时间:"+Time)
		clearTime=setInterval(function count(){
			if(Time > 0){
				Time --;
				$("#register_btn").val("冷却时间:"+Time)
			}else{
				Time = 0
				$("#register_btn").val("注册")
				$("#register_btn").click(c => register());
				clearInterval(clearTime);
			}
		},1000);
		$.ajax({
			type: "POST",
			dataType: "json",
			url: '/api/register',
			contentType: "application/json",
			data:JSON.stringify({
				"username": $("#Rusername").val(),
				"password": $("#Rpassword").val(),
				"phone_number": $("#Rphone").val(),
				"verify_code" : $("#RcodeInput").val(),
				"ip": $("#ip").val()
			}),
			success: function (result) {
				console.log(result)
				if (result.code === 0) {
					alert("注册成功");
				    window.location.href = "/";
				}else {
					alert("注册失败:" + result.msg)
				}
			}
		});
	}