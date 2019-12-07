var Buttons = document.getElementsByClassName('calculator')[0],//按钮
    content1 = document.getElementsByClassName('CalculatorScreen1')[0].getElementsByTagName('p')[0], // 显示器表达式内容
    content2 = document.getElementsByClassName('CalculatorScreen2')[0].getElementsByTagName('p')[0], // 显示器结果
    IsEval = false;//判断是否点击了“=”
    var sock=null;//创建websocket
    var wsuri="ws://127.0.0.1:1234/websocket";//服务器地址
    window.onload=function(){
        console.log("onload");
        sock=new WebSocket(wsuri);
        sock.onopen=function(){//开始
            console.log("connected to"+wsuri);
        }
        sock.onclose=function(e){
            console.log("connected closed("+e.code+")");
        }
        sock.onmessage=function(e){//接收到信息
            console.log("message received:"+e.data);
            content2.innerText = e.data;
        }
    };
    Buttons.onclick = function(e){
    var target = e.target;
    if(target.nodeName.toLowerCase() === 'button'){//判断是否点击到了按钮
        var btnType = target.innerText;
        if(IsEval){//如果已经点击过“=”，任意操作都会将上一个计算结果显示到表达显示屏中，并清零结果屏
            IsEval=false;
            content1.innerText=content2.innerText;
            content2.innerText = '';
            return;
        }

        if(btnType == 'CE'){//删除键
            if(content1.innerText != ''){
                if(content1.innerText.length === 1){
                    content1.innerText = '0';
                }else{
                    content1.innerText = content1.innerText.slice(0,-1);
                }
            }
        }else if(btnType == 'AC'){//清零键
            if(content1.innerText != ''){
                    content1.innerText = '0';
            }
        }else if(btnType == '='){
            IsEval=true;
            var text = content1.innerText;
            if(!text){
                return;
            }else{
                //在前端先进行第一步预处理
                //由于有的字符串不能进行全局替换，这里运用了循环
                for(var i=0;i<=50;i++){
                    text = text.replace('x','*');
                    text = text.replace('%','/100');
                    text = text.replace('lg','g');
                    text = text.replace('ln','n');
                    text = text.replace('sin','s');
                    text = text.replace('cos','c');
                    text = text.replace('tan','t');
                    text = text.replace('^2','p');
                    text = text.replace('√','k');
                    text = text.replace('^(-1)','d');
                }
                if(text[0]=='-')text = text.replace('-','0-');//这里处理第一个数为负数的情况

                //把处理好的字符串发往后台进行计算
                var msg=text;
                sock.send(msg);
            }
        }else if(btnType == 'x²'){
            if(content1.innerText == '0'){
                content1.innerText = '';
            }
            content1.innerText +='^2' 
        }else if(btnType == '√x'){
            if(content1.innerText == '0'){
                content1.innerText = '';
            }
            content1.innerText +='√'
        }else if(btnType == '1/x'){
            if(content1.innerText == '0'){
                content1.innerText = '';
            }
            content1.innerText +='^(-1)'
        }else if(btnType == 'x!'){
            if(content1.innerText == '0'){
                content1.innerText = '';
            }
            content1.innerText +='!'
        }else{
            if(content1.innerText == '0' && (!isNaN(+btnType) ||btnType == '(' || btnType == ')'|| btnType == 'sin'|| btnType == 'cos'|| btnType == 'tan'|| btnType == 'ln'|| btnType == 'lg')){
                content1.innerText = '';
            }
            content1.innerText += btnType;
        }   
    }
}
