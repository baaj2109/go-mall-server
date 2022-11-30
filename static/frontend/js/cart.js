(function($){
    var app={
        init:function(){             
            this.changeCartNum();       
            this.deleteConfirm();
            this.initCheckBox();
            this.isCheckedAll();
        },
        deleteConfirm:function(){
            $('.delete').click(function(){    
                var flag=confirm('您確定要刪除吗?');    
                return flag;    
            })
    
        },    
        initCheckBox(){
            //全選按钮點擊
            $("#checkAll").click(function() {               
                if (this.checked) {
                    $(":checkbox").prop("checked", true);
                    $.get('/cart/changeAllCart?flag=1',function(response){
                        if(response.success){
                            $("#allPrice").html(response.allPrice+"元")                      
                        }
                    })
                  
                }else {
                    $(":checkbox").prop("checked", false);   
                    $.get('/cart/changeAllCart?flag=0',function(response){
                        if(response.success){
                            $("#allPrice").html(response.allPrice+"元")                      
                        }
                    })                    
                }

               
            });    

            // //點擊單個選擇框按钮的時候触發
            var _that=this;  //this是app对象            
            $(".cart_list :checkbox").click(function() {            
                _that.isCheckedAll(); 

                var product_id=$(this).attr('product_id');
                var product_color=$(this).attr('product_color');
                $.get('/cart/changeOneCart?product_id='+product_id+'&product_color='+product_color,function(response){

                    if(response.success){
                        $("#allPrice").html(response.allPrice+"元")                      
                    }
                })

            });
        },
        //判断全選是否選擇
        isCheckedAll(){             
            var chknum = $(".cart_list :checkbox").size();//checkbox總個數
            var chk = 0;  //checkbox checked=true總個數
            $(".cart_list :checkbox").each(function () {  
                if($(this).prop("checked")==true){
                    chk++;
                }
            });
            if(chknum==chk){//全選
                $("#checkAll").prop("checked",true);
            }else{//不全選
                $("#checkAll").prop("checked",false);
            }
        }, 
        changeCartNum(){

            $('.decCart').click(function(){
                
                var product_id=$(this).attr('product_id');
                var product_color=$(this).attr('product_color');
          
                $.get('/cart/decCart?product_id='+product_id+'&product_color='+product_color,function(response){

                    if(response.success){
                        $("#allPrice").html(response.allPrice+"元")
                        //注意this指向
                        $(this).siblings(".input_center").find("input").val(response.num)
                        $(this).parent().parent().siblings(".totalPrice").html(response.currentAllPrice+"元")
                    }
                }.bind(this))

            });

            $('.incCart').click(function(){             
                var product_id=$(this).attr('product_id');
                var product_color=$(this).attr('product_color');
          
                $.get('/cart/incCart?product_id='+product_id+'&product_color='+product_color,function(response){

                    if(response.success){
                        $("#allPrice").html(response.allPrice+"元")
                        //注意this指向
                        $(this).siblings(".input_center").find("input").val(response.num)
                        $(this).parent().parent().siblings(".totalPrice").html(response.currentAllPrice+"元")
                    }
                }.bind(this))
               
            });
        }
    }

    $(function(){
        app.init();
    })    
})($)
