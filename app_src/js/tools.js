function genCostIncomeId(keys){
	var id="TT201503251610SJ";
	var now=new Date();
	id = "CIT" + now.getFullYear().toString() + (now.getMonth() + 1).toString()
		+ now.getDate().toString() + now.getHours().toString()
		+ now.getMinutes().toString() + now.getSeconds().toString()
		+ keys;
	return id;
}

var xhr=null;

function login(uname,passwd){
	if(uname==='' || passwd===''){
		alert("请输入用户名和密码");
		return;
	}
	//var url="http://www.ourtest.com/application/login";
	var url="http://app.1keda.cn:3000/login";
	if(!xhr){
		xhr = new plus.net.XMLHttpRequest();
	}
	xhr.onreadystatechange = function () {
        switch ( xhr.readyState ) {
            case 0:
            	console.log( "xhr请求已初始化" );
            break;
            case 1:
            	console.log( "xhr请求已打开" );
            break;
            case 2:
            	console.log( "xhr请求已发送" );
            break;
            case 3:
                console.log( "xhr请求已响应");
                break;
            case 4:
                console.log( "xhr请求已完成");
                if ( xhr.status == 200 ) {
                	console.log( "xhr请求成功："+xhr.responseText );
                	var result = JSON.parse(xhr.responseText);
                	if(result.returnCode == '1'){
                		//login sucessed
                	}else{
                		alert("登录失败")
                	}
                	
                } else {
                	console.log( "xhr请求失败："+xhr.status );
                }
                break;
            default :
                break;
        }
	}
	
	xhr.open( "POST", url, uname,passwd);
	xhr.send();
}

function checkVersion(){
	var url="http://www.ourtest.com/application/checkVersion";
	if(!xhr){
		xhr = new plus.net.XMLHttpRequest();
	}
	xhr.onreadystatechange = function () {
        switch ( xhr.readyState ) {
            case 0:
            	console.log( "xhr请求已初始化" );
            break;
            case 1:
            	console.log( "xhr请求已打开" );
            break;
            case 2:
            	console.log( "xhr请求已发送" );
            break;
            case 3:
                console.log( "xhr请求已响应");
                break;
            case 4:
                console.log( "xhr请求已完成");
                if ( xhr.status == 200 ) {
                	console.log( "xhr请求成功："+xhr.responseText );
                } else {
                	console.log( "xhr请求失败："+xhr.status );
                }
                break;
            default :
                break;
        }
	}
	
	xhr.open( "GET", url );
	xhr.send();
}



function stepData(){
	if(!xhr){
		xhr = new plus.net.XMLHttpRequest();
	}
	xhr.open( "GET", url );
	

}
