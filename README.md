# Building B2C Web Application In Go

使用go语言搭建一个简单的电商网站平台

## 需求与设计

### 功能介绍

1. 用户可以完成注册、登录然后登录后对商品的浏览。
2. 登录之后，用户可以对相关商品进行选购并添加到购物车。
3. 用户可以对购物车里面的商品进行增加、减少、删除操作。
4. 用户可对购物车商品进行结算操作。

### 技术选型

本项目涉及使用到Golang+MongoDB+Redis对象模型库，详细介绍如下：

Go：go是一门静态强类型语言, 使用它主要是并发性好

MongoDB：MongoDB是一个基于分布式文件存储的一个高性能，开源，无模式的文档型数据库，数据以BSON文档的格式存储在磁盘上。

mgo.v2：一个针对MongoDB操作的对象模型库，封装了MongoDB对文档的的一些增删改查等常用方法。

### 结构划分

项目主要分为以下几大模块：注册模块，登录模块，商品模块、购物车模块、结算模块。

1. 用户注册模块：填写用户名、密码、确认密码后，实现成功注册，然后进行登录。
2. 用户登录模块：填写已注册的用户名称，填写正确的密码，进入商品展示页面。
3. 商品模块：用户选择相关产品加入购物车。
4. 购物车模块：对相关商品进行增加、减少、删除操作。
5. 结算模块：对购物车内已选择商品进行结算。

模块结构如下图所示：

```
                    ----------------
                    | 用户购物模块 |
                    ----------------
                            |
     ----------------------------------------------------
     |              |                  |                |
------------    ------------    --------------    --------------
| 登录注册 |    | 商品展示 |    | 数据库设计 |    | 购物车结算 |
------------    ------------    --------------    --------------

```


### 流程设计

此流程图显示用户可以进行登录和注册操作，如果用户已经注册，则可以直接登录，若未注册则必须先注册成功后才能进行登录，登录成功后可以进入商品页浏览商品，也可以选择相关商品并可加入购物车，在购物车页面内可以对购物车商品进行相关操作，最后结选择相关商品进行结算。

其流程如下图所示：

```
--------    ------------    ------------    --------------
| 登录 | >> | 浏览商品 | >> | 选择商品 | >> | 加入购物车 |
--------    ------------    ------------    --------------
   ^                                               X
--------                        --------    --------------
| 注册 |                        | 结算 | << | 查看购物车 |
--------                        --------    --------------
```


## 数据库设计

### 数据库介绍
数据库的使用呢，我们这里选择了MongoDB。

MongoDB的简单介绍如下：

MongoDB是一个开源的NoSQL数据库，相比MySQL那样的关系型数据库，它更显得轻巧、灵活， 非常适合在数据规模很大、事务性不强的场合下使用。

同时它也是一个对象数据库，没有表、行等概念，也没有固定的模式和结构，所有的数据以Document(以下简称文档)的形式存储(Document，就是一个关联数组式的对象，它的内部由属性组成，一个属性对应的值可能是一个数、字符串、日期、数组，甚至是一个嵌套的文档。)

我们一共要创建三个集合，分别是user(用户)集合、commodity(商品)集合、cart(购物车)集合，后面我们会一一介绍。

### user集合属性值展示

关于user集合，我们设计的属性有name(用户名)、password(密码)， 如下所示：

```json
name: { type: string, required: true },
password: { type: string, required: true },
```

### commodity集合属性值展示

关于commodity集合，我们设计的属性有name(商品名称)、price(商品价格)、imgSrc(商品展示图片路径)， 如下所示：

```json
name: { type: string, required: true },
price: { type: string, required: true },
imagSrc: { type: string, required: false },
```

### carts集合属性值展示

关于cart集合，我们设计的属性有uId(用户ID)、cId(商品ID)、cName(商品名称)、cPrice(商品价格)、cImgSrc(商品展示图片路径)、cQuantity(商品数量)、cStatus(商品结算状态，未结算为false,已结算为true)， 如下所示：


