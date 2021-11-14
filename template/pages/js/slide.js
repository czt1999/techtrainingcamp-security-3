  //加载（事件会在页面加载完成后触发）
        window.onload = function() {
            //获取滑动控件容器,灰色背景
            var dragContainer = document.getElementById("dragContainer");
            //获取滑块左边部分,绿色背景
            var dragBg = document.getElementById("dragBg");
            //获取滑动验证容器文本
            var dragText = document.getElementById("dragText");
            //获取滑块
            var dragHandler = document.getElementById("dragHandler");

            //滑块的最大偏移量 =滑动验证容器文本长度- 滑块长度
            var maxHandlerOffset = dragContainer.clientWidth - dragHandler.clientWidth;
            //是否验证成功的标记
            var isVertifySucc = false;
            
			var HandMaxlen = 0;
			
			var StartIndex 
            initDrag();
 
            function initDrag() {
                //在滑动验证容器文本写入“拖动滑块验证”
                dragText.textContent = "拖动滑块验证";
                //给滑块添加鼠标按下监听
                dragHandler.addEventListener("mousedown", onDragHandlerMouseDown);
            }
             
           //选中滑块
            function onDragHandlerMouseDown() {
                //鼠标移动监听
                document.addEventListener("mousemove", onDragHandlerMouseMove);
                //鼠标松开监听
                document.addEventListener("mouseup",  onDragHandlerMouseUp);
				
				StartIndex = event.clientX 
            }
             
           //滑块移动
            function onDragHandlerMouseMove() {
                /*
                html元素不存在width属性,只有clientWidth
                offsetX是相对当前元素的,clientX和pageX是相对其父元素的
				但是在页面中间，无法确认滑块在页面的初始X，可以在选中时
				以当前鼠标坐标为初始
                */
               //滑块移动量
                var left = event.clientX - StartIndex - dragHandler.clientWidth  / 2;
                //
                if(left < 0) {
                    left = 0;
                 //如果滑块移动量   > 滑块的最大偏移量 ，则调用验证成功函数
                } else if(left > maxHandlerOffset) {
                    left = maxHandlerOffset;
                    verifySucc();
                }
                //滑块移动量
                dragHandler.style.left = left + "px";
                //绿色背景的长度
                dragBg.style.width = dragHandler.style.left;
				
				HandMaxlen = left
            }
            
           //松开滑块函数
            function onDragHandlerMouseUp() {
                //移除鼠标移动监听
                document.removeEventListener("mousemove", onDragHandlerMouseMove);
                //移除鼠标松开监听
                document.removeEventListener("mouseup", onDragHandlerMouseUp);
                //初始化滑块移动量      
				clearLen=setInterval(function count(){
					//当滑块当前长度不为0时则进行计时减少操作
					if(HandMaxlen >= 0){
						HandMaxlen -=2;
						dragHandler.style.left = HandMaxlen + "px";
						//绿色背景的长度
						dragBg.style.width = dragHandler.style.left;
					}else{
						HandMaxlen = 0
						dragHandler.style.left = 0
						dragBg.style.width = 0;
						clearInterval(clearLen);
					}
				},1.5);
				
				
            }

            //验证成功
            function verifySucc() {
                //成功标记，不可回弹
				$("#slideCode").val("true")
                isVertifySucc = false;
                //容器文本的文字改为白色“验证通过”字体
                dragText.textContent = "验证通过";
                dragText.style.color = "white";
                //验证通过的滑块背景
                dragHandler.setAttribute("class", "dragHandlerOkBg");
                //移除鼠标按下监听
                dragHandler.removeEventListener("mousedown", onDragHandlerMouseDown);
                //移除 鼠标移动监听
                document.removeEventListener("mousemove", onDragHandlerMouseMove);
                //移除鼠标松开监听
                document.removeEventListener("mouseup", onDragHandlerMouseUp);
            };
        }