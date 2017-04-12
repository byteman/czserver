var app = angular.module('myApp')

app.filter('myfilter', function() {  
   return function(input, param1) {
       if(input==1) return "单点重量"
	 else if(input==3) return "总重量"
	 else if(input==2) return "排水重量"
	 return input
   };  
 });
 
 app.filter('range', function(){
    return function(n) {
      var res = [];
      for (var i = 0; i < n; i++) {
        res.push(i+1);
      }
      return res;
	}
});

 app.filter('version', function(){
    return function(n) {
      var ver = "ver{0}.{1}.{2}".format([(n<<16)&0xff,(n<<8)&0xff,n&0xff]);;
	  
      return ver;
	}
});