```json
uid: { type: string, required: true },
cid: { type: string, required: true },
cname: { type: string, required: false },
cprice: { type: string, required: true },
cimgsrc: { type: string, required: true },
cquantity: { type: string, required: false },
cstatus: { type: string, required: true },
```

## 模块功能设计

### 注册模块功能设计介绍

功能：本模块主要用于新用户注册，用户通过表单提供用户名和密码信息，系统根据用户提供的注册信息对用户进行具体操作。

输入操作：用户名、密码、确认密码。

对应处理：

1. 输入注册信息：在页面提供的表单处输入用户的用户名和密码信息，点击“注册”按钮提交表单数据信息。已注册用户，可点击“登录”按钮，进入登录页面。

2. 用户注册身份验证：连接数据库，以输入的“用户名”数据为查询条件来查看输入用户名是否已存在，如果用户名未注册，则提示注册成功并转到登录页进行登录，如果用户已注册，则给出用户已存在提示并重新注册。

注册界面
用户注册界面草图如下：



### 登录模块设计与实现

功能：本模块主要用于对用户身份进行鉴别。用户通过表单提供用户名和密码信息，系统根据用户提供的登录信息对用户进行查询鉴别。 如果身份合法，则用户可进入商品页面。

输入操作：用户名、密码。

对应处理：

1. 输入用户的登录信息：在页面提供的表单处输入用户的用户名和密码信息，点击“登录”按钮提交表单数据信息。 也可点击“注册”按钮，注册新用户。

2. 用户身份进行验证：连接数据库，打开user集合，检验用户登录信息。以输入的“用户名”数据为查询条件来查看用户是否存在，如果存在，继续检验其输入密码是否正确，密码和用户名都正确，则进入商品展示页面；如果用户名不存在或密码不正确，则给出相应的提示。

登录界面草图
用户登录界面草图如下：



### 商品页面介绍

商品页呢，我们主要要来展示商品，用户登录成功之后将会跳入商品页，可对所有商品进行查看， 每个商品信息的内容包括：商品名称、商品图片、商品价格，用户可浏览也可将其加入购物车。

页面草图如下所示：



### 商品页面介绍

关于购物车页面，主要展示用户已购买的商品，包括商品的信息、价格、数量，当然用户可以对其中商品进行增加、减少、删除操作，最后，用户可选择对其中商品进行结算，选择结算后，会提示相应的付款金额。

页面草图如下所示：



## 注册
本章课程是在NodeJS、Express、Mongoose、MongoDB的基础之上所建立的，所以，学习者请先了解以上基础，方可更好的完成多本章节的学习和掌握。

此次学习将带领大家来做一个用户注册、登录、商品的购物车加入和结算的简单实战项目，通过逐步的学习来更好的了解和掌握对前面课程的学习和巩固。

通过本章课程循序渐进的学习，将对前面的所学知识进行温习和巩固，达到温故而知新的效果。

简单启动
首先呢，我们先新建一个项目工程目录，然后在目录下创建启动文件app.js，开始我们的第一步。

这里我们会用到Express框架来实现相关功能，所以，需要先安装它，具体安装方法这里就不在做介绍了。

在启动文件添加如下内容，来测试Express框架是否引用成功。

var express = require('express');
var app = express();
app.get('/', function (req, res) {
   res.send('Hello World!');
});

app.listen(80);
浏览器查看结果，如收到响应信息则表明我们项目的第一步已经成功搞定。

创建目录
项目已经启动成功，下面我们开始创建相关目录，用于存储不同的文件。

说明：右侧栏"文件管理"中可添加相应目录和文件。

1.public目录：存放静态文件。

2.routes目录：存放路由文件。

3.views目录： 存放页面文件。

4.common目录：存放公共文件。

5.public目录(可不选)，新建js、css、img三个目录用以存储js、css、img相关文件。

这里我们内置了一些js、css文件来实现简单页面样式和操作，在页面视图中直接使用即可，引用方法如下：

