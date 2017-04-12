var app = angular.module('myApp')


app.controller('OnLineCtrl', function($scope, $http,$interval) {
		
	$scope.admin=false
	$scope.getData=function(){
		$http.get("/online").success(
			function (response) {
				$scope.data = response;
				console.log(response)
			}
		);
		
	};
	$scope.getData();
	var toDo = function () {
		$scope.getData();
	};
	var timer = $interval(toDo, 1000);
	$scope.$on("$destroy",
		function( event ) {
			console.log("controller destroy,cancel timer")
			$interval.cancel( timer );

		}
    );
	$scope.save = function() {
	var par = $scope.dev;
    $http({
        method  : 'POST',
        url     : '/params',
        data    : par,  // pass in data as strings
        headers : { 'Content-Type': 'application/x-www-form-urlencoded' }  // set the headers so angular passing info as form data (not request payload)
    }).success(function(data) {
			console.log(data.result);

			if (data.result) {
				// if not successful, bind errors to error variables
				
			} else {
				// if successful, bind success message to message
				$scope.message = data.message;
				$('#myModal').modal("hide");
			}
		});
	};
	$scope.getparam = function(index){
		
			$('#myModal').modal();
			$scope.dev = $scope.data[index];
		
	
	}
	$scope.login = function(user,password){
		var ok = false;
		if(user=="admin"){
			if(password="123321")
			{
				$scope.admin=true;
				ok = true;
				
			}else if(password="123456"){
				ok = true;
			}
		}
		return ok;
	}

});
