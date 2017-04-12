var app = angular.module('myApp')


app.controller('LoginController',['$scope', '$http', function($scope, $http,$location) {
	$scope.user={
		user:"admin",
		pwd:""
	}
    $scope.isUserAuth=false
   $scope.login = function(){
	    $http.post("/login",$scope.user).success(
				function(data,status,headers,config)
				{
					if(data.result==0)
					{
                        $scope.isUserAuth=true
						$location.path("/static/html/online.html");
					}
				}
		).error(
				function(data,status,headers,config)
				{
					
				}
		);
   }
  
}]);