<link href="example/css/bootstrap.min.css" rel="stylesheet" >
<script src="example/js/jquery.min.js" type="text/javascript"></script>
<script src="example/js/bootstrap.min.js" type="text/javascript"></script>
添加文件
有了目录，我们开始添加文件，先来添加一个登录页面register.html，便于管理和开发我们统一把视图页面放到views目录下。

views目录，添加register.html注册视图页，如下简单效果图：



有了视图页面，我们就可以访问它了，那要如何访问呢，这里就要使用到ejs模板了，安装方法不在阐述，直接如下引用：

app.set( 'view engine', 'html' );
app.engine( '.html', require( 'ejs' ).__express );
使用engine函数注册模板引擎并指定处理后缀名为html的文件。

设定视图存放的目录

app.set('views', require('path').join(__dirname, 'views'));
如果是在本地项目中，我们还要指定本地静态资源访问的路径,如下设置：

app.use(express.static(require('path').join(__dirname, 'public')));
访问注册页
有了视图页面，下面我们就开始访问它，app.js文件部分内容，引入相关模块资源，然后简单访问如下：

app.get('/', function (req, res) {
    res.render('register');
});
app.listen(80);
本节启动访问80端口，如成功看到注册页面则表示项目已经运行成功，如未看到，查看相关错误信息，是否缺少相关模块，安装和引用即可。

定义Schema
首先在common目录内添加models.js文件用来保存各个集合的Schema文件(集合属性)，也便于我们查看和访问，具体内容如下所示：

module.exports = {
    user: {
        name: { type: String, required: true },
        password: { type: String, required: true },
        gender: { type: Boolean, default: true }
    }
};
有了集合的Schema文件，如何访问呢，接着我们会介绍如何使用Model模型操作这些属性。

创建公共方法
还是common目录，我们在新建一个公共方法 —— dbHelper.js文件，来操作这些Schema，因为后面还会涉及此问题，所以我们写成一个公共的方法，dbHelper文件内容如下：

var mongoose = require('mongoose'),
    Schema = mongoose.Schema,
    models = require('./models');

for(var m in models) {
    mongoose.model(m, new Schema(models[m]));
}
module.exports = {
    getModel: function (type) {
        return _getModel(type);
    }
};
var _getModel = function (type) {
    return mongoose.model(type);
};
如上所示我们通过getModel可获取集合的Model模型就可以对数据库有实质性的操作了

关于Model，简单介绍：由Schema构造生成的模型，具有数据库操作的行为。

添加函数
关于dbHelper.js文件里方法的访问很简单，如下所示：

global.dbHelper = require( './common/dbHelper' );
这里我们使用globa来定义全局变量dbHelper，那么dbHelper就可以在任何模块内调用了。

然后我们就开始修改register视图页面，添加单击事件，例如：

<input type="button"  onclick="register()" value="注 册" />
对应register()函数，大致如下：

function register(){
   //通过serialize()方法进行序列化表单值，创建文本字符串。
   var data = $("form").serialize();
   //例如："username=张三&password=12345"
   $.ajax({
       url:'/register',
       type:'POST',
       data:data,
       success:function(data,status){
           if(status == 'success'){
               location.href='register';
           }
       },
       error:function(res,err){
           location.href='register';
       }
   });
}
添加路由
这里我们需要新建一个文件register.js，专门用来处理来之register页面的post请求， 在后面的学习中还会有多个不同处理文件，万所以我们统一管理在routes目录下，在实际开发中我们可能需要针对不同文件请求给出相应文件的处理，所以我们就做分开处理。

这里贴出register.js文件处理get和post请求的相关代码以供参考，如下：

