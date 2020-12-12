(window.webpackJsonp=window.webpackJsonp||[]).push([[0],{12:function(e,t,a){e.exports=a(26)},20:function(e,t,a){},22:function(e,t,a){},24:function(e,t,a){},26:function(e,t,a){"use strict";a.r(t);var r=a(0),n=a.n(r),s=a(11),o=a.n(s),u=a(1),i=a.n(u),c=a(2),l=a(3),m=a(4),p=a(6),d=a(5),f=a(7),h={signIn:"SIGNIN",signUp:"SIGNUP",signedInMain:"SIGNEDIN_MAIN",signedInUpdateName:"SIGNEDIN_UPDATENAME",signedInUpdateAvatar:"SIGNEDIN_UPDATEAVATAR",forgotPassword:"FORGOT_PASSWORD"},g=a(8),E=function(e){var t=e.setField,a=e.submitForm,r=e.values,s=e.fields;return n.a.createElement(n.a.Fragment,null,n.a.createElement("form",{onSubmit:a},s.map(function(e){var a=e.key,s=e.name;return n.a.createElement("div",{key:a},n.a.createElement("span",null,s,": "),n.a.createElement("input",{value:r[a],name:a,onChange:t,type:"password"===a||"passwordConf"===a?"password":""}))}),n.a.createElement("input",{type:"submit",value:"Submit"})))},v={base:"https://api.xutiancheng.me",testbase:"https://localhost:8080",handlers:{users:"/v1/users",myuser:"/v1/users/me",myuserAvatar:"/v1/users/me/avatar",sessions:"/v1/sessions",sessionsMine:"/v1/sessions/mine",resetPasscode:"/v1/resetcodes",passwords:"/v1/passwords/"}},b=(a(20),function(e){var t=e.error,a=e.setError;switch(t){case"":return n.a.createElement(n.a.Fragment,null);default:return n.a.createElement("div",{className:"error"},n.a.createElement("span",{className:"error-hide",onClick:function(){return a("")}},"x"),"Error: ",t)}}),w=function(e){function t(e){var a;return Object(l.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).setField=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.submitForm=function(){var e=Object(c.a)(i.a.mark(function e(t){var r,n,s,o,u,c,l,m,p,d,f,h;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),r=a.state,n=r.email,s=r.userName,o=r.firstName,u=r.lastName,c=r.password,l=r.passwordConf,m={email:n,userName:s,firstName:o,lastName:u,password:c,passwordConf:l},e.next=5,fetch(v.base+v.handlers.users,{method:"POST",body:JSON.stringify(m),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((p=e.sent).status>=300)){e.next=12;break}return e.next=9,p.text();case 9:return d=e.sent,a.setError(d),e.abrupt("return");case 12:return f=p.headers.get("Authorization"),localStorage.setItem("Authorization",f),a.setError(""),a.props.setAuthToken(f),e.next=18,p.json();case 18:h=e.sent,a.props.setUser(h);case 20:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.state={email:"",userName:"",firstName:"",lastName:"",password:"",passwordConf:"",error:""},a.fields=[{name:"Email",key:"email"},{name:"Username",key:"userName"},{name:"First name",key:"firstName"},{name:"Last name",key:"lastName"},{name:"Password",key:"password"},{name:"Password Confirmation",key:"passwordConf"}],a}return Object(f.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=this.state.error;return n.a.createElement(n.a.Fragment,null,n.a.createElement(b,{error:a,setError:this.setError}),n.a.createElement(E,{setField:this.setField,submitForm:this.submitForm,values:t,fields:this.fields}),n.a.createElement("button",{onClick:function(t){return e.props.setPage(t,h.signIn)}},"Sign in instead"))}}]),t}(r.Component),k=function(e){function t(e){var a;return Object(l.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).setField=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.submitForm=function(){var e=Object(c.a)(i.a.mark(function e(t){var r,n,s,o,u,c,l,m;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),r=a.state,n=r.email,s=r.password,o={email:n,password:s},e.next=5,fetch(v.base+v.handlers.sessions,{method:"POST",body:JSON.stringify(o),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((u=e.sent).status>=300)){e.next=12;break}return e.next=9,u.text();case 9:return c=e.sent,a.setError(c),e.abrupt("return");case 12:return l=u.headers.get("Authorization"),localStorage.setItem("Authorization",l),a.setError(""),a.props.setAuthToken(l),e.next=18,u.json();case 18:m=e.sent,a.props.setUser(m);case 20:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.state={email:"",password:"",error:""},a.fields=[{name:"Email",key:"email"},{name:"Password",key:"password"}],a}return Object(f.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=this.state.error;return n.a.createElement(n.a.Fragment,null,n.a.createElement(b,{error:a,setError:this.setError}),n.a.createElement(E,{setField:this.setField,submitForm:this.submitForm,values:t,fields:this.fields}),n.a.createElement("button",{onClick:function(t){return e.props.setPage(t,h.signUp)}},"Sign up instead"),n.a.createElement("button",{onClick:function(t){return e.props.setPage(t,h.forgotPassword)}},"Forgot password"))}}]),t}(r.Component),y=function(e){function t(e){var a;return Object(l.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).sendResetCode=function(){var e=Object(c.a)(i.a.mark(function e(t){var r,n,s,o;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),r=a.state.email,n={email:r},e.next=5,fetch(v.base+v.handlers.resetPasscode,{method:"POST",body:JSON.stringify(n),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((s=e.sent).status>=300)){e.next=12;break}return e.next=9,s.text();case 9:return o=e.sent,a.setError(o),e.abrupt("return");case 12:a.setError(""),alert("A reset code has been sent to your email"),a.setState({resetCodeSent:!0});case 15:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.updatePassword=function(){var e=Object(c.a)(i.a.mark(function e(t){var r,n,s,o,u,c,l,m;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),r=a.state,n=r.email,s=r.password,o=r.passwordConf,u=r.resetCode,c={password:s,passwordConf:o,resetCode:u},e.next=5,fetch(v.base+v.handlers.passwords+n,{method:"PUT",body:JSON.stringify(c),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((l=e.sent).status>=300)){e.next=12;break}return e.next=9,l.text();case 9:return m=e.sent,a.setError(m),e.abrupt("return");case 12:a.setError(""),alert("Your password has been updated");case 14:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.setValue=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.state={email:"",password:"",passwordConf:"",resetCode:"",resetCodeSent:!1,error:""},a}return Object(f.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=t.email,r=t.password,s=t.passwordConf,o=t.resetCode,u=t.resetCodeSent,i=t.error;return n.a.createElement(n.a.Fragment,null,n.a.createElement(b,{error:i,setError:this.setError}),u?n.a.createElement(n.a.Fragment,null,n.a.createElement("form",{onSubmit:function(t){return e.updatePassword(t)}},n.a.createElement("div",null,n.a.createElement("span",null,"Password: "),n.a.createElement("input",{name:"password",type:"password",onChange:this.setValue,value:r})),n.a.createElement("div",null,n.a.createElement("span",null,"Password Confirmation: "),n.a.createElement("input",{name:"passwordConf",type:"password",onChange:this.setValue,value:s})),n.a.createElement("div",null,n.a.createElement("span",null,"Reset Code: "),n.a.createElement("input",{name:"resetCode",onChange:this.setValue,value:o})),n.a.createElement("input",{type:"submit",value:"Update my password"}))):n.a.createElement(n.a.Fragment,null,n.a.createElement("form",{onSubmit:function(t){return e.sendResetCode(t)}},n.a.createElement("div",null,n.a.createElement("span",null,"Email: "),n.a.createElement("input",{name:"email",onChange:this.setValue,value:a})),n.a.createElement("input",{type:"submit",value:"Send me a reset code"}))),n.a.createElement("button",{onClick:function(t){return e.props.setPage(t,h.signIn)}},"Back to sign in"))}}]),t}(r.Component),O=function(e){var t=e.page,a=e.setPage,r=e.setAuthToken,s=e.setUser;switch(t){case h.signUp:return n.a.createElement(w,{setPage:a,setAuthToken:r,setUser:s});case h.signIn:return n.a.createElement(k,{setPage:a,setAuthToken:r,setUser:s});case h.forgotPassword:return n.a.createElement(y,{setPage:a});default:return n.a.createElement(n.a.Fragment,null,"Error, invalid path reached")}},j=a(9),S=(a(22),function(e){var t=e.user,a=e.setPage,s=Object(r.useState)(null),o=Object(j.a)(s,2),u=o[0],l=o[1];function m(){return(m=Object(c.a)(i.a.mark(function e(){var a,r;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,fetch(v.base+v.handlers.myuserAvatar,{method:"GET",headers:new Headers({Authorization:localStorage.getItem("Authorization")})});case 2:if(!((a=e.sent).status>=300)){e.next=6;break}return l(t.photoURL),e.abrupt("return");case 6:return e.next=8,a.blob();case 8:r=e.sent,l(URL.createObjectURL(r));case 10:case"end":return e.stop()}},e)}))).apply(this,arguments)}return Object(r.useEffect)(function(){!function(){m.apply(this,arguments)}()},[]),n.a.createElement(n.a.Fragment,null,n.a.createElement("div",null,"Welcome to my application, ",t.firstName," ",t.lastName),u&&n.a.createElement("img",{className:"avatar",src:u,alt:"".concat(t.firstName,"'s avatar")}),n.a.createElement("div",null,n.a.createElement("button",{onClick:function(e){a(e,h.signedInUpdateName)}},"Update name")),n.a.createElement("div",null,n.a.createElement("button",{onClick:function(e){a(e,h.signedInUpdateAvatar)}},"Update avatar")))}),C=function(e){var t=e.setAuthToken,a=e.setUser,s=Object(r.useState)(""),o=Object(j.a)(s,2),u=o[0],l=o[1];return n.a.createElement(n.a.Fragment,null,n.a.createElement("button",{onClick:function(){var e=Object(c.a)(i.a.mark(function e(r){var n,s;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return r.preventDefault(),e.next=3,fetch(v.base+v.handlers.sessionsMine,{method:"DELETE",headers:new Headers({Authorization:localStorage.getItem("Authorization")})});case 3:if(!((n=e.sent).status>=300)){e.next=10;break}return e.next=7,n.text();case 7:return s=e.sent,l(s),e.abrupt("return");case 10:localStorage.removeItem("Authorization"),l(""),t(""),a(null);case 14:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}()},"Sign out"),u&&n.a.createElement("div",null,n.a.createElement(b,{error:u,setError:l})))},x=function(e){function t(e){var a;return Object(l.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).sendRequest=function(){var e=Object(c.a)(i.a.mark(function e(t){var r,n,s,o,u,c,l;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),r=a.state,n=r.firstName,s=r.lastName,o={firstName:n,lastName:s},e.next=5,fetch(v.base+v.handlers.myuser,{method:"PATCH",body:JSON.stringify(o),headers:new Headers({Authorization:localStorage.getItem("Authorization"),"Content-Type":"application/json"})});case 5:if(!((u=e.sent).status>=300)){e.next=13;break}return e.next=9,u.text();case 9:return c=e.sent,console.log(c),a.setError(c),e.abrupt("return");case 13:return alert("Name changed"),e.next=16,u.json();case 16:l=e.sent,a.props.setUser(l);case 18:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.setValue=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.state={firstName:"",lastName:"",error:""},a}return Object(f.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this.state,t=e.firstName,a=e.lastName,r=e.error;return n.a.createElement(n.a.Fragment,null,n.a.createElement(b,{error:r,setError:this.setError}),n.a.createElement("div",null,"Enter a new name"),n.a.createElement("form",{onSubmit:this.sendRequest},n.a.createElement("div",null,n.a.createElement("span",null,"First name: "),n.a.createElement("input",{name:"firstName",value:t,onChange:this.setValue})),n.a.createElement("div",null,n.a.createElement("span",null,"Last name: "),n.a.createElement("input",{name:"lastName",value:a,onChange:this.setValue})),n.a.createElement("input",{type:"submit",value:"Change name"})))}}]),t}(r.Component),A=function(e){function t(e){var a;return Object(l.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).sendRequest=function(){var e=Object(c.a)(i.a.mark(function e(t){var r,n,s,o;return i.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),r=a.state.file,(n=new FormData).append("uploadfile",r),e.next=6,fetch(v.base+v.handlers.myuserAvatar,{method:"PUT",body:n,headers:new Headers({Authorization:localStorage.getItem("Authorization")})});case 6:if(!((s=e.sent).status>=300)){e.next=14;break}return e.next=10,s.text();case 10:return o=e.sent,console.log(o),a.setError(o),e.abrupt("return");case 14:alert("Avatar changed");case 15:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.handleFile=function(e){a.setState({file:e.target.files[0]})},a.setError=function(e){a.setState({error:e})},a.state={file:null,error:""},a}return Object(f.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this.state.error;return n.a.createElement(n.a.Fragment,null,n.a.createElement(b,{error:e,setError:this.setError}),n.a.createElement("div",null,"Give yourself a new avatar"),n.a.createElement("form",{onSubmit:this.sendRequest},n.a.createElement("div",null,n.a.createElement("span",null,"Upload new avatar "),n.a.createElement("input",{type:"file",onChange:this.handleFile})),n.a.createElement("input",{type:"submit",value:"Change avatar"})))}}]),t}(r.Component),N=function(e){var t=e.page,a=e.setPage,r=e.setAuthToken,s=e.setUser,o=e.user,u=n.a.createElement(n.a.Fragment,null),i=!0;switch(t){case h.signedInMain:u=n.a.createElement(S,{user:o,setPage:a});break;case h.signedInUpdateName:u=n.a.createElement(x,{user:o,setUser:s});break;case h.signedInUpdateAvatar:u=n.a.createElement(A,{user:o,setUser:s});break;default:u=n.a.createElement(n.a.Fragment,null,"Error, invalid path reached"),i=!1}return n.a.createElement(n.a.Fragment,null,u,i&&n.a.createElement("button",{onClick:function(e){return a(e,h.signedInMain)}},"Back to main"),n.a.createElement(C,{setUser:s,setAuthToken:r}))},U=(a(24),function(e){function t(){var e;return Object(l.a)(this,t),(e=Object(p.a)(this,Object(d.a)(t).call(this))).getCurrentUser=Object(c.a)(i.a.mark(function t(){var a,r;return i.a.wrap(function(t){for(;;)switch(t.prev=t.next){case 0:if(e.state.authToken){t.next=2;break}return t.abrupt("return");case 2:return t.next=4,fetch(v.base+v.handlers.myuser,{headers:new Headers({Authorization:e.state.authToken})});case 4:if(!((a=t.sent).status>=300)){t.next=11;break}return alert("Unable to verify login. Logging out..."),localStorage.setItem("Authorization",""),e.setAuthToken(""),e.setUser(null),t.abrupt("return");case 11:return t.next=13,a.json();case 13:r=t.sent,e.setUser(r);case 15:case"end":return t.stop()}},t)})),e.setPageToSignIn=function(t){t.preventDefault(),e.setState({page:h.signIn})},e.setPageToSignUp=function(t){t.preventDefault(),e.setState({page:h.signUp})},e.setPage=function(t,a){t.preventDefault(),e.setState({page:a})},e.setAuthToken=function(t){e.setState({authToken:t,page:""===t?h.signIn:h.signedInMain})},e.setUser=function(t){e.setState({user:t})},e.state={page:localStorage.getItem("Authorization")?h.signedInMain:h.signIn,authToken:localStorage.getItem("Authorization")||null,user:null},e.getCurrentUser(),e}return Object(f.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this.state,t=e.page,a=e.user;return n.a.createElement("div",null,a?n.a.createElement(N,{page:t,setPage:this.setPage,setAuthToken:this.setAuthToken,user:a,setUser:this.setUser}):n.a.createElement(O,{page:t,setPage:this.setPage,setAuthToken:this.setAuthToken,setUser:this.setUser}))}}]),t}(r.Component));Boolean("localhost"===window.location.hostname||"[::1]"===window.location.hostname||window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/));o.a.render(n.a.createElement(U,null),document.getElementById("root")),"serviceWorker"in navigator&&navigator.serviceWorker.ready.then(function(e){e.unregister()})}},[[12,2,1]]]);
//# sourceMappingURL=main.e1a4dd05.chunk.js.map