var app = angular.module('myApp')


app.controller('AppController', function($scope,$http) {
	$scope.isUserAuth=false;
	$scope.role=1

    $scope.user={
        name:"admin",
        pwd:""
    }
    console.log($scope.role)
    $scope.login = function(){
        console.log($scope.user)

        $http.post("/login",$scope.user).success(
            function(data,status,headers,config)
            {
                console.log(data)
                if(data.result==0)
                {
                    console.log($scope.isUserAuth)
                    $scope.isUserAuth=true
                    $scope.role=data.role
                    console.log($scope.isUserAuth)
                }
            }
        ).error(
            function(data,status,headers,config)
            {

            }
        );
    }

});