// app：express对象
module.exports = function ( app ) {
  app.get('/register', function(req, res) {
      res.render('register');
  });
  app.post('/register', function (req, res) {
     var User = global.dbHelper.getModel('user'),
     uname = req.body.uname;
     User.findOne({name: uname}, function (error, doc) {
       if (doc) {
            req.session.error = '用户名已存在！';
            res.send(500);
        } else {
            User.create({
                name: uname,
                password: req.body.upwd
            }, function (error, doc) {
                if (error) {
                    res.send(500);
                } else {
                    req.session.error = '用户名创建成功！';
                    res.send(200);
                }
            });
        }
    });
  });
}
模块的加载和引用
register的post请求处理中，我们使用了session(express-session模块)还有处理post请求数据的body属性(body-parser和multer模块)，需先安装他们，然后引用即可，如下参考：

//引用模块
var bodyParser = require('body-parser');
var multer = require('multer');
var session = require('express-session');

//调用中间件使用
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(multer());
后面我们还会再次添加多个路由记录，所以便于管理和访问，我们可以把他们统一放到一起，比如routes目录下新建index.js文件专门用来存放添加的文件，代码如下：

module.exports = function ( app ) {
    require('./register')(app);
};
那么我们在app.js文件中直接引用index.js文件就可以访问这些文件了，如下所示：

require('./routes')(app); //app:express对象。;
中间件传递信息
这里我们就一步到位，在register的post请求处理中我们使用了express-session模块来保存相关信息，这里我们就使用中间件来传递这些提示信息，中间件内容如下所示：

app.use(function(req, res, next){
    res.locals.user = req.session.user; //保存用户信息
    var err = req.session.error;  //保存结果响应信息
    res.locals.message = '';  // 保存html标签
    if (err) res.locals.message = '<div class="alert alert-danger" style="margin-bottom: 20px;color:red;">' + err + '</div>';
    next();
});
这里注意中间件的安放位置，还有我们设置了变量message并为其简单添加了样式，这里我们在register视图里就用它来作为操作结果的信息提示，直接添加<%- message %>到视图第一个div内即可。

关于注册我们基本已经准备就绪，开始打开连接数据库并设置用户过期时间(注意执行顺序，应放置在首个中间件位置)，app.js条件内容如下：

mongoose.connect("mongodb://127.0.0.1:27017/test");

app.use(session({
    secret:'secret',
    cookie:{
        maxAge:1000*60*30
    }
}));
到这里，注册功能已经完毕，在用户注册的信息录入中，我们没有进行相关的为空、两次密码的不匹配等等验证等等(可自行添加)，赶紧注册试试吧，本地的话可以通过MongoVUE(可视化客户端)来查看数据是否成功写入数据库。

添加视图
前面我们已经实现了注册功能，用户可以成功注册，接着我们就开始让用户登录了，此节我们就实现用户的登录功能，并且登录成功后跳转商品页面查看商品。

首先，我们还是在views目录下添加登录视图页面 —— login.html，效果图如下：



访问视图
有了登录页面，那么注册页面(register)的登录按钮添加指向登陆页面的链接，相应的登陆页的注册按钮也是如此。

这里我们还是添加一个相对应的文件用来处理login页面的请求，routes目录下新建名为login.js的文件，先来增加一个处理get请求的方法，代码参考如下：

module.exports = function ( app ) {
    app.get('/login',function(req,res){
        res.render('login');
    });
}
和register文件一样添加到index.js中，如下：

require('./login')(app);
register视图页的register()函数的回调中，当注册成功时我们就可以跳转到登陆页面了，如下：

  location.href='login';
好，赶紧试试登陆、注册按钮能否成功跳转吧！

实现登陆
我们为登陆按钮增加单击事件和对应函数login()，参考如下：

function login(){
  //通过serialize()方法获取form表单内输入框用户名、密码的值
  //如："username=张三&password=12345"
  var data = $("form").serialize();
  $.ajax({
      url:'/login',
      type:'POST',
      data:data,
      success:function(data,status){
          if(status == 'success'){
              location.href='home';
          }
      },
      error:function(data,status){
           if(status == "error"){
              location.href='login'
           }
       }
   });
}
在相应的login.js文件中，我们还得添加相对应的post请求处理方法，接着我们会讲述。

