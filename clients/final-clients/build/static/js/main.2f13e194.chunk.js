(window.webpackJsonp=window.webpackJsonp||[]).push([[0],{16:function(e,t,a){e.exports=a(37)},24:function(e,t,a){},26:function(e,t,a){},33:function(e,t,a){},35:function(e,t,a){},37:function(e,t,a){"use strict";a.r(t);var n=a(0),r=a.n(n),s=a(11),o=a.n(s),i=a(1),u=a.n(i),l=a(2),c=a(3),m=a(4),p=a(6),d=a(5),h=a(7),f={signIn:"SIGNIN",signUp:"SIGNUP",signedInMain:"SIGNEDIN_MAIN",signedInUpdateName:"SIGNEDIN_UPDATENAME",signedInUpdateAvatar:"SIGNEDIN_UPDATEAVATAR",signedInUploadImage:"SIGNEDIN_UPLOADIMAGE",forgotPassword:"FORGOT_PASSWORD",mainPage:"MAIN_PAGE"},g=a(10),E=function(e){var t=e.setField,a=e.submitForm,n=e.values,s=e.fields;return r.a.createElement(r.a.Fragment,null,r.a.createElement("form",{onSubmit:a},s.map(function(e){var a=e.key,s=e.name;return r.a.createElement("div",{key:a},r.a.createElement("span",null,s,": "),r.a.createElement("input",{value:n[a],name:a,onChange:t,type:"password"===a||"passwordConf"===a?"password":""}))}),r.a.createElement("input",{className:"animated_btn",type:"submit",value:"Submit"})))},b={base:"https://api.xutiancheng.me",testbase:"https://localhost:8080",handlers:{users:"/v1/users",myuser:"/v1/users/me",myuserAvatar:"/v1/users/me/avatar",sessions:"/v1/sessions",sessionsMine:"/v1/sessions/mine",resetPasscode:"/v1/resetcodes",passwords:"/v1/passwords/",upload:"/v1/upload/"}},v=(a(24),function(e){var t=e.error,a=e.setError;switch(t){case"":return r.a.createElement(r.a.Fragment,null);default:return r.a.createElement("div",{className:"error"},r.a.createElement("span",{className:"error-hide",onClick:function(){return a("")}},"x"),"Error: ",t)}}),w=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).setField=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.submitForm=function(){var e=Object(l.a)(u.a.mark(function e(t){var n,r,s,o,i,l,c,m,p,d,h,f;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),n=a.state,r=n.email,s=n.userName,o=n.firstName,i=n.lastName,l=n.password,c=n.passwordConf,m={email:r,userName:s,firstName:o,lastName:i,password:l,passwordConf:c},e.next=5,fetch(b.base+b.handlers.users,{method:"POST",body:JSON.stringify(m),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((p=e.sent).status>=300)){e.next=12;break}return e.next=9,p.text();case 9:return d=e.sent,a.setError(d),e.abrupt("return");case 12:return h=p.headers.get("Authorization"),localStorage.setItem("Authorization",h),a.setError(""),a.props.setAuthToken(h),e.next=18,p.json();case 18:f=e.sent,a.props.setUser(f);case 20:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.state={email:"",userName:"",firstName:"",lastName:"",password:"",passwordConf:"",error:""},a.fields=[{name:"Email",key:"email"},{name:"Username",key:"userName"},{name:"First name",key:"firstName"},{name:"Last name",key:"lastName"},{name:"Password",key:"password"},{name:"Password Confirmation",key:"passwordConf"}],a}return Object(h.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=this.state.error;return r.a.createElement(r.a.Fragment,null,r.a.createElement(v,{error:a,setError:this.setError}),r.a.createElement(E,{setField:this.setField,submitForm:this.submitForm,values:t,fields:this.fields}),r.a.createElement("button",{onClick:function(t){return e.props.setPage(t,f.signIn)}},"Sign in instead"))}}]),t}(n.Component),y=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).setField=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.submitForm=function(){var e=Object(l.a)(u.a.mark(function e(t){var n,r,s,o,i,l,c,m;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),n=a.state,r=n.email,s=n.password,o={email:r,password:s},e.next=5,fetch(b.base+b.handlers.sessions,{method:"POST",body:JSON.stringify(o),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((i=e.sent).status>=300)){e.next=12;break}return e.next=9,i.text();case 9:return l=e.sent,a.setError(l),e.abrupt("return");case 12:return c=i.headers.get("Authorization"),localStorage.setItem("Authorization",c),a.setError(""),a.props.setAuthToken(c),e.next=18,i.json();case 18:m=e.sent,a.props.setUser(m);case 20:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.state={email:"",password:"",error:""},a.fields=[{name:"Email",key:"email"},{name:"Password",key:"password"}],a}return Object(h.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=this.state.error;return r.a.createElement(r.a.Fragment,null,r.a.createElement(v,{error:a,setError:this.setError}),r.a.createElement(E,{setField:this.setField,submitForm:this.submitForm,values:t,fields:this.fields}),r.a.createElement("button",{onClick:function(t){return e.props.setPage(t,f.signUp)}},"Sign up instead"),r.a.createElement("button",{onClick:function(t){return e.props.setPage(t,f.forgotPassword)}},"Forgot password"))}}]),t}(n.Component),S=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).sendResetCode=function(){var e=Object(l.a)(u.a.mark(function e(t){var n,r,s,o;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),n=a.state.email,r={email:n},e.next=5,fetch(b.base+b.handlers.resetPasscode,{method:"POST",body:JSON.stringify(r),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((s=e.sent).status>=300)){e.next=12;break}return e.next=9,s.text();case 9:return o=e.sent,a.setError(o),e.abrupt("return");case 12:a.setError(""),alert("A reset code has been sent to your email"),a.setState({resetCodeSent:!0});case 15:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.updatePassword=function(){var e=Object(l.a)(u.a.mark(function e(t){var n,r,s,o,i,l,c,m;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),n=a.state,r=n.email,s=n.password,o=n.passwordConf,i=n.resetCode,l={password:s,passwordConf:o,resetCode:i},e.next=5,fetch(b.base+b.handlers.passwords+r,{method:"PUT",body:JSON.stringify(l),headers:new Headers({"Content-Type":"application/json"})});case 5:if(!((c=e.sent).status>=300)){e.next=12;break}return e.next=9,c.text();case 9:return m=e.sent,a.setError(m),e.abrupt("return");case 12:a.setError(""),alert("Your password has been updated");case 14:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.setValue=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.state={email:"",password:"",passwordConf:"",resetCode:"",resetCodeSent:!1,error:""},a}return Object(h.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=t.email,n=t.password,s=t.passwordConf,o=t.resetCode,i=t.resetCodeSent,u=t.error;return r.a.createElement(r.a.Fragment,null,r.a.createElement(v,{error:u,setError:this.setError}),i?r.a.createElement(r.a.Fragment,null,r.a.createElement("form",{onSubmit:function(t){return e.updatePassword(t)}},r.a.createElement("div",null,r.a.createElement("span",null,"Password: "),r.a.createElement("input",{name:"password",type:"password",onChange:this.setValue,value:n})),r.a.createElement("div",null,r.a.createElement("span",null,"Password Confirmation: "),r.a.createElement("input",{name:"passwordConf",type:"password",onChange:this.setValue,value:s})),r.a.createElement("div",null,r.a.createElement("span",null,"Reset Code: "),r.a.createElement("input",{name:"resetCode",onChange:this.setValue,value:o})),r.a.createElement("input",{type:"submit",value:"Update my password"}))):r.a.createElement(r.a.Fragment,null,r.a.createElement("form",{onSubmit:function(t){return e.sendResetCode(t)}},r.a.createElement("div",null,r.a.createElement("span",null,"Email: "),r.a.createElement("input",{name:"email",onChange:this.setValue,value:a})),r.a.createElement("input",{type:"submit",class:"animated_btn",value:"Send me a reset code"}))),r.a.createElement("button",{onClick:function(t){return e.props.setPage(t,f.signIn)}},"Back to sign in"))}}]),t}(n.Component),k=function(e){var t=e.page,a=e.setPage,n=e.setAuthToken,s=e.setUser;switch(t){case f.signUp:return r.a.createElement(w,{setPage:a,setAuthToken:n,setUser:s});case f.signIn:return r.a.createElement(y,{setPage:a,setAuthToken:n,setUser:s});case f.forgotPassword:return r.a.createElement(S,{setPage:a});default:return r.a.createElement(r.a.Fragment,null,"Error, invalid path reached")}},T=a(8),O=(a(26),function(e){var t=e.user,a=e.setPage,s=Object(n.useState)(null),o=Object(T.a)(s,2),i=o[0],c=o[1];function m(){return(m=Object(l.a)(u.a.mark(function e(){var a,n;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,fetch(b.base+b.handlers.myuserAvatar,{method:"GET",headers:new Headers({Authorization:localStorage.getItem("Authorization")})});case 2:if(!((a=e.sent).status>=300)){e.next=6;break}return c(t.photoURL),e.abrupt("return");case 6:return e.next=8,a.blob();case 8:n=e.sent,c(URL.createObjectURL(n));case 10:case"end":return e.stop()}},e)}))).apply(this,arguments)}return Object(n.useEffect)(function(){!function(){m.apply(this,arguments)}()},[]),r.a.createElement(r.a.Fragment,null,r.a.createElement("div",null,"Welcome to my application, ",t.firstName," ",t.lastName),i&&r.a.createElement("img",{className:"avatar",src:i,alt:"".concat(t.firstName,"'s avatar")}),r.a.createElement("div",null,r.a.createElement("button",{onClick:function(e){a(e,f.signedInUpdateName)}},"Update name")),r.a.createElement("div",null,r.a.createElement("button",{onClick:function(e){a(e,f.signedInUpdateAvatar)}},"Update avatar")),r.a.createElement("div",null,r.a.createElement("button",{onClick:function(e){a(e,f.signedInUploadImage)}},"Upload an image")),r.a.createElement("div",null,r.a.createElement("button",{onClick:function(e){a(e,f.mainPage)}},"Photo Management")))}),j=function(e){var t=e.setAuthToken,a=e.setUser,s=Object(n.useState)(""),o=Object(T.a)(s,2),i=o[0],c=o[1];return r.a.createElement(r.a.Fragment,null,r.a.createElement("button",{onClick:function(){var e=Object(l.a)(u.a.mark(function e(n){var r,s;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return n.preventDefault(),e.next=3,fetch(b.base+b.handlers.sessionsMine,{method:"DELETE",headers:new Headers({Authorization:localStorage.getItem("Authorization")})});case 3:if(!((r=e.sent).status>=300)){e.next=10;break}return e.next=7,r.text();case 7:return s=e.sent,c(s),e.abrupt("return");case 10:c(""),t(""),a(null);case 13:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}()},"Sign out"),i&&r.a.createElement("div",null,r.a.createElement(v,{error:i,setError:c})))},x=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).sendRequest=function(){var e=Object(l.a)(u.a.mark(function e(t){var n,r,s,o,i,l,c;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),n=a.state,r=n.firstName,s=n.lastName,o={firstName:r,lastName:s},e.next=5,fetch(b.base+b.handlers.myuser,{method:"PATCH",body:JSON.stringify(o),headers:new Headers({Authorization:localStorage.getItem("Authorization"),"Content-Type":"application/json"})});case 5:if(!((i=e.sent).status>=300)){e.next=13;break}return e.next=9,i.text();case 9:return l=e.sent,console.log(l),a.setError(l),e.abrupt("return");case 13:return alert("Name changed"),e.next=16,i.json();case 16:c=e.sent,a.props.setUser(c);case 18:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.setValue=function(e){a.setState(Object(g.a)({},e.target.name,e.target.value))},a.setError=function(e){a.setState({error:e})},a.state={firstName:"",lastName:"",error:""},a}return Object(h.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this.state,t=e.firstName,a=e.lastName,n=e.error;return r.a.createElement(r.a.Fragment,null,r.a.createElement(v,{error:n,setError:this.setError}),r.a.createElement("div",null,"Enter a new name"),r.a.createElement("form",{onSubmit:this.sendRequest},r.a.createElement("div",null,r.a.createElement("span",null,"First name: "),r.a.createElement("input",{name:"firstName",value:t,onChange:this.setValue})),r.a.createElement("div",null,r.a.createElement("span",null,"Last name: "),r.a.createElement("input",{name:"lastName",value:a,onChange:this.setValue})),r.a.createElement("input",{type:"submit",value:"Change name"})))}}]),t}(n.Component),I=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).sendRequest=function(){var e=Object(l.a)(u.a.mark(function e(t){var n,r,s,o;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),n=a.state.file,(r=new FormData).append("uploadfile",n),e.next=6,fetch(b.base+b.handlers.myuserAvatar,{method:"PUT",body:r,headers:new Headers({Authorization:localStorage.getItem("Authorization")})});case 6:if(!((s=e.sent).status>=300)){e.next=14;break}return e.next=10,s.text();case 10:return o=e.sent,console.log(o),a.setError(o),e.abrupt("return");case 14:alert("Avatar changed");case 15:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.handleFile=function(e){a.setState({file:e.target.files[0]})},a.setError=function(e){a.setState({error:e})},a.state={file:null,error:""},a}return Object(h.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this.state.error;return r.a.createElement(r.a.Fragment,null,r.a.createElement(v,{error:e,setError:this.setError}),r.a.createElement("div",null,"Give yourself a new avatar"),r.a.createElement("form",{onSubmit:this.sendRequest},r.a.createElement("div",null,r.a.createElement("span",null,"Upload new avatar "),r.a.createElement("input",{type:"file",onChange:this.handleFile})),r.a.createElement("input",{type:"submit",value:"Change avatar"})))}}]),t}(n.Component),N=a(12),U=function(e){function t(e){var a;return Object(c.a)(this,t),(a=Object(p.a)(this,Object(d.a)(t).call(this,e))).sendRequest=function(){var e=Object(l.a)(u.a.mark(function e(t){var n,r,s,o,i;return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return t.preventDefault(),n=a.state.file,r=n.name,(s=new FormData).append("file",n),e.next=7,fetch(b.base+b.handlers.upload,{method:"POST",body:s,headers:new Headers({Authorization:localStorage.getItem("Authorization")})});case 7:if(!((o=e.sent).status>=300)){e.next=16;break}return e.next=11,o.text();case 11:return i=e.sent,a.setError(i),e.abrupt("return");case 16:200===o.status&&fetch("https://api.xutiancheng.me/v1/photos",{method:"POST",body:JSON.stringify({url:"https://image-441.s3.amazonaws.com/"+r}),headers:new Headers({"Content-Type":"application/json",Authorization:localStorage.getItem("Authorization")})}).then(function(e){if(201==e.status){a.setShow(!0)}else{var t=e.text();a.setError(t)}});case 17:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),a.handleFile=function(e){a.setState({file:e.target.files[0]})},a.setError=function(e){a.setState({error:e})},a.setShow=function(e){a.setState({show:e})},a.state={file:null,error:"",show:!1},a}return Object(h.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=(t.firstName,t.lastName,t.error);return t.show?r.a.createElement(r.a.Fragment,null,r.a.createElement(v,{error:a,setError:this.setError}),r.a.createElement("div",null,"Upload a new image"),r.a.createElement("form",{target:"_blank",encType:"multipart/form-data",method:"post"},r.a.createElement("input",{type:"file",name:"file",onChange:function(t){return e.handleFile(t)}}),r.a.createElement("input",{type:"submit",value:"upload",onClick:function(t){return e.sendRequest(t)}})),r.a.createElement(N.a,{variant:"danger",onClose:function(){return e.setShow(!1)},dismissible:!0},r.a.createElement(N.a.Heading,null,"photo uploaded"),r.a.createElement("p",null,"to check the photo, return to main page, then go to the photo management page."))):r.a.createElement(r.a.Fragment,null,r.a.createElement(v,{error:a,setError:this.setError}),r.a.createElement("div",null,"Upload a new image"),r.a.createElement("form",{target:"_blank",encType:"multipart/form-data",method:"post"},r.a.createElement("input",{type:"file",name:"file",onChange:function(t){return e.handleFile(t)}}),r.a.createElement("input",{type:"submit",value:"upload",onClick:function(t){return e.sendRequest(t)}})))}}]),t}(n.Component),C=(a(33),function(e){var t=e.tags.map(function(t){return r.a.createElement(P,{key:t.id,tag:t,setTag:e.setTag,NotifyTagUpdate:e.NotifyTagUpdate,handleShow:e.handleShow})});return r.a.createElement("div",{style:{display:"flex",flexDirection:"row"}},t)}),A=function(e){return r.a.createElement("input",{type:"text",placeholder:"create new tag",onKeyPress:function(e){if("Enter"===e.key){var t=e.target.value,a={method:"post",headers:new Headers({"Content-Type":"application/json",Authorization:localStorage.getItem("Authorization")}),body:JSON.stringify({name:t})};fetch("https://api.xutiancheng.me/v1/tags",a)}}})},P=function(e){return r.a.createElement("div",{style:{display:"flex",flexWrap:"wrap",flexDirection:"row",margin:"15px"}},r.a.createElement("button",{onClick:function(t){e.setTag(e.tag.id)}},e.tag.name),r.a.createElement(H,{tagID:e.tag.id,NotifyTagUpdate:e.NotifyTagUpdate,handleShow:e.handleShow}))},D=function(e){return-1===e.tag?e.imgDataList.map(function(t){return r.a.createElement(F,{key:t.id,img:t,style:{display:"flex","flex-wrap":"wrap",flexDirection:"row"},tagIDTable:e.tagIDTable,IDTagTable:e.IDTagTable,NotifyTagUpdate:e.NotifyTagUpdate,handleShow:e.handleShow})}):e.imgDataList.map(function(t){var a=[];if(t.tags.forEach(function(e){a.push(e.id)}),a.includes(e.tag))return r.a.createElement(F,{key:t.id,img:t,style:{display:"flex","flex-wrap":"wrap",flexDirection:"row"},tagIDTable:e.tagIDTable,IDTagTable:e.IDTagTable,NotifyTagUpdate:e.NotifyTagUpdate,handleShow:e.handleShow})})},F=function(e){for(var t=[],a=0;a<e.img.tags.length;a++){var n=e.img.tags[a].id,s=e.tagIDTable[n];t.push(s)}var o=t.map(function(e,t){return r.a.createElement("p",{key:t,className:"font-size-0-8"},e)});return r.a.createElement("div",null,r.a.createElement("img",{src:e.img.url,alt:"wrong image url",style:{width:"150px",height:"150px"}}),o,r.a.createElement(z,{imgID:e.img.id,IDTagTable:e.IDTagTable,NotifyTagUpdate:e.NotifyTagUpdate,handleShow:e.handleShow}))},z=function(e){var t=e.imgID;return r.a.createElement("input",{type:"text",placeholder:"bind photo with old tag",onKeyPress:function(a){if("Enter"===a.key){var n=a.target.value,r=e.IDTagTable[n],s={method:"post",headers:new Headers({Authorization:localStorage.getItem("Authorization")})};fetch("https://api.xutiancheng.me/v1/photos/"+t+"/tag/"+r,s)}}})},H=function(e){var t=e.tagID;return r.a.createElement("input",{type:"text",placeholder:"userID",style:{width:"50px"},onKeyPress:function(e){if("Enter"===e.key){var a=e.target.value,n={method:"post",headers:new Headers({"Content-Type":"application/json",Authorization:localStorage.getItem("Authorization")}),body:JSON.stringify({id:a})};fetch("https://api.xutiancheng.me/v1/tags/"+t+"/members",n)}}})},R=function(e){var t=Object(n.useState)(!1),a=Object(T.a)(t,2),s=a[0],o=a[1],i=Object(n.useState)([]),u=Object(T.a)(i,2),l=u[0],c=u[1],m=Object(n.useState)([]),p=Object(T.a)(m,2),d=p[0],h=p[1],f=Object(n.useState)({}),g=Object(T.a)(f,2),E=g[0],b=(g[1],Object(n.useState)({})),v=Object(T.a)(b,2),w=v[0],y=(v[1],Object(n.useState)(-1)),S=Object(T.a)(y,2),k=S[0],O=S[1],j=function(e){O(e)},x=function(){o(!0)},I=Object(n.useState)(!1),U=Object(T.a)(I,2),P=U[0],F=U[1],z=function(){F(!P)};return Object(n.useEffect)(function(){fetch("https://api.xutiancheng.me/v1/photos",{method:"get",headers:new Headers({Authorization:localStorage.getItem("Authorization")})}).then(function(e){return e.json()}).then(function(e){c(e)})},[]),Object(n.useEffect)(function(){fetch("https://api.xutiancheng.me/v1/tags",{method:"get",headers:new Headers({Authorization:localStorage.getItem("Authorization")})}).then(function(e){return e.json()}).then(function(e){for(var t=0;t<e.length;t++){var a=e[t],n={};n[a.id]=a.name,Object.assign(E,n);var r={};r[a.name]=a.id,Object.assign(w,r)}h(e)})},[P,s]),s?r.a.createElement("div",null,r.a.createElement(A,{NotifyTagUpdate:z,handleShow:x}),r.a.createElement(C,{tags:d,setTag:j,NotifyTagUpdate:z,handleShow:x}),r.a.createElement(D,{tag:k,imgDataList:l,tagIDTable:E,IDTagTable:w,NotifyTagUpdate:z,handleShow:x}),r.a.createElement(N.a,{variant:"danger",onClose:function(){return o(!1)},dismissible:!0},r.a.createElement(N.a.Heading,null,"error for this request"),r.a.createElement("p",null,"make sure you have a valid auth session, the image upload isn't duplicate, also the tag is case sensitive. If it's not your problem then maybe there's an internal error."))):r.a.createElement("div",null,r.a.createElement(A,{NotifyTagUpdate:z}),r.a.createElement(C,{tags:d,setTag:j,NotifyTagUpdate:e.NotifyTagUpdate}),r.a.createElement(D,{tag:k,imgDataList:l,tagIDTable:E,IDTagTable:w,NotifyTagUpdate:z}))},M=function(e){var t=e.page,a=e.setPage,n=e.setAuthToken,s=e.setUser,o=e.user,i=r.a.createElement(r.a.Fragment,null),u=!0;switch(t){case f.signedInMain:i=r.a.createElement(O,{user:o,setPage:a});break;case f.signedInUpdateName:i=r.a.createElement(x,{user:o,setUser:s});break;case f.signedInUpdateAvatar:i=r.a.createElement(I,{user:o,setUser:s});break;case f.signedInUploadImage:i=r.a.createElement(U,{user:o,setUser:s});break;case f.mainPage:i=r.a.createElement(R,{user:o});break;default:i=r.a.createElement(r.a.Fragment,null,"Error, invalid path reached"),u=!1}return r.a.createElement(r.a.Fragment,null,i,u&&r.a.createElement("button",{onClick:function(e){return a(e,f.signedInMain)}},"Back to main"),r.a.createElement(j,{setUser:s,setAuthToken:n}))},L=(a(35),function(e){function t(){var e;return Object(c.a)(this,t),(e=Object(p.a)(this,Object(d.a)(t).call(this))).getCurrentUser=Object(l.a)(u.a.mark(function t(){var a,n;return u.a.wrap(function(t){for(;;)switch(t.prev=t.next){case 0:if(e.state.authToken){t.next=2;break}return t.abrupt("return");case 2:return t.next=4,fetch(b.base+b.handlers.myuser,{headers:new Headers({Authorization:e.state.authToken})});case 4:if(!((a=t.sent).status>=300)){t.next=11;break}return alert("Unable to verify login. Logging out..."),localStorage.setItem("Authorization",""),e.setAuthToken(""),e.setUser(null),t.abrupt("return");case 11:return t.next=13,a.json();case 13:n=t.sent,e.setUser(n);case 15:case"end":return t.stop()}},t)})),e.setPageToSignIn=function(t){t.preventDefault(),e.setState({page:f.signIn})},e.setPageToSignUp=function(t){t.preventDefault(),e.setState({page:f.signUp})},e.setPage=function(t,a){t.preventDefault(),e.setState({page:a})},e.setAuthToken=function(t){e.setState({authToken:t,page:""===t?f.signIn:f.signedInMain})},e.setUser=function(t){e.setState({user:t})},e.state={page:localStorage.getItem("Authorization")?f.signedInMain:f.signIn,authToken:localStorage.getItem("Authorization")||null,user:null},e.getCurrentUser(),e}return Object(h.a)(t,e),Object(m.a)(t,[{key:"render",value:function(){var e=this.state,t=e.page,a=e.user;return r.a.createElement("div",null,r.a.createElement("h1",{id:"proj-title"},"441 Photo Management"),a?r.a.createElement(M,{page:t,setPage:this.setPage,setAuthToken:this.setAuthToken,user:a,setUser:this.setUser}):r.a.createElement(k,{page:t,setPage:this.setPage,setAuthToken:this.setAuthToken,setUser:this.setUser}))}}]),t}(n.Component));Boolean("localhost"===window.location.hostname||"[::1]"===window.location.hostname||window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/));o.a.render(r.a.createElement(L,null),document.getElementById("root")),"serviceWorker"in navigator&&navigator.serviceWorker.ready.then(function(e){e.unregister()})}},[[16,2,1]]]);
//# sourceMappingURL=main.2f13e194.chunk.js.map