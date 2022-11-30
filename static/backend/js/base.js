$(function(){	
	app.init()
	$(window).resize(function(){
		app.resizeIframe()
	})	
})

var config={
	adminPath:"backend"
}
var app={
	
	init:function(){
		this.slideToggle();
		this.resizeIframe()
		this.confirmDelete()
		this.changeStatus()
		this.changeNum()
	},	
	slideToggle:function(){
		$('.aside>li:nth-child(1) ul,.aside>li:nth-child(2) ul').hide()
		$('.aside h4').click(function(){
		
			$(this).siblings('ul').slideToggle();
		})
	},
	resizeIframe:function(){		
		$("#rightMain").height($(window).height()-80)
	},
	// 刪除提示
	confirmDelete:function(){
		$(".delete").click(function(){
			var flag=confirm("您確定要刪除吗?")
			return flag
		})
	},		
	changeStatus:function(){
		$(".chStatus").click(function(){
			var id=$(this).attr("data-id");
			var table=$(this).attr("data-table");
			var field=$(this).attr("data-field");
			var el=$(this)
			$.get("/"+config.adminPath+"/main/changestatus",{id:id,table:table,field:field},function(response){				
				if(response.success){

					if(el.attr("src").indexOf("yes")!=-1){
						el.attr("src","/static/backend/images/no.gif")
					}else{
						el.attr("src","/static/backend/images/yes.gif")
					}
				}else{
					console.log(response)
				}
			})
		})
	},
	changeNum:function(){
		/*
		1. 獲取el里面的值  var spanNum=$(this).html()


		2. 創建一個input的dom節點   var input=$("<input value='' />");


		3. 把input放在el里面   $(this).html(input);


		4. 让input獲取焦点  给input赋值    $(input).trigger('banner').val(val);

		5. 點擊input的時候阻止冒泡 

					$(input).click(function(e){
						e.stopPropagation();				
					})

		6. 鼠标离開的時候给span赋值,并触發ajax請求

			$(input).blur(function(){
				var inputNum=$(this).val();
				spanEl.html(inputNum);
				触發ajax請求
				
			})
		*/

		$(".chSpanNum").click(function(){
			var id=$(this).attr("data-id");
			var table=$(this).attr("data-table");
			var field=$(this).attr("data-field");
			var spanNum=$(this).attr("data-num");
			var spanEl=$(this)  //保存span这個dom節點

			var input=$("<input value='' style='width:60px' />");
			$(this).html(input);
			$(input).trigger('banner').val(spanNum);   //让輸入框獲取焦点并設置值
			$(input).click(function(e){
				e.stopPropagation();				
			})
			$(input).blur(function(e){
				var inputNum=$(this).val();
				spanEl.html(inputNum);
				//异步請求修改數量
				$.get("/"+config.adminPath+"/main/editnum",{id:id,table:table,field:field,num:inputNum},function(response){				
					if(!response.success){
						console.log(response)
					}
				})
			})
		})
	}

}