登陆处理
关于login视图页的post请求处理，我们需要判断用户所输入用户名是否存在，密码是否正确，并使用变量保存相应提示信息，当用户名和密码全部正确时，则返回成功并保存用户的个人信息，用作来判断用户的登陆状态，具体可参考register视图页的post请求。

login.js文件的post请求处理代码，参考如下：

app.post('/login', function (req, res) {
  var User = global.dbHelper.getModel('user'),uname = req.body.uname;
  User.findOne({name: uname}, function (error, doc) {
           if (用户不存在) {
                req.session.error = '用户名不存在！';
                res.send(404);
            } else if(用户存在，密码错误) {
                   req.session.error = "密码错误!";
                   res.send(404);
               }else{ //用户名、密码正确
                   req.session.user=doc;
                   res.send(200);
               }
            });
    });
}
还记得我们登陆的本地变量message嘛，用来保存html标签并包含相应提示信息，这里在登陆页面我们也可以使用，用法：<%- message %>，指定到相应位置即可。

商品页视图
用户登录成功之后则跳转至home视图页面(商品主页)，就可以进行对商品的浏览和选择了。

还是views目录，添加home商品视图页，如下简单效果图：



用户成功登录之后跳转至home页，这里我们还是做分开处理，routes目录下新建home.js文件用来处理来之home也的get请求。

这里我们假设如果用户未登录将不能查看商品主页，所以，在请求处理中我们还需要判断用户的登陆状态，这个可以使用我们在登录处理时所保存的用户个人信息。

关于商品页的视图展示我们只需要有其名称、价格、图片，这里使用ejs模板循环展示，可参考如下方式：

注：Commodity：商品集合所有数据，内置图片路径为“/example/img”



请求处理
在home的get请求处理中，我们需要首先判断用户的登陆状态，只有用户登录了方可跳转到商品页，如果为登陆呢则跳转到登录页，而且在进入商品页的时候并传入Commodity集合的所有数据数据在页面展示。

首先呢，在models.js文件中定义Commodity集合的Schema属性，共包括商品名称、商品价格、商品图片，这里简单定义如下：

commodity: {
   name: String,
   price: Number,
   imgSrc: String
}
routes目录下添加home.js文件(index.js文件中引用)。

具体处理方式可参考如下代码：

module.exports = function ( app ) {
    app.get('/home', function (req, res) {
        if(req.session.user){
            var Commodity = global.dbHelper.getModel('commodity');
            Commodity.find({}, function (error, docs) {
               //将Commoditys变量传入home模板
               res.render('home',{Commoditys:docs});
            });
        }else{
            req.session.error = "请先登录"
            res.redirect('/login');
        }
    });
}
添加商品
添加商品，views目录下添加addcommodity视图页面用来对商品的添加，这里简单样式参考如下：



相对应的addcommodity函数参考代码如下：

//imgSrc表示图片路径)，这里内置了5张图片，格式为：xmsz-X.jpg(X为1-5数字)。
var data = $("form").serialize()+"&imgSrc=" + "xmsz-"+Math.floor(Math.random()*5+1)+".jpg";
$.ajax({
     url:'/addcommodity',
     type:'POST',
     data:data,
     success:function(data,status){
         if(status == 'success'){
            alert('添加成功！')
         }
     },
     error:function(data,err){
         alert('添加失败！')
     }
});
商品添加处理
这里我们就直接在home.js文件中添加保存商品的处理方法，如下：

    app.get('/addcommodity', function(req, res) {
        res.render('addcommodity');
    });
    app.post('/addcommodity', function (req, res) {
        var Commodity = global.dbHelper.getModel('commodity');
        Commodity.create({
            name: req.body.name,
            price: req.body.price,
            imgSrc: req.body.imgSrc
        }, function (error, doc) {
            if (doc) {
                res.send(200);
            }else{
                res.send(404);
            }
        });
    });
到这里关于商品页的展示和添加就完成了，在下一节里我们将实现商品页商品加入购物车并结算的功能，继续加油吧！

