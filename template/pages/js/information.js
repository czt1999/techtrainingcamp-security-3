function exit(){
	window.location.href="../school/index.html"
}
function buttonEnter (e) {
	var p = document.getElementById(a[e])
	if(f[e]){return;}
	p.style.borderTopRightRadius = 8+"px";
	p.style.borderBottomRightRadius = 8+"px";
	p.style.borderTop = "solid 1px #89ff02aa inset"
	p.style.borderBottom = "solid 1px #89ff02aa inset"
	p.style.backgroundColor = "#89ff02aa";
	p.style.borderRight = "solid 9px #89ff02aa";
	p.style.color = "white";
}
function buttonOut(e){
	var p = document.getElementById(a[e])
	if(f[e]){return;}
	p.style.borderRadius = 0+"px";
	p.style.backgroundColor = "white";
	p.style.border = "solid 0px #89ff02aa";
	p.style.color = "#888888";
}
function buttonDown(e){
	if(index !=0){
		f[index] = false;
		var p = document.getElementById(a[index])
		p.style.borderRadius = 0+"px";
		p.style.backgroundColor = "white";
		p.style.border = "solid 0px #47a0ffaa ";
		p.style.color = "#888888";
	}
		var d = document.getElementById(b[index]);
			d.style.display = "none";

	/////////////////////////////////////////////////
	var d = document.getElementById(b[e]);
		d.style.display = "flex";
	var close =setInterval(apper,6);
	var opr = 0;
	setTimeout( function(){}, 6);
		function apper(){
			slidClose = true
			if (opr >= 1) {
				slidClose = false
			    clearInterval(close);
			} else {
				opr = opr + 0.01; 
				d.style.opacity = opr;
				d.style.marginTop = ((10-opr*2) * 1.25)+"rem";
			}
		}
	//////////////////////////////////////////////
	var p = document.getElementById(a[e])
	p.style.borderTopRightRadius = 8+"px";
	p.style.borderBottomRightRadius = 8+"px";
	p.style.backgroundColor = "#47a0ffaa";
	p.style.borderRight = "solid 9px #47a0ffaa ";
	p.style.color = "white"
	f[e] = true;
	index = e;
}
function change(){
	var p = document.getElementById('phone')
	p.style.cursor = 'text'
	p.style.background = "white"
	p.placeholder="请输入"
	p.readOnly = false
	var btn = document.getElementById('change_btn')
	btn.style.display = 'none'
	var btn = document.getElementById('save_btn_1')
	btn.style.display = 'block'
	var btn = document.getElementById('save_btn_2')
	btn.style.display = 'block'
}
function cancel(){
	var p = document.getElementById('phone')
	p.style.cursor = 'text'
	p.style.background = "#EEEEEE"
	p.placeholder="无"
	p.readOnly = true
	var btn = document.getElementById('change_btn')
	btn.style.display = 'block'
	var btn = document.getElementById('save_btn_1')
	btn.style.display = 'none'
	var btn = document.getElementById('save_btn_2')
	btn.style.display = 'none'
}

var a=["Main","graduation_project","internship_filling","deal_message"];
var f=[false,false,false,false];
var b=["zero","first","second","thrid"];
var f2=[false,false,false,false];
var index = 0;
var index2 = 0;