添加商品链接
上节课程里我们已经实现了商品的添加和展示，接下来我们开始进行对商品的操作——加入购物车。

首先，商品页的加入购物车按钮、购物车查看按钮添加链接，如下所示：

[">加入购物车](/addToCart/<%=Commoditys[i]._id%)
我们先定义购物车(cart)集合的Schema属性，包含：uId(用户ID)、cId(商品ID)、cName(商品名称)、cPrice(商品价格)、cImgSrc(商品展示图片路径)、cQuantity(商品数量)、cStatus(商品结算状态，默认为false)，参考如下：

 cart:{
        uId: { type: String },
        cId: { type: String },
        cName: { type: String },
        cPrice: { type: String },
        cImgSrc: { type:String } ,
        cQuantity: { type: Number },
        cStatus : { type: Boolean, default: false  }
    }
以上属性定义我们还是统一放到models.js文件中以方便管理和操作。

接着views目录下添加购物车(cart.html)视图页面，参考如下：



购物车页面的展示实现可参考如下贴图：



访问视图
接着在routes目录添加cart.js文件(index.js文件中引用)，来定义路由处理规则。

比如cart路径的处理，如下：

1.1 检测登陆用户是否过期，已过期则跳转login页同时给出提示信息。

1.2 通过global.dbHelper.getModel方法获取cart模型。
1.3 通过find查询用户所有商品并传入模板，条件：用户ID，结算状态。
1.4 贴出部分代码：

app.get('/cart', function(req, res) {
 ......
Cart.find({"uId":用户ID,"cStatus":false},function (error, docs) {
   res.render('cart',{carts:docs});
});
});
加入购物车按钮链接所对应路径的处理，如：/addToCart/:商品id

app.get("/addToCart/:id", function(req, res) {...});
2.1 通过req.params.id获取商品ID号并检测登陆用户状态。

2.2 通过global.dbHelper.getModel方法获取购物车(cart)、商品(commodity)模型。
2.3 商品已存在则进行更新操作，贴出部分参考代码：

Cart.update({"uId":用户ID,"cId":商品ID},{$set:{cQuantity:已有数量+1}}
2.4商品未存在则进行添加操作，贴出部分代码：

Commodity.findOne({"_id": 商品ID}, function (error, doc) {
    Cart.create({
        uId: 用户ID,
        cId: 商品ID,
        cName: doc.name,
        cPrice: doc.price,
        cImgSrc: doc.imgSrc,
        cQuantity : 1
    },function(error,doc){
        if(doc){
            res.redirect('/home');
        }
    });
});
商品的单选和全选
在购物车页面展示还是使用ejs模板(具体实现可参考贴图)来实现用户购物车商品的展示，现期间我们用到了checkbox，这里我们就来简单的实现商品的全选和单选。

首先，简单介绍实现步骤：

用户选中全选按钮时，列表内选框全部变为勾选状态。
用户取消全选按钮选中状态时，列表内选框全部取消勾选状态。
列表内选框为全部勾选状态时，全选按钮为选中状态，反之不勾选。
全选按钮未选中情况下，当列表内按钮全部选中则全选按钮也要被选中。
对于checkbox按钮如下图所示：



注：data-index属性存储索引值用来获取商品数量，data=id存储商品ID，data-price存储商品价格。

实现全选：为全选框(如id为CheckAll)添加单击事件，这里我们使用prop方法设置checkbox状态、is方法判断checkbox状态，如下参考示例：

 $('#CheckAll').click(function(){
       var self = $(this);
       $('input[name="chkItem"]').each(function(){
          $(this).prop("checked",self.is(':checked'));
       });
 });
实现单选：单选框(如name为chkItem)添加单击事件，这里我们使用prop方法设置checkbox状态、not方法判断checkbox状态，如下参考示例：

$('input[name="chkItem"]:checkbox').click(function(){
   var isCheck = $('input[name="chkItem"]:not(:checked)').length?false:true;
   $('#CheckAll').prop("checked",isCheck);
});
到这里我们就简单实现了按钮的全选和单选，以上方法可供参考也可以自行思路去实现。

商品的增加和减少
前面我们实现了对于商品的单选和全选功能，下面我们就开始实现商品数量的加减。

对于+、-按钮如下图所示：



简单实现步骤如下：

定义data-type属性用于区分+、-按钮。
a标签添加单击事件。并同时改变显示框的值。
使用siblings()方法改变同级标签的数量值。
具体实现可参考如下方法：

$('.li-quantity a').click(function(){
    var self = $(this);
    var type = self.attr('data-type'),
        num = parseFloat(self.siblings('input').val());
    if(type == 'add'){
          num += 1;
     }else if(type == 'subtr'){
          if(num > 1){
             num -= 1;
          }else{
             return false;
          }
     }
    self.siblings('input').val(num);
});
到这里对于商品数量的加、减也实现了，接着我们开始实现用户选中商品的总价格计算功能。

计算总金额
商品的状态选择和数量的加减功能我们已经实现了，选择我们就开始实现当用户选择相关商品时总金额的计算功能。

实现步骤如下：

定义公共方法，用于用户不同加、减、全选等操作时皆可调用。
循环列表内所有被选中的checkbox(假设name为chkItem)。
使用checkbox的相应属性值，来获取价格和数量。
定义全局变量(假设为sum)存储总金额并赋值显示。
具体方法实现可参考如下内容：

   sum = 0;
   $('input[name="chkItem"]:checked').each(function(){
       var self = $(this),
        price = self.attr('data-price'),
        index = self.attr('data-index');
        var quantity = $('#Q'+index).val();
        sum +=(parseFloat(price)*parseFloat(quantity));
    });
    $("金额标签").html('￥'+ sum +'.00');
到这里，选中商品的价格总计功能也就简单实现了。

商品结算
关于结算功能，这里就做简单介绍，当用户点击结算按钮时，计算出被选商品的总金额，给予alert出总金额，数据库则更新相应商品的结算状态和数量，当然你也可以在购物车集合中添加一个用户消费金额的属性来保存所消费的金额总数更好，这里就简单处理了。

具体功能实现这里简单介绍如下：

定义公共方法，用于用户不同加、减、全选等操作时皆可调用。
循环列表内所有被选中的checkbox(假设name为chkItem)。
使用checkbox的相应属性值，来获取价格和数量。
定义全局变量(假设为sum)存储总金额并赋值显示。
这里贴出部分实现代码，如下：

//结算方法内容
$('input[name="chkItem"]:checked').each(function(){
     var self = $(this),
         //通过data-index属性，获取索引值
         index = self.attr('data-index'),
        //通过data-id属性，获取对应ID
         ID= self.attr('data-id');
     var 数量= $('#Q'+index).val();
     var data = { "cid": ID, "cnum":数量};
     //发送数据请求
     ...
});
alert('所付金额为：￥'+sum);
location.href = "cart";
cart.js文件添加对应路由处理如下：

app.post("路径",function(req,res){
   var Cart = global.dbHelper.getModel('cart');
   Cart.update({"_id":req.body.ID},{$set : { cQuantity : req.body.数量,cStatus:true }},function(error,doc){
...
商品删除
关于商品的删除功能就简单多了，我们只需获取其ID即可实现对于购物车内商品的删除操作

在购物车商品的展示功能实现时，我们就可以获取其ID，如下参考：

[" >删除](/delFromCart/<%=carts[i]._id%)
cart.js文件，添加对应路径处理方法，这里简单实现参考如下：

app.get("/delFromCart/:id", function(req, res) {
   //req.params.id 获取ID号
   var Cart = global.dbHelper.getModel('cart');
   Cart.remove({"_id":req.params.id},function(error,doc){
       //成功返回1  失败返回0
       if(doc > 0){
           res.redirect('/cart');
       }
   });
});
好，到这里所有功能都已经实现了，赶快好好回味回味所学所悟吧！

附：项目源码下载地址：https://github.com/hubwiz/example-node